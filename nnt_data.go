package gosnowth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// WriteNNT - Write NNT data to a node, data should be a slice of NNTData
// and node is the node to write the data to
func (sc *SnowthClient) WriteNNT(data []NNTData, node *SnowthNode) error {
	buf := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(buf)
	if err := enc.Encode(data); err != nil {
		return errors.Wrap(err, "failed to encode NNTData for write")
	}

	req, err := http.NewRequest("POST", sc.getURL(node, "/write/nnt"), buf)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	resp, err := sc.do(req)
	if err != nil {
		return errors.Wrap(err, "failed to perform request")
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		return fmt.Errorf("non-success status code returned: %s -> %s",
			resp.Status, string(body))
	}

	return nil
}

func (sc *SnowthClient) ReadNNTAllValues(
	node *SnowthNode, start, end time.Time, period int64,
	id, metric string) ([]NNTAllValue, error) {

	var ref = path.Join("/read",
		strconv.FormatInt(start.Unix(), 10),
		strconv.FormatInt(end.Unix(), 10),
		strconv.FormatInt(period, 10), id, "all", metric)

	req, err := http.NewRequest("GET", sc.getURL(node, ref), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	resp, err := sc.do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform request")
	}

	var (
		nntvr = NNTAllValueResponse{}
	)
	if err := decodeJSONFromResponse(&nntvr, resp); err != nil {
		return nil, errors.Wrap(err, "failed to decode")
	}

	return nntvr.Data, nil
}

type NNTAllValueResponse struct {
	Data []NNTAllValue
}

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

type NNTAllValue struct {
	Time              time.Time `json:"-"`
	Count             int64     `json:"count"`
	Value             int64     `json:"value"`
	StdDev            int64     `json:"stddev"`
	Derivitive        int64     `json:"derivative"`
	DerivitiveStdDev  int64     `json:"derivative_stddev"`
	Counter           int64     `json:"counter"`
	CounterStdDev     int64     `json:"counter_stddev"`
	Derivative2       int64     `json:"derivative2"`
	Derivative2StdDev int64     `json:"derivative2_stddev"`
	Counter2          int64     `json:"counter2"`
	Counter2StdDev    int64     `json:"counter2_stddev"`
}

// ReadNNTAverage - Read NNT data from a node
func (sc *SnowthClient) ReadNNTValue(
	node *SnowthNode, start, end time.Time, period int64,
	t, id, metric string) ([]NNTValue, error) {

	var ref = path.Join("/read",
		strconv.FormatInt(start.Unix(), 10),
		strconv.FormatInt(end.Unix(), 10),
		strconv.FormatInt(period, 10), id, t, metric)

	req, err := http.NewRequest("GET", sc.getURL(node, ref), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	resp, err := sc.do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform request")
	}

	var (
		nntvr = NNTValueResponse{}
	)
	if err := decodeJSONFromResponse(&nntvr, resp); err != nil {
		return nil, errors.Wrap(err, "failed to decode")
	}

	return nntvr.Data, nil
}

type NNTValueResponse struct {
	Data []NNTValue
}

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

type NNTValue struct {
	Time  time.Time
	Value int64
}

// ReadNNT - Read NNT data from a node
func (sc *SnowthClient) ReadNNT(data []NNTData, node *SnowthNode) error {

	return nil
}

type ReadNNTData struct {
}

// NNTData - representation of NNT Data for data submission and retrieval
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

// NNTBaseData - representation of NNT Base Data for data
// submission and retrieval
type NNTPartsData struct {
	Count            int64 `json:"count"`
	Value            int64 `json:"value"`
	Derivative       int64 `json:"derivative"`
	Counter          int64 `json:"counter"`
	StdDev           int64 `json:"stddev"`
	DerivativeStdDev int64 `json:"derivative_stddev"`
	CounterStdDev    int64 `json:"counter_stddev"`
}

// Parts - NNTData submission Parts that compose the NNT Rollup
type Parts struct {
	Period int64
	Data   []NNTPartsData
}

// MarshalJSON - cusom marshaller for the parts tuple structure
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
