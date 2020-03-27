// Package main contains code demonstrating how to use the gosnowth client
// library for various operations.
//
// The examples executable should be run with a space separated list of IRONdb
// servers in the format <host>:<port> as its arguments.  It will default to
// using localhost:8112.
package main

import (
	"os"
	"strings"
)

// SnowthServers contains the IRONdb servers to use when running the examples.
var SnowthServers = []string{"http://localhost:8112"}

func main() {
	if len(os.Args) > 1 {
		SnowthServers = []string{}
		for _, svr := range os.Args[1:] {
			if !strings.HasPrefix(svr, "http://") {
				svr = "http://" + svr
			}

			SnowthServers = append(SnowthServers, svr)
		}
	}

	ExampleGetNodeState()
	ExampleGetNodeGossip()
	ExampleGetTopology()
	ExampleSubmitNNT()
	ExampleSubmitHistogram()
	ExampleReadNNT()
	ExampleReadText()
}
