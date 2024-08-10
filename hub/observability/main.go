package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() (err error) {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := setupOTelSDK(ctx)
	if err != nil {
		return
	}
	meterInit()

	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	select {
	case <-ctx.Done():
		stop()
	}

	return
}
