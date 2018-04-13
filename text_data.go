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

func (sc *SnowthClient) ReadTextValues(
	node *SnowthNode, start, end time.Time,
	id, metric string) ([]TextValue, error) {

	var ref = path.Join("/read",
		strconv.FormatInt(start.Unix(), 10),
		strconv.FormatInt(end.Unix(), 10),
		id, metric)

	req, err := http.NewRequest("GET", sc.getURL(node, ref), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	resp, err := sc.do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform request")
	}

	var (
		tvr = TextValueResponse{}
	)
	if err := decodeJSONFromResponse(&tvr, resp); err != nil {
		return nil, errors.Wrap(err, "failed to decode")
	}

	return tvr.Data, nil
}

type TextValueResponse struct {
	Data []TextValue
}

func (tvr *TextValueResponse) UnmarshalJSON(b []byte) error {
	tvr.Data = []TextValue{}
	var values = [][]interface{}{}

	if err := json.Unmarshal(b, &values); err != nil {
		return errors.Wrap(err, "failed to deserialize nnt average response")
	}

	for _, entry := range values {
		var tv = TextValue{}
		tv.Value = entry[1].(string)
		// grab the timestamp
		if v, ok := entry[0].(float64); ok {
			tv.Time = time.Unix(int64(v), 0)
		}
		tvr.Data = append(tvr.Data, tv)
	}
	return nil
}

type TextValue struct {
	Time  time.Time
	Value string
}

// TextData - representation of Text Data for data submission and retrieval
type TextData struct {
	Metric string `json:"metric"`
	ID     string `json:"id"`
	Offset string `json:"offset"`
	Value  string `json:"value"`
}
