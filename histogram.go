package gosnowth

import (
	"bytes"
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
	data ...HistogramData) (err error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(data); err != nil {
		return errors.Wrap(err, "failed to encode HistogramData for write")
	}

	err = sc.do(node, "POST", "/histogram/write", buf, nil, nil)
	return
}
