package gosnowth

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestSnowthNode(t *testing.T) {
	u, err := url.Parse("localhost")
	if err != nil {
		t.Fatal(err)
	}

	sn := &SnowthNode{
		url:             u,
		identifier:      "test",
		currentTopology: "test",
	}

	if sn.GetURL() != u {
		t.Error("Invalid URL returned")
	}

	if sn.GetCurrentTopology() != "test" {
		t.Errorf("Expected topology: test, got: %v", sn.GetCurrentTopology())
	}
}

func TestNewSnowthClient(t *testing.T) {
	// crude test to ensure err is returned for invalid snowth url
	badAddr := "foobar"
	_, err := NewSnowthClient(false, badAddr)
	if err == nil {
		t.Errorf("Error not encountered on invalid snowth addr %v", badAddr)
	}
}

func TestIsNodeActive(t *testing.T) {
	// mock out GetNodeState, GetGossipInfo
}

func TestSnowthClientRequest(t *testing.T) {
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
			if r.Header.Get("X-Test-Header") != "test" {
				t.Error("Expected X-Test-Header:test")
			}

			w.Write([]byte(tagsTestData))
			return
		}
	}))

	defer ms.Close()
	sc, err := NewSnowthClient(false, ms.URL)
	if err != nil {
		t.Fatal("Unable to create snowth client", err)
	}

	sc.SetRequestFunc(func(r *http.Request) error {
		r.Header.Set("X-Test-Header", "test")
		return nil
	})

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

	if res.Items[0].AccountID != 1 {
		t.Errorf("Expected account ID: 1, got: %v", res.Items[0].AccountID)
	}

	body, _, err := sc.DoRequest(node, "GET", "/stats.json", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	r := map[string]map[string]interface{}{}
	err = decodeJSON(body, &r)
	if err != nil {
		t.Fatal(err)
	}

	appValue := r["application"]["_value"]
	if appValue != "snowth" {
		t.Fatalf("Expected application: snowth, got: %v", appValue)
	}
}

func TestSnowthClientDiscoverNodesWatch(t *testing.T) {
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

		if strings.HasPrefix(r.RequestURI, "/topology/xml/") {
			w.Write([]byte(topologyXMLTestData))
			return
		}

		if r.RequestURI == "/gossip/json" {
			if r.Header.Get("ALT") != "" {
				w.Write([]byte(gossipTestAltData))
			}

			w.Write([]byte(gossipTestData))
			return
		}
	}))

	defer ms.Close()
	sc, err := NewSnowthClient(true, ms.URL)
	if err != nil {
		t.Fatal("Unable to create snowth client", err)
	}

	u, err := url.Parse(ms.URL)
	if err != nil {
		t.Fatal("Invalid test URL")
	}

	node := &SnowthNode{
		url:             u,
		identifier:      "bb6f7162-4828-11df-bab8-6bac200dcc2a",
		currentTopology: "294cbd39999c2270964029691e8bc5e231a867d525ccba62181dc8988ff218dc",
	}

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

	if res.Items[0].AccountID != 1 {
		t.Errorf("Expected account ID: 1, got: %v", res.Items[0].AccountID)
	}

	sc.watchInterval = 100 * time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sc.WatchAndUpdate(ctx)
	sc.AddNodes(node)
	sc.ActivateNodes(node)
	if !sc.isNodeActive(node) {
		t.Errorf("Expected node to be active")
	}

	sc.SetRequestFunc(func(r *http.Request) error {
		r.Header.Set("ALT", "true")
		return nil
	})

	time.Sleep(150 * time.Millisecond)
	if sc.isNodeActive(node) {
		t.Errorf("Expected node to be inactive")
	}

	sc.SetRequestFunc(nil)
	time.Sleep(150 * time.Millisecond)
	if !sc.isNodeActive(node) {
		t.Errorf("Expected node to be active")
	}

	cancel()
	time.Sleep(50 * time.Millisecond)
	sc.watchInterval = 0
	sc.WatchAndUpdate(ctx)
}

type mockLog struct {
	last string
}

func (m *mockLog) Debugf(format string, args ...interface{}) {
	m.last = fmt.Sprintf("DEBUG "+format, args...)
}

func (m *mockLog) Errorf(format string, args ...interface{}) {
	m.last = fmt.Sprintf("ERROR "+format, args...)
}
func (m *mockLog) Infof(format string, args ...interface{}) {
	m.last = fmt.Sprintf("INFO "+format, args...)
}

func (m *mockLog) Warnf(format string, args ...interface{}) {
	m.last = fmt.Sprintf("WARN "+format, args...)
}

func TestSnowthClientLog(t *testing.T) {
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

		if strings.HasPrefix(r.RequestURI, "/topology/xml/") {
			w.Write([]byte(topologyXMLTestData))
			return
		}
	}))

	defer ms.Close()
	sc, err := NewSnowthClient(true, ms.URL)
	if err != nil {
		t.Fatal("Unable to create snowth client", err)
	}

	ml := &mockLog{}
	sc.SetLog(ml)
	sc.LogDebugf("test %d", 1)
	exp := "DEBUG test 1"
	if ml.last != exp {
		t.Errorf("Expected log entry: %v, got: %v", exp, ml.last)
	}

	sc.LogErrorf("test %d", 1)
	exp = "ERROR test 1"
	if ml.last != exp {
		t.Errorf("Expected log entry: %v, got: %v", exp, ml.last)
	}

	sc.LogInfof("test %d", 1)
	exp = "INFO test 1"
	if ml.last != exp {
		t.Errorf("Expected log entry: %v, got: %v", exp, ml.last)
	}

	sc.LogWarnf("test %d", 1)
	exp = "WARN test 1"
	if ml.last != exp {
		t.Errorf("Expected log entry: %v, got: %v", exp, ml.last)
	}
}

func TestSnowthClientSetWatch(t *testing.T) {
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
			if r.Header.Get("X-Test-Header") != "test" {
				t.Error("Expected X-Test-Header:test")
			}

			w.Write([]byte(tagsTestData))
			return
		}
	}))

	defer ms.Close()
	sc, err := NewSnowthClient(false, ms.URL)
	if err != nil {
		t.Fatal("Unable to create snowth client", err)
	}

	sc.SetWatchInterval(100 * time.Millisecond)
	sc.SetWatchFunc(func(n *SnowthNode) {
		exp := "0.1.1570000000"
		if n.SemVer() == exp {
			sc.DeactivateNodes(n)
		} else {
			t.Errorf("Expected version: %v, got: %v", exp, n.SemVer())
		}
	})

	nodes := sc.ListActiveNodes()
	if len(nodes) < 1 {
		t.Errorf("Expected length nodes > 1, got: %d", len(nodes))
	}

	sc.WatchAndUpdate(context.Background())
	time.Sleep(150 * time.Millisecond)
	nodes = sc.ListActiveNodes()
	if len(nodes) > 0 {
		t.Fatalf("Expected length nodes: 0, got: %d", len(nodes))
	}
}
