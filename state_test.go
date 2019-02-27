package gosnowth

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestStateDeserialization(t *testing.T) {
	dec := json.NewDecoder(bytes.NewBufferString(stateTestData))
	dec.UseNumber()
	state := new(NodeState)
	err := dec.Decode(state)
	if err != nil {
		t.Errorf("failed to decode node stats, %s\n", err.Error())
	}
	t.Log(state)

	res := state.Identity
	exp := "bb6f7162-4828-11df-bab8-6bac200dcc2a"
	if res != exp {
		t.Errorf("Expected state.Identity: %v, got: %v", exp, res)
	}

	res = state.Current
	exp = "294cbd39999c2270964029691e8bc5e231a867d525ccba62181dc8988ff218dc"
	if res != exp {
		t.Errorf("Expected state.Current: %v, got: %v", exp, res)
	}

	resUI := state.BaseRollup
	expUI := uint64(60)
	if resUI != expUI {
		t.Errorf("Expected state.BaseRollup: %v, got: %v", expUI, resUI)
	}

	if len(state.NNT.RollupEntries) != 4 {
		t.Errorf("Expected length state.NNT.RollupEntries: 4, got: %v",
			len(state.NNT.RollupEntries))
	}
}
