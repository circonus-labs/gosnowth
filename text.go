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

// TextValueResponse values represent text data responses.
type TextValueResponse []TextValue

// UnmarshalJSON decodes a JSON format byte slice into a TextValueResponse.
func (tvr *TextValueResponse) UnmarshalJSON(b []byte) error {
	*tvr = TextValueResponse{}
	values := [][]interface{}{}
	if err := json.Unmarshal(b, &values); err != nil {
		return errors.Wrap(err, "failed to decode JSON response")
	}

	for _, entry := range values {
		var tv = TextValue{}
		tv.Value = entry[1].(string)
		if v, ok := entry[0].(float64); ok {
			tv.Time = time.Unix(int64(v), 0)
		}

		*tvr = append(*tvr, tv)
	}

	return nil
}

// TextValue values represent text data read from IRONdb.
type TextValue struct {
	Time  time.Time
	Value string
}

// ReadTextValues reads text data values from an IRONdb node.
func (sc *SnowthClient) ReadTextValues(node *SnowthNode, start, end time.Time,
	id, metric string) ([]TextValue, error) {
	return sc.ReadTextValuesContext(context.Background(), node, start, end,
		id, metric)
}

// ReadTextValuesContext is the context aware version of ReadTextValues.
func (sc *SnowthClient) ReadTextValuesContext(ctx context.Context,
	node *SnowthNode, start, end time.Time,
	id, metric string) ([]TextValue, error) {
	tvr := new(TextValueResponse)
	err := sc.do(ctx, node, "GET", path.Join("/read",
		strconv.FormatInt(start.Unix(), 10),
		strconv.FormatInt(end.Unix(), 10), id, metric), nil, tvr, decodeJSON)
	if tvr == nil {
		return nil, err
	}

	return *tvr, err
}

// TextData values represent text data to be written to IRONdb.
type TextData struct {
	Metric string `json:"metric"`
	ID     string `json:"id"`
	Offset string `json:"offset"`
	Value  string `json:"value"`
}

// WriteText writes text data to an IRONdb node.
func (sc *SnowthClient) WriteText(node *SnowthNode, data ...TextData) error {
	return sc.WriteTextContext(context.Background(), node, data...)
}

// WriteTextContext is the context aware version of WriteText.
func (sc *SnowthClient) WriteTextContext(ctx context.Context, node *SnowthNode,
	data ...TextData) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return errors.Wrap(err, "failed to encode TextData for write")
	}

	return sc.do(ctx, node, "POST", "/write/text", buf, nil, nil)
}
