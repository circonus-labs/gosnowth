package gosnowth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// WriteNNT - Write NNT data to a node, data should be a slice of NNTData
// and node is the node to write the data to
func (sc *SnowthClient) WriteNNT(node *SnowthNode, data ...NNTData) error {
	buf := new(bytes.Buffer)
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
