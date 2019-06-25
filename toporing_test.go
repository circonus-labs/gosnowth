package gosnowth

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const topoRingTestJSONData = `[
	{
		"id": "1f846f26-0cfd-4df5-b4f1-e0930604e577",
		"idx": 1,
		"location": 11.000000
	},
	{
		"id": "1f846f26-0cfd-4df5-b4f1-e0930604e577",
		"idx": 2,
		"location": 22.000000
	},
	{
		"id": "1f846f26-0cfd-4df5-b4f1-e0930604e577",
		"idx": 3,
		"location": 33.000000
	},
	{
		"id": "765ac4cc-1929-4642-9ef1-d194d08f9538",
		"idx": 1,
		"location": 44.000000
	},
	{
		"id": "765ac4cc-1929-4642-9ef1-d194d08f9538",
		"idx": 2,
		"location": 55.000000
	},
	{
		"id": "765ac4cc-1929-4642-9ef1-d194d08f9538",
		"idx": 3,
		"location": 66.000000
	},
	{
		"id": "8c2fc7b8-c569-402d-a393-db433fb267aa",
		"idx": 1,
		"location": 77.000000
	},
	{
		"id": "8c2fc7b8-c569-402d-a393-db433fb267aa",
		"idx": 2,
		"location": 88.000000
	},
	{
		"id": "8c2fc7b8-c569-402d-a393-db433fb267aa",
		"idx": 3,
		"location": 99.000000
	}
]`

const topoRingXMLTestData = `<vnodes n="2">
	<vnode id="1f846f26-0cfd-4df5-b4f1-e0930604e577"
		idx="1"
		location="11.000000"/>
	<vnode id="1f846f26-0cfd-4df5-b4f1-e0930604e577"
		idx="2"
		location="22.000000"/>
	<vnode id="1f846f26-0cfd-4df5-b4f1-e0930604e577"
		idx="3"
		location="33.000000"/>
	<vnode id="765ac4cc-1929-4642-9ef1-d194d08f9538"
		idx="1"
		location="44.000000"/>
	<vnode id="765ac4cc-1929-4642-9ef1-d194d08f9538"
		idx="2"
		location="55.000000"/>
	<vnode id="765ac4cc-1929-4642-9ef1-d194d08f9538"
		idx="3"
		location="66.000000"/>
	<vnode id="8c2fc7b8-c569-402d-a393-db433fb267aa"
		idx="1"
		location="77.000000"/>
	<vnode id="8c2fc7b8-c569-402d-a393-db433fb267aa"
		idx="2"
		location="88.000000"/>
	<vnode id="8c2fc7b8-c569-402d-a393-db433fb267aa"
		idx="3"
		location="99.000000"/>
</vnodes>`

func TestTopoRingJSONDeserialization(t *testing.T) {
	dec := json.NewDecoder(bytes.NewBufferString(topoRingTestJSONData))
	dec.UseNumber()
	tr := []TopoRingDetail{}
	err := dec.Decode(&tr)
	if err != nil {
		t.Errorf("failed to decode node stats, %s\n", err.Error())
	}

	if len(tr) != 9 {
		t.Error("should be 9 elements")
	}
}

func TestTopoRingXMLDeserialization(t *testing.T) {
	dec := xml.NewDecoder(bytes.NewBufferString(topoRingXMLTestData))
	tr := new(TopoRing)
	err := dec.Decode(tr)
	if err != nil {
		t.Errorf("failed to decode node stats, %s\n", err.Error())
	}

	if len(tr.VirtualNodes) != 9 {
		t.Error("should be 9 elements")
	}
}

func TestGetTopoRingInfo(t *testing.T) {
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
			"/toporing/xml/test") {
			w.Write([]byte(topoRingXMLTestData))
			return
		}

		w.WriteHeader(500)
		return
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
	res, err := sc.GetTopoRingInfo("test", node)
	if err != nil {
		t.Fatal(err)
	}

	if len(res.VirtualNodes) != 9 {
		t.Fatalf("Expected nodes length: 9, got: %v", len(res.VirtualNodes))
	}

	exp := "1f846f26-0cfd-4df5-b4f1-e0930604e577"
	if res.VirtualNodes[0].ID != exp {
		t.Errorf("Expected node ID: %v, got: %v", exp, res.VirtualNodes[0].ID)
	}
}
