package gosnowth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// RollupValues values are individual components of a roll up.
type RollupValues struct {
	Timestamp int64
	Value     float64
}

// UnmarshalJSON decodes a JSON format byte slice into a RollupValues value.
func (rv *RollupValues) UnmarshalJSON(b []byte) error {
	tt := []interface{}{&rv.Timestamp, &rv.Value}
	json.Unmarshal(b, &tt)
	if len(tt) < 2 { // error not enough fields
		return fmt.Errorf("rollup value should contain two entries,"+
			" %d given in payload", len(tt))
	}

	return nil
}

// ReadRollupValues reads rollup data from a node.
func (sc *SnowthClient) ReadRollupValues(
	node *SnowthNode, id, metric string, tags []string, rollup time.Duration,
	start, end time.Time) ([]RollupValues, error) {
	return sc.ReadRollupValuesContext(context.Background(), node, id, metric,
		tags, rollup, start, end)
}

// ReadRollupValuesContext is the context aware version of ReadRollupValues.
func (sc *SnowthClient) ReadRollupValuesContext(ctx context.Context,
	node *SnowthNode, id, metric string, tags []string, rollup time.Duration,
	start, end time.Time) ([]RollupValues, error) {
	startTS := start.Unix() - start.Unix()%int64(rollup/time.Second)
	endTS := end.Unix() - end.Unix()%int64(rollup/time.Second) +
		int64(rollup/time.Second)
	var metricBuilder strings.Builder
	metricBuilder.WriteString(metric)
	if len(tags) > 0 {
		metricBuilder.WriteString("|ST[")
		metricBuilder.WriteString(strings.Join(tags, ","))
		metricBuilder.WriteString("]")
	}

	r := []RollupValues{}
	body, _, err := sc.do(ctx, node, "GET",
		fmt.Sprintf("%s?start_ts=%d&end_ts=%d&rollup_span=%ds",
			path.Join("/rollup", id,
				url.QueryEscape(metricBuilder.String())),
			startTS, endTS, int(rollup/time.Second)), nil)
	if err != nil {
		return nil, err
	}

	if err := decodeJSON(body, &r); err != nil {
		return nil, errors.Wrap(err, "unable to decode IRONdb response")
	}

	return r, nil
}
