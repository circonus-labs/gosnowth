package gosnowth

import (
	"bytes"
	"encoding/xml"
	"testing"
)

const dataLocationXMLTestData = `<nodes n="2">
	<node id="1f846f26-0cfd-4df5-b4f1-e0930604e577"
		address="10.8.20.1"
		port="8112"
		apiport="8112"
		weight="32"/>
	<node id="07fa2237-5744-4c28-a622-a99cfc1ac87e"
		address="10.8.20.4"
		port="8112"
		apiport="8112"
		weight="32"/>
</nodes>`

func TestDataLocationXMLDeserialization(t *testing.T) {
	dec := xml.NewDecoder(bytes.NewBufferString(dataLocationXMLTestData))
	dl := new(Topology)
	err := dec.Decode(dl)
	if err != nil {
		t.Errorf("failed to decode node stats, %s\n", err.Error())
	}

	if len(dl.Nodes) != 2 {
		t.Errorf("Expected number of nodes: 2, got: %v", len(dl.Nodes))
	}
}
