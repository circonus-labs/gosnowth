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

const testPromQLError = `{
    "status": "error",
	"errorType": "test",
	"error": "test"
}`

const testPromQLResponse = `{
    "status": "success",
    "data": {
        "result": [
            {
                "values": [
                    [
                        1676388600,
                        "3568"
                    ]
                ],
                "metric": {
                    "__name__": "bytes",
                    "__check_uuid": "09fc1c4e-8540-49a8-a109-8895553718fc"
                }
            },
            {
                "values": [
                    [
                        1676388600,
                        "0"
                    ]
                ],
                "metric": {
                    "__name__": "bytes",
                    "__check_uuid": "fedbe76c-56df-4f3a-87b6-1eb787b89361"
                }
            }
        ],
        "resulttype": "matrix"
    }
}`

func TestPromQLRangeQuery(t *testing.T) {
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
				_, _ = w.Write([]byte(testPromQLError))

				return
			}

			if len(b) == 0 {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(testPromQLError))

				return
			}

			if strings.Contains(string(b), "127") {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(testCAQLError))

				return
			}

			_, _ = w.Write([]byte(testPromQLResponse))

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

	res, err := sc.PromQLRangeQuery(&PromQLRangeQuery{
		AccountID: "1",
		Query:     "test",
		Start:     "0",
		End:       "900",
		Step:      "300",
	}, node)
	if err != nil {
		t.Fatal(err)
	}

	if res.Data == nil {
		t.Fatalf("Expected data: 2, got: %v", res.Data)
	}

	res, err = sc.PromQLRangeQuery(&PromQLRangeQuery{
		AccountID: "1",
		Query:     "test",
		Start:     "0",
		End:       "900",
		Step:      "127",
	}, node)
	if err == nil {
		t.Fatal("Expected PromQL error response")
	}

	if res.ErrorType != "caql" {
		t.Errorf("Expected error type: caql, got: %v", res.ErrorType)
	}

	exp := "Function not found: histograms"

	if res.Error != exp {
		t.Errorf("Expected error: %v, got: %v", exp, res.Error)
	}
}
