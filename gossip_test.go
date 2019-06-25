package gosnowth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const gossipTestData = `[
	{
		"id": "1f846f26-0cfd-4df5-b4f1-e0930604e577",
		"gossip_time": "1409082055.744880",
		"gossip_age": "0.000000",
		"topo_current": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		"topo_next": "-",
		"topo_state": "n/a",
		"latency": {
			"765ac4cc-1929-4642-9ef1-d194d08f9538": "0",
			"8c2fc7b8-c569-402d-a393-db433fb267aa": "0",
			"07fa2237-5744-4c28-a622-a99cfc1ac87e": "0"
		}
	},
	{
		"id": "765ac4cc-1929-4642-9ef1-d194d08f9538",
		"gossip_time": "1409082055.744880",
		"gossip_age": "0.000000",
		"topo_current": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		"topo_next": "-",
		"topo_state": "n/a",
		"latency": {
			"1f846f26-0cfd-4df5-b4f1-e0930604e577": "0",
			"8c2fc7b8-c569-402d-a393-db433fb267aa": "0",
			"07fa2237-5744-4c28-a622-a99cfc1ac87e": "0"
		}
	},
	{
		"id": "8c2fc7b8-c569-402d-a393-db433fb267aa",
		"gossip_time": "1409082055.744880",
		"gossip_age": "0.000000",
		"topo_current": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		"topo_next": "-",
		"topo_state": "n/a",
		"latency": {
			"765ac4cc-1929-4642-9ef1-d194d08f9538": "0",
			"1f846f26-0cfd-4df5-b4f1-e0930604e577": "0",
			"07fa2237-5744-4c28-a622-a99cfc1ac87e": "0"
		}
	},
	{
		"id": "07fa2237-5744-4c28-a622-a99cfc1ac87e",
		"gossip_time": "1409082055.744880",
		"gossip_age": "0.000000",
		"topo_current": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		"topo_next": "-",
		"topo_state": "n/a",
		"latency": {
			"765ac4cc-1929-4642-9ef1-d194d08f9538": "0",
			"8c2fc7b8-c569-402d-a393-db433fb267aa": "0",
			"1f846f26-0cfd-4df5-b4f1-e0930604e577": "0"
		}
	}
]`

const gossipTestAltData = `[
	{
		"id": "765ac4cc-1929-4642-9ef1-d194d08f9538",
		"gossip_time": "1409082055.744880",
		"gossip_age": "0.000000",
		"topo_current": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		"topo_next": "-",
		"topo_state": "n/a",
		"latency": {
			"1f846f26-0cfd-4df5-b4f1-e0930604e577": "0",
			"8c2fc7b8-c569-402d-a393-db433fb267aa": "0",
			"07fa2237-5744-4c28-a622-a99cfc1ac87e": "0"
		}
	},
	{
		"id": "8c2fc7b8-c569-402d-a393-db433fb267aa",
		"gossip_time": "1409082055.744880",
		"gossip_age": "0.000000",
		"topo_current": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		"topo_next": "-",
		"topo_state": "n/a",
		"latency": {
			"765ac4cc-1929-4642-9ef1-d194d08f9538": "0",
			"1f846f26-0cfd-4df5-b4f1-e0930604e577": "0",
			"07fa2237-5744-4c28-a622-a99cfc1ac87e": "0"
		}
	},
	{
		"id": "07fa2237-5744-4c28-a622-a99cfc1ac87e",
		"gossip_time": "1409082055.744880",
		"gossip_age": "0.000000",
		"topo_current": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		"topo_next": "-",
		"topo_state": "n/a",
		"latency": {
			"765ac4cc-1929-4642-9ef1-d194d08f9538": "0",
			"8c2fc7b8-c569-402d-a393-db433fb267aa": "0",
			"1f846f26-0cfd-4df5-b4f1-e0930604e577": "0"
		}
	}
]`

func TestGossipDeserialization(t *testing.T) {
	dec := json.NewDecoder(bytes.NewBufferString(gossipTestData))
	dec.UseNumber()
	gossip := new(Gossip)
	err := dec.Decode(gossip)
	if err != nil {
		t.Errorf("failed to decode gossip data, %s\n", err.Error())
	}

	if len(*gossip) != 4 {
		t.Error("Should have 4 entries")
	}

	res := []GossipDetail(*gossip)[0].Time
	exp := float64(1409082055.744880)
	if res != exp {
		t.Errorf("Expected time: %v, got: %v", exp, res)
	}

	resA := []GossipDetail(*gossip)[0].Age
	expA := float64(0.0)
	if resA != expA {
		t.Errorf("Expected age: %v, got: %v", exp, res)
	}
}

func TestGetGossipInfo(t *testing.T) {
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

		if r.RequestURI == "/gossip/json" {
			w.Write([]byte(gossipTestData))
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
	res, err := sc.GetGossipInfo(node)
	if err != nil {
		t.Fatal(err)
	}

	if res == nil || len(*res) != 4 {
		t.Fatalf("Expected result length: 4, got: %v", len(*res))
	}

	if (*res)[0].ID != "1f846f26-0cfd-4df5-b4f1-e0930604e577" {
		t.Errorf("Expected ID: 1f846f26-0cfd-4df5-b4f1-e0930604e577, got: %v",
			(*res)[0].ID)
	}

	ctx, cancel := context.WithCancel(context.Background())
	res, err = sc.GetGossipInfoContext(ctx, node)
	if err != nil {
		t.Fatal(err)
	}

	if res == nil || len(*res) != 4 {
		t.Fatalf("Expected result length: 4, got: %v", len(*res))
	}

	if (*res)[0].ID != "1f846f26-0cfd-4df5-b4f1-e0930604e577" {
		t.Errorf("Expected ID: 1f846f26-0cfd-4df5-b4f1-e0930604e577, got: %v",
			(*res)[0].ID)
	}

	cancel()
	res, err = sc.GetGossipInfoContext(ctx, node)
	if err == nil || !strings.Contains(err.Error(), "context") {
		t.Error("Expected context error.", err)
	}
}
