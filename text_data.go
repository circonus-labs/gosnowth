package gosnowth

import (
	"bytes"
	"encoding/json"
	"path"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// WriteText writes text data to an IRONdb node.
func (sc *SnowthClient) WriteText(node *SnowthNode, data ...TextData) (err error) {
	var (
		buf = new(bytes.Buffer)
		enc = json.NewEncoder(buf)
	)
	if err := enc.Encode(data); err != nil {
		return errors.Wrap(err, "failed to encode TextData for write")
	}

	err = sc.do(node, "POST", "/write/text", buf, nil, nil)
	return
}

// ReadTextValues reads text data values from an IRONdb node.
func (sc *SnowthClient) ReadTextValues(node *SnowthNode, start, end time.Time,
	id, metric string) ([]TextValue, error) {
	var (
		tvr = new(TextValueResponse)
		err = sc.do(node, "GET", path.Join("/read",
			strconv.FormatInt(start.Unix(), 10),
			strconv.FormatInt(end.Unix(), 10),
			id, metric), nil, tvr, decodeJSONFromResponse)
	)

	return tvr.Data, err
}

// TextValueResponse values represent text data responses.
type TextValueResponse struct {
	Data []TextValue
}

// UnmarshalJSON decodes a JSON format byte slice into a TextValueResponse.
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

// TextValue values represent text data read from IRONdb.
type TextValue struct {
	Time  time.Time
	Value string
}

// TextData values represent text data to be written to IRONdb.
type TextData struct {
	Metric string `json:"metric"`
	ID     string `json:"id"`
	Offset string `json:"offset"`
	Value  string `json:"value"`
}
