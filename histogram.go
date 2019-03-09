package gosnowth

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/circonus-labs/circonusllhist"
	"github.com/pkg/errors"
)

// HistogramData values represent histogram data records in IRONdb.
type HistogramData struct {
	Metric    string                    `json:"metric"`
	ID        string                    `json:"id"`
	Offset    int64                     `json:"offset"`
	Period    int64                     `json:"period"`
	Histogram *circonusllhist.Histogram `json:"histogram"`
}

// WriteHistogram sends a variadic list of histogram data values to be written
// to an IRONdb node.
func (sc *SnowthClient) WriteHistogram(node *SnowthNode,
	data ...HistogramData) error {
	return sc.WriteHistogramContext(context.Background(), node, data...)
}

// WriteHistogramContext is the context aware version of WriteHistogram.
func (sc *SnowthClient) WriteHistogramContext(ctx context.Context,
	node *SnowthNode, data ...HistogramData) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return errors.Wrap(err, "failed to encode HistogramData for write")
	}

	return sc.do(ctx, node, "POST", "/histogram/write", buf, nil, nil)
}
