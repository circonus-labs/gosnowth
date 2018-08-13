package gosnowth

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
	"time"
)

type mockHTTPClient struct {
	mockDo func(req *http.Request) (*http.Response, error)
}

func (mhc *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if mhc.mockDo != nil {
		mhc.mockDo(req)
	}
	return &http.Response{
		StatusCode: 200,
		Body: nooprc{
			buf: bytes.NewBufferString("[[1,1]]"),
		},
	}, nil
}

type nooprc struct {
	buf *bytes.Buffer
}

func (nrc nooprc) Read(p []byte) (n int, err error) {
	return nrc.buf.Read(p)
}

func (nrc nooprc) Close() error {
	return nil
}

func TestGetRollup(t *testing.T) {
	body := &nooprc{
		buf: bytes.NewBufferString("[[1, 1]]"),
	}

	client := &SnowthClient{
		c: &mockHTTPClient{
			mockDo: func(req *http.Request) (*http.Response, error) {
				t.Logf("%+v", req.URL)
				return &http.Response{StatusCode: 200, Body: body}, nil
			},
		},
	}

	nodeURL, _ := url.Parse("http://localhost/")

	_, err := client.ReadRollupValues(
		&SnowthNode{
			url: nodeURL,
		}, "123", "metric.name",
		[]string{"abc:123", "foo:bar"}, time.Second,
		time.Now(), time.Now())
	if err != nil {
		t.Error(err.Error())
	}
}
