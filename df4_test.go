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
		"error": ["test", "test"],
		"warning": "test",
		"explain":{
			"info":{
				"putype":["none","number"]
			}
		}
	},
	"meta": [
		{
			"kind": "numeric",
			"label": "test_numeric",
			"tags": [
				"__check_uuid:11223344-5566-7788-9900-aabbccddeeff",
				"__name:test_numeric"
			]
		},
		{
			"kind": "text",
			"label": "test_text",
			"tags": [
				"__check_uuid:11223344-5566-7788-9900-aabbccddeeff",
				"__name:test_text"
			]
		},
		{
			"kind": "histogram",
			"label": "test_histogram",
			"tags": [
				"__check_uuid:11223344-5566-7788-9900-aabbccddeeff",
				"__name:test_histogram"
			]
		}
	],
	"data": [
		[
			1,
			null,
			2
		],
		[
			[
				[
					6866,
					"test1"
				]
			],
			[],
			[
				[
					6866,
					"test2"
				]
			]
		],
		[
			{"+12e-004": 1},
			null,
			{"+12e-004": 1}
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

func TestDF4ResponseMarshaling(t *testing.T) {
	t.Parallel()

	var v *DF4Response

	if err := json.NewDecoder(bytes.NewBufferString(
		testDF4Response)).Decode(&v); err != nil {
		t.Fatal(err)
	}

	if len(v.Head.Error) != 2 {
		t.Fatalf("Expected length: 2, got: %v", len(v.Head.Error))
	}

	if len(v.Head.Warning) != 1 {
		t.Fatalf("Expected length: 1, got: %v", len(v.Head.Warning))
	}

	if v.Head.Warning[0] != "test" {
		t.Errorf("Expected warning: test, got: %v", v.Head.Warning)
	}

	if len(v.Data) != 3 {
		t.Fatalf("Expected length: 3, got: %v", len(v.Data))
	}

	if v.Data[0][2] != 2.0 {
		t.Errorf("Expected value: 2.0, got: %v", v.Data[0][2])
	}

	if v.Head.Start != 0 {
		t.Errorf("Expected time start: 0, got: %v", v.Head.Start)
	}

	if v.Head.Period != 300 {
		t.Errorf("Expected time period: 300, got: %v", v.Head.Period)
	}

	exp := `{"info":{"putype":["none","number"]}}`
	if string(v.Head.Explain) != exp {
		t.Errorf("Expected explain: %v, got: %v", exp, string(v.Head.Explain))
	}

	buf := &bytes.Buffer{}

	if err := json.NewEncoder(buf).Encode(&v); err != nil {
		t.Fatal(err)
	}

	exp = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
		testDF4Response, "\n", ""), " ", ""), "\t", "") + "\n"
	if buf.String() != exp {
		t.Errorf("Expected JSON: %s, got: %s", exp, buf.String())
	}
}

func TestDF4Data(t *testing.T) {
	t.Parallel()

	var v *DF4Response

	if err := json.NewDecoder(bytes.NewBufferString(
		testDF4Response)).Decode(&v); err != nil {
		t.Fatal(err)
	}

	if len(v.Data) != 3 {
		t.Fatalf("Expected length: 3, got: %v", len(v.Data))
	}

	num := v.Data[0].Numeric()

	if len(num) != 3 {
		t.Fatalf("Expected length: 3, got: %v", len(num))
	}

	if num[1] != nil {
		t.Errorf("Expected value: nil, got: %v", num[1])
	}

	if *num[2] != 2.0 {
		t.Errorf("Expected value: 2.0, got: %v", *num[2])
	}

	text := v.Data[1].Text()

	if len(text) != 3 {
		t.Fatalf("Expected length: 3, got: %v", len(text))
	}

	if text[1] != nil {
		t.Errorf("Expected value: nil, got: %v", text[1])
	}

	if *text[2] != "test2" {
		t.Errorf("Expected value: test2, got: %v", *text[2])
	}

	hist := v.Data[2].Histogram()

	if len(hist) != 3 {
		t.Fatalf("Expected length: 3, got: %v", len(hist))
	}

	if hist[1] != nil {
		t.Errorf("Expected value: nil, got: %v", hist[1])
	}

	if (*hist[2])["+12e-004"] != 1 {
		t.Errorf("Expected value: 1, got: %v", (*hist[2])["+12e-004"])
	}
}

func TestDF4DataNullEmpty(t *testing.T) {
	t.Parallel()

	var v *DF4Response

	if err := json.NewDecoder(bytes.NewBufferString(
		testDF4Response)).Decode(&v); err != nil {
		t.Fatal(err)
	}

	if len(v.Data) != 3 {
		t.Fatalf("Expected length: 3, got: %v", len(v.Data))
	}

	tv, ok := v.Data[1][1].([]interface{})

	if !ok || len(tv) != 0 {
		t.Errorf("Expected value: [], got: %v", v.Data[1][1])
	}

	for _, dv := range v.Data {
		dv.NullEmpty()
	}

	if v.Data[1][1] != nil {
		t.Errorf("Expected value: nil, got: %v", v.Data[1][1])
	}
}
