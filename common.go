package gosnowth

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
}

func resolveURL(baseURL *url.URL, ref string) string {
	refURL, _ := url.Parse(ref)
	return baseURL.ResolveReference(refURL).String()
}

type multiError struct {
	errs []error
}

func newMultiError() *multiError {
	return &multiError{
		errs: []error{},
	}
}

func (me *multiError) Add(err error) {
	if err != nil {
		me.errs = append(me.errs, err)
	}
}

func (me *multiError) HasError() bool {
	if len(me.errs) > 0 {
		return true
	}
	return false
}

func (me *multiError) Error() string {
	var errStrs []string
	for _, err := range me.errs {
		errStrs = append(errStrs, err.Error())
	}
	return strings.Join(errStrs, "; ")
}

// moveNode - move a url from a slice to a new slice, if this is used for
// SnowthInstances' active or inactive slices wrap in a write lock
func moveNode(from, dest *[]*SnowthNode, u *SnowthNode) {
	// put this url in active
	*dest = append(*dest, u)

	// find the item index in the deactive list
	var index = -1
	for i, v := range *from {
		if v.url.String() == u.url.String() {
			index = i
		}
	}
	if index != -1 {
		// remove from deactive
		*from = removeNode(*from, index)
	}
}

// removeNode - remove a url from a slice, if this is used for
// SnowthInstances' active or inactive slices wrap in a write lock
func removeNode(a []*SnowthNode, index int) []*SnowthNode {
	copy(a[index:], a[index+1:])
	a[len(a)-1] = nil // or the zero value of T
	a = a[:len(a)-1]
	return a
}

func decodeJSONFromResponse(v interface{}, resp *http.Response) error {
	defer closeBody(resp)
	dec := json.NewDecoder(resp.Body)

	if err := dec.Decode(v); err != nil {
		return errors.Wrap(err, "failed to decode response body")
	}
	return nil
}

func encodeXML(v interface{}) (io.Reader, error) {
	buf := bytes.NewBuffer([]byte{})
	dec := xml.NewEncoder(buf)

	if err := dec.Encode(v); err != nil {
		return nil, errors.Wrap(err, "failed to encode")
	}
	return buf, nil
}

func decodeXMLFromResponse(v interface{}, resp *http.Response) error {
	defer closeBody(resp)
	dec := xml.NewDecoder(resp.Body)

	if err := dec.Decode(v); err != nil {
		return errors.Wrap(err, "failed to decode response body")
	}
	return nil
}
