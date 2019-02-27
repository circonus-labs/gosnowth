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
	exp := "http://localhost:1234/a/resource/path"
	if result != exp {
		t.Errorf("Expected result: %v, got: %v", exp, result)
	}
}

func TestMultiError(t *testing.T) {
	me := newMultiError()
	if me.HasError() {
		t.Error("Should have no errors yet")
	}

	me.Add(errors.New("error 1"))
	me.Add(errors.New("error 2"))
	me.Add(nil)
	if !me.HasError() {
		t.Error("Should have errors")
	}

	res := me.Error()
	exp := "error 1; error 2"
	if res != exp {
		t.Errorf("Expected result: %v, got: %v", exp, res)
	}
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

	if len(to) != 1 {
		t.Error("Length of to should be 1")
	}

	if len(from) != 1 {
		t.Error("Length of from should be 1")
	}
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

	if decoded["something"] != 1 {
		t.Error("something should be 1")
	}

	if decoded["something_else"] != 2 {
		t.Error("something_else should be 2")
	}
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

	if decoded.Something != 1 {
		t.Error("something should be 1")
	}

	if decoded.SomethingElse != 2 {
		t.Error("something else should be 2")
	}
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

	if !strings.Contains(string(b), "somethingelse") {
		t.Error("Should contain somethingelse")
	}
}
