package gosnowth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

const histogramTestData = `[
	[
	  1556290800,
	  300,
	  {
		"+23e-004": 1,
		"+85e-004": 1
	  }
	],
	[
	  1556291100,
	  300,
	  {
		"+22e-004": 1,
		"+23e-004": 2,
		"+30e-004": 1,
		"+39e-003": 1
	  }
	]
]`

const histTestData = `[
	{
		"account_id": 1,
		"metric": "example1",
		"id": "ae0f7f90-2a6b-481c-9cf5-21a31837020e",
		"check_name": "test",
		"offset": 1408724400,
		"period": 60,
		"histogram": "AAA="
	}
]`

func TestHistogramValueMarshaling(t *testing.T) {
	v := []HistogramValue{}
	err := json.NewDecoder(bytes.NewBufferString(histogramTestData)).Decode(&v)
	if err != nil {
		t.Fatal(err)
	}

	if len(v) != 2 {
		t.Fatalf("Expected length: 2, got %v", len(v))
	}

	if v[0].Timestamp() != "1556290800" {
		t.Errorf("Expected timestamp: 1556290800, got: %v", v[0].Timestamp())
	}

	if v[0].Period.Seconds() != 300.0 {
		t.Errorf("Expected seconds: 300, got: %v", v[0].Period.Seconds())
	}

	if v[0].Data["+23e-004"] != 1 {
		t.Errorf("Expected data: 1, got: %v", v[0].Data["+23e-004"])
	}

	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(&v)
	if err != nil {
		t.Fatal(err)
	}

	exp := strings.Replace(strings.Replace(strings.Replace(histogramTestData,
		" ", "", -1), "\n", "", -1), "\t", "", -1) + "\n"
	if buf.String() != exp {
		t.Errorf("Expected JSON: %v, got: %v", exp, buf.String())
	}
}

func TestReadHistogramValues(t *testing.T) {
	ms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		if r.RequestURI == "/state" {
			w.Write([]byte(stateTestData))
			return
		}

		if r.RequestURI == "/stats.json" {
			w.Write([]byte(statsTestData))
			return
		}

		u := "/histogram/1556290800/1556291400/300/" +
			"ae0f7f90-2a6b-481c-9cf5-21a31837020e/example1"
		if strings.HasPrefix(r.RequestURI, u) {
			w.Write([]byte(histogramTestData))
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
	res, err := sc.ReadHistogramValues(node,
		"ae0f7f90-2a6b-481c-9cf5-21a31837020e", "example1",
		300*time.Second, time.Unix(1556290800, 0),
		time.Unix(1556291200, 0))
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 2 {
		t.Fatalf("Expected length: 1, got: %v", len(res))
	}

	if res[0].Timestamp() != "1556290800" {
		t.Errorf("Expected timestamp: 1556290800, got: %v", res[0].Timestamp())
	}

	if res[0].Period.Seconds() != 300.0 {
		t.Errorf("Expected seconds: 300, got: %v", res[0].Period.Seconds())
	}

	if res[0].Data["+23e-004"] != 1 {
		t.Errorf("Expected data: 1, got: %v", res[0].Data["+23e-004"])
	}
}

func TestWriteHistogram(t *testing.T) {
	ms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		if r.RequestURI == "/state" {
			w.Write([]byte(stateTestData))
			return
		}

		if r.RequestURI == "/stats.json" {
			w.Write([]byte(statsTestData))
			return
		}

		if r.RequestURI == "/histogram/write" {
			rb := []HistogramData{}
			if err := json.NewDecoder(r.Body).Decode(&rb); err != nil {
				w.WriteHeader(500)
				t.Error("Unable to decode JSON data")
				return
			}

			if len(rb) < 1 {
				w.WriteHeader(500)
				t.Error("Invalid request")
				return
			}

			exp := "ae0f7f90-2a6b-481c-9cf5-21a31837020e"
			if rb[0].ID != exp {
				w.WriteHeader(500)
				t.Errorf("Expected UUID: %v, got: %v", exp, rb[0].ID)
				return
			}

			w.Write([]byte(histTestData))
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

	v := []HistogramData{}
	err = json.NewDecoder(bytes.NewBufferString(histTestData)).Decode(&v)
	if err != nil {
		t.Fatalf("Unable to encode JSON %v", err)
	}

	node := &SnowthNode{url: u}
	err = sc.WriteHistogram(node, v...)
	if err != nil {
		t.Fatal(err)
	}
}
