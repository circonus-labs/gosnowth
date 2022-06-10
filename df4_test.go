package gosnowth

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

const testDF4Response = `{
	"version": "DF4",
	"head": {
		"count": 3,
		"start": 0,
		"period": 300,
		"explain":{
			"info":{
				"putype":["none","number"]
			}
		}
	},
	"meta": [
		{
			"kind": "numeric",
			"label": "test",
			"tags": [
				"__check_uuid:11223344-5566-7788-9900-aabbccddeeff",
				"__name:test"
			]
		}
	],
	"data": [
		[
			1,
			2,
			3
		]
	]
}`

//nolint shortif
func TestDF4ResponseCopy(t *testing.T) {
	t.Parallel()

	var v *DF4Response
	err := json.NewDecoder(bytes.NewBufferString(testDF4Response)).Decode(&v)
	if err != nil {
		t.Fatal(err)
	}

	b := v.Copy()
	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(&b)
	if err != nil {
		t.Fatal(err)
	}

	s1 := buf.String()
	buf = &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(&v)
	if err != nil {
		t.Fatal(err)
	}

	s2 := buf.String()
	if s1 != s2 {
		t.Errorf("Expected JSON: %v, got: %v", s2, s1)
	}
}

func TestMarshalDF4Response(t *testing.T) {
	t.Parallel()

	v := &DF4Response{
		Data: [][]interface{}{{1, 2, 3}},
		Meta: []DF4Meta{{
			Tags: []string{
				"__check_uuid:11223344-5566-7788-9900-aabbccddeeff",
				"__name:test",
			},
			Label: "test",
			Kind:  "numeric",
		}},
		Ver: "DF4",
		Head: DF4Head{
			Count:  3,
			Start:  0,
			Period: 300,
			Explain: json.RawMessage([]byte(
				`{"info":{"putype":["none","number"]}}`)),
		},
	}

	buf := &bytes.Buffer{}

	if err := json.NewEncoder(buf).Encode(&v); err != nil {
		t.Fatal(err)
	}

	exp := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
		testDF4Response, "\n", ""), " ", ""), "\t", "") + "\n"
	if buf.String() != exp {
		t.Errorf("Expected JSON: %s, got: %s", exp, buf.String())
	}
}

func TestUnmarshalDF4Timeseries(t *testing.T) {
	t.Parallel()

	var v *DF4Response

	if err := json.NewDecoder(bytes.NewBufferString(
		testDF4Response)).Decode(&v); err != nil {
		t.Fatal(err)
	}

	if len(v.Data) != 1 {
		t.Fatalf(`Expected length: 1. got %d`, len(v.Data))
	}

	if v.Data[0][1] != 2.0 {
		t.Errorf(`Expected value: 2.0. got %f`, v.Data[1][1])
	}

	if v.Head.Start != 0 {
		t.Errorf(`Expected time start: 0. got %d`, v.Head.Start)
	}

	if v.Head.Period != 300 {
		t.Errorf(`Expected time period: 300. got %d`, v.Head.Period)
	}

	if len(v.Head.Explain) != 53 {
		t.Errorf("Expected length explain: 53, got: %v", len(v.Head.Explain))
	}
}
