package main

import (
	"log"
	"strconv"
	"time"

	"github.com/circonus-labs/circonusllhist"
	"github.com/google/uuid"

	"github.com/circonus-labs/gosnowth"
)

// ExampleSubmitText demonstrates how to submit a text metric to a node.
func ExampleSubmitText() {
	// Create a new client.
	cfg, err := gosnowth.NewConfig(SnowthServers...)
	if err != nil {
		log.Fatalf("failed to create snowth configuration: %v", err)
	}

	cfg.SetDiscover(true)
	client, err := gosnowth.NewClient(cfg)
	if err != nil {
		log.Fatalf("failed to create snowth client: %v", err)
	}

	// Write text data.
	for _, node := range client.ListActiveNodes() {
		id := uuid.New().String()
		// WriteText takes in a node and variadic of TextData entries.
		if err := client.WriteText(node, gosnowth.TextData{
			Metric: "test-text-metric2",
			ID:     id,
			Offset: strconv.FormatInt(time.Now().Unix(), 10),
			Value:  "a_text_data_value",
		}); err != nil {
			log.Fatalf("failed to write text data: %v", err)
		}
	}
}

// ExampleSubmitNNT demonstrates how to submit an NNT metric to a node.
func ExampleSubmitNNT() {
	// Create a new client.
	cfg, err := gosnowth.NewConfig(SnowthServers...)
	if err != nil {
		log.Fatalf("failed to create snowth configuration: %v", err)
	}

	cfg.SetDiscover(true)
	client, err := gosnowth.NewClient(cfg)
	if err != nil {
		log.Fatalf("failed to create snowth client: %v", err)
	}

	// Write NNT data to the node.
	for _, node := range client.ListActiveNodes() {
		id := uuid.New().String()
		if err := client.WriteNNT(node, gosnowth.NNTData{
			Metric: "test-metric",
			ID:     id,
			Offset: time.Now().Unix(),
			Count:  5, Value: 100,
			Parts: gosnowth.Parts{
				Period: 60,
				Data: []gosnowth.NNTPartsData{
					gosnowth.NNTPartsData{Count: 1, Value: 100},
					gosnowth.NNTPartsData{Count: 1, Value: 100},
					gosnowth.NNTPartsData{Count: 1, Value: 100},
					gosnowth.NNTPartsData{Count: 1, Value: 100},
					gosnowth.NNTPartsData{Count: 1, Value: 100},
				},
			},
		}); err != nil {
			log.Fatalf("failed to write nnt data: %v", err)
		}
	}
}

// ExampleSubmitHistogram demonstrates how to submit histogram data to a node.
func ExampleSubmitHistogram() {
	// Create a new client.
	cfg, err := gosnowth.NewConfig(SnowthServers...)
	if err != nil {
		log.Fatalf("failed to create snowth configuration: %v", err)
	}

	cfg.SetDiscover(true)
	client, err := gosnowth.NewClient(cfg)
	if err != nil {
		log.Fatalf("failed to create snowth client: %v", err)
	}

	// Using the circonusllhist library, create a new histogram from a string
	// representation of a histogram.
	hist, err := circonusllhist.NewFromStrings([]string{
		"H[0.0e+00]=1",
		"H[1.0e+01]=1",
		"H[2.0e+02]=1",
		"H[3.0e+03]=1",
	}, false)

	// Write histogram data.
	for _, node := range client.ListActiveNodes() {
		id := uuid.New().String()
		if err := client.WriteHistogram(node, gosnowth.HistogramData{
			AccountID: 1,
			Metric:    "test-hist-metric",
			ID:        id,
			CheckName: "test",
			Offset:    time.Now().Unix(),
			Histogram: hist,
			Period:    60,
		}); err != nil {
			log.Fatalf("failed to write histogram data: %v", err)
		}
	}
}
