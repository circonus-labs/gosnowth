package gosnowth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/circonus-labs/circonusllhist"
	"github.com/pkg/errors"
)

// WriteHistogram - Write Histogram data to a node, data should be a slice of
// Histogram Data and node is the node to write the data to
func (sc *SnowthClient) WriteHistogram(node *SnowthNode, data ...HistogramData) error {

	buf := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(buf)
	if err := enc.Encode(data); err != nil {
		return errors.Wrap(err, "failed to encode HistogramData for write")
	}

	req, err := http.NewRequest("POST", sc.getURL(node, "/histogram/write"), buf)
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

// HistogramData - representation of Text Data for data submission and retrieval
type HistogramData struct {
	Metric    string                    `json:"metric"`
	ID        string                    `json:"id"`
	Offset    int64                     `json:"offset"`
	Period    int64                     `json:"period"`
	Histogram *circonusllhist.Histogram `json:"histogram"`
}
