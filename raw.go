package gosnowth

import (
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
	r, err := http.NewRequest("POST", sc.getURL(node, "/raw"), data)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	r.Header.Add("X-Snowth-Datapoints", strconv.FormatUint(dataPoints, 10))
	if fb { // is flatbuffer?
		r.Header.Add("Content-Type", FlatbufferContentType)
	}

	sc.LogDebugf("Snowth Request: %+v", r)
	var start = time.Now()
	resp, err := sc.c.Do(r)
	if err != nil {
		return errors.Wrap(err, "failed to perform request")
	}

	defer resp.Body.Close()
	sc.LogDebugf("Snowth Response: %+v", resp)
	sc.LogDebugf("Snowth Response Latency: %+v", time.Now().Sub(start))
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		sc.LogWarnf("status code not 200: %+v", resp)
		return fmt.Errorf("non-success status code returned: %s -> %s",
			resp.Status, string(body))
	}

	return nil
}
