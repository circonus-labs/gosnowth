package main

import (
	"log"
	"time"

	"github.com/circonus-labs/circonusllhist"
	"github.com/satori/go.uuid"

	"github.com/circonus/gosnowth"
)

// main - example program that uses the snowth client
func main() {
	// create a client, with a seed of nodes
	client, err := gosnowth.NewSnowthClient(
		"http://localhost:8112",
		"http://localhost:8113",
	)
	if err != nil {
		log.Fatalf("failed to create snowth client: %v", err)
	}

	// get the node state from the node
	for _, node := range client.ListActiveNodes() {
		state, err := client.GetNodeState(node)
		if err != nil {
			log.Fatalf("failed to get state: %v", err)
		}
		log.Printf("%+v", state)
	}
	// get the gossip data from the node
	for _, node := range client.ListActiveNodes() {
		gossip, err := client.GetGossipInfo(node)
		if err != nil {
			log.Fatalf("failed to get gossip: %v", err)
		}
		log.Printf("%+v", gossip)
	}
	// get the topology from the node
	for _, node := range client.ListActiveNodes() {
		topology, err := client.GetTopologyInfo(node)
		if err != nil {
			log.Fatalf("failed to get topology: %v", err)
		}
		log.Printf("%+v", topology)
	}
	// get the toporing from the node
	for _, node := range client.ListActiveNodes() {
		toporing, err := client.GetTopoRingInfo(node.GetCurrentTopology(), node)
		if err != nil {
			log.Fatalf("failed to get topology: %v", err)
		}
		log.Printf("%+v", toporing)
	}
	/*
		// write nnt data to node
		for _, node := range client.ListActiveNodes() {
			guid, _ := uuid.NewV4()

			err := client.WriteNNT(
				[]gosnowth.NNTData{
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
					}}, node)
			if err != nil {
				log.Fatalf("failed to write nnt data: %v", err)
			}
		}
		// write text data
		for _, node := range client.ListActiveNodes() {
			guid, _ := uuid.NewV4()
			err := client.WriteText(
				[]gosnowth.TextData{
					gosnowth.TextData{
						Metric: "test-text-metric2", ID: guid.String(),
						Offset: strconv.FormatInt(time.Now().Unix(), 10),
						Value:  "a_text_data_value",
					}}, node)
			if err != nil {
				log.Fatalf("failed to write text data: %v", err)
			}
		}
	*/

	newHistogram, err := circonusllhist.NewFromStrings([]string{
		"H[0.0e+00]=1",
		"H[1.0e+01]=1",
		"H[2.0e+02]=1",
		"H[3.0e+03]=1",
	}, false)

	// write histogram data
	for _, node := range client.ListActiveNodes() {
		guid, _ := uuid.NewV4()
		err := client.WriteHistogram(
			[]gosnowth.HistogramData{
				gosnowth.HistogramData{
					Metric: "test-text-metric2", ID: guid.String(),
					Offset:    time.Now().Unix(),
					Histogram: newHistogram, Period: 60,
				}}, node)
		if err != nil {
			log.Fatalf("failed to write histogram data: %v", err)
		}
	}
}
