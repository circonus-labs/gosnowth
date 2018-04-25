// Package example - this package contains examples on how to use the gosnowth
// client for various operations.
package main

import "github.com/circonus-labs/gosnowth/cmd/example"

func main() {
	example.ExampleGetNodeState()
	example.ExampleGetNodeGossip()
	example.ExampleGetTopology()
	example.ExampleGetTopoRing()

	// Perform Example Write NNT
	example.ExampleSubmitNNT()
	// Perform Example Write Histogram
	example.ExampleSubmitHistogram()
	// Perform Example Read NNT
	example.ExampleReadNNT()
	// Perform Example Read Text
	example.ExampleReadText()
}
