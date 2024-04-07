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

	resp, err := client.Write(context.Background(), &promwrite.WriteRequest{
		TimeSeries: []promwrite.TimeSeries{
			{
				Labels: []promwrite.Label{
					{
						Name:  "__name__",
						Value: "my_metric_name",
					},
				},
				Sample: promwrite.Sample{
					Time:  time.Now(),
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
