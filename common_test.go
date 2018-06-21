package gosnowth

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type noOpReadCloser struct {
	*bytes.Buffer
	WasClosed bool
}

func (norc *noOpReadCloser) Close() error {
	norc.WasClosed = true
	return nil
}

func TestResolveURL(t *testing.T) {
	base, _ := url.Parse("http://localhost:1234")
	result := resolveURL(base, "/a/resource/path")
	assert.Equal(t,
		"http://localhost:1234/a/resource/path", result, "should equal")
}

func TestMultiError(t *testing.T) {
	merr := newMultiError()
	assert.True(t, !merr.HasError(), "should have no errors yet")
	merr.Add(errors.New("error 1"))
	merr.Add(errors.New("error 2"))
	merr.Add(nil)

	assert.True(t, merr.HasError(), "should have errors")
	assert.Equal(t, "error 1; error 2", merr.Error(), "errors should be joined")
}

func TestMoveNode(t *testing.T) {
	urlA, _ := url.Parse("http://localhost:1")
	urlB, _ := url.Parse("http://localhost:2")

	from := []*SnowthNode{
		&SnowthNode{
			identifier: "a",
			url:        urlA,
		},
		&SnowthNode{
			identifier: "b",
			url:        urlB,
		},
	}
	to := []*SnowthNode{}

	moveNode(&from, &to, &SnowthNode{
		identifier: "a",
		url:        urlA,
	})

	assert.True(t, len(to) == 1, "length of to should be 1")
	assert.True(t, len(from) == 1, "length of from should be 1")
}

func TestDecodeJSONFromResponse(t *testing.T) {
	resp := &http.Response{
		Body: &noOpReadCloser{
			bytes.NewBufferString(`{
				"something": 1,
				"something_else": 2
			}`),
			false},
	}

	decoded := make(map[string]int)
	err := decodeJSONFromResponse(&decoded, resp.Body)
	if err != nil {
		t.Error("error encountered from decode function: ", err)
	}
	assert.Equal(t, 1, decoded["something"], "something should be 1")
	assert.Equal(t, 2, decoded["something_else"], "something_else should be 1")
}

func TestDecodeXMLFromResponse(t *testing.T) {
	resp := &http.Response{
		Body: &noOpReadCloser{
			bytes.NewBufferString(`<data><something>1</something><somethingelse>2</somethingelse></data>`),
			false},
	}
	type data struct {
		XMLName       xml.Name `xml:"data"`
		Something     int      `xml:"something"`
		SomethingElse int      `xml:"somethingelse"`
	}
	decoded := &data{}

	err := decodeXMLFromResponse(decoded, resp.Body)
	if err != nil {
		t.Error("error encountered from decode function: ", err)
	}
	assert.Equal(t, 1, decoded.Something, "something should be 1")
	assert.Equal(t, 2, decoded.SomethingElse, "something_else should be 1")
}

func TestEncodeXML(t *testing.T) {
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
	b, _ := ioutil.ReadAll(reader)

	assert.True(t, strings.Contains(string(b), "somethingelse"), "should contain somethingelse")
}
