package gosnowth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

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
