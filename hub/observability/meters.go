package main

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const name = "hub-monitoring"

var meter = otel.Meter(name)

func meterInit() error {
	sipConn, err := meter.Int64ObservableGauge(
		"hub.connection.sip",
		metric.WithDescription("HUB connection status with SIP"),
		metric.WithUnit("{bool}"),
	)
	capillaryConn, err := meter.Int64ObservableGauge(
		"hub.connection.capillary",
		metric.WithDescription("HUB connection status with Capillary"),
		metric.WithUnit("{bool}"),
	)
	ghlConn, err := meter.Int64ObservableGauge(
		"hub.connection.ghl",
		metric.WithDescription("HUB connection status with GHL"),
		metric.WithUnit("{bool}"),
	)
	apiConn, err := meter.Int64ObservableGauge(
		"hub.shellapi.health",
		metric.WithDescription("Shell API health status"),
		metric.WithUnit("{count}"),
	)

	if err != nil {
		return err
	}

	_, err = meter.RegisterCallback(
		func(_ context.Context, o metric.Observer) error {
			resp, err := getHub()
			if err != nil {
				return err
			}

			sipConnVal := 0
			capillaryConnVal := 0
			ghlConnVal := 0

			if resp.SIP {
				sipConnVal = 1
			}

			if resp.CAPILLARY {
				capillaryConnVal = 1
			}

			if resp.GHL {
				ghlConnVal = 1
			}

			o.ObserveInt64(sipConn, int64(sipConnVal))
			o.ObserveInt64(capillaryConn, int64(capillaryConnVal))
			o.ObserveInt64(ghlConn, int64(ghlConnVal))

			attributes := []attribute.KeyValue{
				attribute.Bool("sip.health", resp.SIP),
				attribute.Bool("capillary.health", resp.CAPILLARY),
				attribute.Bool("ghl.health", resp.GHL),
			}
			opts := make([]metric.ObserveOption, len(attributes))
			for i, attr := range attributes {
				opts[i] = metric.WithAttributes(attr)
			}
			apiConnVal := sipConnVal + capillaryConnVal + ghlConnVal
			o.ObserveInt64(apiConn, int64(apiConnVal), opts...)

			return nil
		},
		sipConn,
		capillaryConn,
		ghlConn,
		apiConn,
	)

	if err != nil {
		return err
	}
	return nil
}
