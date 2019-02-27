// Package example - this package contains examples on how to use the gosnowth
// client for various operations.
package main

func main() {
	ExampleGetNodeState()
	ExampleGetNodeGossip()
	ExampleGetTopology()
	ExampleGetTopoRing()

	// Perform Example Write NNT
	ExampleSubmitNNT()
	// Perform Example Write Histogram
	ExampleSubmitHistogram()
	// Perform Example Read NNT
	ExampleReadNNT()
	// Perform Example Read Text
	ExampleReadText()
}
