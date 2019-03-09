package gosnowth

import (
	"context"
	"path"
)

// LocateMetric locates which nodes contain specified metric data.
func (sc *SnowthClient) LocateMetric(uuid string, metric string,
	node *SnowthNode) (location *Topology, err error) {
	return sc.LocateMetricContext(context.Background(), uuid, metric, node)
}

// LocateMetricContext is the context aware version of LocateMetric.
func (sc *SnowthClient) LocateMetricContext(ctx context.Context, uuid string,
	metric string, node *SnowthNode) (location *Topology, err error) {
	location = new(Topology)
	err = sc.do(ctx, node, "GET", path.Join("/locate/xml", uuid, metric),
		nil, location, decodeXML)
	return
}
