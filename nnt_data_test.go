package gosnowth

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestNNTValue(t *testing.T) {
	var (
		data = "[[1380000000,50],[1380000300,60]]"
		nnta = NNTValueResponse{}
	)
	if err := json.Unmarshal([]byte(data), &nnta); err != nil {
		t.Error("error unmarshalling: ", err)
	}

	fmt.Println(nnta)

	if nnta.Data[0].Time != time.Unix(1380000000, 0) {
		t.Error("invalid time parsing")
	}
	if nnta.Data[1].Time != time.Unix(1380000300, 0) {
		t.Error("invalid time parsing")
	}

	if nnta.Data[0].Value != 50 {
		t.Error("invalid value parsing")
	}
	if nnta.Data[1].Value != 60 {
		t.Error("invalid value parsing")
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
