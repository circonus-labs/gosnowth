package gosnowth

import (
	"encoding/xml"
	"net/http"
	"path"

	"github.com/pkg/errors"
)

// GetTopologyInfo - Get the topology information from the node.
func (sc *SnowthClient) GetTopologyInfo(node *SnowthNode) (*Topology, error) {
	var resource = path.Join("/topology/xml", node.currentTopology)
	req, err := http.NewRequest("GET", sc.getURL(node, resource), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	resp, err := sc.do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform request")
	}

	var topology = new(Topology)
	if err := decodeXMLFromResponse(topology, resp); err != nil {
		return nil, errors.Wrap(err, "failed to decode")
	}
	topology.Hash = node.currentTopology

	return topology, nil
}

// Topology - the topology structure from the API
type Topology struct {
	XMLName     xml.Name       `xml:"nodes" json:"-"`
	NumberNodes int            `xml:"n,attr" json:"-"`
	Hash        string         `xml:"-"`
	Nodes       []TopologyNode `xml:"node"`
}

// TopologyNode - the topology node structure from the API
type TopologyNode struct {
	XMLName     xml.Name `xml:"node" json:"-"`
	ID          string   `xml:"id,attr" json:"id"`
	Address     string   `xml:"address,attr" json:"address"`
	Port        uint16   `xml:"port,attr" json:"port"`
	APIPort     uint16   `xml:"apiport,attr" json:"apiport"`
	Weight      int      `xml:"weight,attr" json:"weight"`
	NumberNodes int      `xml:"-" json:"n"`
}
