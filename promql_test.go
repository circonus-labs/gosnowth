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

const testPromQLInstantQueryResponse = `{
    "status": "success",
    "data": {
        "result": [
            {
                "value": [
                    [
                        1676388600,
                        "3568"
                    ]
                ],
                "metric": {
                    "__name__": "bytes",
                    "__check_uuid": "09fc1c4e-8540-49a8-a109-8895553718fc"
                }
            }
        ],
        "resulttype": "vector"
    }
}`

const testPromQLRangeQueryResponse = `{
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

func TestPromQLInstantQuery(t *testing.T) {
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

			_, _ = w.Write([]byte(testPromQLInstantQueryResponse))

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

	res, err := sc.PromQLInstantQuery(&PromQLInstantQuery{
		AccountID: "1",
		Query:     "test",
		Time:      "300.123",
	}, node)
	if err != nil {
		t.Fatal(err)
	}

	if res.Data == nil {
		t.Fatalf("Expected data, got: %v", res.Data)
	}

	res, err = sc.PromQLInstantQuery(&PromQLInstantQuery{
		AccountID: "1",
		Query:     "test",
		Time:      "127",
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

			_, _ = w.Write([]byte(testPromQLRangeQueryResponse))

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
		t.Fatalf("Expected data, got: %v", res.Data)
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

func TestConvertSeriesSelector(t *testing.T) {
	tests := []struct {
		name string
		sel  string
		exp  string
	}{{
		name: "series_selector",
		sel:  `test{test="test"}`,
		exp:  `and(__name:b"dGVzdA==",and(b"dGVzdA==":b"dGVzdA=="))`,
	}, {
		name: "series_selector_two",
		sel:  `test{test="test",test1="test1"}`,
		exp: `and(__name:b"dGVzdA==",and(b"dGVzdA==":b"dGVzdA=="),` +
			`and(b"dGVzdDE=":b"dGVzdDE="))`,
	}, {
		name: "series_selector_not",
		sel:  `test{test="test",test1!="test1"}`,
		exp: `and(__name:b"dGVzdA==",and(b"dGVzdA==":b"dGVzdA=="),` +
			`not(b"dGVzdDE=":b"dGVzdDE="))`,
	}, {
		name: "series_selector_regex",
		sel:  `test{test=~"test",test1!~"test1"}`,
		exp: `and(__name:b"dGVzdA==",and(b"dGVzdA==":b/dGVzdA==/),` +
			`not(b"dGVzdDE=":b/dGVzdDE=/))`,
	}}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r, err := ConvertSeriesSelector(tt.sel)
			if err != nil {
				t.Fatal(err)
			}

			if r != tt.exp {
				t.Errorf("Expected query: %v, got: %v", tt.exp, r)
			}
		})
	}
}

const promQLSeriesTestData = `[
	{
		"uuid": "3aa57ac2-28de-4ec4-aa3d-ed0ddd48fa4d",
		"check_tags": [
			"test:test",
			"__check_id:1"
		],
		"metric_name": "test|ST[test1:test1]",
		"type": "numeric,histogram",
		"activity": [
			[
				1555610400,
				1556588400
			],
			[
				1556625600,
				1561848300
			]
		],
		"latest": {
			"numeric": [
				[1561848300000, 1]
			]
		},
		"account_id": 1
	}
]`

func TestPromQLSeriesQuery(t *testing.T) {
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

		w.Header().Set("X-Snowth-Search-Result-Count", "1")
		_, _ = w.Write([]byte(promQLSeriesTestData))

		return
	}))

	defer ms.Close()

	sc, err := NewClient(context.Background(),
		&Config{Servers: []string{ms.URL}})
	if err != nil {
		t.Fatal("Unable to create snowth client", err)
	}

	u, err := url.Parse(ms.URL)
	if err != nil {
		t.Fatal("Invalid test URL")
	}

	node := &SnowthNode{url: u}

	res, err := sc.PromQLSeriesQuery(&PromQLSeriesQuery{
		Match:     []string{"test"},
		Start:     "0",
		End:       "300",
		AccountID: "1",
	}, node)
	if err != nil {
		t.Fatal(err)
	}

	if res.Data == nil {
		t.Fatalf("Expected data, got: %v", res.Data)
	}

	dm, ok := res.Data.([]map[string]string)
	if !ok {
		t.Fatalf("Invalid type for result data: %T", res.Data)
	}

	if len(dm) < 1 {
		t.Fatalf("Expected data length: 1, got: %v", len(dm))
	}

	if dm[0]["__name__"] != "test" {
		t.Errorf("Expected __name__: test, got: %v", dm[0]["__name__"])
	}

	if dm[0]["test"] != "test" {
		t.Errorf("Expected test: test, got: %v", dm[0]["test"])
	}

	if dm[0]["test1"] != "test1" {
		t.Errorf("Expected test1: test1, got: %v", dm[0]["test1"])
	}
}

func TestPromQLLabelQuery(t *testing.T) {
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

		w.Header().Set("X-Snowth-Search-Result-Count", "1")
		_, _ = w.Write([]byte(tagCatsValsTestData))

		return
	}))

	defer ms.Close()

	sc, err := NewClient(context.Background(),
		&Config{Servers: []string{ms.URL}})
	if err != nil {
		t.Fatal("Unable to create snowth client", err)
	}

	u, err := url.Parse(ms.URL)
	if err != nil {
		t.Fatal("Invalid test URL")
	}

	node := &SnowthNode{url: u}

	res, err := sc.PromQLLabelQuery(&PromQLLabelQuery{
		Match:     []string{"test"},
		Start:     "0",
		End:       "300",
		AccountID: "1",
	}, node)
	if err != nil {
		t.Fatal(err)
	}

	if res.Data == nil {
		t.Fatalf("Expected data, got: %v", res.Data)
	}

	ds, ok := res.Data.([]string)
	if !ok {
		t.Fatalf("Invalid type for result data: %T", res.Data)
	}

	if len(ds) != 3 {
		t.Fatalf("Expected data length: 3, got: %v", len(ds))
	}

	if ds[0] != "__name__" {
		t.Errorf("Expected data: __name__, got: %v", ds[0])
	}

	if ds[1] != "test" {
		t.Errorf("Expected data: test, got: %v", ds[1])
	}
}

func TestPromQLLabelValuesQuery(t *testing.T) {
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

		w.Header().Set("X-Snowth-Search-Result-Count", "1")
		_, _ = w.Write([]byte(tagCatsValsTestData))

		return
	}))

	defer ms.Close()

	sc, err := NewClient(context.Background(),
		&Config{Servers: []string{ms.URL}})
	if err != nil {
		t.Fatal("Unable to create snowth client", err)
	}

	u, err := url.Parse(ms.URL)
	if err != nil {
		t.Fatal("Invalid test URL")
	}

	node := &SnowthNode{url: u}

	res, err := sc.PromQLLabelValuesQuery("test",
		&PromQLLabelQuery{
			Match:     []string{"test"},
			Start:     "0",
			End:       "300",
			AccountID: "1",
		}, node)
	if err != nil {
		t.Fatal(err)
	}

	if res.Data == nil {
		t.Fatalf("Expected data, got: %v", res.Data)
	}

	ds, ok := res.Data.([]string)
	if !ok {
		t.Fatalf("Invalid type for result data: %T", res.Data)
	}

	if len(ds) != 2 {
		t.Fatalf("Expected data length: 2, got: %v", len(ds))
	}

	if ds[0] != "test" {
		t.Errorf("Expected data: test, got: %v", ds[0])
	}
}

func TestPromQLMetadataQuery(t *testing.T) {
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

		w.Header().Set("X-Snowth-Search-Result-Count", "1")
		_, _ = w.Write([]byte(promQLSeriesTestData))

		return
	}))

	defer ms.Close()

	sc, err := NewClient(context.Background(),
		&Config{Servers: []string{ms.URL}})
	if err != nil {
		t.Fatal("Unable to create snowth client", err)
	}

	u, err := url.Parse(ms.URL)
	if err != nil {
		t.Fatal("Invalid test URL")
	}

	node := &SnowthNode{url: u}

	res, err := sc.PromQLMetadataQuery(&PromQLMetadataQuery{
		Limit:     "2",
		Metric:    "test",
		AccountID: "1",
	}, node)
	if err != nil {
		t.Fatal(err)
	}

	if res.Data == nil {
		t.Fatalf("Expected data, got: %v", res.Data)
	}

	dm, ok := res.Data.(map[string][]map[string]string)
	if !ok {
		t.Fatalf("Invalid type for result data: %T", res.Data)
	}

	if len(dm) != 1 {
		t.Fatalf("Expected data length: 1, got: %v", len(dm))
	}

	if dm["test"][0]["type"] != "histogram" {
		t.Errorf("Expected type: histogram, got: %v", dm["test"][0]["type"])
	}
}
