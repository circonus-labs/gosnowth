package gosnowth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// FlatbufferContentType is the content type header for flatbuffer data.
const FlatbufferContentType = "application/x-circonus-metric-list-flatbuffer"

// RawNumericValueResponse values represent raw numeric data responses from IRONdb.
type RawNumericValueResponse struct {
	Data []RawNumericValue
}

// UnmarshalJSON decodes a JSON format byte slice into a RawNumericValueResponse.
func (rv *RawNumericValueResponse) UnmarshalJSON(b []byte) error {
	rv.Data = []RawNumericValue{}
	values := [][]interface{}{}
	if err := json.Unmarshal(b, &values); err != nil {
		return errors.Wrap(err, "failed to deserialize raw numeric response")
	}

	for _, entry := range values {
		var rnv = RawNumericValue{}
		if m, ok := entry[1].(float64); ok {
			rnv.Value = m
		}

		// grab the timestamp
		if v, ok := entry[0].(float64); ok {
			rnv.Time = time.Unix(int64(v/1000), 0)
		}

		rv.Data = append(rv.Data, rnv)
	}

	return nil
}

// RawNumericValue values represent raw numeric data.
type RawNumericValue struct {
	Time  time.Time
	Value float64
}

// ReadRawNumericValues reads raw numeric data from a node.
func (sc *SnowthClient) ReadRawNumericValues(node *SnowthNode,
	start time.Time, end time.Time, uuid string, metric string) ([]RawNumericValue, error) {
	return sc.ReadRawNumericValuesContext(context.Background(), node, start, end, uuid, metric)
}

// ReadRawNumericValuesContext is the context aware version of ReadRawNumericValues.
func (sc *SnowthClient) ReadRawNumericValuesContext(ctx context.Context,
	node *SnowthNode, start, end time.Time, uuid string, metric string) ([]RawNumericValue, error) {
	qp := url.Values{}
	qp.Add("start_ts", formatTimestamp(start))
	qp.Add("end_ts", formatTimestamp(end))

	r := &RawNumericValueResponse{}
	body, _, err := sc.do(ctx, node, "GET", path.Join("/raw",
		uuid, metric)+"?"+qp.Encode(), nil, nil)
	if err != nil {
		return nil, err
	}

	if err := decodeJSON(body, &r); err != nil {
		return nil, errors.Wrap(err, "unable to decode IRONdb response")
	}
	return r.Data, nil
}

// WriteRawResponse values represent raw IRONdb data write responses.
type WriteRawResponse struct {
	Errors      uint64 `json:"errors"`
	Misdirected uint64 `json:"misdirected"`
	Records     uint64 `json:"records"`
	Updated     uint64 `json:"updated"`
}

// WriteRaw writes raw IRONdb data to a node.
func (sc *SnowthClient) WriteRaw(node *SnowthNode, data io.Reader,
	fb bool, dataPoints uint64) (*WriteRawResponse, error) {
	return sc.WriteRawContext(context.Background(), node, data, fb, dataPoints)
}

// WriteRawContext is the context aware version of WriteRaw.
func (sc *SnowthClient) WriteRawContext(ctx context.Context, node *SnowthNode,
	data io.Reader, fb bool, dataPoints uint64) (*WriteRawResponse, error) {

	hdrs := http.Header{"X-Snowth-Datapoints": {strconv.FormatUint(dataPoints, 10)}}
	if fb { // is flatbuffer?
		hdrs["Content-Type"] = []string{FlatbufferContentType}
	}

	body, _, err := sc.do(ctx, node, "POST", "/raw", data, hdrs)
	if err != nil {
		return nil, err
	}

	r := &WriteRawResponse{}
	if err := decodeJSON(body, &r); err != nil {
		return nil, errors.Wrap(err, "unable to decode IRONdb response")
	}

	return r, nil
}
