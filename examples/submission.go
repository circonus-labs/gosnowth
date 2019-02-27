package main

import (
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/circonus-labs/circonusllhist"
	"github.com/circonus-labs/gosnowth"
)

// ExampleSubmitText - this example shows how you
// can submit a text metric to a particular snowth
// node.  In this example you need snowth nodes running
// at http://localhost:8112 and http://localhost:8113
func ExampleSubmitText() {
	// create a client, with a seed of nodes
	client, err := gosnowth.NewSnowthClient(
		true,
		"http://localhost:8112",
		"http://localhost:8113",
	)
	if err != nil {
		log.Fatalf("failed to create snowth client: %v", err)
	}

	// write text data
	for _, node := range client.ListActiveNodes() {
		// create a new metric ID, a UUIDv4
		guid := uuid.New()
		// WriteText takes in a node and variadic of
		// gosnowth.TextData entries
		err := client.WriteText(
			node,
			gosnowth.TextData{
				Metric: "test-text-metric2", ID: guid.String(),
				Offset: strconv.FormatInt(time.Now().Unix(), 10),
				Value:  "a_text_data_value",
			})
		if err != nil {
			log.Fatalf("failed to write text data: %v", err)
		}
	}
}

// ExampleSubmitNNT - this example shows how you
// can submit an NNT metric to a particular snowth
// node.  In this example you need snowth nodes running
// at http://localhost:8112 and http://localhost:8113
func ExampleSubmitNNT() {
	// create a client, with a seed of nodes
	client, err := gosnowth.NewSnowthClient(
		true,
		"http://localhost:8112",
		"http://localhost:8113",
	)
	if err != nil {
		log.Fatalf("failed to create snowth client: %v", err)
	}

	// write nnt data to node
	for _, node := range client.ListActiveNodes() {
		// create a new UUID for our NNT metric
		guid := uuid.New()
		err := client.WriteNNT(
			node,
			gosnowth.NNTData{
				Metric: "test-metric", ID: guid.String(),
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
			})
		if err != nil {
			log.Fatalf("failed to write nnt data: %v", err)
		}
	}
}

// ExampleSubmitHistogram - this example shows how you
// can submit a histogram metric to a particular snowth
// node.  In this example you need snowth nodes running
// at http://localhost:8112 and http://localhost:8113
func ExampleSubmitHistogram() {
	// create a client, with a seed of nodes
	client, err := gosnowth.NewSnowthClient(
		true,
		"http://localhost:8112",
		"http://localhost:8113",
	)
	if err != nil {
		log.Fatalf("failed to create snowth client: %v", err)
	}

	// using the circonusllhist library, we are creating a new
	// histogram from the string representation of a histogram.
	newHistogram, err := circonusllhist.NewFromStrings([]string{
		"H[0.0e+00]=1",
		"H[1.0e+01]=1",
		"H[2.0e+02]=1",
		"H[3.0e+03]=1",
	}, false)

	// write histogram data
	for _, node := range client.ListActiveNodes() {
		guid := uuid.New()
		err := client.WriteHistogram(
			node,
			gosnowth.HistogramData{
				Metric: "test-text-metric2", ID: guid.String(),
				Offset: time.Now().Unix(),
				// our histogram is of circonusllhist.Histogram type
				Histogram: newHistogram, Period: 60,
			})
		if err != nil {
			log.Fatalf("failed to write histogram data: %v", err)
		}
	}
}
