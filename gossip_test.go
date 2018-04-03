package gosnowth

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGossipDeserialization(t *testing.T) {
	dec := json.NewDecoder(bytes.NewBufferString(gossipTestData))
	dec.UseNumber()
	gossip := new(Gossip)
	err := dec.Decode(gossip)
	if err != nil {
		t.Errorf("failed to decode gossip data, %s\n", err.Error())
	}
	t.Log(gossip)

	assert.Equal(t, 4, len(*gossip), "should have 4 entries")
	assert.Equal(t, 1409082055.744880, []GossipDetail(*gossip)[0].Time, "time should be")
	assert.Equal(t, 0.0, []GossipDetail(*gossip)[0].Age, "age should be")
}
