package gosnowth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// WriteText - Write Text data to a node, data should be a slice of TextData
// and node is the node to write the data to
func (sc *SnowthClient) WriteText(data []TextData, node *SnowthNode) error {

	buf := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(buf)
	if err := enc.Encode(data); err != nil {
		return errors.Wrap(err, "failed to encode TextData for write")
	}

	req, err := http.NewRequest("POST", sc.getURL(node, "/write/text"), buf)
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

// TextData - representation of Text Data for data submission and retrieval
type TextData struct {
	Metric string `json:"metric"`
	ID     string `json:"id"`
	Offset string `json:"offset"`
	Value  string `json:"value"`
}
