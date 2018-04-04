package gosnowth

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, "bb6f7162-4828-11df-bab8-6bac200dcc2a", state.Identity, "should equal")
	assert.Equal(t, "294cbd39999c2270964029691e8bc5e231a867d525ccba62181dc8988ff218dc", state.Current, "should equal")
	assert.Equal(t, uint64(60), state.BaseRollup, "should equal")
	assert.Equal(t, 4, len(state.NNT.RollupEntries), "should equal")
}
