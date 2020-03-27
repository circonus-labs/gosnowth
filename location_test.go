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
	<node id="1f846f26-0cfd-4df5-b4f1-e0930604e577"
		address="10.8.20.1"
		port="8112"
		apiport="8112"
		weight="32"/>
	<node id="07fa2237-5744-4c28-a622-a99cfc1ac87e"
		address="10.8.20.4"
		port="8112"
		apiport="8112"
		weight="32"/>
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
			w.Write([]byte(stateTestData))
			return
		}

		if r.RequestURI == "/stats.json" {
			w.Write([]byte(statsTestData))
			return
		}

		if strings.HasPrefix(r.RequestURI,
			"/locate/xml/1f846f26-0cfd-4df5-b4f1-e0930604e577/test") {
			w.Write([]byte(locateXMLTestData))
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

	exp := "1f846f26-0cfd-4df5-b4f1-e0930604e577"
	if res[0].ID != exp {
		t.Errorf("Expected ID: %v, got: %v", exp, res[0].ID)
	}
}
