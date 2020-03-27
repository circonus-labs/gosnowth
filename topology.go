package gosnowth

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"path"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

type TopologyNodeSlot struct {
	Location [sha256.Size]byte
	Node     *TopologyNode
	Idx      uint16
}

// Topology values represent IRONdb topology structure.
type Topology struct {
	XMLName        xml.Name       `xml:"nodes" json:"-"`
	OldWriteCopies uint8          `xml:"n,attr" json:"-"`
	WriteCopies    uint8          `xml:"write_copies,attr" json:"-"`
	Hash           string         `xml:"-"`
	Nodes          []TopologyNode `xml:"node"`
	use_side       bool
	ring           []TopologyNodeSlot
}

func (topo *Topology) Len() int { return len(topo.ring) }
func (topo *Topology) Swap(i, j int) {
	topo.ring[i], topo.ring[j] = topo.ring[j], topo.ring[i]
}
func (topo *Topology) Less(i, j int) bool {
	return bytes.Compare(topo.ring[i].Location[:], topo.ring[j].Location[:]) < 0
}

type TopoSide uint8

func (i *TopoSide) UnmarshalXMLAttr(attr xml.Attr) error {
	switch strings.ToLower(attr.Value) {
	default:
		*i = 0
	case "a":
		*i = 1
	case "b":
		*i = 2
	}
	return nil
}
func (i TopoSide) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	var s string
	switch i {
	default:
		s = "both"
	case 1:
		s = "a"
	case 2:
		s = "b"
	}
	return xml.Attr{Name: name, Value: s}, nil
}

// TopologyNode represent a node in the IRONdb topology structure.
type TopologyNode struct {
	XMLName     xml.Name `xml:"node" json:"-"`
	ID          string   `xml:"id,attr" json:"id"`
	Address     string   `xml:"address,attr" json:"address"`
	Port        uint16   `xml:"port,attr" json:"port"`
	APIPort     uint16   `xml:"apiport,attr" json:"apiport"`
	Weight      uint16   `xml:"weight,attr" json:"weight"`
	Side        TopoSide `xml:"side,attr" json:"side"`
	WriteCopies uint8    `xml:"-" json:"n"`
}

func (topo *Topology) compile() error {
	nslots := 0
	if topo.WriteCopies == 0 {
		topo.WriteCopies = topo.OldWriteCopies
	}
	topo.OldWriteCopies = topo.WriteCopies
	for _, node := range topo.Nodes {
		node.ID = strings.ToLower(node.ID)
		if node.Side != 0 {
			topo.use_side = true
		}
		nslots += int(node.Weight)
	}
	hash := sha256.New()
	for _, node := range topo.Nodes {
		hash.Write([]byte(node.ID))
		hash.Write([]byte{0, 0})
		netshort := make([]byte, 2)
		binary.BigEndian.PutUint16(netshort, node.Weight)
		hash.Write(netshort)
		if topo.use_side {
			binary.BigEndian.PutUint16(netshort, uint16(node.Side))
			hash.Write(netshort)
		}
	}
	// This matches the horrible backware compatibility requirements in the C version
	if topo.WriteCopies != 2 {
		hash.Write(bytes.Repeat([]byte{0}, 38))
		netshort := make([]byte, 2)
		binary.BigEndian.PutUint16(netshort, uint16(topo.WriteCopies))
		hash.Write(netshort)
		if topo.use_side {
			binary.BigEndian.PutUint16(netshort, 0)
			hash.Write(netshort)
		}
	}
	sum := hex.EncodeToString(hash.Sum(nil))
	if topo.Hash == "" {
		topo.Hash = sum
	}
	if topo.Hash != sum {
		return errors.New("bad topology hash")
	}

	topo.ring = make([]TopologyNodeSlot, nslots)
	i := 0
	for node_idx, node := range topo.Nodes {
		for w := uint16(0); w < node.Weight; w++ {
			location := make([]byte, 38)
			copy(location, node.ID)
			binary.BigEndian.PutUint16(location[36:], w)
			sum := sha256.Sum256(location)
			copy(topo.ring[i].Location[:], sum[:])
			topo.ring[i].Node = &topo.Nodes[node_idx]
			topo.ring[i].Idx = w
			if node.Side == 1 {
				topo.ring[i].Location[0] &= 0x7f
			} else if node.Side == 2 {
				topo.ring[i].Location[0] |= 0x80
			}
			i++
		}
	}
	sort.Sort(topo)
	return nil
}

func (topo *Topology) binSearchNext(location [sha256.Size]byte) (match, next int) {
	start := 0
	end := len(topo.ring) - 1
	mid := len(topo.ring) / 2
	cmp := -1
	for start <= end {
		cmp = bytes.Compare(location[:], topo.ring[mid].Location[:])

		if cmp == 0 {
			return mid, mid
		}
		if cmp < 0 {
			end = mid - 1
			mid = (start + end) / 2
			continue
		}
		start = mid + 1
		mid = (start + end) / 2
	}
	cmp = bytes.Compare(location[:], topo.ring[mid].Location[:])
	if cmp > 0 {
		mid++
	}
	if mid >= len(topo.ring) {
		mid = 0
	}
	return -1, mid
}
func nodeListContains(list []TopologyNode, item TopologyNode) bool {
	for _, node := range list {
		if node.ID == item.ID {
			return true
		}
	}
	return false
}
func (topo *Topology) FindNext(location [sha256.Size]byte, found []TopologyNode) *TopologyNode {
	side := 0
	if topo.use_side {
		side = 1
		if (location[0] & 0x80) != 0 {
			side = 2
		}
	}
	idx, next := topo.binSearchNext(location)
	if idx < 0 {
		idx = next
	}
	for attempts := 0; attempts < 2; attempts++ {
		for idx < len(topo.ring) {
			slot := &topo.ring[idx]
			if side == 0 ||
				((slot.Location[0]&0x80) == 0 && side == 1) ||
				((slot.Location[0]&0x80) != 0 && side == 2) {
				// We're unsided or this side matches
				if len(found) == 0 || !nodeListContains(found, *slot.Node) {
					return slot.Node
				}
			}
			idx++
		}
		idx = 0
	}
	return nil
}
func (topo *Topology) FindMetric(uuid, metric string) ([]TopologyNode, error) {
	return topo.FindN(strings.ToLower(uuid)+"-"+metric, int(topo.WriteCopies))
}
func (topo *Topology) FindMetricN(uuid, metric string, n int) ([]TopologyNode, error) {
	return topo.FindN(strings.ToLower(uuid)+"-"+metric, n)
}
func (topo *Topology) Find(s string) ([]TopologyNode, error) {
	return topo.FindN(s, int(topo.WriteCopies))
}
func (topo *Topology) FindN(s string, n int) ([]TopologyNode, error) {
	if topo.ring == nil || len(topo.ring) < 1 {
		return nil, errors.New("Empty topology")
	}
	location := sha256.Sum256([]byte(s))
	nodes := make([]TopologyNode, 0)
	for i := 0; i < n; i++ {
		node := topo.FindNext(location, nodes)
		if node == nil {
			break
		}
		nodes = append(nodes, *node)
		location[0] ^= 0x80
	}
	return nodes, nil
}

// GetTopologyInfo retrieves topology information from a node.
func (sc *SnowthClient) GetTopologyInfo(node *SnowthNode) (*Topology, error) {
	return sc.GetTopologyInfoContext(context.Background(), node)
}

func TopologyLoadXML(xml string) (*Topology, error) {
	r := &Topology{}

	if err := decodeXML(strings.NewReader(xml), &r); err != nil {
		return nil, errors.Wrap(err, "unable to decode IRONdb response")
	}
	if err := r.compile(); err != nil {
		return nil, err
	}
	return r, nil
}

// GetTopologyInfoContext is the context aware version of GetTopologyInfo.
func (sc *SnowthClient) GetTopologyInfoContext(ctx context.Context,
	node *SnowthNode) (*Topology, error) {
	r := &Topology{}
	if node == nil {
		nodes := sc.ListActiveNodes()
		if len(nodes) == 0 {
			return nil, errors.New("no active nodes")
		}
		node = nodes[0]
	}
	topology_id := node.GetCurrentTopology()
	if topology_id == "" {
		return nil, errors.New("no active topology")
	}
	if topology_id == sc.currentTopology && sc.currentTopologyCompiled != nil {
		return sc.currentTopologyCompiled, nil
	}
	body, _, err := sc.do(ctx, node, "GET",
		path.Join("/topology/xml", node.GetCurrentTopology()), nil, nil)
	if err != nil {
		return nil, err
	}

	if err := decodeXML(body, &r); err != nil {
		return nil, errors.Wrap(err, "unable to decode IRONdb response")
	}
	if err = r.compile(); err != nil {
		return nil, err
	}
	sc.currentTopology = topology_id
	sc.currentTopologyCompiled = r

	return r, nil
}

// LoadTopology loads a new topology on a node without activating it.
func (sc *SnowthClient) LoadTopology(hash string, t *Topology,
	node *SnowthNode) error {
	return sc.LoadTopologyContext(context.Background(), hash, t, node)
}

// LoadTopologyContext is the context aware version of LoadTopology.
func (sc *SnowthClient) LoadTopologyContext(ctx context.Context, hash string,
	t *Topology, node *SnowthNode) error {
	b, err := encodeXML(t)
	if err != nil {
		return errors.Wrap(err, "failed to encode request data")
	}

	_, _, err = sc.do(ctx, node, "POST", path.Join("/topology", hash), b, nil)
	return err
}

// ActivateTopology activates a new topology on the node.
// WARNING THIS IS DANGEROUS.
func (sc *SnowthClient) ActivateTopology(hash string, node *SnowthNode) error {
	return sc.ActivateTopologyContext(context.Background(), hash, node)
}

// ActivateTopologyContext is the context aware version of ActivateTopology.
// WARNING THIS IS DANGEROUS.
func (sc *SnowthClient) ActivateTopologyContext(ctx context.Context,
	hash string, node *SnowthNode) error {
	_, _, err := sc.do(ctx, node, "GET", path.Join("/activate", hash), nil, nil)
	return err
}
