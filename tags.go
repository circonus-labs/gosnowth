package gosnowth

import (
	"fmt"
	"net/url"
)

// FindTagsItem values represent results returned from IRONdb tag queries.
type FindTagsItem struct {
	UUID       string
	CheckName  string `json:"check_name"`
	MetricName string `json:"metric_name"`
	Category   string
	Type       string
	AccountID  int32 `json:"account_id"`
}

// FindTags retrieves metrics that are associated with the provided tag query.
func (sc *SnowthClient) FindTags(node *SnowthNode, accountID int32,
	query string, start, end string) ([]FindTagsItem, error) {
	u := fmt.Sprintf("%s?query=%s",
		sc.getURL(node, fmt.Sprintf("/find/%d/tags", accountID)),
		url.QueryEscape(query))
	if start != "" && end != "" {
		u += fmt.Sprintf("&activity_start_secs=%s&activity_end_secs=%s",
			url.QueryEscape(start), url.QueryEscape(end))
	}

	r := []FindTagsItem{}
	err := sc.do(node, "GET", u, nil, &r, decodeJSONFromResponse)
	return r, err
}
