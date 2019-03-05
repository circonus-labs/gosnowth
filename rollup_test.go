package gosnowth

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

const rollupTestData = "[[1,1]]"

const rollupAllTestData = `[
    [
        1529509020,
        {
            "count": 1,
            "value": 0,
            "stddev": 0,
            "derivative": 0,
            "derivative_stddev": 0,
            "counter": 0,
            "counter_stddev": 0,
            "derivative2": 0,
            "derivative2_stddev": 0,
            "counter2": 0,
            "counter2_stddev": 0
        }
    ],
    [
        1529509080,
        {
            "count": 1,
            "value": 0,
            "stddev": 0,
            "derivative": 0,
            "derivative_stddev": 0,
            "counter": 0,
            "counter_stddev": 0,
            "derivative2": 0,
            "derivative2_stddev": 0,
            "counter2": 0,
            "counter2_stddev": 0
        }
    ],
    [
        1529509140,
        {
            "count": 1,
            "value": 0,
            "stddev": 0,
            "derivative": 0,
            "derivative_stddev": 0,
            "counter": 0,
            "counter_stddev": 0,
            "derivative2": 0,
            "derivative2_stddev": 0,
            "counter2": 0,
            "counter2_stddev": 0
        }
    ]
]`

func TestReadRollupValues(t *testing.T) {
	ms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		if r.RequestURI == "/state" {
			w.Write([]byte(stateTestData))
			return
		}

		u := "/rollup/fc85e0ab-f568-45e6-86ee-d7443be8277d/" +
			"online%7CST%5Btest%3Atest%5D?start_ts=1529509020" +
			"&end_ts=1529509201&rollup_span=1s"
		if strings.HasPrefix(r.RequestURI, u) {
			w.Write([]byte(rollupTestData))
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
	res, err := sc.ReadRollupValues(node,
		"fc85e0ab-f568-45e6-86ee-d7443be8277d", "online", []string{"test:test"},
		time.Second, time.Unix(1529509020, 0), time.Unix(1529509200, 0))
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 1 {
		t.Fatalf("Expected results: 1, got: %v", len(res))
	}

	if res[0].Value != 1 {
		t.Errorf("Expected value: 1, got: %v", res[0].Value)
	}
}
