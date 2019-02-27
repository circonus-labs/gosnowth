package gosnowth

import (
	"bytes"
	"encoding/xml"
	"testing"
)

func TestDataLocationXMLDeserialization(t *testing.T) {
	dec := xml.NewDecoder(bytes.NewBufferString(dataLocationXMLTestData))
	dl := new(DataLocation)
	err := dec.Decode(dl)
	if err != nil {
		t.Errorf("failed to decode node stats, %s\n", err.Error())
	}
	t.Log(dl)

	if len(dl.Nodes) != 2 {
		t.Errorf("Expected number of nodes: 2, got: %v", len(dl.Nodes))
	}
}
