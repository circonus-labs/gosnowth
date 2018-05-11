package gosnowth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type FindTagsResponse []FindTagsItem

type FindTagsItem struct {
	UUID       string
	CheckName  string `json:"check_name"`
	MetricName string `json:"metric_name"`
	Category   string
	Type       string
	AccountID  int32 `json:"account_id"`
}

// FindTags - Find metrics that are associated with tags
func (sc *SnowthClient) FindTags(node *SnowthNode, accountID int32, query string) (FindTagsResponse, error) {
	url := fmt.Sprintf("%s?query=%s",
		sc.getURL(node, fmt.Sprintf("/find/%d/tags", accountID)),
		url.QueryEscape(query),
	)

	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	resp, err := sc.c.Do(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform request")
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-success status code returned: %s -> %s",
			resp.Status, string(body))
	}

	var ftr = FindTagsResponse{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&ftr); err != nil {
		return nil, errors.Wrap(err, "failed to decode json response body")
	}

	return ftr, nil
}
