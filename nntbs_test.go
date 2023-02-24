package gosnowth

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/circonus-labs/gosnowth/fb/nntbs"
	flatbuffers "github.com/google/flatbuffers/go"
)

func TestWriteNNTBSFlatbuffer(t *testing.T) {
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

		if strings.HasPrefix(r.RequestURI, "/nntbs") {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				t.Error("Unable to read request body")
			}

			if string(b)[4:8] == "CINN" {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{ "records": 1, "updated": 1, "misdirected": 0, "errors": 0 }`))

				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("invalid request body:"))
			_, _ = w.Write(b)

			return
		}

		t.Errorf("Unexpected request: %v", r)
		w.WriteHeader(http.StatusInternalServerError)
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

	builder := flatbuffers.NewBuilder(1024 * 1024)

	merge := &nntbs.NNTMergeT{
		Ops: []*nntbs.NNTMergeOpT{
			{
				Metric: &nntbs.MetricInfoT{
					MetricLocator: &nntbs.MetricLocatorT{
						CheckUuid:  []byte("11223344-5566-7788-9900-aabbccddeeff"),
						MetricName: "test.metric",
					},
					AccountId:     int32(1),
					CheckName:     "test.check",
					CheckCategory: metricSourceGraphite,
				},
				Nnt: []*nntbs.NNTT{
					{
						Epoch:      uint64(0),
						Apocalypse: uint64(300),
						Period:     uint32(60),
						Blocks: []*nntbs.NNTBlockT{
							{
								Data: []*nntbs.NNTDatumT{{
									Zero:             false,
									Count:            0,
									Stddev:           0.0,
									Derivative:       0.0,
									DerivativeStddev: 0.0,
									Counter:          0.0,
									CounterStddev:    0.0,
									Value: &nntbs.NumericValueT{
										Type: nntbs.NumericValueDoubleValue,
										Value: &nntbs.DoubleValueT{
											V: 1,
										},
									},
								}},
							},
						},
					},
				},
			},
		},
	}

	err = sc.WriteNNTBSFlatbuffer(merge, builder, node)
	if err != nil {
		t.Fatal(err)
	}
}
