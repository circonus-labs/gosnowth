package gosnowth

// DF4Response values represent time series data in the DF4 format.
type DF4Response struct {
	Data [][]interface{} `json:"data,omitempty"`
	Meta []DF4Meta       `json:"meta,omitempty"`
	Ver  string          `json:"version,omitempty"`
	Head DF4Head         `json:"head"`
}

// DF4Meta values contain information and metadata about the metrics in a DF4
// time series data response.
type DF4Meta struct {
	Kind  string   `json:"kind"`
	Label string   `json:"label"`
	Tags  []string `json:"tags,omitempty"`
}

// DF4Head values contain information about the time range of the data elements
// in a DF4 time series data response.
type DF4Head struct {
	Count  int64 `json:"count"`
	Start  int64 `json:"start"`
	Period int64 `json:"period"`
}

// Copy returns a deep copy of the base DF4 response.
func (dr *DF4Response) Copy() *DF4Response {
	b := &DF4Response{
		Data: make([][]interface{}, len(dr.Data)),
		Meta: make([]DF4Meta, len(dr.Meta)),
		Ver:  dr.Ver,
		Head: DF4Head{
			Count:  dr.Head.Count,
			Start:  dr.Head.Start,
			Period: dr.Head.Period,
		},
	}

	copy(b.Meta, dr.Meta)
	for i, v := range dr.Data {
		b.Data[i] = make([]interface{}, len(v))
		copy(b.Data[i], v)
	}

	return b
}
