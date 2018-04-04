package gosnowth

import (
	"net/http"

	"github.com/pkg/errors"
)

// GetGossipInfo - Get the gossip information from the client.
func (sc *SnowthClient) GetGossipInfo(node *SnowthNode) (*Gossip, error) {
	req, err := http.NewRequest("GET", sc.getURL(node, "/gossip/json"), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	resp, err := sc.do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform request")
	}

	var gossip = new(Gossip)
	if err := decodeJSONFromResponse(gossip, resp); err != nil {
		return nil, errors.Wrap(err, "failed to decode")
	}

	return gossip, nil
}

// Gossip - the gossip information from a node.  This structure includes
// information on how the nodes are communicating with each other, and if an
// nodes are behind with each other with regards to data replication.
type Gossip []GossipDetail

// GossipDetail - Gossip information about a node identified by ID
type GossipDetail struct {
	ID          string        `json:"id"`
	Time        float64       `json:"gossip_time,string"`
	Age         float64       `json:"gossip_age,string"`
	CurrentTopo string        `json:"topo_current"`
	NextTopo    string        `json:"topo_next"`
	TopoState   string        `json:"topo_state"`
	Latency     GossipLatency `json:"latency"`
}

// GossipLatency - a map of the uuid of the node to the latency in seconds
type GossipLatency map[string]string
