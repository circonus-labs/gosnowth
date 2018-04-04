package gosnowth

import (
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataLocationXMLDeserialization(t *testing.T) {
	dec := xml.NewDecoder(bytes.NewBufferString(dataLocationXMLTestData))
	dl := new(DataLocation)
	err := dec.Decode(dl)
	if err != nil {
		t.Errorf("failed to decode node stats, %s\n", err.Error())
	}
	t.Log(dl)

	assert.Equal(t, 2, len(dl.Nodes), "number of nodes wrong")
}
