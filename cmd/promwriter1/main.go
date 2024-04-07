package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/m3db/prometheus_remote_client_golang/promremote"
)

func main() {
	userAgent := "Ovat User Agent"
	writeURLFlag := "http://localhost:9090/api/v1/write"

	// create config and client
	cfg := promremote.NewConfig(
		promremote.WriteURLOption(writeURLFlag),
		promremote.HTTPClientTimeoutOption(60*time.Second),
		promremote.UserAgent(userAgent),
	)

	client, err := promremote.NewClient(cfg)
	if err != nil {
		fmt.Printf("unable to construct client: %v\n", err)
	}

	timeSeriesList := []promremote.TimeSeries{
		{
			Labels: []promremote.Label{
				{
					Name:  "__name__",
					Value: "foo_bar",
				},
				{
					Name:  "biz",
					Value: "baz",
				},
			},
			Datapoint: promremote.Datapoint{
				Timestamp: time.Now(),
				Value:     1415.92,
			},
		},
	}

	ctx := context.Background()

	writeResult, err := client.WriteTimeSeries(ctx, timeSeriesList, promremote.WriteOptions{})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Write Result: %+v\n", writeResult)
}
