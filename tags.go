package gosnowth

import (
	"context"
	"fmt"
	"net/url"
)

// FindTagsItem values represent results returned from IRONdb tag queries.
type FindTagsItem struct {
	UUID       string    `json:"uuid"`
	CheckName  string    `json:"check_name"`
	CheckTags  []string  `json:"check_tags,omitempty"`
	MetricName string    `json:"metric_name"`
	Category   string    `json:"category"`
	Type       string    `type:"type"`
	AccountID  int32     `json:"account_id"`
	Activity   [][]int32 `json:"activity,omitempty"`
}

// FindTags retrieves metrics that are associated with the provided tag query.
func (sc *SnowthClient) FindTags(node *SnowthNode, accountID int32,
	query string, start, end string) ([]FindTagsItem, error) {
	return sc.FindTagsContext(context.Background(), node, accountID, query,
		start, end)
}

// FindTagsContext is the context aware version of FindTags.
func (sc *SnowthClient) FindTagsContext(ctx context.Context, node *SnowthNode,
	accountID int32, query string, start, end string) ([]FindTagsItem, error) {
	u := fmt.Sprintf("%s?query=%s",
		sc.getURL(node, fmt.Sprintf("/find/%d/tags", accountID)),
		url.QueryEscape(query))
	if start != "" && end != "" {
		u += fmt.Sprintf("&activity_start_secs=%s&activity_end_secs=%s",
			url.QueryEscape(start), url.QueryEscape(end))
	}

	r := []FindTagsItem{}
	err := sc.do(ctx, node, "GET", u, nil, &r, decodeJSON)
	return r, err
}
