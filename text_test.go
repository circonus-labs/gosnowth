package gosnowth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

const textTestData = `[[1380000000,"hello"],[1380000300,"world"]]`

func TestTextValue(t *testing.T) {
	t.Parallel()

	tvr := TextValueResponse{}
	if err := json.Unmarshal([]byte(textTestData), &tvr); err != nil {
		t.Error("error unmarshaling: ", err)
	}
}

func TestReadTextValuesFindMetricNode(t *testing.T) {
	t.Parallel()

	ms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request,
	) {
		if r.RequestURI == "/state" {
			_, _ = w.Write([]byte(stateTestData))

			return
		}

		if r.RequestURI == "/stats.json" {
			_, _ = w.Write([]byte(statsTestData))

			return
		}

		if strings.HasPrefix(r.RequestURI,
			"/read/1/2/3aa57ac2-28de-4ec4-aa3d-ed0ddd48fa4d/test") {
			_, _ = w.Write([]byte(textTestData))

			return
		}

		w.WriteHeader(500)
	}))

	defer ms.Close()
	sc, err := NewSnowthClient(false, ms.URL)
	if err != nil {
		t.Fatal("Unable to create snowth client", err)
	}

	res, err := sc.ReadTextValues("3aa57ac2-28de-4ec4-aa3d-ed0ddd48fa4d",
		"test", time.Unix(1, 0), time.Unix(2, 0))
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 2 {
		t.Fatalf("Expected result length: 2, got: %v", len(res))
	}

	if res[0].Value == nil {
		t.Fatal("Expected value: not nil, got: nil")
	}

	if *res[0].Value != "hello" {
		t.Errorf("Expected value: hello, got: %v", *res[0].Value)
	}
}

func TestReadTextValues(t *testing.T) {
	t.Parallel()

	ms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request,
	) {
		if r.RequestURI == "/state" {
			_, _ = w.Write([]byte(stateTestData))

			return
		}

		if r.RequestURI == "/stats.json" {
			_, _ = w.Write([]byte(statsTestData))

			return
		}

		if strings.HasPrefix(r.RequestURI,
			"/read/1/2/3aa57ac2-28de-4ec4-aa3d-ed0ddd48fa4d/test") {
			_, _ = w.Write([]byte(textTestData))

			return
		}

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
	res, err := sc.ReadTextValues("3aa57ac2-28de-4ec4-aa3d-ed0ddd48fa4d",
		"test", time.Unix(1, 0), time.Unix(2, 0), node)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 2 {
		t.Fatalf("Expected result length: 2, got: %v", len(res))
	}

	if res[0].Value == nil {
		t.Fatal("Expected value: not nil, got: nil")
	}

	if *res[0].Value != "hello" {
		t.Errorf("Expected value: hello, got: %v", *res[0].Value)
	}
}

func TestWriteText(t *testing.T) {
	t.Parallel()

	ms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request,
	) {
		if r.RequestURI == "/state" {
			_, _ = w.Write([]byte(stateTestData))

			return
		}

		if r.RequestURI == "/stats.json" {
			_, _ = w.Write([]byte(statsTestData))

			return
		}

		if strings.HasPrefix(r.RequestURI,
			"/write/text") {
			w.WriteHeader(200)

			return
		}

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
	err = sc.WriteText([]TextData{{
		Metric: "test",
		ID:     "3aa57ac2-28de-4ec4-aa3d-ed0ddd48fa4d",
		Offset: "1",
		Value:  "test",
	}}, node)
	if err != nil {
		t.Fatal(err)
	}
}
