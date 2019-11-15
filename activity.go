package gosnowth

import (
	"context"

	"github.com/pkg/errors"
)

// RebuildActivityRequest values represent a request to rebuild activity tracking data.
type RebuildActivityRequest struct {
	UUID   string `json:"check_uuid"`
	Metric string `json:"metric_name"`
}

// RebuildActivityResponse values represent a response to activity tracking rebuilds.
type RebuildActivityResponse WriteRawResponse

// RebuildActivity rebuilds IRONdb activity tracking data for a list of metrics.
func (sc *SnowthClient) RebuildActivity(node *SnowthNode, rebuildRequest []RebuildActivityRequest) (*RebuildActivityResponse, error) {
	return sc.RebuildActivityContext(context.Background(), node, rebuildRequest)
}

// RebuildActivityContext is the context aware version of RebuildActivity.
func (sc *SnowthClient) RebuildActivityContext(ctx context.Context, node *SnowthNode,
	rebuildRequest []RebuildActivityRequest) (*RebuildActivityResponse, error) {

	data, err := encodeJSON(rebuildRequest)
	if err != nil {
		return nil, err
	}
	body, _, err := sc.do(ctx, node, "POST", "/surrogate/activity_rebuild", data, nil)
	if err != nil {
		return nil, err
	}

	r := &RebuildActivityResponse{}
	if err := decodeJSON(body, &r); err != nil {
		return nil, errors.Wrap(err, "unable to decode IRONdb response")
	}

	return r, nil
}
