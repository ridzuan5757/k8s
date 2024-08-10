package main

import (
	"context"

	"go.opentelemetry.io/otel"
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

	if err != nil {
		return err
	}

	_, err = meter.RegisterCallback(
		func(_ context.Context, o metric.Observer) error {
			resp, err := getHub()
			if err != nil {
				return err
			}

			if resp.SIP {
				o.ObserveInt64(sipConn, int64(1))
			} else {
				o.ObserveInt64(sipConn, int64(0))
			}

			if resp.CAPILLARY {
				o.ObserveInt64(capillaryConn, int64(1))
			} else {
				o.ObserveInt64(capillaryConn, int64(0))
			}

			if resp.GHL {
				o.ObserveInt64(ghlConn, int64(1))
			} else {
				o.ObserveInt64(ghlConn, int64(0))
			}

			return nil
		},
		sipConn,
		capillaryConn,
		ghlConn,
	)

	if err != nil {
		return err
	}
	return nil
}
