package gosnowth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/circonus-labs/circonusllhist"
	"github.com/pkg/errors"
)

// HistogramValue values are individual data points of a rollup.
type HistogramValue struct {
	Time   time.Time
	Period time.Duration
	Data   map[string]int64
}

// MarshalJSON encodes a HistogramValue value into a JSON format byte slice.
func (hv *HistogramValue) MarshalJSON() ([]byte, error) {
	v := []interface{}{}
	tn := float64(0)
	fv, err := strconv.ParseFloat(formatTimestamp(hv.Time), 64)
	if err == nil {
		tn = float64(fv)
	}

	v = append(v, tn)
	v = append(v, hv.Period.Seconds())
	v = append(v, hv.Data)
	return json.Marshal(v)
}

// UnmarshalJSON decodes a JSON format byte slice into a HistogramValue value.
func (hv *HistogramValue) UnmarshalJSON(b []byte) error {
	v := []interface{}{}
	json.Unmarshal(b, &v)
	if len(v) != 3 {
		return fmt.Errorf("histogram value should contain three entries: %s",
			string(b))
	}

	if fv, ok := v[0].(float64); ok {
		tv, err := parseTimestamp(strconv.FormatFloat(fv, 'f', 3, 64))
		if err != nil {
			return err
		}

		hv.Time = tv
	}

	if fv, ok := v[1].(float64); ok {
		hv.Period = time.Duration(fv) * time.Second
	}

	if m, ok := v[2].(map[string]interface{}); ok {
		hv.Data = make(map[string]int64, len(m))
		for k, iv := range m {
			if fv, ok := iv.(float64); ok {
				hv.Data[k] = int64(fv)
			}
		}
	}

	return nil
}

// Timestamp returns the HistogramValue time as a string in the IRONdb
// timestamp format.
func (hv *HistogramValue) Timestamp() string {
	return formatTimestamp(hv.Time)
}

// ReadHistogramValues reads histogram data from a node.
func (sc *SnowthClient) ReadHistogramValues(
	node *SnowthNode, id, metric string, tags []string, period time.Duration,
	start, end time.Time) ([]HistogramValue, error) {
	return sc.ReadHistogramValuesContext(context.Background(), node, id, metric,
		tags, period, start, end)
}

// ReadHistogramValuesContext is the context aware version of
// ReadHistogramValues.
func (sc *SnowthClient) ReadHistogramValuesContext(ctx context.Context,
	node *SnowthNode, uuid, metric string, tags []string, period time.Duration,
	start, end time.Time) ([]HistogramValue, error) {
	startTS := start.Unix() - start.Unix()%int64(period.Seconds())
	endTS := end.Unix() - end.Unix()%int64(period.Seconds()) +
		int64(period.Seconds())
	var metricBuilder strings.Builder
	metricBuilder.WriteString(metric)
	if len(tags) > 0 {
		metricBuilder.WriteString("|ST[")
		metricBuilder.WriteString(strings.Join(tags, ","))
		metricBuilder.WriteString("]")
	}

	r := []HistogramValue{}
	body, _, err := sc.do(ctx, node, "GET",
		path.Join("/histogram", strconv.FormatInt(startTS, 10),
			strconv.FormatInt(endTS, 10),
			strconv.FormatInt(int64(period.Seconds()), 10), uuid,
			url.QueryEscape(metricBuilder.String())), nil)
	if err != nil {
		return nil, err
	}

	if err := decodeJSON(body, &r); err != nil {
		return nil, errors.Wrap(err, "unable to decode IRONdb response")
	}

	return r, nil
}

// HistogramData values represent histogram data records in IRONdb.
type HistogramData struct {
	AccountID int64                     `json:"account_id"`
	Metric    string                    `json:"metric"`
	ID        string                    `json:"id"`
	CheckName string                    `json:"check_name"`
	Offset    int64                     `json:"offset"`
	Period    int64                     `json:"period"`
	Histogram *circonusllhist.Histogram `json:"histogram"`
}

// WriteHistogram sends a variadic list of histogram data values to be written
// to an IRONdb node.
func (sc *SnowthClient) WriteHistogram(node *SnowthNode,
	data ...HistogramData) error {
	return sc.WriteHistogramContext(context.Background(), node, data...)
}

// WriteHistogramContext is the context aware version of WriteHistogram.
func (sc *SnowthClient) WriteHistogramContext(ctx context.Context,
	node *SnowthNode, data ...HistogramData) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return errors.Wrap(err, "failed to encode HistogramData for write")
	}

	_, _, err := sc.do(ctx, node, "POST", "/histogram/write", buf)
	return err
}
