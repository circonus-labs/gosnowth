package gosnowth

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const statsTestData = `{
	"application": {
		"_type": "s",
		"_value": "snowth"
	},
	"identity": {
		"_type": "s",
		"_value": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
	},
	"version": {
		"_type": "s",
		"_value": "test"
	},
	"topology": {
		"next": {
			"_type": "s",
			"_value": "-"
		},
		"current": {
			"_type": "s",
			"_value": "test"
		}
	},
	"semver": {
		"_type": "s",
		"_value": "0.1.1570000000"
	}
}`

func TestGetStats(t *testing.T) {
	ms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		if r.RequestURI == "/state" {
			w.Write([]byte(stateTestData))
			return
		}

		if r.RequestURI == "/stats.json" {
			w.Write([]byte(statsTestData))
			return
		}

		if strings.HasPrefix(r.RequestURI, "/find/1/tags?query=test") {
			w.Write([]byte(tagsTestData))
			return
		}
	}))

	defer ms.Close()
	sc, err := NewSnowthClient(false, ms.URL)
	if err != nil {
		t.Fatal("Unable to create snowth client", err)
	}

	u, err := url.Parse(ms.URL)
	if err != nil {
		t.Fatal("Invalid test URL")
	}

	node := &SnowthNode{url: u}
	res, err := sc.GetStats(node)
	if err != nil {
		t.Fatal(err)
	}

	exp := "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
	if res.Identity() != exp {
		t.Fatalf("Expected identity: %v, got: %v", exp, res.Identity())
	}

	exp = "0.1.1570000000"
	if res.SemVer() != exp {
		t.Fatalf("Expected version: %v, got: %v", exp, res.SemVer())
	}

	exp = "test"
	if res.CurrentTopology() != exp {
		t.Fatalf("Expected current: %v, got: %v", exp, res.CurrentTopology())
	}

	exp = "-"
	if res.NextTopology() != exp {
		t.Fatalf("Expected next: %v, got: %v", exp, res.NextTopology())
	}
}
