package example

import (
	"log"
	"strconv"
	"time"

	"github.com/circonus-labs/gosnowth"
	uuid "github.com/satori/go.uuid"
)

// ExampleReadNNT - this example shows how you are
// able to read NNT values from a given snowth node.
// In this example you need snowth nodes running
// at http://localhost:8112 and http://localhost:8113
func ExampleReadNNT() {
	// create a client, with a seed of nodes
	client, err := gosnowth.NewSnowthClient(
		true,
		"http://localhost:8112",
		"http://localhost:8113",
	)
	if err != nil {
		log.Fatalf("failed to create snowth client: %v", err)
	}
	// write text data in order to read back the data
	guid, _ := uuid.NewV4()
	for _, node := range client.ListActiveNodes() {
		// create a new metric ID, a UUIDv4
		// WriteText takes in a node and variadic of
		// gosnowth.TextData entries
		err := client.WriteNNT(node,
			gosnowth.NNTData{
				Metric: "test-metric", ID: guid.String(),
				Offset: (time.Now().Unix() / 60) * 60,
				Count:  5, Value: 100,
				Parts: gosnowth.Parts{
					Period: 60,
					Data: []gosnowth.NNTPartsData{
						gosnowth.NNTPartsData{Count: 1, Value: 100},
						gosnowth.NNTPartsData{Count: 1, Value: 100},
						gosnowth.NNTPartsData{Count: 1, Value: 100},
						gosnowth.NNTPartsData{Count: 1, Value: 100},
						gosnowth.NNTPartsData{Count: 1, Value: 100},
					}},
			})

		if err != nil {
			log.Fatalf("failed to write text data: %v", err)
		}

		data, err := client.ReadNNTValues(node,
			time.Now().Add(-60*time.Second), time.Now().Add(60*time.Second), 60,
			"count", guid.String(), "test-metric")

		if err != nil {
			log.Fatalf("failed to read nnt data: %v", err)
		}
		log.Printf("%+v\n", data)
		data1, err := client.ReadNNTAllValues(node,
			time.Now().Add(-60*time.Second), time.Now().Add(60*time.Second), 60,
			guid.String(), "test-metric")
		log.Printf("%+v\n", data1)
	}
}

// ExampleReadText - this example shows how you are
// able to read Text values from a given snowth node.
// In this example you need snowth nodes running
// at http://localhost:8112 and http://localhost:8113
func ExampleReadText() {
	// create a client, with a seed of nodes
	client, err := gosnowth.NewSnowthClient(
		true,
		"http://localhost:8112",
		"http://localhost:8113",
	)
	if err != nil {
		log.Fatalf("failed to create snowth client: %v", err)
	}
	// write text data in order to read back the data
	for _, node := range client.ListActiveNodes() {
		guid, _ := uuid.NewV4()

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

		data, err := client.ReadTextValues(node,
			time.Now().Add(-60*time.Second), time.Now().Add(60*time.Second),
			guid.String(), "test-text-metric2")

		if err != nil {
			log.Fatalf("failed to read TEXT data: %v", err)
		}
		log.Printf("%+v\n", data)
	}
}
