// Package gosnowth contains an IRONdb client library written in Go.
package gosnowth

import (
	"fmt"
	"bytes"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const locateXMLTestData = `<nodes n="2">
	<node id="00000000-0000-0000-0000-000000000001"
		address="10.1.2.1"
		port="8112"
		apiport="8112"
		weight="32"/>
	<node id="00000000-0000-0000-0000-000000000002"
		address="10.1.2.2"
		port="8112"
		apiport="8112"
		weight="32"/>
	<node id="00000000-0000-0000-0000-000000000003"
		address="10.1.2.3"
		port="8112"
		apiport="8112"
		weight="32"/>
	<node id="00000000-0000-0000-0000-000000000004"
		address="10.1.2.4"
		port="8112"
		apiport="8112"
		weight="32"/>
	<node id="00000000-0000-0000-0000-000000000005"
		address="10.1.2.5"
		port="8112"
		apiport="8112"
		weight="32"/>
	<node id="00000000-0000-0000-0000-000000000006"
		address="10.1.2.6"
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

	if len(dl.Nodes) != 6 {
		t.Errorf("Expected number of nodes: 6, got: %v", len(dl.Nodes))
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
			"/locate/xml/d76b5011-ded0-4523-8310-6132d863d02d/test") {
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
	res, err := sc.LocateMetric("d76b5011-ded0-4523-8310-6132d863d02d",
		"test", node)
	fmt.Println(res)
	if err != nil {
		t.Fatal(err)
	}

	// probably need to test that we get ALL the expected owning nodes,
	// not just the first one.
	exp := "00000000-0000-0000-0000-000000000001"
	if res[0].ID != exp {
		t.Errorf("Expected ID: %v, got: %v", exp, res[0].ID)
	}
}
