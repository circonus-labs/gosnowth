package gosnowth

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

// CAQLQuery values represent CAQL queries and associated parameters.
type CAQLQuery struct {
	Query   string
	Start   int64
	End     int64
	Period  int64
	Timeout int64
}

// GetCAQLQuery retrieves data values for metrics matching a CAQL format.
func (sc *SnowthClient) GetCAQLQuery(node *SnowthNode, accountID int64,
	q *CAQLQuery) (*DF4Response, error) {
	return sc.GetCAQLQueryContext(context.Background(), node, accountID, q)
}

// GetCAQLQueryContext is the context aware version of GetCAQLQuery.
func (sc *SnowthClient) GetCAQLQueryContext(ctx context.Context,
	node *SnowthNode, accountID int64, q *CAQLQuery) (*DF4Response, error) {
	u := sc.getURL(node, "/extension/lua/public/caql_v1") +
		"?query=" + url.PathEscape(q.Query)
	if q.Start != 0 {
		u += fmt.Sprintf("&start=%d", q.Start)
	}

	if q.End != 0 {
		u += fmt.Sprintf("&end=%d", q.End)
	}

	if q.Period != 0 {
		u += fmt.Sprintf("&period=%d", q.Period)
	}

	if q.Timeout != 0 {
		u += fmt.Sprintf("&_timeout=%d", q.Timeout)
	}

	u += fmt.Sprintf("&format=DF4&account_id=%d", accountID)
	r := &DF4Response{}
	body, _, err := sc.do(ctx, node, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	if err := decodeJSON(body, &r); err != nil {
		return nil, errors.Wrap(err, "unable to decode IRONdb response")
	}

	return r, err
}
