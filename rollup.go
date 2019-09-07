package gosnowth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// RollupValue values are individual data points of a rollup.
type RollupValue struct {
	Time  time.Time
	Value float64
}

// MarshalJSON encodes a RollupValue value into a JSON format byte slice.
func (rv *RollupValue) MarshalJSON() ([]byte, error) {
	v := []interface{}{}
	tn := float64(0)
	fv, err := strconv.ParseFloat(formatTimestamp(rv.Time), 64)
	if err == nil {
		tn = float64(fv)
	}

	v = append(v, tn)
	v = append(v, rv.Value)
	return json.Marshal(v)
}

// UnmarshalJSON decodes a JSON format byte slice into a RollupValue value.
func (rv *RollupValue) UnmarshalJSON(b []byte) error {
	v := []interface{}{}
	json.Unmarshal(b, &v)
	if len(v) != 2 {
		return fmt.Errorf("rollup value should contain two entries: %s",
			string(b))
	}

	if fv, ok := v[0].(float64); ok {
		tv, err := parseTimestamp(strconv.FormatFloat(fv, 'f', 3, 64))
		if err != nil {
			return err
		}

		rv.Time = tv
	}

	if fv, ok := v[1].(float64); ok {
		rv.Value = fv
	}

	return nil
}

// Timestamp returns the RollupValue time as a string in the IRONdb timestamp
// format.
func (rv *RollupValue) Timestamp() string {
	return formatTimestamp(rv.Time)
}

// RollupAllValue values contain all parts of an individual rollup data point.
type RollupAllValue struct {
	Time              time.Time
	Count             int64
	Counter           float64
	Counter2          float64
	CounterStddev     float64
	Counter2Stddev    float64
	Derivative        float64
	Derivative2       float64
	DerivativeStddev  float64
	Derivative2Stddev float64
	Stddev            float64
	Value             float64
}

// MarshalJSON encodes a RollupValue value into a JSON format byte slice.
func (rv *RollupAllValue) MarshalJSON() ([]byte, error) {
	v := []interface{}{}
	tn := float64(0)
	fv, err := strconv.ParseFloat(formatTimestamp(rv.Time), 64)
	if err == nil {
		tn = float64(fv)
	}

	v = append(v, tn)
	v = append(v, map[string]interface{}{
		"count":              rv.Count,
		"value":              rv.Value,
		"stddev":             rv.Stddev,
		"derivative":         rv.Derivative,
		"derivative_stddev":  rv.DerivativeStddev,
		"counter":            rv.Counter,
		"counter_stddev":     rv.CounterStddev,
		"derivative2":        rv.Derivative2,
		"derivative2_stddev": rv.Derivative2Stddev,
		"counter2":           rv.Counter2,
		"counter2_stddev":    rv.Counter2Stddev,
	})

	return json.Marshal(v)
}

// UnmarshalJSON decodes a JSON format byte slice into a RollupValue value.
func (rv *RollupAllValue) UnmarshalJSON(b []byte) error {
	v := []interface{}{}
	json.Unmarshal(b, &v)
	if len(v) != 2 {
		return fmt.Errorf("rollup value should contain two entries: %s",
			string(b))
	}

	if fv, ok := v[0].(float64); ok {
		tv, err := parseTimestamp(strconv.FormatFloat(fv, 'f', 3, 64))
		if err != nil {
			return err
		}

		rv.Time = tv
	}

	if m, ok := v[1].(map[string]interface{}); ok {
		for key, val := range m {
			if fv := val.(float64); ok {
				switch key {
				case "count":
					rv.Count = int64(fv)
				case "value":
					rv.Value = fv
				case "stddev":
					rv.Stddev = fv
				case "derivative":
					rv.Derivative = fv
				case "derivative_stddev":
					rv.DerivativeStddev = fv
				case "counter":
					rv.Counter = fv
				case "counter_stddev":
					rv.CounterStddev = fv
				case "derivative2":
					rv.Derivative2 = fv
				case "derivative2_stddev":
					rv.Derivative2Stddev = fv
				case "counter2":
					rv.Counter2 = fv
				case "counter2_stddev":
					rv.Counter2Stddev = fv
				}
			}
		}
	}

	return nil
}

// Timestamp returns the RollupAllValue time as a string in the IRONdb
// timestamp format.
func (rv *RollupAllValue) Timestamp() string {
	return formatTimestamp(rv.Time)
}

// ReadRollupValues reads rollup data from a node.
func (sc *SnowthClient) ReadRollupValues(
	node *SnowthNode, id, metric string, tags []string, rollup time.Duration,
	start, end time.Time) ([]RollupValue, error) {
	return sc.ReadRollupValuesContext(context.Background(), node, id, metric,
		tags, rollup, start, end)
}

// ReadRollupValuesContext is the context aware version of ReadRollupValues.
func (sc *SnowthClient) ReadRollupValuesContext(ctx context.Context,
	node *SnowthNode, id, metric string, tags []string, rollup time.Duration,
	start, end time.Time) ([]RollupValue, error) {
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

	r := []RollupValue{}
	body, _, err := sc.do(ctx, node, "GET",
		fmt.Sprintf("%s?start_ts=%d&end_ts=%d&rollup_span=%ds&type=average",
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

// ReadRollupAllValues reads rollup data from a node.
func (sc *SnowthClient) ReadRollupAllValues(
	node *SnowthNode, id, metric string, tags []string, rollup time.Duration,
	start, end time.Time) ([]RollupAllValue, error) {
	return sc.ReadRollupAllValuesContext(context.Background(), node, id, metric,
		tags, rollup, start, end)
}

// ReadRollupAllValuesContext is the context aware version of ReadRollupValues.
func (sc *SnowthClient) ReadRollupAllValuesContext(ctx context.Context,
	node *SnowthNode, id, metric string, tags []string, rollup time.Duration,
	start, end time.Time) ([]RollupAllValue, error) {
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

	r := []RollupAllValue{}
	body, _, err := sc.do(ctx, node, "GET",
		fmt.Sprintf("%s?start_ts=%d&end_ts=%d&rollup_span=%ds&type=all",
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
