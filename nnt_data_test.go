package gosnowth

import (
	"encoding/json"
	"testing"
)

func TestNNTValue(t *testing.T) {
	var (
		data = "[[1380000000,50],[1380000300,60]]"
		nnta = NNTValueResponse{}
	)
	if err := json.Unmarshal([]byte(data), &nnta); err != nil {
		t.Error("error unmarshalling: ", err)
	}
}

func TestNNTAllValue(t *testing.T) {
	var (
		data = `[
			[
				1379998800,
				{"count":60,"value":10,"stddev":0,"derivative":0,"derivative_stddev":0,"counter":0,"counter_stddev":0,"derivative2":0,"derivative2_stddev":0,"counter2":0,"counter2_stddev":0}
			],
			[
				1380002400,
				{"count":60,"value":10,"stddev":0,"derivative":0,"derivative_stddev":0,
		"counter":0,"counter_stddev":0,"derivative2":0,"derivative2_stddev":0,"counter2":0,"counter2_stddev":0}
			],
			[
				1380006000,
				{"count":60,"value":10,"stddev":1,"derivative":1,"derivative_stddev":1,"counter":1,"counter_stddev":1,"derivative2":1,"derivative2_stddev":1,"counter2":1,"counter2_stddev":1}
			]
		]`
		nntavr = NNTAllValueResponse{}
	)
	if err := json.Unmarshal([]byte(data), &nntavr); err != nil {
		t.Error("error unmarshalling: ", err)
	}
}
