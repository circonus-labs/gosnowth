package example

import (
	"fmt"
	"log"

	"github.com/circonus-labs/gosnowth"
)

// ExampleGetNodeState - this example shows how you can get
// the snowth node's state from a particular node.  In this
// example you need a snowth instance running at
// http://localhost:8112 and http://localhost:8113
func ExampleGetNodeState() {
	// create a client, with a seed of nodes
	client, err := gosnowth.NewSnowthClient(
		true,
		"http://localhost:8112",
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
		fmt.Println(state)
	}
}

// ExampleGetNodeGossip - this example shows how you can get
// the snowth node's gossip details from a particular node.  In this
// example you need a snowth instance running at
// http://localhost:8112 and http://localhost:8113
func ExampleGetNodeGossip() {
	// create a client, with a seed of nodes
	client, err := gosnowth.NewSnowthClient(
		true,
		"http://localhost:8112",
		"http://localhost:8113",
	)
	if err != nil {
		log.Fatalf("failed to create snowth client: %v", err)
	}

	// get the gossip data from the node
	for _, node := range client.ListActiveNodes() {
		gossip, err := client.GetGossipInfo(node)
		if err != nil {
			log.Fatalf("failed to get gossip: %v", err)
		}
		fmt.Println(gossip)
	}
}

// ExampleGetNodeTopology - this example shows how you can get
// the snowth node's topology details from a particular node.  In this
// example you need a snowth instance running at
// http://localhost:8112 and http://localhost:8113
func ExampleGetTopology() {
	// create a client, with a seed of nodes
	client, err := gosnowth.NewSnowthClient(
		true,
		"http://localhost:8112",
		"http://localhost:8113",
	)
	if err != nil {
		log.Fatalf("failed to create snowth client: %v", err)
	}

	// get the topology from the node
	for _, node := range client.ListActiveNodes() {
		topology, err := client.GetTopologyInfo(node)
		if err != nil {
			log.Fatalf("failed to get topology: %v", err)
		}
		fmt.Println(topology)
	}
}

// ExampleGetNodeTopoRing - this example shows how you can get
// the snowth node's toporing details from a particular node.  In this
// example you need a snowth instance running at
// http://localhost:8112 and http://localhost:8113
func ExampleGetTopoRing() {
	// create a client, with a seed of nodes
	client, err := gosnowth.NewSnowthClient(
		true,
		"http://localhost:8112",
		"http://localhost:8113",
	)
	if err != nil {
		log.Fatalf("failed to create snowth client: %v", err)
	}

	// get the toporing from the node
	for _, node := range client.ListActiveNodes() {
		toporing, err := client.GetTopoRingInfo(node.GetCurrentTopology(), node)
		if err != nil {
			log.Fatalf("failed to get toporing: %v", err)
		}
		fmt.Println(toporing)
	}
}
