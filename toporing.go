package gosnowth

import (
	"encoding/xml"
	"path"
)

// TopoRing values represent IRONdb topology ring data.
type TopoRing struct {
	XMLName      xml.Name         `xml:"vnodes" json:"-"`
	VirtualNodes []TopoRingDetail `xml:"vnode"`
	NumberNodes  int              `xml:"n,attr" json:"-"`
}

// TopoRingDetail values represent IRONdb topology ring node details.
type TopoRingDetail struct {
	XMLName  xml.Name `xml:"vnode" json:"-"`
	ID       string   `xml:"id,attr" json:"id"`
	IDX      int      `xml:"idx,attr" json:"idx"`
	Location float64  `xml:"location,attr" json:"location"`
}

// GetTopoRingInfo retrieves topology ring information from a node.
func (sc *SnowthClient) GetTopoRingInfo(hash string,
	node *SnowthNode) (*TopoRing, error) {
	tr := new(TopoRing)
	err := sc.do(node, "GET", path.Join("/toporing/xml", hash),
		nil, tr, decodeXML)
	return tr, err
}
