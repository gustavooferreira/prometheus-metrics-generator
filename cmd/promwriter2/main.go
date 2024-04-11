package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/castai/promwrite"
)

func main() {
	client := promwrite.NewClient("http://localhost:9090/api/v1/write")

	// Time: time.Now().Add(500000 * time.Hour),
	// sampleTime := time.Now().Add(-500000 * time.Hour)
	sampleTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	fmt.Printf("Sample time: %+v\n", sampleTime)

	resp, err := client.Write(context.Background(), &promwrite.WriteRequest{
		TimeSeries: []promwrite.TimeSeries{
			{
				Labels: []promwrite.Label{
					{
						Name:  "__name__",
						Value: "my_metric_name_8",
					},
				},
				Sample: promwrite.Sample{
					Time:  sampleTime,
					Value: 123,
				},
			},
		},
	})

	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Response: %+v\n", resp)
}
