package gosnowth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

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

// FindTagsOption values contain optional parameters to be passed to the
// IRONdb find tags call by a FindTags operation.
type FindTagsOption struct {
	Name  string `json:"name"`
	Value string `json:"value"`
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

	if fv, ok := v[1].(float64); ok {
		ftl.Value = &fv
	} else {
		return fmt.Errorf("unable to decode latest numeric value, "+
			"invalid value: %v", string(b))
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

	if sv, ok := v[1].(string); ok {
		ftl.Value = &sv
	} else {
		return fmt.Errorf("unable to decode latest text value, "+
			"invalid value: %v", string(b))
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

	if sv, ok := v[1].(string); ok {
		ftl.Value = &sv
	} else {
		return fmt.Errorf("unable to decode latest histogram value, "+
			"invalid value: %v", string(b))
	}

	return nil
}

// NewFindTagsOption creates a new find tags option with any name and value.
func NewFindTagsOption(name, value string) *FindTagsOption {
	return &FindTagsOption{
		Name:  name,
		Value: value,
	}
}

// FindTagsOptionStart creates a new find tags option for activity start.
func FindTagsOptionStart(value string) *FindTagsOption {
	return &FindTagsOption{
		Name:  "activity_start_secs",
		Value: value,
	}
}

// FindTagsOptionEnd creates a new find tags option for activity end.
func FindTagsOptionEnd(value string) *FindTagsOption {
	return &FindTagsOption{
		Name:  "activity_end_secs",
		Value: value,
	}
}

// FindTagsOptionActivity creates a new find tags option for retrieving
// activity window data.
func FindTagsOptionActivity(value string) *FindTagsOption {
	return &FindTagsOption{
		Name:  "activity",
		Value: value,
	}
}

// FindTagsOptionLatest creates a new find tags option for retrieving latest
// metric values.
func FindTagsOptionLatest(value string) *FindTagsOption {
	return &FindTagsOption{
		Name:  "latest",
		Value: value,
	}
}

// FindTagsOptionCountOnly creates a new find tags option for retrieving only
// the count of matching metrics.
func FindTagsOptionCountOnly(value string) *FindTagsOption {
	return &FindTagsOption{
		Name:  "count_only",
		Value: value,
	}
}

// FindTags retrieves metrics that are associated with the provided tag query.
func (sc *SnowthClient) FindTags(node *SnowthNode, accountID int64,
	query string, opts ...*FindTagsOption) (*FindTagsResult, error) {
	return sc.FindTagsContext(context.Background(), node, accountID, query,
		opts...)
}

// FindTagsContext is the context aware version of FindTags.
func (sc *SnowthClient) FindTagsContext(ctx context.Context, node *SnowthNode,
	accountID int64, query string,
	opts ...*FindTagsOption) (*FindTagsResult, error) {
	u := fmt.Sprintf("%s?query=%s",
		sc.getURL(node, fmt.Sprintf("/find/%d/tags", accountID)),
		url.QueryEscape(query))
	start, end := "", ""
	for _, opt := range opts {
		if opt != nil {
			switch opt.Name {
			case "start", "activity_start_secs":
				start = opt.Value
			case "end", "activity_end_secs":
				end = opt.Value
			default:
				u += fmt.Sprintf("&%s=%s", url.QueryEscape(opt.Name),
					url.QueryEscape(opt.Value))
			}
		}
	}

	if start != "" && end != "" {
		u += fmt.Sprintf("&activity_start_secs=%s&activity_end_secs=%s",
			url.QueryEscape(start), url.QueryEscape(end))
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
