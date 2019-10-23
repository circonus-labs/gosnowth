package gosnowth

import (
	"context"
	"path"

	"github.com/pkg/errors"
)

// LocateMetric locates which nodes contain specified metric data.
func (sc *SnowthClient) LocateMetric(uuid string, metric string,
	node *SnowthNode) (location *Topology, err error) {
	return sc.LocateMetricContext(context.Background(), uuid, metric, node)
}

// LocateMetricContext is the context aware version of LocateMetric.
func (sc *SnowthClient) LocateMetricContext(ctx context.Context, uuid string,
	metric string, node *SnowthNode) (*Topology, error) {
	r := &Topology{}
	body, _, err := sc.do(ctx, node, "GET",
		path.Join("/locate/xml", uuid, metric), nil, nil)
	if err != nil {
		return nil, err
	}

	if err := decodeXML(body, &r); err != nil {
		return nil, errors.Wrap(err, "unable to decode IRONdb response")
	}

	return r, nil
}
