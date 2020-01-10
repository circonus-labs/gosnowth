package gosnowth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// FindTagsItem values represent results returned from IRONdb tag queries.
type FindTagsItem struct {
	UUID       string          `json:"uuid"`
	CheckName  string          `json:"check_name"`
	CheckTags  []string        `json:"check_tags,omitempty"`
	MetricName string          `json:"metric_name"`
	Category   string          `json:"category"`
	Type       string          `type:"type"`
	AccountID  int64           `json:"account_id"`
	Activity   [][]int64       `json:"activity,omitempty"`
	Latest     *FindTagsLatest `json:"latest,omitempty"`
}

// FindTagsResult values contain the results of a find tags request.
type FindTagsResult struct {
	Items []FindTagsItem
	Count int64
}

// FindTagsOptions values contain optional parameters to be passed to the
// IRONdb find tags call by a FindTags operation.
type FindTagsOptions struct {
	Start     time.Time `json:"activity_start_secs"`
	End       time.Time `json:"activity_end_secs"`
	Activity  int64     `json:"activity"`
	Latest    int64     `json:"latest"`
	CountOnly int64     `json:"count_only"`
}

// FindTagsLatest values contain the most recent data values for a metric.
type FindTagsLatest struct {
	Numeric   []FindTagsLatestNumeric   `json:"numeric,omitempty"`
	Text      []FindTagsLatestText      `json:"text,omitempty"`
	Histogram []FindTagsLatestHistogram `json:"histogram,omitempty"`
}

// FindTagsLatestNumeric values contain recent metric numeric data.
type FindTagsLatestNumeric struct {
	Time  int64
	Value *float64
}

// MarshalJSON encodes a FindTagsLatestNumeric value into a JSON format byte
// slice.
func (ftl *FindTagsLatestNumeric) MarshalJSON() ([]byte, error) {
	v := []interface{}{ftl.Time, ftl.Value}
	return json.Marshal(v)
}

// UnmarshalJSON decodes a JSON format byte slice into a FindTagsLatestNumeric
// value.
func (ftl *FindTagsLatestNumeric) UnmarshalJSON(b []byte) error {
	v := []interface{}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if len(v) != 2 {
		return fmt.Errorf("unable to decode latest numeric value, "+
			"invalid length: %v: %v", string(b), err)
	}

	if fv, ok := v[0].(float64); ok {
		ftl.Time = int64(fv)
	} else {
		return fmt.Errorf("unable to decode latest numeric value, "+
			"invalid timestamp: %v", string(b))
	}

	if v[1] != nil {
		if fv, ok := v[1].(float64); ok {
			ftl.Value = &fv
		} else {
			return fmt.Errorf("unable to decode latest numeric value, "+
				"invalid value: %v", string(b))
		}
	}

	return nil
}

// FindTagsLatestText values contain recent metric text data.
type FindTagsLatestText struct {
	Time  int64
	Value *string
}

// MarshalJSON encodes a FindTagsLatestText value into a JSON format byte slice.
func (ftl *FindTagsLatestText) MarshalJSON() ([]byte, error) {
	v := []interface{}{ftl.Time, ftl.Value}
	return json.Marshal(v)
}

// UnmarshalJSON decodes a JSON format byte slice into a FindTagsLatestText
// value.
func (ftl *FindTagsLatestText) UnmarshalJSON(b []byte) error {
	v := []interface{}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if len(v) != 2 {
		return fmt.Errorf("unable to decode latest text value, "+
			"invalid length: %v: %v", string(b), err)
	}

	if fv, ok := v[0].(float64); ok {
		ftl.Time = int64(fv)
	} else {
		return fmt.Errorf("unable to decode latest text value, "+
			"invalid timestamp: %v", string(b))
	}

	if v[1] != nil {
		if sv, ok := v[1].(string); ok {
			ftl.Value = &sv
		} else {
			return fmt.Errorf("unable to decode latest text value, "+
				"invalid value: %v", string(b))
		}
	}

	return nil
}

// FindTagsLatestHistogram values contain recent metric histogram data.
type FindTagsLatestHistogram struct {
	Time  int64
	Value *string
}

// MarshalJSON encodes a FindTagsLatestHistogram value into a JSON format byte
// slice.
func (ftl *FindTagsLatestHistogram) MarshalJSON() ([]byte, error) {
	v := []interface{}{ftl.Time, ftl.Value}
	return json.Marshal(v)
}

// UnmarshalJSON decodes a JSON format byte slice into a
// FindTagsLatestHistogram value.
func (ftl *FindTagsLatestHistogram) UnmarshalJSON(b []byte) error {
	v := []interface{}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if len(v) != 2 {
		return fmt.Errorf("unable to decode latest histogram value, "+
			"invalid length: %v: %v", string(b), err)
	}

	if fv, ok := v[0].(float64); ok {
		ftl.Time = int64(fv)
	} else {
		return fmt.Errorf("unable to decode latest histogram value, "+
			"invalid timestamp: %v", string(b))
	}

	if v[1] != nil {
		if sv, ok := v[1].(string); ok {
			ftl.Value = &sv
		} else {
			return fmt.Errorf("unable to decode latest histogram value, "+
				"invalid value: %v", string(b))
		}
	}

	return nil
}

// FindTags retrieves metrics that are associated with the provided tag query.
func (sc *SnowthClient) FindTags(node *SnowthNode, accountID int64,
	query string, options *FindTagsOptions) (*FindTagsResult, error) {
	return sc.FindTagsContext(context.Background(), node, accountID, query,
		options)
}

// FindTagsContext is the context aware version of FindTags.
func (sc *SnowthClient) FindTagsContext(ctx context.Context, node *SnowthNode,
	accountID int64, query string,
	options *FindTagsOptions) (*FindTagsResult, error) {
	u := fmt.Sprintf("%s?query=%s",
		sc.getURL(node, fmt.Sprintf("/find/%d/tags", accountID)),
		url.QueryEscape(query))
	if !options.Start.IsZero() && !options.End.IsZero() &&
		options.Start.Unix() != 0 && options.End.Unix() != 0 {
		u += fmt.Sprintf("&activity_start_secs=%s&activity_end_secs=%s",
			formatTimestamp(options.Start), formatTimestamp(options.End))
	}

	u += fmt.Sprintf("&activity=%d", options.Activity)
	u += fmt.Sprintf("&latest=%d", options.Latest)
	if options.CountOnly != 0 {
		u += fmt.Sprintf("&count_only=%d", options.CountOnly)
	}

	r := &FindTagsResult{}
	body, header, err := sc.do(ctx, node, "GET", u, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := decodeJSON(body, &r.Items); err != nil {
		return nil, errors.Wrap(err, "unable to decode IRONdb response")
	}

	// Return a results count and capture it from the header , if provided.
	r.Count = int64(len(r.Items))
	if header != nil {
		c := header.Get("X-Snowth-Search-Result-Count")
		if c != "" {
			if cv, err := strconv.ParseInt(c, 10, 64); err == nil {
				r.Count = cv
			}
		}
	}

	return r, err
}
