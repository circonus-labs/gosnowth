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
	tvr := TextValueResponse{}
	if err := json.Unmarshal([]byte(textTestData), &tvr); err != nil {
		t.Error("error unmarshalling: ", err)
	}
}

func TestReadTextValues(t *testing.T) {
	ms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		if r.RequestURI == "/state" {
			w.Write([]byte(stateTestData))
			return
		}

		if strings.HasPrefix(r.RequestURI,
			"/read/1/2/3aa57ac2-28de-4ec4-aa3d-ed0ddd48fa4d/test") {
			w.Write([]byte(textTestData))
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
	res, err := sc.ReadTextValues(node, time.Unix(1, 0), time.Unix(2, 0),
		"3aa57ac2-28de-4ec4-aa3d-ed0ddd48fa4d", "test")
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 2 {
		t.Fatalf("Expected result length: 2, got: %v", len(res))
	}

	if res[0].Value != "hello" {
		t.Errorf("Expected value: hello, got: %v", res[0].Value)
	}
}

func TestWriteText(t *testing.T) {
	ms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		if r.RequestURI == "/state" {
			w.Write([]byte(stateTestData))
			return
		}

		if strings.HasPrefix(r.RequestURI,
			"/write/text") {
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
	err = sc.WriteText(node, TextData{
		Metric: "test",
		ID:     "3aa57ac2-28de-4ec4-aa3d-ed0ddd48fa4d",
		Offset: "1",
		Value:  "test",
	})
	if err != nil {
		t.Fatal(err)
	}
}