package gosnowth

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
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
	if ctx == nil {
		ctx = context.Background()
	}

	r, err := http.NewRequest("POST", sc.getURL(node, "/raw"), data)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	r.Close = true
	r.Header.Add("X-Snowth-Datapoints", strconv.FormatUint(dataPoints, 10))
	if fb { // is flatbuffer?
		r.Header.Add("Content-Type", FlatbufferContentType)
	}

	r = r.WithContext(ctx)
	sc.RLock()
	rf := sc.request
	sc.RUnlock()
	if rf != nil {
		if err := rf(r); err != nil {
			return errors.Wrap(err, "unable to process request")
		}

		if r == nil {
			return errors.New("invalid request after processing")
		}
	}

	sc.LogDebugf("snowth Request: %+v", r)
	var start = time.Now()
	resp, err := sc.c.Do(r)
	if err != nil {
		return errors.Wrap(err, "failed to perform request")
	}

	defer resp.Body.Close()
	sc.LogDebugf("snowth response: %+v", resp)
	sc.LogDebugf("snowth latency: %+v", time.Since(start))
	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "context terminated")
	default:
		break
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		sc.LogWarnf("error returned from IRONdb: [%d] %s",
			resp.StatusCode, string(body))
		return fmt.Errorf("error returned from IRONdb: [%d] %s",
			resp.StatusCode, string(body))
	}

	return nil
}
