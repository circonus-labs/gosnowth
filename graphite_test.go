package gosnowth

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const graphiteMetricsTestData = `[
  {
    "leaf": true,
    "name": "11223344-5566-7788-9900-aabbccddeeff.test;test=test",
    "leaf_data": {
      "uuid": "11223344-5566-7788-9900-aabbccddeeff",
      "name": "test|ST[test:test]",
      "egress_function": "avg"
    }
  }
]`

const graphiteDatapointsTestData = `{
  "from": 0,
  "to": 900,
  "step": 300,
  "series": {
    "11223344-5566-7788-9900-aabbccddeeff.test": [
      null,
      0.1,
      null
    ]
  }
}`

func TestGraphiteFindMetrics(t *testing.T) {
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
			"/graphite/1/test/metrics/find?query=test") {
			w.Header().Set("X-Snowth-Search-Result-Count", "1")
			_, _ = w.Write([]byte(graphiteMetricsTestData))

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

	res, err := sc.GraphiteFindMetrics(1, "test", "test", nil, node)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 1 {
		t.Fatalf("Expected result length: 1, got: %v", len(res))
	}

	exp := "11223344-5566-7788-9900-aabbccddeeff.test;test=test"
	if res[0].Name != exp {
		t.Errorf("Expected metric name: %v, got: %v", exp, res[0].Name)
	}
}

func TestGraphiteFindTags(t *testing.T) {
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
			"/graphite/1/test/tags/find?query=test") {
			w.Header().Set("X-Snowth-Search-Result-Count", "1")
			_, _ = w.Write([]byte(graphiteMetricsTestData))

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

	res, err := sc.GraphiteFindTags(1, "test", "test", nil, node)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 1 {
		t.Fatalf("Expected result length: 1, got: %v", len(res))
	}

	exp := "11223344-5566-7788-9900-aabbccddeeff.test;test=test"
	if res[0].Name != exp {
		t.Errorf("Expected metric name: %v, got: %v", exp, res[0].Name)
	}
}

func TestGraphiteGetDatapoints(t *testing.T) {
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
			"/graphite/1/test/series_multi") {
			w.Header().Set("X-Snowth-Search-Result-Count", "1")
			_, _ = w.Write([]byte(graphiteDatapointsTestData))

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

	res, err := sc.GraphiteGetDatapoints(1, "test", &GraphiteLookup{
		Start: 0,
		End:   900,
		Names: []string{"11223344-5566-7788-9900-aabbccddeeff.test"},
	}, nil, node)
	if err != nil {
		t.Fatal(err)
	}

	if res.From != 0 {
		t.Errorf("Expected from: 0, got: %v", res.From)
	}

	if res.To != 900 {
		t.Errorf("Expected to: 900, got: %v", res.To)
	}

	if res.Step != 300 {
		t.Errorf("Expected step: 300, got: %v", res.Step)
	}

	rv := res.Series["11223344-5566-7788-9900-aabbccddeeff.test"][0]
	if rv != nil {
		t.Errorf("Expected null value, got: %v", rv)
	}

	rv = res.Series["11223344-5566-7788-9900-aabbccddeeff.test"][1]
	if *rv != 0.1 {
		t.Errorf("Expected value: 0.1, got: %v", *rv)
	}
}
