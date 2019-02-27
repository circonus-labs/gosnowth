package gosnowth

import (
	"bytes"
	"encoding/json"
	"testing"
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

	if len(*gossip) != 4 {
		t.Error("Should have 4 entries")
	}

	res := []GossipDetail(*gossip)[0].Time
	exp := float64(1409082055.744880)
	if res != exp {
		t.Errorf("Expected time: %v, got: %v", exp, res)
	}

	resA := []GossipDetail(*gossip)[0].Age
	expA := float64(0.0)
	if resA != expA {
		t.Errorf("Expected age: %v, got: %v", exp, res)
	}
}
