package gosnowth

import (
	"bytes"
	"encoding/json"
	"path"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// WriteNNT writes NNT data to a node.
func (sc *SnowthClient) WriteNNT(node *SnowthNode, data ...NNTData) (err error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(data); err != nil {
		return errors.Wrap(err, "failed to encode NNTData for write")
	}
	err = sc.do(node, "POST", "/write/nnt", buf, nil, nil)
	return
}

// ReadNNTAllValues reads all NNT data form a node.
func (sc *SnowthClient) ReadNNTAllValues(
	node *SnowthNode, start, end time.Time, period int64,
	id, metric string) ([]NNTAllValue, error) {

	var (
		nntvr = new(NNTAllValueResponse)
		err   = sc.do(node, "GET", path.Join("/read",
			strconv.FormatInt(start.Unix(), 10),
			strconv.FormatInt(end.Unix(), 10),
			strconv.FormatInt(period, 10), id, "all", metric),
			nil, nntvr, decodeJSONFromResponse)
	)
	return nntvr.Data, err
}

// NNTAllValueResponse values represent NNT data responses from IRONdb.
type NNTAllValueResponse struct {
	Data []NNTAllValue
}

// UnmarshalJSON decodes a JSON format byte slice into an NNTAllValueResponse.
func (nntvr *NNTAllValueResponse) UnmarshalJSON(b []byte) error {
	nntvr.Data = []NNTAllValue{}
	var values = [][]interface{}{}

	if err := json.Unmarshal(b, &values); err != nil {
		return errors.Wrap(err, "failed to deserialize nnt average response")
	}

	for _, entry := range values {
		var nntavr = NNTAllValue{}
		if m, ok := entry[1].(map[string]interface{}); ok {
			valueBytes, err := json.Marshal(m)
			if err != nil {
				return errors.Wrap(err,
					"failed to marshal intermediate value from tuple")
			}
			if err := json.Unmarshal(valueBytes, &nntavr); err != nil {
				return errors.Wrap(err,
					"failed to unmarshal value from tuple")
			}
		}
		// grab the timestamp
		if v, ok := entry[0].(float64); ok {
			nntavr.Time = time.Unix(int64(v), 0)
		}
		nntvr.Data = append(nntvr.Data, nntavr)
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

// ReadNNTValues reads NNT data from a node.
func (sc *SnowthClient) ReadNNTValues(
	node *SnowthNode, start, end time.Time, period int64,
	t, id, metric string) ([]NNTValue, error) {

	var (
		nntvr = new(NNTValueResponse)
		err   = sc.do(node, "GET", path.Join("/read",
			strconv.FormatInt(start.Unix(), 10),
			strconv.FormatInt(end.Unix(), 10),
			strconv.FormatInt(period, 10), id, t, metric),
			nil, nntvr, decodeJSONFromResponse)
	)
	return nntvr.Data, err
}

// NNTValueResponse values represent responses containing NNT data.
type NNTValueResponse struct {
	Data []NNTValue
}

// UnmarshalJSON decodes a JSON format byte slice into an NNTValueResponse.
func (nntvr *NNTValueResponse) UnmarshalJSON(b []byte) error {
	nntvr.Data = []NNTValue{}
	var values = [][]int64{}

	if err := json.Unmarshal(b, &values); err != nil {
		return errors.Wrap(err, "failed to deserialize nnt average response")
	}

	for _, tuple := range values {
		nntvr.Data = append(nntvr.Data, NNTValue{
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

// ReadNNT reads NNT data from a node.
func (sc *SnowthClient) ReadNNT(data []NNTData, node *SnowthNode) error {

	return nil
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
	Period int64
	Data   []NNTPartsData
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
