package gosnowth

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRebuildActivity(t *testing.T) {
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

		if strings.HasPrefix(r.RequestURI, "/surrogate/activity_rebuild") {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Error("Unable to read request body")
			}

			if string(b) == "[]\n" {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`{ "records": 0, "updated": 0, "misdirected": 0, "errors": 0 }`))
				return
			}

			w.WriteHeader(500)
			_, _ = w.Write([]byte("invalid request body"))
			return
		}

		t.Errorf("Unexpected request: %v", r)
		w.WriteHeader(500)
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
	_, err = sc.RebuildActivity(node, []RebuildActivityRequest{})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = sc.RebuildActivityContext(ctx, node, []RebuildActivityRequest{})
	if err == nil {
		t.Fatal("Expected error response")
	}

	if !strings.Contains(err.Error(), "context") {
		t.Errorf("Expected context error, got: %v", err.Error())
	}
}
