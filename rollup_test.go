package gosnowth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

const rollupTestData = "[[1,1]]"

const rollupAllTestData = `[
	[
		1529509020,
		{
			"count": 1,
			"counter": 0,
			"counter2": 0,
			"counter2_stddev": 0,
			"counter_stddev": 0,
			"derivative": 0,
			"derivative2": 0,
			"derivative2_stddev": 0,
			"derivative_stddev": 0,
			"stddev": 0,
			"value": 0
		}
	],
	[
		1529509080,
		{
			"count": 1,
			"counter": 0,
			"counter2": 0,
			"counter2_stddev": 0,
			"counter_stddev": 0,
			"derivative": 0,
			"derivative2": 0,
			"derivative2_stddev": 0,
			"derivative_stddev": 0,
			"stddev": 0,
			"value": 0
		}
	],
	[
		1529509140,
		{
			"count": 1,
			"counter": 0,
			"counter2": 0,
			"counter2_stddev": 0,
			"counter_stddev": 0,
			"derivative": 0,
			"derivative2": 0,
			"derivative2_stddev": 0,
			"derivative_stddev": 0,
			"stddev": 0,
			"value": 0
		}
	]
]`

func TestRollupValueMarshaling(t *testing.T) {
	v := []RollupValue{}
	err := json.NewDecoder(bytes.NewBufferString(rollupTestData)).Decode(&v)
	if err != nil {
		t.Fatal(err)
	}

	if len(v) != 1 {
		t.Fatalf("Expected length: 1, got %v", len(v))
	}

	if v[0].Timestamp() != "1" {
		t.Errorf("Expected timestamp: 1, got: %v", v[0].Timestamp())
	}

	if v[0].Value == nil {
		t.Fatal("Expected value: not nil, got: nil")
	}

	if *v[0].Value != 1 {
		t.Errorf("Expected value: 1, got: %v", *v[0].Value)
	}

	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(&v)
	if err != nil {
		t.Fatal(err)
	}

	if buf.String() != rollupTestData+"\n" {
		t.Errorf("Expected JSON: %v, got: %v", rollupTestData, buf.String())
	}
}

func TestRollupAllValueMarshaling(t *testing.T) {
	v := []RollupAllValue{}
	err := json.NewDecoder(bytes.NewBufferString(`[
		[
			1529509020,
			null
		]
	]`)).Decode(&v)
	if err != nil {
		t.Fatal(err)
	}

	if v[0].Timestamp() != "1529509020" {
		t.Errorf("Expected timestamp: 1529509020, got: %v", v[0].Timestamp())
	}

	if v[0].Data != nil {
		t.Fatalf("Expected data: nil, got: %v", *v[0].Data)
	}

	err = json.NewDecoder(bytes.NewBufferString(rollupAllTestData)).Decode(&v)
	if err != nil {
		t.Fatal(err)
	}

	if len(v) != 3 {
		t.Fatalf("Expected length: 3, got %v", len(v))
	}

	if v[0].Timestamp() != "1529509020" {
		t.Errorf("Expected timestamp: 1529509020, got: %v", v[0].Timestamp())
	}

	if v[0].Data == nil {
		t.Fatal("Expected data: not nil, got: nil")
	}

	if v[0].Data.Value != 0.0 {
		t.Errorf("Expected value: 0, got: %v", v[0].Data.Value)
	}

	if v[0].Data.Count != 1 {
		t.Errorf("Expected count: 1, got: %v", v[0].Data.Count)
	}

	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(&v)
	if err != nil {
		t.Fatal(err)
	}

	exp := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
		rollupAllTestData, " ", ""), "\n", ""), "\t", "") + "\n"
	if buf.String() != exp {
		t.Errorf("Expected JSON: %v, got: %v", exp, buf.String())
	}
}

func TestReadRollupValues(t *testing.T) {
	ms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		if r.RequestURI == "/state" {
			_, _ = w.Write([]byte(stateTestData))
			return
		}

		if r.RequestURI == "/stats.json" {
			_, _ = w.Write([]byte(statsTestData))
			return
		}

		u := "/rollup/fc85e0ab-f568-45e6-86ee-d7443be8277d/" +
			"online?start_ts=1529509020&end_ts=1529509201&rollup_span=1s" +
			"&type=average"
		if strings.HasPrefix(r.RequestURI, u) {
			_, _ = w.Write([]byte(rollupTestData))
			return
		}
	}))

	defer ms.Close()
	sc, err := NewSnowthClient(false, ms.URL)
	if err != nil {
		t.Fatal("Unable to create snowth client", err)
	}

	u, err := url.Parse(ms.URL)
	if err != nil {
		t.Fatal("Invalid test URL")
	}

	node := &SnowthNode{url: u}
	res, err := sc.ReadRollupValues(
		"fc85e0ab-f568-45e6-86ee-d7443be8277d", "online", time.Second,
		time.Unix(1529509020, 0), time.Unix(1529509200, 0), "average", node)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 1 {
		t.Fatalf("Expected length: 1, got: %v", len(res))
	}

	if res[0].Value == nil {
		t.Fatal("Expected value: not nil, got: nil")
	}

	if *res[0].Value != 1 {
		t.Errorf("Expected value: 1, got: %v", *res[0].Value)
	}
}

func TestReadRollupAllValues(t *testing.T) {
	ms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		if r.RequestURI == "/state" {
			_, _ = w.Write([]byte(stateTestData))
			return
		}

		if r.RequestURI == "/stats.json" {
			_, _ = w.Write([]byte(statsTestData))
			return
		}

		u := "/rollup/fc85e0ab-f568-45e6-86ee-d7443be8277d/" +
			"online?start_ts=1529509020&end_ts=1529509201&rollup_span=1s" +
			"&type=all"
		if strings.HasPrefix(r.RequestURI, u) {
			_, _ = w.Write([]byte(rollupAllTestData))
			return
		}
	}))

	defer ms.Close()
	sc, err := NewSnowthClient(false, ms.URL)
	if err != nil {
		t.Fatal("Unable to create snowth client", err)
	}

	u, err := url.Parse(ms.URL)
	if err != nil {
		t.Fatal("Invalid test URL")
	}

	node := &SnowthNode{url: u}
	res, err := sc.ReadRollupAllValues(
		"fc85e0ab-f568-45e6-86ee-d7443be8277d", "online", time.Second,
		time.Unix(1529509020, 0), time.Unix(1529509200, 0), node)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 3 {
		t.Fatalf("Expected length: 3, got: %v", len(res))
	}

	if res[0].Data == nil {
		t.Fatal("Expected data: not nil, got: nil")
	}

	if res[0].Data.Count != 1 {
		t.Errorf("Expected count: 1, got: %v", res[0].Data.Count)
	}

	if res[0].Data.Value != 0 {
		t.Errorf("Expected value: 0, got: %v", res[0].Data.Value)
	}
}

func TestReadRollupValuesOpts(t *testing.T) {
	ms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		if r.RequestURI == "/state" {
			_, _ = w.Write([]byte(stateTestData))
			return
		}

		if r.RequestURI == "/stats.json" {
			_, _ = w.Write([]byte(statsTestData))
			return
		}

		u := "/rollup/fc85e0ab-f568-45e6-86ee-d7443be8277d/" +
			"online?start_ts=now-1d&end_ts=now&rollup_span=1s" +
			"&type=average"
		if strings.HasPrefix(r.RequestURI, u) {
			_, _ = w.Write([]byte(rollupTestData))
			return
		}

		u = "/rollup/fc85e0ab-f568-45e6-86ee-d7443be8277d/" +
			"online?start_ts=now-1d&end_ts=now&rollup_span=1s" +
			"&type=all"
		if strings.HasPrefix(r.RequestURI, u) {
			_, _ = w.Write([]byte(rollupAllTestData))
			return
		}
	}))

	defer ms.Close()
	sc, err := NewSnowthClient(false, ms.URL)
	if err != nil {
		t.Fatal("Unable to create snowth client", err)
	}

	u, err := url.Parse(ms.URL)
	if err != nil {
		t.Fatal("Invalid test URL")
	}

	node := &SnowthNode{url: u}
	res, err := sc.ReadRollupValuesOpts(context.Background(), &RollupOptions{
		UUID:   "fc85e0ab-f568-45e6-86ee-d7443be8277d",
		Metric: "online",
		Type:   "all",
		Period: "1s",
		Start:  "now-1d",
		End:    "now",
	}, node)
	if err != nil {
		t.Fatal(err)
	}

	if res.AllValues == nil {
		t.Fatal("Expected all values: not nil, got: nil")
	}

	if len(res.AllValues) != 3 {
		t.Fatalf("Expected length: 3, got: %v", len(res.AllValues))
	}

	if res.AllValues[0].Data == nil {
		t.Fatal("Expected data: not nil, got: nil")
	}

	if res.AllValues[0].Data.Count != 1 {
		t.Errorf("Expected count: 1, got: %v", res.AllValues[0].Data.Count)
	}

	if res.AllValues[0].Data.Value != 0 {
		t.Errorf("Expected value: 0, got: %v", res.AllValues[0].Data.Value)
	}

	res, err = sc.ReadRollupValuesOpts(context.Background(), &RollupOptions{
		UUID:   "fc85e0ab-f568-45e6-86ee-d7443be8277d",
		Metric: "online",
		Type:   "average",
		Period: "1s",
		Start:  "now-1d",
		End:    "now",
	}, node)
	if err != nil {
		t.Fatal(err)
	}

	if res.Values == nil {
		t.Fatal("Expected values: not nil, got: nil")
	}

	if len(res.Values) != 1 {
		t.Fatalf("Expected length: 1, got: %v", len(res.Values))
	}

	if res.Values[0].Value == nil {
		t.Fatal("Expected value: not nil, got: nil")
	}

	if *res.Values[0].Value != 1 {
		t.Errorf("Expected value: 1, got: %v", *res.Values[0].Value)
	}
}
