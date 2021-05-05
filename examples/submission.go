package main

import (
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/google/uuid"
	"github.com/openhistogram/circonusllhist"

	"github.com/circonus-labs/gosnowth"
	"github.com/circonus-labs/gosnowth/fb/noit"
)

// ExampleSubmitText demonstrates how to submit a text metric to a node.
func ExampleSubmitText() {
	// Create a new client.
	cfg, err := gosnowth.NewConfig(SnowthServers...)
	if err != nil {
		log.Fatalf("failed to create snowth configuration: %v", err)
	}

	cfg.SetDiscover(true)
	client, err := gosnowth.NewClient(cfg)
	if err != nil {
		log.Fatalf("failed to create snowth client: %v", err)
	}

	// Write text data.
	id := uuid.New().String()
	// WriteText takes in a node and variadic of TextData entries.
	if err := client.WriteText([]gosnowth.TextData{{
		Metric: "test-text-metric2",
		ID:     id,
		Offset: strconv.FormatInt(time.Now().Unix(), 10),
		Value:  "a_text_data_value",
	}}); err != nil {
		log.Fatalf("failed to write text data: %v", err)
	}
}

// ExampleSubmitNNT demonstrates how to submit an NNT metric to a node.
func ExampleSubmitNNT() {
	// Create a new client.
	cfg, err := gosnowth.NewConfig(SnowthServers...)
	if err != nil {
		log.Fatalf("failed to create snowth configuration: %v", err)
	}

	cfg.SetDiscover(true)
	client, err := gosnowth.NewClient(cfg)
	if err != nil {
		log.Fatalf("failed to create snowth client: %v", err)
	}

	// Write NNT data to the node.
	id := uuid.New().String()
	if err := client.WriteNNT([]gosnowth.NNTData{{
		Metric: "test-metric",
		ID:     id,
		Offset: time.Now().Unix(),
		Count:  5, Value: 100,
		Parts: gosnowth.Parts{
			Period: 60,
			Data: []gosnowth.NNTPartsData{
				{Count: 1, Value: 100},
				{Count: 1, Value: 100},
				{Count: 1, Value: 100},
				{Count: 1, Value: 100},
				{Count: 1, Value: 100},
			},
		},
	}}); err != nil {
		log.Fatalf("failed to write nnt data: %v", err)
	}
}

// ExampleSubmitHistogram demonstrates how to submit histogram data to a node.
func ExampleSubmitHistogram() {
	// Create a new client.
	cfg, err := gosnowth.NewConfig(SnowthServers...)
	if err != nil {
		log.Fatalf("failed to create snowth configuration: %v", err)
	}

	cfg.SetDiscover(true)
	client, err := gosnowth.NewClient(cfg)
	if err != nil {
		log.Fatalf("failed to create snowth client: %v", err)
	}

	// Using the circonusllhist library, create a new histogram from a string
	// representation of a histogram.
	hist, err := circonusllhist.NewFromStrings([]string{
		"H[0.0e+00]=1",
		"H[1.0e+01]=1",
		"H[2.0e+02]=1",
		"H[3.0e+03]=1",
	}, false)

	// Write histogram data.
	id := uuid.New().String()
	if err := client.WriteHistogram([]gosnowth.HistogramData{{
		AccountID: 1,
		Metric:    "test-hist-metric",
		ID:        id,
		CheckName: "test",
		Offset:    time.Now().Unix(),
		Histogram: hist,
		Period:    60,
	}}); err != nil {
		log.Fatalf("failed to write histogram data: %v", err)
	}
}

func ExampleWriteRawMetricList(b *testing.B) {
	host := os.Getenv("SNOWTH_URL")
	if host == "" {
		return
	}

	sc, err := gosnowth.NewSnowthClient(false, host)
	if err != nil {
		log.Fatal("Unable to create snowth client", err)
	}

	builder := flatbuffers.NewBuilder(1024)

	list := &noit.MetricListT{
		Metrics: []*noit.MetricT{{
			Timestamp: 1589198300149,
			CheckName: "zmon.check.20406",
			CheckUuid: "e312a0cb-dbe9-445d-8346-13b0ae6a3382",
			AccountId: 1,
			Value: &noit.MetricValueT{
				Name:      "containers.recommender.restarts",
				Timestamp: 1589198300149,
				Value: &noit.MetricValueUnionT{
					Type: noit.MetricValueUnionIntValue,
					Value: &noit.IntValueT{
						Value: 0,
					},
				},
				Generation: 1,
				StreamTags: []string{
					"alias:gift-cards",
					"application:vertical-pod-autoscaler",
					"component:recommender",
					"namespace:kube-system",
					"version:v0.6.1-internal.7",
					"entity:pod-vpa-recommender-5f4fbfdc48-8gbnc-kube-system",
				},
			},
		}},
	}

	_, err = sc.WriteRawMetricList(list, builder)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleBulkWriteRawMetricList(b *testing.B) {
	host := os.Getenv("SNOWTH_URL")
	if host == "" {
		return
	}

	sc, err := gosnowth.NewSnowthClient(false, host)
	if err != nil {
		log.Fatal("Unable to create snowth client", err)
	}

	builder := flatbuffers.NewBuilder(1024)

	metrics := []*noit.MetricT{}

	for i := 0; i < 100; i++ {
		metrics = append(metrics, &noit.MetricT{
			Timestamp: 1589198300149 + uint64(i),
			CheckName: "zmon.check.20406",
			CheckUuid: "e312a0cb-dbe9-445d-8346-13b0ae6a3382",
			AccountId: 1,
			Value: &noit.MetricValueT{
				Name:      "containers.recommender.restarts",
				Timestamp: 1589198300149 + uint64(i),
				Value: &noit.MetricValueUnionT{
					Type: noit.MetricValueUnionIntValue,
					Value: &noit.IntValueT{
						Value: int32(i),
					},
				},
				Generation: 1,
				StreamTags: []string{
					"alias:gift-cards",
					"application:vertical-pod-autoscaler",
					"component:recommender",
					"namespace:kube-system",
					"version:v0.6.1-internal.7",
					"entity:pod-vpa-recommender-5f4fbfdc48-8gbnc-kube-system",
				},
			},
		})
	}

	list := &noit.MetricListT{Metrics: metrics}

	_, err = sc.WriteRawMetricList(list, builder)
	if err != nil {
		log.Fatal(err)
	}
}
