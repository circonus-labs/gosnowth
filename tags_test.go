// Package gosnowth contains an IRONdb client library written in Go.
package gosnowth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

const tagsCountTestData = `{"count":22,"estimate":false}`

const tagsTestData = `[
	{
		"uuid": "3aa57ac2-28de-4ec4-aa3d-ed0ddd48fa4d",
		"check_tags": [
			"test:test",
			"__check_id:1"
		],
		"metric_name": "test",
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

func TestFindTagsJSON(t *testing.T) {
	fti := &FindTagsItem{
		UUID:       "11223344-5566-7788-9900-aabbccddeeff",
		CheckTags:  []string{"test:test"},
		MetricName: "test|ST[test:test]",
		Type:       "numeric",
		AccountID:  1,
		Activity:   [][]int64{{1, 1}, {2, 1}},
		Latest: &FindTagsLatest{
			Numeric:   []FindTagsLatestNumeric{{1, float64Ptr(1)}},
			Text:      []FindTagsLatestText{{1, nil}},
			Histogram: []FindTagsLatestHistogram{{1, stringPtr("AAEoAgAB")}},
		},
	}

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(&fti)
	if err != nil {
		t.Fatal(err)
	}

	var r *FindTagsItem
	err = json.NewDecoder(buf).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if fti.Latest == nil {
		t.Fatal("Expected latest data to not be nil")
	}

	if *fti.Latest.Numeric[0].Value != *r.Latest.Numeric[0].Value {
		t.Errorf("Expected numeric latest value: %v, got: %v",
			*fti.Latest.Numeric[0].Value, *r.Latest.Numeric[0].Value)
	}

	if r.Latest.Text[0].Value != nil {
		t.Errorf("Expected text latest value: nil, got: %v",
			r.Latest.Text[0].Value)
	}

	if *fti.Latest.Histogram[0].Value != *r.Latest.Histogram[0].Value {
		t.Errorf("Expected histogram latest value: %v, got: %v",
			*fti.Latest.Histogram[0].Value, *r.Latest.Histogram[0].Value)
	}
}

func TestFindTags(t *testing.T) {
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

		if strings.Contains(r.RequestURI, "&count_only=1") {
			w.Header().Set("X-Snowth-Search-Result-Count", "1")
			_, _ = w.Write([]byte(tagsCountTestData))
			return
		}

		if strings.HasPrefix(r.RequestURI, "/find/1/tags?query=test") {
			w.Header().Set("X-Snowth-Search-Result-Count", "1")
			_, _ = w.Write([]byte(tagsTestData))
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
	res, err := sc.FindTags(1, "test", &FindTagsOptions{
		Start:     time.Unix(1, 0),
		End:       time.Unix(2, 0),
		Activity:  0,
		Latest:    0,
		CountOnly: 1,
		Limit:     -1,
	}, node)
	if err != nil {
		t.Fatal(err)
	}

	if res.Count != 1 {
		t.Fatalf("Expected result count: 1, got: %v", res.Count)
	}

	res, err = sc.FindTags(1, "test", &FindTagsOptions{
		Start:     time.Unix(1, 0),
		End:       time.Unix(2, 0),
		Activity:  1,
		Latest:    1,
		CountOnly: 0,
		Limit:     -1,
	}, node)
	if err != nil {
		t.Fatal(err)
	}

	if res.Count != 1 {
		t.Fatalf("Expected result count: 1, got: %v", res.Count)
	}

	if len(res.Items) != 1 {
		t.Fatalf("Expected result length: 1, got: %v", len(res.Items))
	}

	if res.Items[0].MetricName != "test" {
		t.Errorf("Expected metric name: test, got: %v", res.Items[0].MetricName)
	}

	if len(res.Items[0].CheckTags) != 2 {
		t.Fatalf("Expected tags length: 2, got: %v",
			len(res.Items[0].CheckTags))
	}

	if res.Items[0].CheckTags[0] != "test:test" {
		t.Errorf("Expected check tag: test:test, got: %v",
			res.Items[0].CheckTags[0])
	}

	if len(res.Items[0].Activity) != 2 {
		t.Fatalf("Expected activity length: 2, got %v",
			len(res.Items[0].Activity))
	}

	if len(res.Items[0].Activity[1]) != 2 {
		t.Fatalf("Expected activity[1] length: 2, got %v",
			len(res.Items[0].Activity[1]))
	}

	if res.Items[0].Activity[1][1] != 1561848300 {
		t.Fatalf("Expected activity timestamp: 1561848300, got %v",
			res.Items[0].Activity[1][1])
	}

	res, err = sc.FindTags(1, "test", &FindTagsOptions{}, node)
	if err != nil {
		t.Fatal(err)
	}

	if res.Count != 1 {
		t.Fatalf("Expected result count: 1, got: %v", res.Count)
	}

	if len(res.Items) != 1 {
		t.Fatalf("Expected result length: 1, got: %v", len(res.Items))
	}

	if res.Items[0].AccountID != 1 {
		t.Errorf("Expected account ID: 1, got: %v", res.Items[0].AccountID)
	}

	if res.Items[0].MetricName != "test" {
		t.Errorf("Expected metric name: test, got: %v", res.Items[0].MetricName)
	}

	if len(res.Items[0].CheckTags) != 2 {
		t.Fatalf("Expected tags length: 2, got: %v",
			len(res.Items[0].CheckTags))
	}

	if res.Items[0].CheckTags[0] != "test:test" {
		t.Errorf("Expected check tag: test:test, got: %v",
			res.Items[0].CheckTags[0])
	}

	if len(res.Items[0].Activity) != 2 {
		t.Fatalf("Expected activity length: 2, got %v",
			len(res.Items[0].Activity))
	}

	if len(res.Items[0].Activity[1]) != 2 {
		t.Fatalf("Expected activity[1] length: 2, got %v",
			len(res.Items[0].Activity[1]))
	}

	if res.Items[0].Activity[1][1] != 1561848300 {
		t.Fatalf("Expected activity timestamp: 1561848300, got %v",
			res.Items[0].Activity[1][1])
	}
}
