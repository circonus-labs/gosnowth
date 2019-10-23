package gosnowth

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// FindTagsItem values represent results returned from IRONdb tag queries.
type FindTagsItem struct {
	UUID       string    `json:"uuid"`
	CheckName  string    `json:"check_name"`
	CheckTags  []string  `json:"check_tags,omitempty"`
	MetricName string    `json:"metric_name"`
	Category   string    `json:"category"`
	Type       string    `type:"type"`
	AccountID  int64     `json:"account_id"`
	Activity   [][]int64 `json:"activity,omitempty"`
}

// FindTagsResult values contain the results of a find tags request.
type FindTagsResult struct {
	Items []FindTagsItem
	Count int64
}

// FindTags retrieves metrics that are associated with the provided tag query.
func (sc *SnowthClient) FindTags(node *SnowthNode, accountID int64,
	query string, start, end string) (*FindTagsResult, error) {
	return sc.FindTagsContext(context.Background(), node, accountID, query,
		start, end)
}

// FindTagsContext is the context aware version of FindTags.
func (sc *SnowthClient) FindTagsContext(ctx context.Context, node *SnowthNode,
	accountID int64, query string, start, end string) (*FindTagsResult, error) {
	u := fmt.Sprintf("%s?query=%s",
		sc.getURL(node, fmt.Sprintf("/find/%d/tags", accountID)),
		url.QueryEscape(query))
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
