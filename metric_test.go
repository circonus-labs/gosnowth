package gosnowth

import (
	"bytes"
	"testing"
)

func TestScanMetricName(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input string    // input
		tok   scanToken // token
		lit   string    // metric name literal
	}{
		{
			input: "testing",
			tok:   tokenMetric,
			lit:   "testing",
		},
		{
			input: "testing*",
			tok:   tokenMetric,
			lit:   "testing*",
		},
		{
			input: "testing|ST[blah:blah]",
			tok:   tokenMetric,
			lit:   "testing",
		},
		{
			input: "t|esting|ST[blah:blah]",
			tok:   tokenMetric,
			lit:   "t|esting",
		},
		{
			input: "t|Sesting|ST[blah:blah]",
			tok:   tokenMetric,
			lit:   "t|Sesting",
		},
		{
			input: "tel|a|M|Ssting|ST[blah:blah]",
			tok:   tokenMetric,
			lit:   "tel|a|M|Ssting",
		},
		{
			input: "testing|ST[blah:blah]",
			tok:   tokenMetric,
			lit:   "testing",
		},
		{
			input: "testing|MT{test:test}",
			tok:   tokenMetric,
			lit:   "testing",
		},
		{
			input: "testing|ST[blah:blah]|MT{test:test}",
			tok:   tokenMetric,
			lit:   "testing",
		},
	}

	for _, c := range cases {
		buf := bytes.NewBufferString(c.input)
		s := newMetricScanner(buf)

		tok, lit, err := s.scanMetricName()
		if err != nil {
			t.Fatal(err)
		}

		if tok != c.tok {
			t.Error("failed to find the metric ident")
		}

		if lit != c.lit {
			t.Error("incorrect literal scanned: ", c.lit, lit)
		}
	}
}

func TestScanTagSep(t *testing.T) {
	t.Parallel()

	stBuf := bytes.NewBufferString("|ST[blah:blah]")
	s := newMetricScanner(stBuf)

	tok, lit, err := s.peekTagSep()
	if err != nil {
		t.Fatal(err)
	}

	if tok != tokenStreamTag {
		t.Error("expected stream tag separator")
	}

	if lit != "|ST" {
		t.Error("incorrect literal scanned: ", lit)
	}

	tok, lit, err = s.scanTagSep()
	if err != nil {
		t.Fatal(err)
	}

	if tok != tokenStreamTag {
		t.Error("expected stream tag separator")
	}

	if lit != "|ST" {
		t.Error("incorrect literal scanned: ", lit)
	}

	mtBuf := bytes.NewBufferString("|MT{blah:blah}")
	s = newMetricScanner(mtBuf)

	tok, lit, err = s.peekTagSep()
	if err != nil {
		t.Fatal(err)
	}

	if tok != tokenMeasurementTag {
		t.Error("expected measurement tag separator")
	}

	if lit != "|MT" {
		t.Error("incorrect literal scanned: ", lit)
	}

	tok, lit, err = s.scanTagSep()
	if err != nil {
		t.Fatal(err)
	}

	if tok != tokenMeasurementTag {
		t.Error("expected measurement tag separator")
	}

	if lit != "|MT" {
		t.Error("incorrect literal scanned: ", lit)
	}
}

func TestMetricParser(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input string // input
		numST int    // number of derived stream tags
		numMT int    // number of derived measurement tags
		lit   string // metric name literal
	}{
		{
			input: "testing",
			numST: 0,
			numMT: 0,
			lit:   "testing",
		},
		{
			input: "testing*",
			numST: 0,
			numMT: 0,
			lit:   "testing*",
		},
		{
			input: "testing|ST[blah:blah]",
			numST: 1,
			numMT: 0,
			lit:   "testing",
		},
		{
			input: "t|esting|ST[blah:blah]",
			numST: 1,
			numMT: 0,
			lit:   "t|esting",
		},
		{
			input: "t|Sesting|ST[blah:blah]",
			numST: 1,
			numMT: 0,
			lit:   "t|Sesting",
		},
		{
			input: "tel|a|M|Ssting|ST[blah:blah]",
			numST: 1,
			numMT: 0,
			lit:   "tel|a|M|Ssting",
		},
		{
			input: "testing|ST[blah:blah]",
			numST: 1,
			numMT: 0,
			lit:   "testing",
		},
		{
			input: "testing|ST[blah:blah,blah1:blah,blah2:blah]",
			numST: 3,
			numMT: 0,
			lit:   "testing",
		},
		{
			input: `testing|ST[b"QUFB":blah,blah1:b"QkJCQgo=",b"YWFh":b"YmJi"]`,
			numST: 3,
			numMT: 0,
			lit:   "testing",
		},
		{
			input: "testing|ST[blah:blah:blah:blah]",
			numST: 1,
			numMT: 0,
			lit:   "testing",
		},
		{
			input: "testing|ST[blah:blah]|ST[blah:blah]",
			numST: 2,
			numMT: 0,
			lit:   "testing",
		},
		{
			input: "testing|MT{blah:blah}",
			numST: 0,
			numMT: 1,
			lit:   "testing",
		},
		{
			input: "testing|ST[blah:blah]|MT{blah:blah}",
			numST: 1,
			numMT: 1,
			lit:   "testing",
		},
		{
			input: "testing|ST[blah:blah]|MT{blah:blah}|ST[blah:blah]",
			numST: 2,
			numMT: 1,
			lit:   "testing",
		},
		{
			input: `testing|ST["blah:|ST[]":blah]|MT{blah:",}:|MTblah"}`,
			numST: 1,
			numMT: 1,
			lit:   "testing",
		},
		{
			input: `testing|ST["blah:|ST[]":blah]|MT{blah:",}:|MTblah"}`,
			numST: 1,
			numMT: 1,
			lit:   "testing",
		},
		{
			input: `testing|ST["b:|ST[]":b]|MT{b:",}:|MTb"}|ST[a:b]|MT{c:d}`,
			numST: 2,
			numMT: 2,
			lit:   "testing",
		},
		{
			input: `testing|ST["quote\"slash\\":bar]`,
			numST: 1,
			numMT: 0,
			lit:   "testing",
		},
		{
			input: `testing|ST["quote\"slash\\":bar]|MT{"q\\\"\\:":bar}`,
			numST: 1,
			numMT: 1,
			lit:   "testing",
		},
	}

	for _, c := range cases {
		buf := bytes.NewBufferString(c.input)
		p := NewMetricParser(buf)

		metricName, err := p.Parse()
		if err != nil {
			t.Fatal("failed to parse the metric name", err)
		}

		if metricName.Name != c.lit {
			t.Error("incorrect literal scanned: ", c.lit, metricName.Name)
		}

		if metricName.CanonicalName != c.input {
			t.Error("incorrect canonical: ", c.input, metricName.CanonicalName)
		}

		if len(metricName.StreamTags) != c.numST {
			t.Error("incorrect number of stream tags: ", c.numST,
				len(metricName.StreamTags))
		}

		if len(metricName.MeasurementTags) != c.numMT {
			t.Error("incorrect number of measurement tags: ", c.numMT,
				len(metricName.MeasurementTags))
		}
	}
}

func TestMetricParserComplex(t *testing.T) {
	t.Parallel()

	mn, err := ParseMetricName("test|ST[sri:sri:spc:compute::" +
		"stackable-cloud:vm/0bce9aef-8186-44c9-a73d-33723173b2c6]")
	if err != nil {
		t.Fatal("failed to parse the metric name")
	}

	if mn.Name != "test" {
		t.Errorf("Expected name: test, got: %v", mn.Name)
	}

	exp := "sri:spc:compute::stackable-cloud:vm" +
		"/0bce9aef-8186-44c9-a73d-33723173b2c6"
	if mn.StreamTags[0].Value != exp {
		t.Errorf("Expected stream tag value: %v, got: %v", exp,
			mn.StreamTags[0].Value)
	}
}
