package gosnowth

import (
	"bytes"
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

		if strings.HasPrefix(r.RequestURI, "/raw") {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Error("Unable to read request body")
			}

			if string(b) != "test" {
				t.Errorf("Expected request body: test, got: %v", string(b))
			}

			w.WriteHeader(200)
			return
		}

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
	err = sc.WriteRaw(node, bytes.NewBufferString("test"), true, 1)
	if err != nil {
		t.Fatal(err)
	}
}
