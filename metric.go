package gosnowth

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

// scanToken values represent individual tokens found by query scanners.
type scanToken uint

// scanToken values used by the scanners to identify tokens.
const (
	tokenEOF scanToken = iota
	tokenIllegal
	tokenB64 // Base64 encoding marker
	tokenWS  // Whitespace
	tokenOP  // Open parenthesis
	tokenCP  // Close parenthesis
	tokenOB  // Open bracket
	tokenCB  // Close bracket
	tokenColon
	tokenComma
	tokenQuote
	tokenKeyword
	tokenTagCat
	tokenTagVal
	tokenMetric
	tokenStreamTag
	tokenMeasurementTag
)

// tagType value represent whether a tag is a stream or measurment tag.
type tagType uint

const (
	tagStreamTag = iota
	tagMeasurementTag
)

// Tag values represent stream or measurment tags.
type Tag struct {
	Key   string
	Value string
}

// MetricName values represent metric names.
type MetricName struct {
	CanonicalName   string
	Name            string
	StreamTags      []Tag
	MeasurementTags []Tag
}

// NewMetricName creates and initializes a new metric name value.
func NewMetricName() *MetricName {
	return &MetricName{}
}

// metricScanner represents a lexical scanner for metric names.
type metricScanner struct {
	r *bufio.Reader
}

// newMetricScanner returns a new metric name scanner value.
func newMetricScanner(r io.Reader) *metricScanner {
	return &metricScanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (ms *metricScanner) read() rune {
	ch, _, err := ms.r.ReadRune()
	if err != nil {
		return rune(0) // EOF
	}

	return ch
}

// unread places the previously read rune back on the reader.
func (ms *metricScanner) unread() error { return ms.r.UnreadRune() }

// Scan returns the next token and literal value.
func (ms *metricScanner) scan() (tok scanToken, lit string) {
	ch := ms.read()
	switch ch {
	case rune(0):
		return tokenEOF, ""
	case '[':
		return tokenOB, string(ch)
	case ']':
		return tokenCB, string(ch)
	case '"':
		return tokenQuote, string(ch)
	case 'b':
		return tokenB64, string(ch)
	case ':':
		return tokenColon, string(ch)
	case ',':
		return tokenComma, string(ch)
	}

	return tokenIllegal, string(ch)
}

// scanTagSep attempts to scan a tag separator from the scan buffer.
func (ms *metricScanner) scanTagSep() (scanToken, string, error) {
	var buf bytes.Buffer
	if ch := ms.read(); ch == '|' {
		if _, err := buf.WriteRune(ch); err != nil {
			return tokenIllegal, "", fmt.Errorf(
				"unable to write to tag separator buffer: %w", err)
		}

		for i := 0; i < 2; i++ {
			ch := ms.read()
			if _, err := buf.WriteRune(ch); err != nil {
				return tokenIllegal, "", fmt.Errorf(
					"unable to write to tag separator buffer: %w", err)
			}
		}

		switch buf.String() {
		case "|ST":
			return tokenStreamTag, buf.String(), nil
		case "|MT":
			return tokenMeasurementTag, buf.String(), nil
		default:
			return tokenIllegal, "", nil
		}
	} else if ch == rune(0) {
		return tokenEOF, "", nil
	}

	return tokenIllegal, "", nil
}

// peekTagSep checks for a tag separator next in the scan buffer.
func (ms *metricScanner) peekTagSep() (scanToken, string, error) {
	if ch := ms.read(); ch == '|' {
		if err := ms.unread(); err != nil {
			return tokenIllegal, "", fmt.Errorf(
				"unable to unread to scan buffer: %w", err)
		}

		if b, err := ms.r.Peek(3); err == nil {
			switch string(b) {
			case "|ST":
				return tokenStreamTag, string(b), nil
			case "|MT":
				return tokenMeasurementTag, string(b), nil
			default:
				return tokenIllegal, "", nil
			}
		}
	}

	return tokenIllegal, "", nil
}

// scanMetricName consumes the current rune and all contiguous ident runes.
func (ms *metricScanner) scanMetricName() (scanToken, string, error) {
	var buf bytes.Buffer
	for {
		ch := ms.read()
		if ch == '|' {
			if err := ms.unread(); err != nil {
				return tokenIllegal, "", fmt.Errorf(
					"unable to unread to scan buffer: %w", err)
			}

			tok, _, err := ms.peekTagSep()
			if err != nil {
				return tokenIllegal, "", fmt.Errorf(
					"unable to peek tag separator from scan buffer: %w", err)
			}

			if tok != tokenIllegal {
				// we have a valid separator done scanning name
				break
			}

			// otherwise write the character
			ch := ms.read()
			if _, err := buf.WriteRune(ch); err != nil {
				return tokenIllegal, "", fmt.Errorf(
					"unable to write to metric name buffer: %w", err)
			}
		} else if ch == rune(0) { // EOF
			break
		} else if _, err := buf.WriteRune(ch); err != nil {
			return tokenIllegal, "", fmt.Errorf(
				"unable to write to metric name buffer: %w", err)
		}
	}

	return tokenMetric, buf.String(), nil
}

// scanTagName attempts to read a tag name token from the scan buffer.
func (ms *metricScanner) scanTagName() (scanToken, string, error) {
	var buf bytes.Buffer
	quoted := false

loop:
	for {
		ch := ms.read()
		switch ch {
		case '"':
			quoted = !quoted

			if _, err := buf.WriteRune(ch); err != nil {
				return tokenIllegal, "", fmt.Errorf(
					"unable to write to tag name buffer: %w", err)
			}
		case ':':
			if !quoted {
				if err := ms.unread(); err != nil {
					return tokenIllegal, "", fmt.Errorf(
						"unable to unread to scan buffer: %w", err)
				}

				break loop
			} else {
				if _, err := buf.WriteRune(ch); err != nil {
					return tokenIllegal, "", fmt.Errorf(
						"unable to write to tag name buffer: %w", err)
				}
			}
		case rune(0): // EOF
			break loop
		default:
			if _, err := buf.WriteRune(ch); err != nil {
				return tokenIllegal, "", fmt.Errorf(
					"unable to write to tag name buffer: %w", err)
			}
		}
	}

	return tokenTagCat, buf.String(), nil
}

// scanTagValue attempts to read a tag value token from the scan buffer.
func (ms *metricScanner) scanTagValue() (scanToken, string, error) {
	var buf bytes.Buffer
	quoted := false

loop:
	for {
		ch := ms.read()
		switch ch {
		case '"':
			quoted = !quoted

			if _, err := buf.WriteRune(ch); err != nil {
				return tokenIllegal, "", fmt.Errorf(
					"unable to write to tag name buffer: %w", err)
			}
		case ',', ']':
			if !quoted {
				if err := ms.unread(); err != nil {
					return tokenIllegal, "", fmt.Errorf(
						"unable to unread to scan buffer: %w", err)
				}

				break loop
			} else {
				if _, err := buf.WriteRune(ch); err != nil {
					return tokenIllegal, "", fmt.Errorf(
						"unable to write to tag name buffer: %w", err)
				}
			}
		case rune(0): // EOF
			break loop
		default:
			if _, err := buf.WriteRune(ch); err != nil {
				return tokenIllegal, "", fmt.Errorf(
					"unable to write to tag name buffer: %w", err)
			}
		}
	}

	return tokenTagVal, buf.String(), nil
}

// MetricParser values are used to parse metric names and stream tags.
type MetricParser struct {
	s *metricScanner
}

// NewMetricParser returns a new instance of MetricParser.
func NewMetricParser(r io.Reader) *MetricParser {
	return &MetricParser{s: newMetricScanner(r)}
}

// parseTagSet performs the functionality to parse a tag set.
func (mp *MetricParser) parseTagSet(tt tagType) (string, []Tag, error) {
	canonical := strings.Builder{}
	tags := []Tag{}

	switch tt {
	case tagStreamTag:
		canonical.WriteString("|ST[")
	case tagMeasurementTag:
		canonical.WriteString("|MT[")
	default:
		return "", nil, fmt.Errorf("invalid tag type: %v", tt)
	}

	var tok scanToken
	var lit string
	var err error

	if tok, lit = mp.s.scan(); tok != tokenOB {
		return "", nil, fmt.Errorf("parse failure, expecting '[' got: %s", lit)
	}

	for {
		var tag = Tag{}

		tok, lit, err = mp.s.scanTagName()
		if err != nil {
			return "", nil,
				fmt.Errorf("unable to parse stream tag name: %w", err)
		}

		if tok != tokenTagCat {
			return "", nil, fmt.Errorf("expected stream tag name, got: %s", lit)
		}

		tag.Key = lit
		if strings.HasPrefix(tag.Key, `b"`) &&
			strings.HasSuffix(tag.Key, `"`) {
			val := strings.Trim(tag.Key[1:], `"`)
			b, err := ioutil.ReadAll(base64.NewDecoder(base64.StdEncoding,
				bytes.NewBufferString(val)))
			if err != nil {
				return "", nil, fmt.Errorf(
					"unable to parse base64 stream tag category: %w", err)
			}

			tag.Key = string(b)
		}

		if strings.HasPrefix(tag.Key, `"`) && strings.HasSuffix(tag.Key, `"`) {
			tag.Key = strings.Trim(tag.Key, `"`)
		}

		canonical.WriteString(lit)

		if tok, lit = mp.s.scan(); tok != tokenColon {
			return "", nil,
				fmt.Errorf("parse failure, expecting ':' got: %s", lit)
		}

		canonical.WriteString(":")

		tok, lit, err = mp.s.scanTagValue()
		if err != nil {
			return "", nil,
				fmt.Errorf("unable to parse stream tag value: %w", err)
		}

		if tok != tokenTagVal {
			return "", nil,
				fmt.Errorf("expected stream tag value, got: %s", lit)
		}

		tag.Value = lit
		if strings.HasPrefix(tag.Value, `b"`) &&
			strings.HasSuffix(tag.Value, `"`) {
			val := strings.Trim(tag.Value[1:], `"`)
			b, err := ioutil.ReadAll(base64.NewDecoder(base64.StdEncoding,
				bytes.NewBufferString(val)))
			if err != nil {
				return "", nil, fmt.Errorf(
					"unable to parse base64 stream tag value: %w", err)
			}

			tag.Value = string(b)
		}

		if strings.HasPrefix(tag.Value, `"`) &&
			strings.HasSuffix(tag.Value, `"`) {
			tag.Value = strings.Trim(tag.Value, `"`)
		}

		canonical.WriteString(lit)

		tags = append(tags, tag)

		tok, lit = mp.s.scan()
		if tok == tokenComma {
			// there are additional tags
			canonical.WriteString(",")
			continue
		}

		if tok == tokenCB {
			// done with tags
			canonical.WriteString("]")
			break
		}

		return "", nil, fmt.Errorf("should have , or ], got: %s", lit)
	}

	return canonical.String(), tags, nil
}

// Parse scans and parses a metric name value.
func (mp *MetricParser) Parse() (*MetricName, error) {
	canonical := strings.Builder{}
	metricName := NewMetricName()

	tok, lit, err := mp.s.scanMetricName()
	if err != nil {
		return nil, fmt.Errorf("unable to scan metric name token: %w", err)
	}

	if tok != tokenMetric {
		return nil, fmt.Errorf("expected metric identifier, got: %s ", lit)
	}

	canonical.WriteString(lit)
	metricName.Name = lit

	for {
		// Get any tags in the metric name.
		tok, _, err = mp.s.scanTagSep()
		if err != nil {
			return nil, fmt.Errorf("unable to scan metric tag token: %w", err)
		}

		if tok == tokenEOF {
			break
		}

		if tok == tokenStreamTag {
			can, tags, err := mp.parseTagSet(tagStreamTag)
			if err != nil {
				return nil, err
			}

			canonical.WriteString(can)
			metricName.StreamTags = append(metricName.StreamTags, tags...)
		} else if tok == tokenMeasurementTag {
			can, tags, err := mp.parseTagSet(tagMeasurementTag)
			if err != nil {
				return nil, err
			}

			canonical.WriteString(can)
			metricName.MeasurementTags = append(metricName.MeasurementTags,
				tags...)
		}
	}

	metricName.CanonicalName = canonical.String()

	return metricName, nil
}

// ParseMetricName takes a canonical metric name as a string and parses it
// into a MetricName value containing separated stream tags.
func ParseMetricName(name string) (*MetricName, error) {
	p := NewMetricParser(bytes.NewBufferString(name))

	return p.Parse()
}