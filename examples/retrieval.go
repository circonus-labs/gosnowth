package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/circonus-labs/gosnowth"
	"github.com/google/uuid"
)

// ExampleReadNNT demonstrates how to read NNT values from a given snowth node.
func ExampleReadNNT() {
	// Create a new client.
	cfg := gosnowth.NewConfig(SnowthServers...)

	client, err := gosnowth.NewClient(context.Background(), cfg)
	if err != nil {
		log.Fatalf("failed to create snowth client: %v", err)
	}

	// Write text data in order to read back the data.
	id := uuid.New().String()
	// WriteNNT takes in a node and variadic of NNTPartsData entries.
	if err := client.WriteNNT([]gosnowth.NNTData{{
		Metric: "test-metric",
		ID:     id,
		Offset: (time.Now().Unix() / 60) * 60,
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
		log.Fatalf("failed to write text data: %v", err)
	}

	data, err := client.ReadNNTValues(time.Now().Add(-60*time.Second),
		time.Now().Add(60*time.Second), 60, "count", id, "test-metric")
	if err != nil {
		log.Fatalf("failed to read nnt data: %v", err)
	}

	log.Printf("%+v\n", data)
	data1, err := client.ReadNNTAllValues(time.Now().Add(-60*time.Second),
		time.Now().Add(60*time.Second), 60, id, "test-metric")
	log.Printf("%+v\n", data1)
}

// ExampleReadText demonstrates how to read text values from a given snowth
// node.
func ExampleReadText() {
	// Create a new client.
	cfg := gosnowth.NewConfig(SnowthServers...)

	client, err := gosnowth.NewClient(context.Background(), cfg)
	if err != nil {
		log.Fatalf("failed to create snowth client: %v", err)
	}

	// Write text data in order to read back the data.
	id := uuid.New().String()
	if err := client.WriteText([]gosnowth.TextData{{
		Metric: "test-text-metric2",
		ID:     id,
		Offset: strconv.FormatInt(time.Now().Unix(), 10),
		Value:  "a_text_data_value",
	}}); err != nil {
		log.Fatalf("failed to write text data: %v", err)
	}

	data, err := client.ReadTextValues(id, "test-text-metric2",
		time.Now().Add(-60*time.Second), time.Now().Add(60*time.Second))
	if err != nil {
		log.Fatalf("failed to read text data: %v", err)
	}

	log.Printf("%+v\n", data)
}

func ExampleFetchQuery() {
	host := os.Getenv("SNOWTH_URL")
	if host == "" {
		return
	}

	sc, err := gosnowth.NewClient(context.Background(),
		&gosnowth.Config{Servers: []string{host}})
	if err != nil {
		log.Fatal("Unable to create snowth client", err)
	}

	metrics, err := sc.FindTags(1, "data-service.*.oneMinuteRate",
		&gosnowth.FindTagsOptions{
			Start:     time.Unix(1611696384, 0),
			End:       time.Unix(1611696584, 0),
			Activity:  0,
			Latest:    0,
			CountOnly: 0,
			Limit:     -1,
		})
	if err != nil {
		log.Fatal(err)
	}

	responses := []*gosnowth.DF4Response{}
	for _, metric := range metrics.Items {
		res, err := sc.FetchValues(&gosnowth.FetchQuery{
			Start:  time.Unix(1611696384, 0),
			Period: time.Second,
			Count:  200,
			Streams: []gosnowth.FetchStream{{
				UUID:      metric.UUID,
				Name:      metric.MetricName,
				Kind:      metric.Type,
				Label:     metric.MetricName,
				Transform: "none",
			}},
			Reduce: []gosnowth.FetchReduce{{
				Label:  "test",
				Method: "pass",
			}},
		})
		if err != nil {
			log.Fatal(err)
		}

		responses = append(responses, res)
	}
}

func ExampleFetchQueryMultiStream() {
	host := os.Getenv("SNOWTH_URL")
	if host == "" {
		return
	}

	sc, err := gosnowth.NewClient(context.Background(),
		&gosnowth.Config{Servers: []string{host}})
	if err != nil {
		log.Fatal("Unable to create snowth client", err)
	}

	_, err = sc.FetchValues(&gosnowth.FetchQuery{
		Start:  time.Unix(1611696384, 0),
		Period: time.Second,
		Count:  200,
		Streams: []gosnowth.FetchStream{{
			UUID:      "11223344-5566-7788-9900-aabbccddeeff",
			Name:      "cpu.usage",
			Kind:      "numeric",
			Label:     "cpu.usage",
			Transform: "none",
		}, {
			UUID:      "11223344-5566-7788-9900-aabbccddeeff",
			Name:      "cpu_usage",
			Kind:      "numeric",
			Label:     "cpu_usage",
			Transform: "none",
		}},
		Reduce: []gosnowth.FetchReduce{{
			Label:  "test",
			Method: "average",
		}},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCAQLQuery() {
	host := os.Getenv("SNOWTH_URL")
	if host == "" {
		return
	}

	sc, err := gosnowth.NewClient(context.Background(),
		&gosnowth.Config{Servers: []string{host}})
	if err != nil {
		log.Fatal("Unable to create snowth client", err)
	}

	_, err = sc.GetCAQLQuery(&gosnowth.CAQLQuery{
		AccountID: 1,
		Query:     `(find("orders_per_second", "and(check_name:zmon.check.123)") | aggregate:sum() ) / 60`,
		Start:     1611742469,
		End:       1611753269,
		Period:    60,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = sc.GetCAQLQuery(&gosnowth.CAQLQuery{
		AccountID: 1,
		Query:     `find("orders_per_second", "and(check_name:zmon.check.123)") | aggregate:mean(2m)`,
		Start:     1611742469,
		End:       1611753269,
		Period:    60,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleGetCheckTags() {
	host := os.Getenv("SNOWTH_URL")
	if host == "" {
		return
	}

	sc, err := gosnowth.NewClient(context.Background(),
		&gosnowth.Config{Servers: []string{host}})
	if err != nil {
		log.Fatal("Unable to create snowth client", err)
	}

	_, err = sc.GetCheckTags("e312a0cb-dbe9-445d-8346-13b0ae6a3382")
	if err != nil {
		log.Fatal(err)
	}
}
