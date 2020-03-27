package gosnowth

import (
	"context"
	"path"

	"github.com/pkg/errors"
)

func (sc *SnowthClient) LocateMetric(uuid string, metric string) ([]TopologyNode, error) {
	topo, err := sc.Topology()
	if err != nil {
		return nil, err
	}
	return topo.FindMetric(uuid, metric)
}

// LocateMetric locates which nodes contain specified metric data.
func (sc *SnowthClient) LocateMetricRemote(uuid string, metric string,
	node *SnowthNode) ([]TopologyNode, error) {
	return sc.LocateMetricRemoteContext(context.Background(), uuid, metric, node)
}

// LocateMetricContext is the context aware version of LocateMetric.
func (sc *SnowthClient) LocateMetricRemoteContext(ctx context.Context, uuid string,
	metric string, node *SnowthNode) ([]TopologyNode, error) {
	r := &Topology{}
	if node == nil {
		nodes := sc.ListActiveNodes()
		if len(nodes) == 0 {
			return nil, errors.New("no active nodes")
		}
		node = nodes[0]
	}
	body, _, err := sc.do(ctx, node, "GET",
		path.Join("/locate/xml", uuid, metric), nil, nil)
	if err != nil {
		return nil, err
	}

	if err := decodeXML(body, &r); err != nil {
		return nil, errors.Wrap(err, "unable to decode IRONdb response")
	}
	if r.WriteCopies == 0 {
		r.WriteCopies = r.OldWriteCopies
	}
	r.OldWriteCopies = r.WriteCopies

	return r.Nodes, nil
}
