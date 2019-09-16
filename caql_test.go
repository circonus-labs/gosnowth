package gosnowth

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestGetCAQLQuery(t *testing.T) {
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
			"/extension/lua/public/caql_v1?query=test") {
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
	res, err := sc.GetCAQLQuery(node, 1, &CAQLQuery{
		Query:  "test",
		Start:  0,
		End:    900,
		Period: 300,
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
