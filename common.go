package gosnowth

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// resolveURL resolves the address of a URL plus a string reference.
func resolveURL(baseURL *url.URL, ref string) string {
	refURL, _ := url.Parse(ref)
	return baseURL.ResolveReference(refURL).String()
}

// multiError values keep track of multiple errors.
type multiError struct {
	errs []error
}

// newMultiError initializes a new multiError value.
func newMultiError() *multiError {
	return &multiError{
		errs: []error{},
	}
}

// Add appends an error to the list of errors.
func (me *multiError) Add(err error) {
	if err != nil {
		me.errs = append(me.errs, err)
	}
}

// HasError returns whether the value contains any errors.
func (me multiError) HasError() bool {
	return len(me.errs) > 0
}

// String returns a string representation of the error value.
func (me multiError) String() string {
	es := []string{}
	for _, err := range me.errs {
		es = append(es, err.Error())
	}

	return strings.Join(es, "; ")
}

// Error implements the error interface for multiError values.
func (me multiError) Error() string {
	return me.String()
}

// decodeJSON decodes JSON from a reader into an interface.
func decodeJSON(r io.Reader, v interface{}) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return errors.Wrap(err, "failed to decode JSON")
	}

	return nil
}

// encodeXML create a reader of XML data representing an interface.
func encodeXML(v interface{}) (io.Reader, error) {
	buf := bytes.NewBuffer([]byte{})
	if err := xml.NewEncoder(buf).Encode(v); err != nil {
		return nil, errors.Wrap(err, "failed to encode XML")
	}

	return buf, nil
}

// decodeXML decodes XML from a reader into an interface.
func decodeXML(r io.Reader, v interface{}) error {
	if err := xml.NewDecoder(r).Decode(v); err != nil {
		return errors.Wrap(err, "failed to decode XML")
	}

	return nil
}
