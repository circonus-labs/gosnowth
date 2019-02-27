package gosnowth

import "testing"

func TestNewSnowthClient(t *testing.T) {
	// crude test to ensure err is returned for invalid snowth url
	badAddr := "foobar"
	_, err := NewSnowthClient(false, badAddr)
	if err == nil {
		t.Errorf("Error not encountered on invalid snowth addr %v", badAddr)
	}
}

func TestIsNodeActive(t *testing.T) {
	// mock out GetNodeState, GetGossipInfo
}
