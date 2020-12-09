// Package gosnowth contains an IRONdb client library written in Go.
package gosnowth

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const locateXMLTestData = `<nodes n="2">
	<node id="fac5f5e4-b9f0-49ca-870d-992c5c2a0c79"
		address="10.0.0.1"
		port="8112"
		apiport="8112"
		weight="170"
		side="both"/>
	<node id="97df761e-7c56-40dd-9344-ec064a4ed68d"
		address="10.0.0.2"
		port="8112"
		apiport="8112"
		weight="170"
		side="both"/>
</nodes>`

const locateTopologyXMLTestData = `<nodes n="3">
	<node id="5c32c076-ffeb-cfdd-a541-97e25c028dd6"
		address="10.0.0.100"
		port="8112"
		apiport="8112"
		weight="51"
		side="a"/>
	<node id="1533fc6b-de08-6eac-eb46-d3920a1a18a3"
		address="10.0.0.101"
		port="8112"
		apiport="8112"
		weight="51"
		side="b"/>
	<node id="18111a24-5832-42c8-e780-bcbf88f47215"
		address="10.0.0.102"
		port="8112"
		apiport="8112"
		weight="51"
		side="a"/>
	<node id="4ec7bd67-f279-6f9a-fbe7-be9a0dee4c39"
		address="10.0.0.103"
		port="8112"
		apiport="8112"
		weight="51"
		side="b"/>
	<node id="0475df4e-ee2d-c96c-b6d7-e9d1b0239c2c"
		address="10.0.0.104"
		port="8112"
		apiport="8112"
		weight="51"
		side="a"/>
	<node id="9d1a34cd-b150-4c19-a894-e20280b42b62"
		address="10.0.0.105"
		port="8112"
		apiport="8112"
		weight="51"
		side="b"/>
	<node id="3d8ae36d-3d4d-4eda-ab53-c58538985062"
		address="10.0.0.106"
		port="8112"
		apiport="8112"
		weight="51"
		side="a"/>
	<node id="d2b9a8aa-9503-6cb3-dfdd-e407c1a6bee7"
		address="10.0.0.107"
		port="8112"
		apiport="8112"
		weight="51"
		side="b"/>
	<node id="15e35e06-4069-ecb8-c7a4-93e4c540693d"
		address="10.0.0.108"
		port="8112"
		apiport="8112"
		weight="51"
		side="a"/>
	<node id="8f0073e1-5d52-67da-bd59-e8017e5b5aa1"
		address="10.0.0.109"
		port="8112"
		apiport="8112"
		weight="51"
		side="b"/>
</nodes>`

func TestDataLocationXMLDeserialization(t *testing.T) {
	dec := xml.NewDecoder(bytes.NewBufferString(locateXMLTestData))
	dl := new(Topology)
	err := dec.Decode(dl)
	if err != nil {
		t.Errorf("failed to decode node stats, %s\n", err.Error())
	}

	if len(dl.Nodes) != 2 {
		t.Errorf("Expected number of nodes: 2, got: %v", len(dl.Nodes))
	}
}

func TestLocateMetric(t *testing.T) {
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

		if strings.HasPrefix(r.RequestURI,
			"/locate/xml/1f846f26-0cfd-4df5-b4f1-e0930604e577/test") {
			_, _ = w.Write([]byte(locateXMLTestData))
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
	res, err := sc.LocateMetric("1f846f26-0cfd-4df5-b4f1-e0930604e577",
		"test", node)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 2 {
		t.Errorf("Expected length: 2, got: %v", len(res))
	}

	exp := "fac5f5e4-b9f0-49ca-870d-992c5c2a0c79"
	if res[0].ID != exp {
		t.Errorf("Expected ID: %v, got: %v", exp, res[0].ID)
	}
}

func TestLocateMetricFindMetric(t *testing.T) {
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

		if strings.HasPrefix(r.RequestURI,
			"/topology/xml") {
			_, _ = w.Write([]byte(topologyXMLTestData))
			return
		}

		if strings.HasPrefix(r.RequestURI,
			"/locate/xml/1f846f26-0cfd-4df5-b4f1-e0930604e577/test") {
			_, _ = w.Write([]byte(locateXMLTestData))
			return
		}
	}))

	defer ms.Close()
	sc, err := NewSnowthClient(false, ms.URL)
	if err != nil {
		t.Fatal("Unable to create snowth client", err)
	}

	res, err := sc.LocateMetric("1f846f26-0cfd-4df5-b4f1-e0930604e577",
		"test")
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 3 {
		t.Errorf("Expected length: 3, got: %v", len(res))
	}

	exp := "9d1a34cd-b150-4c19-a894-e20280b42b62"
	if res[0].ID != exp {
		t.Errorf("Expected primary node ID: %v, got: %v", exp, res[0].ID)
	}

	exp = "3d8ae36d-3d4d-4eda-ab53-c58538985062"
	if res[1].ID != exp {
		t.Errorf("Expected secondary node ID: %v, got: %v", exp, res[1].ID)
	}

	exp = "1533fc6b-de08-6eac-eb46-d3920a1a18a3"
	if res[2].ID != exp {
		t.Errorf("Expected ID: %v, got: %v", exp, res[2].ID)
	}
}
