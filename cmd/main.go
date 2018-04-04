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
}
