package gosnowth

import "context"

// Gossip values contain gossip information from a node. This structure includes
// information on how the nodes are communicating with each other, and if any
// nodes are behind with each other with regards to data replication.
type Gossip []GossipDetail

// GossipDetail values represent gossip information about a node.
type GossipDetail struct {
	ID          string        `json:"id"`
	Time        float64       `json:"gossip_time,string"`
	Age         float64       `json:"gossip_age,string"`
	CurrentTopo string        `json:"topo_current"`
	NextTopo    string        `json:"topo_next"`
	TopoState   string        `json:"topo_state"`
	Latency     GossipLatency `json:"latency"`
}

// GossipLatency values contain a map of node UUID's to latencies in seconds.
type GossipLatency map[string]string

// GetGossipInfo fetches the gossip information from an IRONdb node. The gossip
// response body will include a list of "GossipDetail" which provide
// the identifier of the node, the node's gossip_time, gossip_age, as well
// as topology state, current and next topology.
func (sc *SnowthClient) GetGossipInfo(
	node *SnowthNode) (gossip *Gossip, err error) {
	return sc.GetGossipInfoContext(context.Background(), node)
}

// GetGossipInfoContext is the context aware version of GetGossipInfo.
func (sc *SnowthClient) GetGossipInfoContext(ctx context.Context,
	node *SnowthNode) (gossip *Gossip, err error) {
	gossip = new(Gossip)
	err = sc.do(ctx, node, "GET", "/gossip/json", nil, gossip, decodeJSON)
	return
}
