package gosnowth

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// FetchStream values represent queries for individual data streams in an
// IRONdb fetch request.
type FetchStream struct {
	UUID            string   `json:"uuid"`
	Name            string   `json:"name"`
	Kind            string   `json:"kind"`
	Label           string   `json:"label,omitempty"`
	Transform       string   `json:"transform"`
	TransformParams []string `json:"transform_params,omitempty"`
}

// FetchReduce values represent reduce operations to perform on specified
// data streams in an IRONdb fetch request.
type FetchReduce struct {
	Label        string   `json:"label"`
	Method       string   `json:"method"`
	MethodParams []string `json:"method_params,omitempty"`
}

// FetchQuery values represent queries used to fetch IRONdb data.
type FetchQuery struct {
	Start   time.Time     `json:"start"`
	Period  time.Duration `json:"period"`
	Count   int64         `json:"count"`
	Streams []FetchStream `json:"streams"`
	Reduce  []FetchReduce `json:"reduce"`
}

// MarshalJSON encodes a FetchQuery value into a JSON format byte slice.
func (fq *FetchQuery) MarshalJSON() ([]byte, error) {
	v := struct {
		Start   float64       `json:"start"`
		Period  float64       `json:"period"`
		Count   int64         `json:"count"`
		Streams []FetchStream `json:"streams"`
		Reduce  []FetchReduce `json:"reduce"`
	}{}
	fv, err := strconv.ParseFloat(formatTimestamp(fq.Start), 64)
	if err != nil {
		return nil, errors.New("invalid fetch start value: " +
			formatTimestamp(fq.Start))
	}

	v.Start = fv
	v.Period = fq.Period.Seconds()
	v.Count = fq.Count
	if len(fq.Streams) > 0 {
		v.Streams = fq.Streams
	}

	if len(fq.Reduce) > 0 {
		v.Reduce = fq.Reduce
	}

	return json.Marshal(v)
}

// UnmarshalJSON decodes a JSON format byte slice into a HistogramValue value.
func (fq *FetchQuery) UnmarshalJSON(b []byte) error {
	v := struct {
		Start   float64       `json:"start"`
		Period  float64       `json:"period"`
		Count   int64         `json:"count"`
		Streams []FetchStream `json:"streams"`
		Reduce  []FetchReduce `json:"reduce"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if v.Start == 0 {
		return errors.New("fetch query missing start: " + string(b))
	}

	fq.Start, err = parseTimestamp(strconv.FormatFloat(v.Start, 'f', 3, 64))
	if err != nil {
		return err
	}

	if v.Period == 0 {
		return errors.New("fetch query missing period: " + string(b))
	}

	fq.Period = time.Duration(v.Period*1000) * time.Millisecond
	if v.Count == 0 {
		return errors.New("fetch query missing count: " + string(b))
	}

	fq.Count = v.Count
	if len(v.Streams) < 1 {
		return errors.New("fetch query requires at least one stream: " +
			string(b))
	}

	fq.Streams = v.Streams
	if len(v.Reduce) < 1 {
		return errors.New("fetch query requires at least one reduce: " +
			string(b))
	}

	fq.Reduce = v.Reduce
	return nil
}

// Timestamp returns the FetchQuery start time as a string in the IRONdb
// timestamp format.
func (fq *FetchQuery) Timestamp() string {
	return formatTimestamp(fq.Start)
}

// FetchValues retrieves data values using the IRONdb fetch API.
func (sc *SnowthClient) FetchValues(node *SnowthNode,
	q *FetchQuery) (*DF4Response, error) {
	return sc.FetchValuesContext(context.Background(), node, q)
}

// FetchValuesContext is the context aware version of FetchValuesValues.
func (sc *SnowthClient) FetchValuesContext(ctx context.Context,
	node *SnowthNode, q *FetchQuery) (*DF4Response, error) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(&q); err != nil {
		return nil, err
	}

	body, _, err := sc.do(ctx, node, "POST", "/fetch", buf)
	if err != nil {
		return nil, err
	}

	r := &DF4Response{}
	if err := decodeJSON(body, &r); err != nil {
		return nil, errors.Wrap(err, "unable to decode IRONdb response")
	}

	return r, nil
}
