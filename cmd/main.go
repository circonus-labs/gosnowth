package main

import (
	"log"

	"github.com/circonus/gosnowth"
)

// main - example program that uses the snowth client
func main() {
	client, err := gosnowth.NewSnowthClient(
		"http://localhost:8112",
		"http://localhost:8113",
	)
	if err != nil {
		log.Fatalf("failed to create snowth client: %v", err)
	}

	for _, node := range client.ListActiveNodes() {
		state, err := client.GetNodeState(node)
		if err != nil {
			log.Fatalf("failed to get state: %v", err)
		}
		log.Printf("%+v", state)
	}
	for _, node := range client.ListActiveNodes() {
		gossip, err := client.GetGossipInfo(node)
		if err != nil {
			log.Fatalf("failed to get gossip: %v", err)
		}
		log.Printf("%+v", gossip)
	}
	for _, node := range client.ListActiveNodes() {
		topology, err := client.GetTopologyInfo(node)
		if err != nil {
			log.Fatalf("failed to get topology: %v", err)
		}
		log.Printf("%+v", topology)
	}
	for _, node := range client.ListActiveNodes() {
		toporing, err := client.GetTopoRingInfo(node.GetCurrentTopology(), node)
		if err != nil {
			log.Fatalf("failed to get topology: %v", err)
		}
		log.Printf("%+v", toporing)
	}
}
