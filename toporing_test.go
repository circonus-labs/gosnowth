package gosnowth

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"testing"
)

func TestTopoRingJSONDeserialization(t *testing.T) {
	dec := json.NewDecoder(bytes.NewBufferString(topoRingJSONTestData))
	dec.UseNumber()
	toporing := []TopoRingDetail{}
	err := dec.Decode(&toporing)
	if err != nil {
		t.Errorf("failed to decode node stats, %s\n", err.Error())
	}

	t.Log(toporing)
	if len(toporing) != 9 {
		t.Error("should be 9 elements")
	}
}

func TestTopoRingXMLDeserialization(t *testing.T) {
	dec := xml.NewDecoder(bytes.NewBufferString(topoRingXMLTestData))
	toporing := new(TopoRing)
	err := dec.Decode(toporing)
	if err != nil {
		t.Errorf("failed to decode node stats, %s\n", err.Error())
	}

	t.Log(toporing)
	if len(toporing.VirtualNodes) != 9 {
		t.Error("should be 9 elements")
	}
}
