package gosnowth

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func float64Ptr(f float64) *float64 {
	return &f
}

func stringPtr(s string) *string {
	return &s
}

type noOpReadCloser struct {
	*bytes.Buffer
	WasClosed bool
}

func (n *noOpReadCloser) Close() error {
	n.WasClosed = true

	return nil
}

func TestResolveURL(t *testing.T) {
	t.Parallel()

	base, _ := url.Parse("http://localhost:1234")
	result := resolveURL(base, "/a/resource/path")

	exp := "http://localhost:1234/a/resource/path"
	if result != exp {
		t.Errorf("Expected result: %v, got: %v", exp, result)
	}
}

func TestDecodeJSON(t *testing.T) {
	t.Parallel()

	resp := &http.Response{
		Body: &noOpReadCloser{
			bytes.NewBufferString(`{
				"something": 1,
				"something_else": 2
			}`),
			false,
		},
	}

	v := make(map[string]int)

	if err := decodeJSON(resp.Body, &v); err != nil {
		t.Error("error encountered from decode function: ", err)
	}

	if v["something"] != 1 {
		t.Error("something should be 1")
	}

	if v["something_else"] != 2 {
		t.Error("something_else should be 2")
	}
}

func TestDecodeXML(t *testing.T) {
	t.Parallel()

	resp := &http.Response{
		Body: &noOpReadCloser{
			bytes.NewBufferString(`<data><something>1</something><somethingelse>2</somethingelse></data>`),
			false,
		},
	}

	type data struct {
		XMLName       xml.Name `xml:"data"`
		Something     int      `xml:"something"`
		SomethingElse int      `xml:"somethingelse"`
	}

	decoded := &data{}

	if err := decodeXML(resp.Body, decoded); err != nil {
		t.Error("error encountered from decode function: ", err)
	}

	if decoded.Something != 1 {
		t.Error("something should be 1")
	}

	if decoded.SomethingElse != 2 {
		t.Error("something else should be 2")
	}
}

func TestEncodeXML(t *testing.T) {
	t.Parallel()

	type data struct {
		XMLName       xml.Name `xml:"data"`
		Something     int      `xml:"something"`
		SomethingElse int      `xml:"somethingelse"`
	}

	d := &data{
		Something:     1,
		SomethingElse: 2,
	}

	reader, err := encodeXML(d)
	if err != nil {
		t.Error("error encountered encoding: ", err)
	}

	b, _ := io.ReadAll(reader)
	if !strings.Contains(string(b), "somethingelse") {
		t.Error("Should contain somethingelse")
	}
}

func TestFormatTimestamp(t *testing.T) {
	t.Parallel()

	tm := time.Unix(123456789, int64(time.Millisecond))
	exp := "123456789.001"
	res := formatTimestamp(tm)

	if res != exp {
		t.Errorf("Expected string: %v, got: %v", exp, res)
	}

	tm = time.Unix(123456789, 0)
	exp = "123456789"
	res = formatTimestamp(tm)

	if res != exp {
		t.Errorf("Expected string: %v, got: %v", exp, res)
	}
}

func TestParseTimestamp(t *testing.T) {
	t.Parallel()

	res, err := parseTimestamp("123456789.001")
	if err != nil {
		t.Fatal(err)
	}

	exp := time.Unix(123456789, int64(time.Millisecond))
	if !res.Equal(exp) {
		t.Errorf("Expected time: %v, got: %v", exp, res)
	}
}

func TestParseDuration(t *testing.T) {
	t.Parallel()

	res, err := parseDuration("1")
	if err != nil {
		t.Fatal(err)
	}

	if exp := time.Second; res != exp {
		t.Errorf("Expected duration: %v, got: %v", exp, res)
	}
}
