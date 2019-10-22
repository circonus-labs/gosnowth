package gosnowth

import (
	"context"
	"io"
	"strconv"
)

// FlatbufferContentType is the content type header for flatbuffer data.
const FlatbufferContentType = "application/x-circonus-metric-list-flatbuffer"

// WriteRaw writes raw IRONdb data to a node.
func (sc *SnowthClient) WriteRaw(node *SnowthNode, data io.Reader,
	fb bool, dataPoints uint64) error {
	return sc.WriteRawContext(context.Background(), node, data, fb, dataPoints)
}

// WriteRawContext is the context aware version of WriteRaw.
func (sc *SnowthClient) WriteRawContext(ctx context.Context, node *SnowthNode,
	data io.Reader, fb bool, dataPoints uint64) error {

	hdrs := map[string]string{"X-Snowth-Datapoints": strconv.FormatUint(dataPoints, 10)}
	if fb { // is flatbuffer?
		hdrs["Content-Type"] = FlatbufferContentType
	}

	_, _, err := sc.do(ctx, node, "POST", "/raw", data, hdrs)
	if err != nil {
		return err
	}

	return nil
}
