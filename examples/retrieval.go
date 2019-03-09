package main

import (
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/circonus-labs/gosnowth"
)

// ExampleReadNNT demonstrates how to read NNT values from a given snowth node.
func ExampleReadNNT() {
	// Create a new client.
	client, err := gosnowth.NewClient(gosnowth.NewConfig(SnowthServers...).
		WithDiscover(true))
	if err != nil {
		log.Fatalf("failed to create snowth client: %v", err)
	}

	// Write text data in order to read back the data.
	id := uuid.New().String()
	for _, node := range client.ListActiveNodes() {
		// WriteNNT takes in a node and variadic of NNTPartsData entries.
		if err := client.WriteNNT(node, gosnowth.NNTData{
			Metric: "test-metric",
			ID:     id,
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
		}); err != nil {
			log.Fatalf("failed to write text data: %v", err)
		}

		data, err := client.ReadNNTValues(node,
			time.Now().Add(-60*time.Second), time.Now().Add(60*time.Second),
			60, "count", id, "test-metric")
		if err != nil {
			log.Fatalf("failed to read nnt data: %v", err)
		}

		log.Printf("%+v\n", data)
		data1, err := client.ReadNNTAllValues(node,
			time.Now().Add(-60*time.Second), time.Now().Add(60*time.Second),
			60, id, "test-metric")
		log.Printf("%+v\n", data1)
	}
}

// ExampleReadText demonstrates how to read text values from a given snowth
// node.
func ExampleReadText() {
	// Create a new client.
	client, err := gosnowth.NewClient(gosnowth.NewConfig(SnowthServers...).
		WithDiscover(true))
	if err != nil {
		log.Fatalf("failed to create snowth client: %v", err)
	}

	// Write text data in order to read back the data.
	for _, node := range client.ListActiveNodes() {
		id := uuid.New().String()
		if err := client.WriteText(node, gosnowth.TextData{
			Metric: "test-text-metric2",
			ID:     id,
			Offset: strconv.FormatInt(time.Now().Unix(), 10),
			Value:  "a_text_data_value",
		}); err != nil {
			log.Printf("failed to write text data: %v", err)
		}

		data, err := client.ReadTextValues(node,
			time.Now().Add(-60*time.Second), time.Now().Add(60*time.Second),
			id, "test-text-metric2")
		if err != nil {
			log.Printf("failed to read text data: %v", err)
		}

		log.Printf("%+v\n", data)
	}
}
