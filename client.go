// Package gosnowth contains an IRONdb client library written in Go.
package gosnowth

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// Logger values implement the behavior used by SnowthClient for logging,
// if the client has been assigned a logger with this interface.
type Logger interface {
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}

// SnowthNode values represent a snowth node. An IRONdb cluster consists of
// several nodes.  A SnowthNode has a URL to the API of that Node, an identifier,
// and a current topology.  The identifier is how the node is identified within
// the cluster, and the topology is the current topology that the node falls
// within.  A topology is a set of nodes that distribute data amongst each other.
type SnowthNode struct {
	url             *url.URL
	identifier      string
	currentTopology string
	semVer          string
}

// GetURL returns the *url.URL for a given SnowthNode. This is useful if you
// need the raw connection string of a given snowth node, such as when making a
// proxy for a snowth node.
func (sn *SnowthNode) GetURL() *url.URL {
	return sn.url
}

// SemVer returns a string containing the semantic version of IRONdb the node
// is currently running.
func (sn *SnowthNode) SemVer() string {
	return sn.semVer
}

// GetCurrentTopology return the hash string representation of the
// node's current topology.
func (sn *SnowthNode) GetCurrentTopology() string {
	return sn.currentTopology
}

// httpClient values are used to define the behavior needed from HTTP client
// values.
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// SnowthClient values provide client functionality for accessing IRONdb.
// Operations for the client can be broken down into 6 major sections:
//		1.) State and Topology
// Within the state and topology APIs, there are several useful apis, including
// apis to retrieve Node state, Node gossip information, topology information,
// and topo ring information.  Each of these operations is implemented as a method
// on this client.
//		2.) Rebalancing APIs
// In order to add or remove nodes within an IRONdb cluster you will have to use
// the rebalancing APIs.  Implemented within this package you will be able to
// load a new topology, rebalance nodes to the new topology, as well as check
// load state information and abort a topology change.
//		3.) Data Retrieval APIs
// IRONdb houses data, and the data retrieval APIs allow for accessing of that
// stored data.  Data types implemented include NNT, Text, and Histogram data
// element types.
//		4.) Data Submission APIs
// IRONdb houses data, to which you can use to submit data to the cluster.  Data
// types supported include the same for the retrieval APIs, NNT, Text and
// Histogram data types.
//		5.) Data Deletion APIs
// Data sometimes needs to be deleted, and that is performed with the data
// deletion APIs.  This client implements the data deletion apis to remove data
// from the nodes.
//		6.) Lua Extensions APIs
type SnowthClient struct {
	sync.RWMutex
	c httpClient

	// in order to keep track of healthy nodes within the cluster,
	// we have two lists of SnowthNode types, active and inactive.
	activeNodes   []*SnowthNode
	inactiveNodes []*SnowthNode

	// watchInterval is the duration between checks to tell if a node is active
	// or inactive.
	watchInterval time.Duration

	// If log output is desired, a value matching the Logger interface can be
	// assigned.  If this is nil, no log output will be attempted.
	log Logger

	// request is an assignable middleware function which modifies the request
	// before it is used by SnowthClient to connect with IRONdb. Tracing headers
	// or other context information can be added by this function.
	request func(r *http.Request) error

	// watch is an assignable middleware function which can plugin functionality
	// to activate or deactivate snowth cluster nodes during the watch and
	// update process, using custom logic.
	watch func(n *SnowthNode)
}

// NewSnowthClient initializes a new SnowthClient value, constructing all the
// required state to communicate with a cluster of IRONdb nodes.
// The discover parameter, when true, will allow the client to discover new
// nodes from the topology.
func NewSnowthClient(discover bool, addrs ...string) (*SnowthClient, error) {
	cfg, err := NewConfig()
	if err != nil {
		return nil, err
	}

	cfg.SetDiscover(discover)
	if err := cfg.SetServers(addrs...); err != nil {
		return nil, err
	}

	return NewClient(cfg)
}

// NewClient creates and performs initial setup of a new SnowthClient.
func NewClient(cfg *Config) (*SnowthClient, error) {
	client := &http.Client{
		Timeout: cfg.Timeout(),
		Transport: &http.Transport{
			Dial: (&net.Dialer{Timeout: cfg.DialTimeout()}).Dial,
		},
	}

	sc := &SnowthClient{
		c:             client,
		activeNodes:   []*SnowthNode{},
		inactiveNodes: []*SnowthNode{},
		watchInterval: cfg.WatchInterval(),
	}

	// For each of the addrs we need to parse the connection string,
	// then create a node for that connection string, poll the state
	// of that node, and populate the identifier and topology of that
	// node.  Finally we will add the node and activate it.
	numActiveNodes := 0
	nErr := newMultiError()
	for _, addr := range cfg.Servers() {
		url, err := url.Parse(addr)
		if err != nil {
			// This node had an error, put on inactive list.
			nErr.Add(errors.Wrap(err, "unable to parse server url"))
			continue
		}

		// Call get stats to populate the id of this node.
		node := &SnowthNode{url: url}
		stats, err := sc.GetStats(node)
		if err != nil {
			// This node had an error, put on inactive list.
			nErr.Add(errors.Wrap(err, "unable to get status of node"))
			continue
		}

		node.identifier = stats.Identity()
		node.currentTopology = stats.CurrentTopology()
		node.semVer = stats.SemVer()
		sc.AddNodes(node)
		sc.ActivateNodes(node)
		numActiveNodes++
	}

	if numActiveNodes == 0 {
		if nErr.HasError() {
			return nil, errors.Wrap(nErr, "no snowth nodes could be activated")
		}

		return nil, errors.New("no snowth nodes could be activated")
	}

	if cfg.Discover() {
		// For robustness, we will perform a discovery of associated nodes
		// this works by pulling the topology information for given nodes
		// and adding nodes discovered within the topology into the client.
		if err := sc.discoverNodes(); err != nil {
			return nil, errors.Wrap(err,
				"failed to perform discovery of new nodes")
		}
	}

	return sc, nil
}

// SetRequestFunc sets an optional middleware function that is used to modify
// the HTTP request before it is used by SnowthClient to connect with IRONdb.
// Tracing headers or other context information provided by the user of this
// library can be added by this function.
func (sc *SnowthClient) SetRequestFunc(f func(r *http.Request) error) {
	sc.Lock()
	defer sc.Unlock()
	sc.request = f
}

// SetWatchFunc sets an optional middleware function that can be used to
// inspect and activate or deactivate IRONdb cluster nodes during the watch and
// update process.
func (sc *SnowthClient) SetWatchFunc(f func(n *SnowthNode)) {
	sc.Lock()
	defer sc.Unlock()
	sc.watch = f
}

// SetWatchInterval sets the interval at which the watch process executes.
func (sc *SnowthClient) SetWatchInterval(d time.Duration) {
	sc.Lock()
	defer sc.Unlock()
	sc.watchInterval = d
}

// SetLog assigns a logger to the snowth client.
func (sc *SnowthClient) SetLog(log Logger) {
	sc.Lock()
	defer sc.Unlock()
	sc.log = log
}

// LogInfof writes a log entry at the information level.
func (sc *SnowthClient) LogInfof(format string, args ...interface{}) {
	if sc.log != nil {
		sc.log.Infof(format, args...)
	}
}

// LogWarnf writes a log entry at the warning level.
func (sc *SnowthClient) LogWarnf(format string, args ...interface{}) {
	if sc.log != nil {
		sc.log.Warnf(format, args...)
	}
}

// LogErrorf writes a log entry at the error level.
func (sc *SnowthClient) LogErrorf(format string, args ...interface{}) {
	if sc.log != nil {
		sc.log.Errorf(format, args...)
	}
}

// LogDebugf writes a log entry at the debug level.
func (sc *SnowthClient) LogDebugf(format string, args ...interface{}) {
	if sc.log != nil {
		sc.log.Debugf(format, args...)
	}
}

// isNodeActive checks to see if a given node is active or not taking into
// account the ability to get the node state, gossip information and the gossip
// age of the node. If the age is larger than 10 the node is considered
// inactive.
func (sc *SnowthClient) isNodeActive(node *SnowthNode) bool {
	if node.identifier == "" || node.semVer == "" {
		// go get state to figure out identity
		stats, err := sc.GetStats(node)
		if err != nil {
			// error means we failed, node is not active
			sc.LogWarnf("unable to get the state of the node: %s",
				err.Error())
			return false
		}

		node.identifier = stats.Identity()
		node.semVer = stats.SemVer()
		sc.LogDebugf("retrieved state of node: %s -> %s",
			node.GetURL().Host, node.identifier)
	}

	gossip, err := sc.GetGossipInfo(node)
	if err != nil {
		sc.LogWarnf("unable to get the gossip info of the node: %s",
			err.Error())
		return false
	}

	age := float64(100)
	for _, entry := range []GossipDetail(*gossip) {
		if entry.ID == node.identifier {
			age = entry.Age
			break
		}
	}

	if age > 10.0 {
		sc.LogWarnf("gossip age expired: %s -> %d", node.GetURL().Host, age)
		return false
	}

	return true
}

// WatchAndUpdate watches gossip data for all nodes, and move the nodes to
// the active or inactive pools as required.  It returns a function to cancel
// the operation if needed. It accepts a context value as an argument which
// will also cancel the operation if the context is cancelled or expired. If
// context cancellation is not needed, nil can be passed as the argument.
func (sc *SnowthClient) WatchAndUpdate(ctx context.Context) {
	sc.RLock()
	defer sc.RUnlock()
	if sc.watchInterval <= time.Duration(0) {
		return
	}

	go func() {
		tick := time.NewTicker(sc.watchInterval)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				sc.LogDebugf("firing watch and update")
				for _, node := range sc.ListInactiveNodes() {
					sc.LogDebugf("checking node for inactive -> active: %s",
						node.GetURL().Host)
					if sc.isNodeActive(node) {
						// Move to active.
						sc.LogDebugf("active, moving to active list: %s",
							node.GetURL().Host)
						sc.ActivateNodes(node)
					}

					sc.RLock()
					if sc.watch != nil {
						sc.watch(node)
					}

					sc.RUnlock()
				}

				for _, node := range sc.ListActiveNodes() {
					sc.LogDebugf("checking node for active -> inactive: %s",
						node.GetURL().Host)
					if !sc.isNodeActive(node) {
						// Move to inactive.
						sc.LogWarnf("inactive, moving to inactive list: %s",
							node.GetURL().Host)
						sc.DeactivateNodes(node)
					}

					sc.RLock()
					if sc.watch != nil {
						sc.watch(node)
					}

					sc.RUnlock()
				}
			}
		}
	}()
}

// discoverNodes attempts to discover peer nodes related to the topology.
// This function will go through the active nodes and get the topology
// information which shows all other nodes included in the cluster, then adds
// them as nodes to this client's active node pool.
func (sc *SnowthClient) discoverNodes() error {
	success := false
	mErr := newMultiError()
	for _, node := range sc.ListActiveNodes() {
		// lookup the topology
		topology, err := sc.GetTopologyInfo(node)
		if err != nil {
			mErr.Add(errors.Wrap(err, "error getting topology info: %+v"))
			continue
		}

		// populate all the nodes with the appropriate topology information
		for _, topoNode := range topology.Nodes {
			sc.populateNodeInfo(node.GetCurrentTopology(), topoNode)
		}

		success = true
	}

	if !success {
		// we didn't get any topology information, therefore we didn't
		// discover correctly, return the multitude of errors
		return mErr
	}

	return nil
}

// populateNodeInfo populates an existing node with details from the topology.
// If a node doesn't exist, it will be added to the list of active nodes.
func (sc *SnowthClient) populateNodeInfo(hash string, topology TopologyNode) {
	sc.Lock()
	found := false
	for i := 0; i < len(sc.activeNodes); i++ {
		if sc.activeNodes[i].identifier == topology.ID {
			found = true
			url := url.URL{
				Scheme: "http",
				Host: fmt.Sprintf("%s:%d", topology.Address,
					topology.APIPort),
			}

			sc.activeNodes[i].url = &url
			sc.activeNodes[i].currentTopology = hash
			continue
		}
	}

	for i := 0; i < len(sc.inactiveNodes); i++ {
		found = true
		if sc.inactiveNodes[i].identifier == topology.ID {
			url := url.URL{
				Scheme: "http",
				Host: fmt.Sprintf("%s:%d", topology.Address,
					topology.APIPort),
			}

			sc.inactiveNodes[i].url = &url
			sc.inactiveNodes[i].currentTopology = hash
			continue
		}
	}

	sc.Unlock()
	if !found {
		newNode := &SnowthNode{
			identifier: topology.ID,
			url: &url.URL{
				Scheme: "http",
				Host: fmt.Sprintf("%s:%d", topology.Address,
					topology.APIPort),
			},
			currentTopology: hash,
		}

		sc.AddNodes(newNode)
		sc.ActivateNodes(newNode)
	}
}

// ActivateNodes makes provided nodes active.
func (sc *SnowthClient) ActivateNodes(nodes ...*SnowthNode) {
	sc.Lock()
	defer sc.Unlock()
	in := []*SnowthNode{}
	match := false
	for _, iv := range sc.inactiveNodes {
		match = false
		for _, v := range nodes {
			if v.GetURL().String() == iv.GetURL().String() {
				match = true
				break
			}
		}

		if !match {
			in = append(in, iv)
		}
	}

	sc.inactiveNodes = in
	an := []*SnowthNode{}
	for _, v := range nodes {
		match = false
		for _, av := range sc.activeNodes {
			if v.GetURL().String() == av.GetURL().String() {
				match = true
				break
			}
		}

		if !match {
			an = append(an, v)
		}
	}

	sc.activeNodes = append(sc.activeNodes, an...)
}

// DeactivateNodes makes provided nodes inactive.
func (sc *SnowthClient) DeactivateNodes(nodes ...*SnowthNode) {
	sc.Lock()
	defer sc.Unlock()
	an := []*SnowthNode{}
	match := false
	for _, av := range sc.activeNodes {
		match = false
		for _, v := range nodes {
			if v.GetURL().String() == av.GetURL().String() {
				match = true
				break
			}
		}

		if !match {
			an = append(an, av)
		}
	}

	sc.activeNodes = an
	in := []*SnowthNode{}
	for _, v := range nodes {
		match = false
		for _, iv := range sc.inactiveNodes {
			if v.GetURL().String() == iv.GetURL().String() {
				match = true
				break
			}
		}

		if !match {
			in = append(in, v)
		}
	}

	sc.inactiveNodes = append(sc.inactiveNodes, in...)
}

// AddNodes adds node values to the inactive node list.
func (sc *SnowthClient) AddNodes(nodes ...*SnowthNode) {
	sc.Lock()
	defer sc.Unlock()
	in := []*SnowthNode{}
	match := false
	for _, v := range nodes {
		match = false
		for _, iv := range sc.inactiveNodes {
			if v.GetURL().String() == iv.GetURL().String() {
				match = true
				break
			}
		}

		if !match {
			in = append(in, v)
		}
	}

	sc.inactiveNodes = append(sc.inactiveNodes, in...)
}

// ListInactiveNodes lists all of the currently inactive nodes.
func (sc *SnowthClient) ListInactiveNodes() []*SnowthNode {
	sc.RLock()
	defer sc.RUnlock()
	result := []*SnowthNode{}
	for _, url := range sc.inactiveNodes {
		result = append(result, url)
	}

	return result
}

// ListActiveNodes lists all of the currently active nodes.
func (sc *SnowthClient) ListActiveNodes() []*SnowthNode {
	sc.RLock()
	defer sc.RUnlock()
	result := []*SnowthNode{}
	for _, url := range sc.activeNodes {
		result = append(result, url)
	}

	return result
}

// do sends a request to IRONdb.
func (sc *SnowthClient) do(ctx context.Context, node *SnowthNode,
	method, url string, body io.Reader, respValue interface{},
	decodeFunc func(io.Reader, interface{}) error) error {
	if ctx == nil {
		ctx = context.Background()
	}

	r, err := http.NewRequest(method, sc.getURL(node, url), body)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	r.Close = true
	r = r.WithContext(ctx)
	sc.RLock()
	rf := sc.request
	sc.RUnlock()
	if rf != nil {
		if err := rf(r); err != nil {
			return errors.Wrap(err, "unable to process request")
		}

		if r == nil {
			return errors.New("invalid request after processing")
		}
	}

	sc.LogDebugf("snowth request: %+v", r)
	var start = time.Now()
	resp, err := sc.c.Do(r)
	if err != nil {
		return errors.Wrap(err, "failed to perform request")
	}

	defer resp.Body.Close()
	sc.LogDebugf("snowth response: %+v", resp)
	sc.LogDebugf("snowth latency: %+v", time.Since(start))
	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "context terminated")
	default:
		break
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		sc.LogWarnf("error returned from IRONdb: [%d] %s",
			resp.StatusCode, string(body))
		return fmt.Errorf("error returned from IRONdb: [%d] %s",
			resp.StatusCode, string(body))
	}

	if respValue != nil {
		if err := decodeFunc(resp.Body, respValue); err != nil {
			return errors.Wrap(err, "unable to decode IRONdb response")
		}
	}

	return nil
}

// getURL resolves the URL with a reference for a particular node.
func (sc *SnowthClient) getURL(node *SnowthNode, ref string) string {
	return resolveURL(node.url, ref)
}
