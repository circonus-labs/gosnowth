package gosnowth

import (
	"context"
	"io"
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
		"_debug": 0,
		"_id": 33545929,
		"_timeout": 15,
		"account_id": "1",
		"end_time": 1567566000,
		"expansion": [],
		"format": "DF4",
		"ignore_duration_limits": false,
		"method": "fetch",
		"min_prefill": 0,
		"period": 300,
		"prepare_results": "JSON",
		"q": "find:histograms(\"latency\",\"and(service:api)\")",
		"start_time": 1567500000
	},
	"success": false
}`

func TestGetCAQLQuery(t *testing.T) {
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

		if r.Method == "POST" && strings.HasPrefix(r.RequestURI,
			"/extension/lua/public/caql_v1") {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(testCAQLError))

				return
			}

			if len(b) == 0 {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(testCAQLError))

				return
			}

			if strings.Contains(string(b), "histograms") {
				w.WriteHeader(http.StatusBadGateway)
				_, _ = w.Write([]byte(testCAQLError))

				return
			}

			_, _ = w.Write([]byte(testFetchDF4Response))

			return
		}
	}))

	defer ms.Close()

	sc, err := NewClient(context.Background(),
		&Config{Servers: []string{ms.URL}})
	if err != nil {
		t.Fatal("Unable to create snowth client", err)
	}

	sc.SetRetries(1)
	sc.SetConnectRetries(1)

	u, err := url.Parse(ms.URL)
	if err != nil {
		t.Fatal("Invalid test URL")
	}

	node := &SnowthNode{url: u}

	res, err := sc.GetCAQLQuery(&CAQLQuery{
		AccountID: 1,
		Query:     "test",
		Start:     0,
		End:       900,
		Period:    300,
	}, node)
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

	_, err = sc.GetCAQLQuery(&CAQLQuery{
		AccountID: 1,
		Query:     "find:histograms()",
		Start:     0,
		End:       900,
		Period:    300,
	}, node)
	if err == nil {
		t.Fatal("Expected CAQL error response")
	}

	vErr, ok := err.(*CAQLError)
	if !ok {
		t.Fatalf("Expected error type: CAQLError, got: %T: %v", err, err)
	}

	exp := "Function not found: histograms"
	if vErr.Message() != exp {
		t.Errorf("Expected user error: %v, got: %v", exp,
			vErr.Message())
	}

	exp = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
		testCAQLError, " ", ""), "\t", ""), "\n", "")
	val := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
		vErr.Error(), " ", ""), "\t", ""), "\n", "")

	if val != exp {
		t.Errorf("Expected error JSON: %v, got: %v", exp, val)
	}
}
