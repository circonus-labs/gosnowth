package gosnowth

import (
	"path"
)

// LocateMetric locates which nodes contain a metric.
func (sc *SnowthClient) LocateMetric(uuid string, metric string,
	node *SnowthNode) (location *Topology, err error) {
	location = new(Topology)
	err = sc.do(node, "GET", path.Join("/locate/xml", uuid, metric),
		nil, location, decodeXMLFromResponse)
	return
}
