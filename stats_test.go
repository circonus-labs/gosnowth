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
		"_value": "bb6f7162-4828-11df-bab8-6bac200dcc2a"
	},
	"version": {
		"_type": "s",
		"_value": "v52bcc96a9a1a41acd96352b9b63e59cba2b6a8a9\/65ab82cb7281e76e96b2fedafdc6594d50437d91"
	},
	"topology": {
		"next": {
			"_type": "s",
			"_value": "-"
		},
		"current": {
			"_type": "s",
			"_value": "294cbd39999c2270964029691e8bc5e231a867d525ccba62181dc8988ff218dc"
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
			_, _ = w.Write([]byte(stateTestData))
			return
		}

		if r.RequestURI == "/stats.json" {
			_, _ = w.Write([]byte(statsTestData))
			return
		}

		if strings.HasPrefix(r.RequestURI, "/find/1/tags?query=test") {
			_, _ = w.Write([]byte(tagsTestData))
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

	exp := "bb6f7162-4828-11df-bab8-6bac200dcc2a"
	if res.Identity() != exp {
		t.Errorf("Expected identity: %v, got: %v", exp, res.Identity())
	}

	exp = "0.1.1570000000"
	if res.SemVer() != exp {
		t.Errorf("Expected version: %v, got: %v", exp, res.SemVer())
	}

	exp = "294cbd39999c2270964029691e8bc5e231a867d525ccba62181dc8988ff218dc"
	if res.CurrentTopology() != exp {
		t.Errorf("Expected current: %v, got: %v", exp, res.CurrentTopology())
	}

	exp = "-"
	if res.NextTopology() != exp {
		t.Errorf("Expected next: %v, got: %v", exp, res.NextTopology())
	}
}
