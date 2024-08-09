package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	metermetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type Response struct {
	SIP       bool `json:"sip"`
	CAPILLARY bool `json:"capillary"`
	GHL       bool `json:"ghl"`
}

func getHub() (Response, error) {

	var result Response
	url := "http://localhost:8080/hub"

	resp, err := http.Get(url)
	if err != nil {
		return result, fmt.Errorf("Error performing GET request:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("Error reading response body: ", err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, fmt.Errorf("Error parsing JSON: ", err)
	}
	return result, nil
}

func newResource() (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("hub-proactive-monitoring"),
			semconv.ServiceVersion("0.0.1"),
		),
	)
}

func newMeterProvider(res *resource.Resource) (*metric.MeterProvider, error) {
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(
			metric.NewPeriodicReader(
				metricExporter,
				metric.WithInterval(time.Second),
			),
		),
	)
	return meterProvider, nil
}

func main() {
	res, err := newResource()
	if err != nil {
		panic(err)
	}

	meterProvider, err := newMeterProvider(res)
	if err != nil {
		panic(err)
	}

	defer func() {
		err := meterProvider.Shutdown(context.Background())
		if err != nil {
			log.Println(err)
		}
	}()

	otel.SetMeterProvider(meterProvider)

	var meter = otel.Meter("hub.connection")

	sip, err := meter.Int64ObservableGauge(
		"sip",
		metermetric.WithDescription("Hub connection status with SIP"),
	)
	if err != nil {
		panic(err)
	}

	meter.RegisterCallback(
		func(_ context.Context, o metermetric.Observer) error {
			response, err := getHub()
			if err != nil {
				panic(err)
			}

			var sipval int64
			if response.SIP {
				sipval = 1
			} else {
				sipval = 0
			}
			o.ObserveInt64(sip, sipval)
			return nil
		},
		sip,
	)
}
