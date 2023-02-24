package gosnowth

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
)

const topologyTestData = `[
	{
		"id": "1f846f26-0cfd-4df5-b4f1-e0930604e577",
		"address": "10.8.20.1",
		"port": 8112,
		"apiport": 8112,
		"weight": 32,
		"n": 2
	},
	{
		"id": "765ac4cc-1929-4642-9ef1-d194d08f9538",
		"address": "10.8.20.2",
		"port": 8112,
		"apiport": 8112,
		"weight": 32,
		"n": 2
	},
	{
		"id": "8c2fc7b8-c569-402d-a393-db433fb267aa",
		"address": "10.8.20.3",
		"port": 8112,
		"apiport": 8112,
		"weight": 32,
		"n": 2
	},
	{
		"id": "07fa2237-5744-4c28-a622-a99cfc1ac87e",
		"address": "10.8.20.4",
		"port": 8112,
		"apiport": 8112,
		"weight": 32,
		"n": 2
	}
]`

const topologyXMLTestData = `<nodes n="3">
<node id="5c32c076-ffeb-cfdd-a541-97e25c028dd6" address="10.128.0.100" port="8112" apiport="8112" weight="51" side="a"/>
<node id="1533fc6b-de08-6eac-eb46-d3920a1a18a3" address="10.128.0.101" port="8112" apiport="8112" weight="51" side="b"/>
<node id="18111a24-5832-42c8-e780-bcbf88f47215" address="10.128.0.102" port="8112" apiport="8112" weight="51" side="a"/>
<node id="4ec7bd67-f279-6f9a-fbe7-be9a0dee4c39" address="10.128.0.103" port="8112" apiport="8112" weight="51" side="b"/>
<node id="0475df4e-ee2d-c96c-b6d7-e9d1b0239c2c" address="10.128.0.104" port="8112" apiport="8112" weight="51" side="a"/>
<node id="9d1a34cd-b150-4c19-a894-e20280b42b62" address="10.128.0.105" port="8112" apiport="8112" weight="51" side="b"/>
<node id="3d8ae36d-3d4d-4eda-ab53-c58538985062" address="10.128.0.106" port="8112" apiport="8112" weight="51" side="a"/>
<node id="d2b9a8aa-9503-6cb3-dfdd-e407c1a6bee7" address="10.128.0.107" port="8112" apiport="8112" weight="51" side="b"/>
<node id="15e35e06-4069-ecb8-c7a4-93e4c540693d" address="10.128.0.108" port="8112" apiport="8112" weight="51" side="a"/>
<node id="8f0073e1-5d52-67da-bd59-e8017e5b5aa1" address="10.128.0.109" port="8112" apiport="8112" weight="51" side="b"/>
</nodes>`

func TestTopologyJSONDeserialization(t *testing.T) {
	t.Parallel()

	dec := json.NewDecoder(bytes.NewBufferString(topologyTestData))
	dec.UseNumber()

	topo := []TopologyNode{}

	if err := dec.Decode(&topo); err != nil {
		t.Errorf("failed to decode topology, %s\n", err.Error())
	}

	if len(topo) != 4 {
		t.Error("should be 4 elements")
	}
}

func TestTopologyXMLDeserialization(t *testing.T) {
	t.Parallel()

	dec := xml.NewDecoder(bytes.NewBufferString(topologyXMLTestData))
	topo := new(Topology)

	if err := dec.Decode(topo); err != nil {
		t.Errorf("failed to decode topology, %s\n", err.Error())
	}

	if len(topo.Nodes) != 10 {
		t.Error("should be 10 elements")
	}
}

func TestTopologyXMLSerialization(t *testing.T) {
	t.Parallel()

	buf := bytes.NewBuffer([]byte{})
	enc := xml.NewEncoder(buf)
	topo := Topology{
		WriteCopies: 2,
		Nodes: []TopologyNode{
			{
				ID:      "1f846f26-0cfd-4df5-b4f1-e0930604e577",
				Address: "10.8.20.1",
				Port:    8112,
				APIPort: 8112,
				Weight:  32,
			},
			{
				ID:      "765ac4cc-1929-4642-9ef1-d194d08f9538",
				Address: "10.8.20.2",
				Port:    8112,
				APIPort: 8112,
				Weight:  32,
			},
			{
				ID:      "8c2fc7b8-c569-402d-a393-db433fb267aa",
				Address: "10.8.20.3",
				Port:    8112,
				APIPort: 8112,
				Weight:  32,
			},
			{
				ID:      "07fa2237-5744-4c28-a622-a99cfc1ac87e",
				Address: "10.8.20.4",
				Port:    8112,
				APIPort: 8112,
				Weight:  32,
			},
		},
	}

	if err := enc.Encode(topo); err != nil {
		t.Errorf("failed to encode node stats, %s\n", err.Error())
	}

	if strings.Count(buf.String(), "id=") != 4 {
		t.Error("should have 4 nodes")
	}
}

/*
func checkLocationAgainstNode(t *testing.T, sc *SnowthClient, uuid string, metric string) {
	intnodes, err := sc.LocateMetric(uuid, metric)
	if err != nil {
		t.Error(err)
		return
	}
	extnodes, err := sc.LocateMetricRemote(uuid, metric, nil)
	if err != nil {
		t.Error(err)
		return
	}
	if len(intnodes) != len(extnodes) {
		t.Errorf("%s-%s locate list length mismatch", uuid, metric)
	} else {
		for i, node := range intnodes {
			if node.ID != extnodes[i].ID {
				t.Errorf("%s-%s [%d] i[%s] != e[%s]", uuid, metric, i, node.ID, extnodes[i].ID)
			}
		}
	}
}
*/

/*
func TestLiveNode(t *testing.T) {
	base := os.Getenv("SNOWTH_URL")
	if base == "" {
		return
	}

	nids := 10
	nmetrics := 10
	sc, err := NewClient(context.Background(),
		&Config{Servers: []string{ms.URL}})
	if err != nil {
		t.Fatal("Unable to create snowth client", err)
	}

	for i := 0; i < nids; i++ {
		id := uuid.New().String()
		for j := 0; j < nmetrics; j++ {
			checkLocationAgainstNode(t, sc, id, "foo|ST[bar:baz"+strconv.FormatInt(int64(j), 10)+"]")
		}
	}
}
*/

func BenchmarkLookup1(b *testing.B) {
	b.StopTimer()

	id := uuid.New().String()

	topo, err := TopologyLoadXML(topologyXMLTestData)
	if err != nil {
		b.Fatal("cannot load topology for benchmark")
	}

	b.StartTimer()

	for n := 0; n < b.N; n++ {
		_, _ = topo.FindMetric(id,
			"this.is.a.metric|ST[nice:andhappy,with:tags]")
	}
}

func TestTopology(t *testing.T) {
	t.Parallel()

	ms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request,
	) {
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
			"/topology/test") {
			w.WriteHeader(http.StatusOK)

			return
		}

		if strings.HasPrefix(r.RequestURI,
			"/activate/test") {
			w.WriteHeader(http.StatusOK)

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
	}))

	defer ms.Close()

	sc, err := NewClient(context.Background(),
		&Config{Servers: []string{ms.URL}})
	if err != nil {
		t.Fatal("Unable to create snowth client", err)
	}

	u, err := url.Parse(ms.URL)
	if err != nil {
		t.Fatal("Invalid test URL")
	}

	node := &SnowthNode{url: u}

	res, err := sc.GetTopologyInfo(nil)
	if err != nil {
		t.Fatal(err)
	}

	if res.WriteCopies != 3 {
		t.Fatalf("Expected nodes length: 3, got: %v", res.WriteCopies)
	}

	exp := "5c32c076-ffeb-cfdd-a541-97e25c028dd6"
	if res.Nodes[0].ID != exp {
		t.Errorf("Expected node ID: %v, got: %v", exp, res.Nodes[0].ID)
	}

	if res.Hash != "6c5f3aefde5c1f32d088b450fb3f0d9f33dedaaf8bed9cf5f77906f13fd65fc8" {
		t.Errorf("Unexpected topo hash: %v", res.Hash)
	}

	err = sc.LoadTopology("test", res, node)
	if err != nil {
		t.Fatal(err)
	}

	err = sc.ActivateTopology("test", node)
	if err != nil {
		t.Fatal(err)
	}
}
