package gosnowth

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/pkg/errors"
)

// SnowthNode - representation of a snowth node, contains identifier and
// base url for connecting to this node
type SnowthNode struct {
	url             *url.URL
	identifier      string
	currentTopology string
}

// SnowthClient - client functionality for talking with a snowth topology
type SnowthClient struct {
	c *http.Client

	activeNodesMu *sync.RWMutex
	activeNodes   []*SnowthNode

	inactiveNodesMu *sync.RWMutex
	inactiveNodes   []*SnowthNode
}

// NewSnowthClient - given a variadic addrs parameter, the client will
// construct all the needed state to communicate with a group of nodes
// which constitute a cluster
func NewSnowthClient(addrs ...string) (*SnowthClient, error) {
	sc := &SnowthClient{
		c:               http.DefaultClient,
		activeNodesMu:   new(sync.RWMutex),
		activeNodes:     []*SnowthNode{},
		inactiveNodesMu: new(sync.RWMutex),
		inactiveNodes:   []*SnowthNode{},
	}

	for _, addr := range addrs {
		url, err := url.Parse(addr)
		if err != nil {
			// this node had an error, put on inactive list
			log.Printf("failed to bootstrap state of node: %+v", err)
			continue
		}
		node := &SnowthNode{url: url}
		// call get state to populate the id of this node
		state, err := sc.GetNodeState(node)
		if err != nil {
			// this node had an error, put on inactive list
			log.Printf("failed to bootstrap state of node: %+v", err)
			continue
		}
		node.identifier = state.Identity
		node.currentTopology = state.Current
		sc.AddNodes(node)
		sc.ActivateNodes(node)
	}

	// TODO: run background watching task to manage active/inactive nodes

	if err := sc.discoverNodes(); err != nil {
		return nil, errors.Wrap(err, "failed to discover nodes")
	}

	return sc, nil
}

// discoverNodes - private method for the client to discover peer nodes
// related to the topology.  This function will go through the active nodes
// get the topology information which shows all other nodes included in
// the topology, and adds them as snowth nodes to this client's active pool
// of nodesh
func (sc *SnowthClient) discoverNodes() error {
	// take our list of active nodes, interrogate gossipinfo
	// get more nodes from the gossip info
	for _, node := range sc.ListActiveNodes() {
		// lookup the topology
		topology, err := sc.GetTopologyInfo(node)
		if err != nil {
			log.Printf("error getting topology info: %+v", err)
			continue
		}

		// populate all the nodes with the appropriate topology information
		for _, topoNode := range topology.Nodes {
			sc.populateNodeInfo(topology.Hash, topoNode)
		}
	}

	return nil
}

// populateNodeInfo - this helper method populates an existing node with the
// details from the topology.  If a node doesn't exist, it will be added
// to the list of active nodes in the client.
func (sc *SnowthClient) populateNodeInfo(hash string, topology TopologyNode) {
	var found = false

	sc.activeNodesMu.Lock()
	for i := 0; i < len(sc.activeNodes); i++ {
		if sc.activeNodes[i].identifier == topology.ID {
			found = true
			url := url.URL{
				Scheme: "http",
				Host:   fmt.Sprintf("%s:%d", topology.Address, topology.APIPort),
			}
			sc.activeNodes[i].url = &url
			sc.activeNodes[i].currentTopology = hash
			continue
		}
	}
	sc.activeNodesMu.Unlock()
	sc.inactiveNodesMu.Lock()
	for i := 0; i < len(sc.inactiveNodes); i++ {
		found = true
		if sc.inactiveNodes[i].identifier == topology.ID {
			url := url.URL{
				Scheme: "http",
				Host:   fmt.Sprintf("%s:%d", topology.Address, topology.APIPort),
			}
			sc.inactiveNodes[i].url = &url
			sc.inactiveNodes[i].currentTopology = hash
			continue
		}
	}
	sc.inactiveNodesMu.Unlock()
	if !found {
		newNode := &SnowthNode{
			identifier: topology.ID,
			url: &url.URL{
				Scheme: "http",
				Host:   fmt.Sprintf("%s:%d", topology.Address, topology.APIPort),
			},
			currentTopology: hash,
		}
		sc.AddNodes(newNode)
		sc.ActivateNodes(newNode)
	}
}

// doChangeActivation - perform an activation state change
func (sc *SnowthClient) doChangeActivation(from, to *[]*SnowthNode, nodes []*SnowthNode) {
	sc.activeNodesMu.Lock()
	defer sc.activeNodesMu.Unlock()
	sc.inactiveNodesMu.Lock()
	defer sc.inactiveNodesMu.Unlock()
	for _, v := range nodes {
		moveNode(from, to, v)
	}
}

// ActivateNodes - given a list of nodes, make said nodes active for the client
func (sc *SnowthClient) ActivateNodes(nodes ...*SnowthNode) {
	sc.doChangeActivation(&sc.inactiveNodes, &sc.activeNodes, nodes)
}

// DeactivateNodes - given a list of nodes, make said nodes inactive
func (sc *SnowthClient) DeactivateNodes(nodes ...*SnowthNode) {
	sc.doChangeActivation(&sc.activeNodes, &sc.inactiveNodes, nodes)
}

// AddNodes - add nodes parameters to the inactive node list
func (sc *SnowthClient) AddNodes(nodes ...*SnowthNode) {
	sc.inactiveNodesMu.Lock()
	defer sc.inactiveNodesMu.Unlock()
	sc.inactiveNodes = append(sc.inactiveNodes, nodes...)
}

// doListNodes - helper to list the nodes, active or inactive
func doListNodes(nodes *[]*SnowthNode, mu *sync.RWMutex) []*SnowthNode {
	mu.RLock()
	defer mu.RUnlock()
	var result = []*SnowthNode{}
	for _, url := range *nodes {
		result = append(result, url)
	}
	return result
}

// ListInactiveNodes - list all of the currently inactive nodes
func (sc *SnowthClient) ListInactiveNodes() []*SnowthNode {
	return doListNodes(&sc.inactiveNodes, sc.inactiveNodesMu)
}

// ListActiveNodes - list all of the currently active nodes
func (sc *SnowthClient) ListActiveNodes() []*SnowthNode {
	return doListNodes(&sc.activeNodes, sc.activeNodesMu)
}

// do - helper to perform the request for the client
func (sc *SnowthClient) do(r *http.Request) (*http.Response, error) {
	return sc.c.Do(r)
}

// getURL - helper to resolve a reference against a particular node
func (sc *SnowthClient) getURL(node *SnowthNode, ref string) string {
	return resolveURL(node.url, ref)
}
