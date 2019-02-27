package gosnowth

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"testing"
)

const topoRingJSONTestData = `[
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
	dec := json.NewDecoder(bytes.NewBufferString(topoRingJSONTestData))
	dec.UseNumber()
	toporing := []TopoRingDetail{}
	err := dec.Decode(&toporing)
	if err != nil {
		t.Errorf("failed to decode node stats, %s\n", err.Error())
	}

	if len(toporing) != 9 {
		t.Error("should be 9 elements")
	}
}

func TestTopoRingXMLDeserialization(t *testing.T) {
	dec := xml.NewDecoder(bytes.NewBufferString(topoRingXMLTestData))
	toporing := new(TopoRing)
	err := dec.Decode(toporing)
	if err != nil {
		t.Errorf("failed to decode node stats, %s\n", err.Error())
	}

	if len(toporing.VirtualNodes) != 9 {
		t.Error("should be 9 elements")
	}
}
