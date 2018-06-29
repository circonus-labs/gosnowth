package gosnowth

import (
	"fmt"
	"net/url"
)

type FindTagsItem struct {
	UUID       string
	CheckName  string `json:"check_name"`
	MetricName string `json:"metric_name"`
	Category   string
	Type       string
	AccountID  int32 `json:"account_id"`
}

// FindTags - Find metrics that are associated with tags
func (sc *SnowthClient) FindTags(node *SnowthNode, accountID int32, query string, start, end string) ([]FindTagsItem, error) {
	url := fmt.Sprintf("%s?query=%s&activity_start_secs=%s&activity_end_secs=%s",
		sc.getURL(node, fmt.Sprintf("/find/%d/tags", accountID)),
		url.QueryEscape(query), url.QueryEscape(start), url.QueryEscape(end),
	)
	var (
		r   = []FindTagsItem{}
		err = sc.do(node, "GET", url, nil, &r, decodeJSONFromResponse)
	)
	return r, err
}
