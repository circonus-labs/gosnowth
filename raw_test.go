package gosnowth

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/circonus-labs/gosnowth/fb/noit"
)

func TestWriteRaw(t *testing.T) {
	ms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		if r.RequestURI == "/state" {
			_, _ = w.Write([]byte(stateTestData))
			return
		}

		if r.RequestURI == "/stats.json" {
			_, _ = w.Write([]byte(statsTestData))
			return
		}

		if strings.HasPrefix(r.RequestURI, "/raw") {
			buf, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Error("Unable to read request body")
			}

			if string(buf) == "test" {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`{ "records": 0, "updated": 0, "misdirected": 0, "errors": 0 }`))
				return
			}

			w.WriteHeader(500)
			_, _ = w.Write([]byte("invalid request body"))
			return
		}

		t.Errorf("Unexpected request: %v", r)
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
	_, err = sc.WriteRaw(bytes.NewBufferString("test"), true, 1, node)
	if err != nil {
		t.Fatal(err)
	}

	sc.SetRequestFunc(func(r *http.Request) error { return nil })
	_, err = sc.WriteRaw(bytes.NewBufferString("error"), true, 1, node)
	if err == nil {
		t.Fatal("Expected error response")
	}

	if !strings.Contains(err.Error(), "invalid request body") {
		t.Errorf("Unexpected error returned: %v", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = sc.WriteRawContext(ctx, bytes.NewBufferString("test"), true, 1, node)
	if err == nil {
		t.Fatal("Expected error response")
	}

	if !strings.Contains(err.Error(), "context") {
		t.Errorf("Expected context error, got: %v", err.Error())
	}
}

func TestReadRawNumericValues(t *testing.T) {
	ms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		if r.RequestURI == "/state" {
			_, _ = w.Write([]byte(stateTestData))
			return
		}

		if r.RequestURI == "/stats.json" {
			_, _ = w.Write([]byte(statsTestData))
			return
		}

		if strings.HasPrefix(r.RequestURI, "/raw") {
			w.WriteHeader(200)
			_, _ = w.Write([]byte(
				`[[1529509063064,0],[1529509122985,0],[1529509183764,0]]`))
			return
		}

		t.Errorf("Unexpected request: %v", r)
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
	_, err = sc.ReadRawNumericValues(
		time.Unix(1529509020, 0),
		time.Unix(1529509200, 0),
		"11223344-5566-7788-9900-aabbccddeeff",
		"test",
		node)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWriteRawMetricList(t *testing.T) {
	ms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		if r.RequestURI == "/state" {
			_, _ = w.Write([]byte(stateTestData))
			return
		}

		if r.RequestURI == "/stats.json" {
			_, _ = w.Write([]byte(statsTestData))
			return
		}

		if strings.HasPrefix(r.RequestURI, "/raw") {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Error("Unable to read request body")
			}

			if string(b)[4:8] == "CIML" {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`{ "records": 0, "updated": 0, "misdirected": 0, "errors": 0 }`))
				return
			}

			w.WriteHeader(500)
			_, _ = w.Write([]byte("invalid request body"))

			return
		}

		t.Errorf("Unexpected request: %v", r)
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

	builder := flatbuffers.NewBuilder(1024 * 1024)

	list := &noit.MetricListT{
		Metrics: []*noit.MetricT{{
			Timestamp: 1,
			CheckName: "test",
			CheckUuid: "11223344-5566-7788-9900-aabbccddeeff",
			AccountId: 1,
			Value: &noit.MetricValueT{
				Name:      "test",
				Timestamp: 1,
				Value: &noit.MetricValueUnionT{
					Type: noit.MetricValueUnionIntValue,
					Value: &noit.IntValueT{
						Value: 1,
					},
				},
				Generation: 1,
				StreamTags: []string{"test:test"},
			},
		}},
	}

	_, err = sc.WriteRawMetricList(list, builder, node)
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkWriteRawFlatbuffer(b *testing.B) {
	host := os.Getenv("SNOWTH_URL")
	if host == "" {
		return
	}

	b.StopTimer()

	sc, err := NewSnowthClient(false, host)
	if err != nil {
		b.Fatal("Unable to create snowth client", err)
	}

	builder := flatbuffers.NewBuilder(1024)

	list := &noit.MetricListT{
		Metrics: []*noit.MetricT{{
			Timestamp: uint64(time.Now().Unix()) * 1000,
			CheckName: "gosnowth-benchmark",
			CheckUuid: "e312a0cb-dbe9-445d-8346-13b0ae6a3382",
			AccountId: 1,
			Value: &noit.MetricValueT{
				Name:      "gosnowth-benchmark",
				Timestamp: uint64(time.Now().Unix()) * 1000,
				Value: &noit.MetricValueUnionT{
					Type: noit.MetricValueUnionIntValue,
					Value: &noit.IntValueT{
						Value: 1,
					},
				},
				Generation: 1,
				StreamTags: []string{"test:test"},
			},
		}},
	}

	builder.Reset()

	offset := noit.MetricListPack(builder, list)
	builder.FinishWithFileIdentifier(offset, []byte("CIML"))

	b.StartTimer()

	for n := 0; n < b.N; n++ {
		reader := bytes.NewReader(builder.FinishedBytes())

		_, err = sc.WriteRaw(reader, true, 1)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteRawMetricList(b *testing.B) {
	host := os.Getenv("SNOWTH_URL")
	if host == "" {
		return
	}

	b.StopTimer()

	sc, err := NewSnowthClient(false, host)
	if err != nil {
		b.Fatal("Unable to create snowth client", err)
	}

	builder := flatbuffers.NewBuilder(1024)

	list := &noit.MetricListT{
		Metrics: []*noit.MetricT{{
			Timestamp: uint64(time.Now().Unix()) * 1000,
			CheckName: "gosnowth-benchmark",
			CheckUuid: "e312a0cb-dbe9-445d-8346-13b0ae6a3382",
			AccountId: 1,
			Value: &noit.MetricValueT{
				Name:      "gosnowth-benchmark",
				Timestamp: uint64(time.Now().Unix()) * 1000,
				Value: &noit.MetricValueUnionT{
					Type: noit.MetricValueUnionIntValue,
					Value: &noit.IntValueT{
						Value: 1,
					},
				},
				Generation: 1,
				StreamTags: []string{"test:test"},
			},
		}},
	}

	b.StartTimer()

	for n := 0; n < b.N; n++ {
		_, err = sc.WriteRawMetricList(list, builder)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteRawMetricListLocal(b *testing.B) {
	b.StopTimer()

	ms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		if r.RequestURI == "/state" {
			_, _ = w.Write([]byte(stateTestData))
			return
		}

		if r.RequestURI == "/stats.json" {
			_, _ = w.Write([]byte(statsTestData))
			return
		}

		if strings.HasPrefix(r.RequestURI, "/raw") {
			buf, err := ioutil.ReadAll(r.Body)
			if err != nil {
				b.Error("Unable to read request body")
			}

			if string(buf)[4:8] == "CIML" {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`{ "records": 0, "updated": 0, "misdirected": 0, "errors": 0 }`))
				return
			}

			w.WriteHeader(500)
			_, _ = w.Write([]byte("invalid request body"))

			return
		}

		b.Errorf("Unexpected request: %v", r)
		w.WriteHeader(500)
	}))

	defer ms.Close()
	sc, err := NewSnowthClient(false, ms.URL)
	if err != nil {
		b.Fatal("Unable to create snowth client", err)
	}

	u, err := url.Parse(ms.URL)
	if err != nil {
		b.Fatal("Invalid test URL")
	}

	node := &SnowthNode{url: u}

	builder := flatbuffers.NewBuilder(1024)

	list := &noit.MetricListT{
		Metrics: []*noit.MetricT{{
			Timestamp: uint64(time.Now().Unix()) * 1000,
			CheckName: "gosnowth-benchmark",
			CheckUuid: "e312a0cb-dbe9-445d-8346-13b0ae6a3382",
			AccountId: 1,
			Value: &noit.MetricValueT{
				Name:      "gosnowth-benchmark",
				Timestamp: uint64(time.Now().Unix()) * 1000,
				Value: &noit.MetricValueUnionT{
					Type: noit.MetricValueUnionIntValue,
					Value: &noit.IntValueT{
						Value: 1,
					},
				},
				Generation: 1,
				StreamTags: []string{"test:test"},
			},
		}},
	}

	b.StartTimer()

	for n := 0; n < b.N; n++ {
		_, err = sc.WriteRawMetricList(list, builder, node)
		if err != nil {
			b.Fatal(err)
		}
	}
}
