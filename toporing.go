package gosnowth

import (
	"encoding/xml"
	"net/http"
	"path"

	"github.com/pkg/errors"
)

// GetTopoRingInfo - Get the toporing information from the node.
func (sc *SnowthClient) GetTopoRingInfo(hash string, node *SnowthNode) (*TopoRing, error) {
	var resource = path.Join("/toporing/xml", hash)
	req, err := http.NewRequest("GET", sc.getURL(node, resource), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	resp, err := sc.do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform request")
	}

	var toporing = new(TopoRing)
	if err := decodeXMLFromResponse(toporing, resp); err != nil {
		return nil, errors.Wrap(err, "failed to decode")
	}

	return toporing, nil
}

// TopoRing - structure for the response of the toporing api calls
type TopoRing struct {
	XMLName      xml.Name         `xml:"vnodes" json:"-"`
	VirtualNodes []TopoRingDetail `xml:"vnode"`
	NumberNodes  int              `xml:"n,attr" json:"-"`
}

// TopoRingDetail - detail node information from toporing api call
type TopoRingDetail struct {
	XMLName  xml.Name `xml:"vnode" json:"-"`
	ID       string   `xml:"id,attr" json:"id"`
	IDX      int      `xml:"idx,attr" json:"idx"`
	Location float64  `xml:"location,attr" json:"location"`
}
