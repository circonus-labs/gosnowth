package main

import (
	"log"

	"github.com/circonus-labs/gosnowth"
)

// ExampleGetNodeState demonstrates how to get the snowth node's state from
// a particular node.
func ExampleGetNodeState() {
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

	// Get the node state.
	for _, node := range client.ListActiveNodes() {
		state, err := client.GetNodeState(node)
		if err != nil {
			log.Fatalf("failed to get state: %v", err)
		}

		log.Println(state)
	}
}

// ExampleGetNodeGossip demontrates how to get gossip details from a node.
func ExampleGetNodeGossip() {
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

	// Get the gossip data from the node.
	for _, node := range client.ListActiveNodes() {
		gossip, err := client.GetGossipInfo(node)
		if err != nil {
			log.Fatalf("failed to get gossip: %v", err)
		}

		log.Println(gossip)
	}
}

// ExampleGetTopology demonstrates how to get topology details from a node.
func ExampleGetTopology() {
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

	// Get the topology from the node.
	for _, node := range client.ListActiveNodes() {
		topology, err := client.GetTopologyInfo(node)
		if err != nil {
			log.Fatalf("failed to get topology: %v", err)
		}

		log.Println(topology)
	}
}

// ExampleGetTopoRing demonstrates how to get topology ring details from a
// node.
func ExampleGetTopoRing() {
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

	// Get the topology ring data from the node.
	for _, node := range client.ListActiveNodes() {
		tr, err := client.GetTopoRingInfo(node.GetCurrentTopology(), node)
		if err != nil {
			log.Fatalf("failed to get topology ring: %v", err)
		}

		log.Println(tr)
	}
}
