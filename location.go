package gosnowth

import (
	"net/http"
	"path"

	"github.com/pkg/errors"
)

// LocateMetric - locate which nodes a metric lives on
func (sc *SnowthClient) LocateMetric(uuid string, metric string, node *SnowthNode) (*DataLocation, error) {
	var resource = path.Join("/locate/xml", uuid, metric)

	req, err := http.NewRequest("GET", sc.getURL(node, resource), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	resp, err := sc.do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform request")
	}

	var locations = new(DataLocation)
	if err := decodeXMLFromResponse(locations, resp); err != nil {
		return nil, errors.Wrap(err, "failed to decode")
	}

	return locations, nil
}

// DataLocation is from the location api and mimics the topology response
type DataLocation Topology
