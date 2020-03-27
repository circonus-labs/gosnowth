package gosnowth

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestWriteRaw(t *testing.T) {
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

		if strings.HasPrefix(r.RequestURI, "/raw") {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Error("Unable to read request body")
			}

			if string(b) == "test" {
				w.WriteHeader(200)
				w.Write([]byte(`{ "records": 0, "updated": 0, "misdirected": 0, "errors": 0 }`))
				return
			}

			w.WriteHeader(500)
			w.Write([]byte("invalid request body"))
			return
		}

		t.Errorf("Unexpected request: %v", r)
		w.WriteHeader(500)
		return
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
	_, err = sc.WriteRaw(bytes.NewBufferString("test"), true, 1, node)
	if err != nil {
		t.Fatal(err)
	}

	sc.SetRequestFunc(func(r *http.Request) error { return nil })
	_, err = sc.WriteRaw(bytes.NewBufferString("error"), true, 1, node)
	if err == nil {
		t.Fatal("Expected error response")
	}

	if !strings.Contains(err.Error(), "invalid request body") {
		t.Errorf("Unexpected error returned: %v", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = sc.WriteRawContext(ctx, bytes.NewBufferString("test"), true, 1, node)
	if err == nil {
		t.Fatal("Expected error response")
	}

	if !strings.Contains(err.Error(), "context") {
		t.Errorf("Expected context error, got: %v", err.Error())
	}
}
