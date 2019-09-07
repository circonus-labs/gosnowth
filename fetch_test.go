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

const fetchTestQuery = `{
	"start":1555616700,
	"period":300,
	"count":10,
	"streams":[
		{
			"uuid":"11223344-5566-7788-9900-aabbccddeeff",
			"name":"test",
			"kind":"numeric",
			"label":"test",
			"transform":"average"
		}
	],
	"reduce": [{"label":"test","method":"test"}]
}`

func TestFetchQueryMarshaling(t *testing.T) {
	v := &FetchQuery{}
	err := json.NewDecoder(bytes.NewBufferString(fetchTestQuery)).Decode(&v)
	if err != nil {
		t.Fatal(err)
	}

	if v.Timestamp() != "1555616700" {
		t.Errorf("Expected timestamp: 1555616700, got: %v", v.Timestamp())
	}

	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(&v)
	if err != nil {
		t.Fatal(err)
	}

	exp := strings.Replace(strings.Replace(strings.Replace(fetchTestQuery,
		" ", "", -1), "\n", "", -1), "\t", "", -1) + "\n"
	if buf.String() != exp {
		t.Errorf("Expected JSON: %v, got: %v", exp, buf.String())
	}
}

func TestFetchQuery(t *testing.T) {
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

		if strings.HasPrefix(r.RequestURI,
			"/fetch") {
			w.Write([]byte(testDF4Response))
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
	res, err := sc.FetchValues(node, &FetchQuery{
		Start:  time.Unix(0, 0),
		Period: 300 * time.Second,
		Count:  3,
		Streams: []FetchStream{{
			UUID:      "11223344-5566-7788-9900-aabbccddeeff",
			Name:      "test",
			Kind:      "numeric",
			Label:     "test",
			Transform: "none",
		}},
		Reduce: []FetchReduce{{
			Label:  "test",
			Method: "average",
		}},
	})
	if err != nil {
		t.Fatal(err)
	}

	if res.Head.Count != 3 {
		t.Fatalf("Expected header count: 3, got: %v", res.Head.Count)
	}

	if len(res.Data) != 1 {
		t.Fatalf("Expected data length: 1, got: %v", len(res.Data))
	}

	v, ok := res.Data[0][0].(float64)
	if !ok {
		t.Fatal("Unexpected data type")
	}

	if v != 1.0 {
		t.Errorf("Expected value: 1, got: %v", v)
	}

	if len(res.Meta) != 1 {
		t.Fatalf("Expected meta length: 1, got: %v", len(res.Meta))
	}

	if res.Meta[0].Label != "test" {
		t.Errorf("Expected meta label: test, got: %v", res.Meta[0].Label)
	}
}
