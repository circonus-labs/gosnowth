package gosnowth

import (
	"bytes"
	"context"
	"encoding/json"
	"path"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// NNTAllValueResponse values represent NNT data responses from IRONdb.
type NNTAllValueResponse struct {
	Data []NNTAllValue
}

// UnmarshalJSON decodes a JSON format byte slice into an NNTAllValueResponse.
func (nv *NNTAllValueResponse) UnmarshalJSON(b []byte) error {
	nv.Data = []NNTAllValue{}
	values := [][]interface{}{}
	if err := json.Unmarshal(b, &values); err != nil {
		return errors.Wrap(err, "failed to deserialize nnt average response")
	}

	for _, entry := range values {
		var nav = NNTAllValue{}
		if m, ok := entry[1].(map[string]interface{}); ok {
			valueBytes, err := json.Marshal(m)
			if err != nil {
				return errors.Wrap(err,
					"failed to marshal intermediate value from tuple")
			}

			if err := json.Unmarshal(valueBytes, &nav); err != nil {
				return errors.Wrap(err,
					"failed to unmarshal value from tuple")
			}
		}

		// grab the timestamp
		if v, ok := entry[0].(float64); ok {
			nav.Time = time.Unix(int64(v), 0)
		}

		nv.Data = append(nv.Data, nav)
	}

	return nil
}

// NNTAllValue values represent NNT data.
type NNTAllValue struct {
	Time              time.Time `json:"-"`
	Count             int64     `json:"count"`
	Value             int64     `json:"value"`
	StdDev            int64     `json:"stddev"`
	Derivative        int64     `json:"derivative"`
	DerivativeStdDev  int64     `json:"derivative_stddev"`
	Counter           int64     `json:"counter"`
	CounterStdDev     int64     `json:"counter_stddev"`
	Derivative2       int64     `json:"derivative2"`
	Derivative2StdDev int64     `json:"derivative2_stddev"`
	Counter2          int64     `json:"counter2"`
	Counter2StdDev    int64     `json:"counter2_stddev"`
}

// NNTValueResponse values represent responses containing NNT data.
type NNTValueResponse struct {
	Data []NNTValue
}

// UnmarshalJSON decodes a JSON format byte slice into an NNTValueResponse.
func (nv *NNTValueResponse) UnmarshalJSON(b []byte) error {
	nv.Data = []NNTValue{}
	values := [][]int64{}
	if err := json.Unmarshal(b, &values); err != nil {
		return errors.Wrap(err, "failed to deserialize nnt average response")
	}

	for _, tuple := range values {
		nv.Data = append(nv.Data, NNTValue{
			Time:  time.Unix(tuple[0], 0),
			Value: tuple[1],
		})
	}

	return nil
}

// NNTValue values represent individual NNT data values.
type NNTValue struct {
	Time  time.Time
	Value int64
}

// NNTData values represent NNT data.
type NNTData struct {
	Count            int64  `json:"count"`
	Value            int64  `json:"value"`
	Derivative       int64  `json:"derivative"`
	Counter          int64  `json:"counter"`
	StdDev           int64  `json:"stddev"`
	DerivativeStdDev int64  `json:"derivative_stddev"`
	CounterStdDev    int64  `json:"counter_stddev"`
	Metric           string `json:"metric"`
	ID               string `json:"id"`
	Offset           int64  `json:"offset"`
	Parts            Parts  `json:"parts"`
}

// NNTPartsData values represent NNT base data parts.
type NNTPartsData struct {
	Count            int64 `json:"count"`
	Value            int64 `json:"value"`
	Derivative       int64 `json:"derivative"`
	Counter          int64 `json:"counter"`
	StdDev           int64 `json:"stddev"`
	DerivativeStdDev int64 `json:"derivative_stddev"`
	CounterStdDev    int64 `json:"counter_stddev"`
}

// Parts values contain the NNTData submission parts of an NNT rollup.
type Parts struct {
	Period int64          `json:"period"`
	Data   []NNTPartsData `json:"data"`
}

// MarshalJSON marshals a Parts value into a JSON format byte slice.
func (p *Parts) MarshalJSON() ([]byte, error) {
	tuple := []interface{}{}
	tuple = append(tuple, p.Period, p.Data)
	buf := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(buf)
	if err := enc.Encode(tuple); err != nil {
		return buf.Bytes(), err
	}

	return buf.Bytes(), nil
}

// WriteNNT writes NNT data to a node.
func (sc *SnowthClient) WriteNNT(node *SnowthNode, data ...NNTData) error {
	return sc.WriteNNTContext(context.Background(), node, data...)
}

// WriteNNTContext is the context aware version of WriteNNT.
func (sc *SnowthClient) WriteNNTContext(ctx context.Context, node *SnowthNode,
	data ...NNTData) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return errors.Wrap(err, "failed to encode NNTData for write")
	}

	return sc.do(ctx, node, "POST", "/write/nnt", buf, nil, nil)
}

// ReadNNTValues reads NNT data from a node.
func (sc *SnowthClient) ReadNNTValues(node *SnowthNode, start, end time.Time,
	period int64, t, id, metric string) ([]NNTValue, error) {
	return sc.ReadNNTValuesContext(context.Background(), node, start, end,
		period, t, id, metric)
}

// ReadNNTValuesContext is the context aware version of ReadNNTValues.
func (sc *SnowthClient) ReadNNTValuesContext(ctx context.Context,
	node *SnowthNode, start, end time.Time, period int64,
	t, id, metric string) ([]NNTValue, error) {
	nv := new(NNTValueResponse)
	err := sc.do(ctx, node, "GET", path.Join("/read",
		strconv.FormatInt(start.Unix(), 10),
		strconv.FormatInt(end.Unix(), 10),
		strconv.FormatInt(period, 10), id, t, metric),
		nil, nv, decodeJSON)
	return nv.Data, err
}

// ReadNNTAllValues reads all NNT data from a node.
func (sc *SnowthClient) ReadNNTAllValues(node *SnowthNode,
	start, end time.Time, period int64,
	id, metric string) ([]NNTAllValue, error) {
	return sc.ReadNNTAllValuesContext(context.Background(), node, start, end,
		period, id, metric)
}

// ReadNNTAllValuesContext is the context aware version of ReadNNTAllValues.
func (sc *SnowthClient) ReadNNTAllValuesContext(ctx context.Context,
	node *SnowthNode, start, end time.Time, period int64,
	id, metric string) ([]NNTAllValue, error) {
	nv := new(NNTAllValueResponse)
	err := sc.do(ctx, node, "GET", path.Join("/read",
		strconv.FormatInt(start.Unix(), 10),
		strconv.FormatInt(end.Unix(), 10),
		strconv.FormatInt(period, 10), id, "all", metric),
		nil, nv, decodeJSON)
	return nv.Data, err
}