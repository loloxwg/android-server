package main

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"time"
)

func main() {
	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient("http://124.222.47.219:8086", "root")
	// Use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking("my-org", "my-bucket")
	// Create point using full params constructor
	p := influxdb2.NewPoint("stat",
		map[string]string{"unit": "temperature"},
		map[string]interface{}{"avg": 24.5, "max": 45.0},
		time.Now())
	// write point immediately
	writeAPI.WritePoint(context.Background(), p)
	// Create point using fluent style
	p = influxdb2.NewPointWithMeasurement("sensor_data").
		AddTag("location", "bedroom").
		AddField("humidity", 46.4).
		AddField("pm2_5", 0.5).
		AddField("temperature", 25).
		SetTime(time.Now())
	writeAPI.WritePoint(context.Background(), p)

	// Or write directly line protocol
	line := fmt.Sprintf("stat,unit=temperature avg=%f,max=%f", 23.5, 45.0)
	writeAPI.WriteRecord(context.Background(), line)

	// Get query client
	queryAPI := client.QueryAPI("my-org")
	// Get parser flux query result
	result, err := queryAPI.Query(context.Background(), `from(bucket:"my-bucket")|> range(start: -1h) |> filter(fn: (r) => r._measurement == "stat")`)
	if err == nil {
		// Use Next() to iterate over query result lines
		for result.Next() {
			// Observe when there is new grouping key producing new table
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			// read result
			fmt.Printf("row: %s\n", result.Record().String())
		}
		if result.Err() != nil {
			fmt.Printf("Query error: %s\n", result.Err().Error())
		}
	}
	fmt.Println("error:", err.Error())
	// Ensures background processes finishes
	client.Close()
}
