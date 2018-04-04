package gosnowth

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, 9, len(toporing), "should be 9 elements")
}

func TestTopoRingXMLDeserialization(t *testing.T) {
	dec := xml.NewDecoder(bytes.NewBufferString(topoRingXMLTestData))
	toporing := new(TopoRing)
	err := dec.Decode(toporing)
	if err != nil {
		t.Errorf("failed to decode node stats, %s\n", err.Error())
	}
	t.Log(toporing)
	assert.Equal(t, 9, len(toporing.VirtualNodes), "should be 9 elements")
}
