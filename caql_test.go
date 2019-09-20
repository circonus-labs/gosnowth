package gosnowth

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const testCAQLError = `{
	"locals": [],
	"method": "caql_v1",
	"trace": [],
	"user_error": {
		"message": "Function not found: histograms"
	},
	"status": "520 (User facing error)",
	"arguments": {
		"format": "DF4",
		"q": "find:histograms(\"latency\",\"and(service:api)\")",
		"period": 300,
		"_id": 33545929,
		"ignore_duration_limits": false,
		"prepare_results": "JSON",
		"account_id": "1",
		"method": "fetch",
		"start_time": 1567500000,
		"_timeout": 15,
		"min_prefill": 0,
		"_debug": 0,
		"expansion": [],
		"end_time": 1567566000
	},
	"success": false
}`

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

		if r.Method == "POST" && strings.HasPrefix(r.RequestURI,
			"/extension/lua/public/caql_v1") {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(testCAQLError))
				return
			}

			if strings.Contains(string(b), "histograms") {
				w.WriteHeader(502)
				w.Write([]byte(testCAQLError))
				return
			}

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
	res, err := sc.GetCAQLQuery(node, &CAQLQuery{
		AccountID: 1,
		Query:     "test",
		Start:     0,
		End:       900,
		Period:    300,
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

	_, err = sc.GetCAQLQuery(node, &CAQLQuery{
		AccountID: 1,
		Query:     "find:histograms()",
		Start:     0,
		End:       900,
		Period:    300,
	})
	if err == nil {
		t.Fatal("Expected CAQL error response")
	}

	vErr, ok := err.(*CAQLError)
	if !ok {
		t.Fatalf("Expected error type: CAQLError, got: %T", err)
	}

	exp := "Function not found: histograms"
	if vErr.UserError.Message != exp {
		t.Errorf("Expected user error: %v, got: %v", exp,
			vErr.UserError.Message)
	}

	exp = strings.Replace(strings.Replace(strings.Replace(testCAQLError,
		" ", "", -1), "\t", "", -1), "\n", "", -1)
	val := strings.Replace(strings.Replace(strings.Replace(vErr.Error(),
		" ", "", -1), "\t", "", -1), "\n", "", -1)
	if val != exp {
		t.Errorf("Expected error JSON: %v, got: %v", exp, val)
	}
}
