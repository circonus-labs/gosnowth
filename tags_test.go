package gosnowth

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const tagsTestData = `[
	{
		"uuid": "3aa57ac2-28de-4ec4-aa3d-ed0ddd48fa4d",
		"check_name": "test",
		"check_tags": [
			"test:test",
			"__check_id:1"
		],
		"metric_name": "test",
		"category": "reconnoiter",
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
		"account_id": 1
	}
]`

func TestFindTags(t *testing.T) {
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

		if strings.HasPrefix(r.RequestURI, "/find/1/tags?query=test") {
			w.Write([]byte(tagsTestData))
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
	res, err := sc.FindTags(node, 1, "test", "1", "1")
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

	res, err = sc.FindTags(node, 1, "test", "", "")
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
