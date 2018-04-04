package gosnowth

import (
	"encoding/xml"
	"fmt"
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

// LoadTopology - Load a new topology. Will not activate, just load and store.
func (sc *SnowthClient) LoadTopology(hash string, topology *Topology, node *SnowthNode) error {

	reqBody, err := encodeXML(topology)
	if err != nil {
		return errors.Wrap(err, "failed to encode request data")
	}

	var resource = path.Join("/topology", hash)
	req, err := http.NewRequest("POST", sc.getURL(node, resource), reqBody)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	resp, err := sc.do(req)
	if err != nil {
		return errors.Wrap(err, "failed to perform request")
	}
	defer closeBody(resp)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid response status: %s", resp.Status)
	}
	return nil
}

// ActivateTopology - Switch to a new topology.  THIS IS DANGEROUS.
func (sc *SnowthClient) ActivateTopology(hash string, node *SnowthNode) error {
	var resource = path.Join("/activate", hash)
	req, err := http.NewRequest("GET", sc.getURL(node, resource), nil)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	resp, err := sc.do(req)
	if err != nil {
		return errors.Wrap(err, "failed to perform request")
	}
	defer closeBody(resp)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid response status: %s", resp.Status)
	}
	return nil
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
