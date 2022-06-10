package gosnowth

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const testLuaExtensionData = `{
	"test": {
		"documentation": "# test\nReturns parsed payload ",
		"method": null,
		"PARSE_JSON_PAYLOAD": true,
		"params": {
			"tests": {
				"type": "list",
				"default": [],
				"optional": true,
				"alias_list": [
					"alias"
				],
				"description": "comma separated list of tests to run",
				"name": "tests"
			},
			"bc": {
				"type": "boolean",
				"default": false,
				"optional": true,
				"description": "print bytecode of tests"
			},
			"print_output": {
				"type": "boolean",
				"default": false,
				"optional": true,
				"description": "print output of test function"
			},
			"iterations": {
				"type": "number",
				"default": 1,
				"optional": true,
				"description": "Repeat each test multiple times"
			}
		},
		"man": "",
		"name": "test",
		"description": "Returns parsed payload (for testing)."
	},
	"cor": [],
	"registration": {
		"documentation": "",
		"method": null,
		"params": [],
		"name": "registration",
		"description": "Returns the extension register."
	}
}`

func TestGetLuaExtension(t *testing.T) {
	t.Parallel()

	ms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request,
	) {
		if r.RequestURI == "/state" {
			_, _ = w.Write([]byte(stateTestData))

			return
		}

		if r.RequestURI == "/stats.json" {
			_, _ = w.Write([]byte(statsTestData))

			return
		}

		if strings.HasPrefix(r.RequestURI, "/extension/lua") {
			_, _ = w.Write([]byte(testLuaExtensionData))

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
	res, err := sc.GetLuaExtensions(node)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 3 {
		t.Fatalf("Expected length: 3, got: %v", len(res))
	}

	if res["test"].Name != "test" {
		t.Errorf("Expected name: test, got: %v", res["test"].Name)
	}

	exp := "Returns parsed payload (for testing)."
	if res["test"].Description != exp {
		t.Errorf("Expected name: %v, got: %v", exp, res["test"].Description)
	}

	if len(res["test"].Params) != 4 {
		t.Fatalf("Expected params length: 4, got: %v", len(res["test"].Params))
	}

	if len(res["test"].Params["tests"].AliasList) != 1 {
		t.Fatalf("Expected alias list length: 1, got: %v",
			len(res["test"].Params["tests"].AliasList))
	}

	if res["test"].Params["tests"].AliasList[0] != "alias" {
		t.Fatalf("Expected alias: alias, got: %v",
			res["test"].Params["tests"].AliasList[0])
	}

	if v, ok := res["test"].Params["iterations"].Default.(float64); !ok ||
		v != 1.0 {
		t.Fatalf("Expected default: 1, got: %v", v)
	}

	if len(res["registration"].Params) != 0 {
		t.Errorf("Expected params length: 0, got: %v",
			len(res["registration"].Params))
	}

	if res["cor"].Description != "" {
		t.Errorf("Expected empty description, got: %v", res["cor"].Description)
	}
}

func TestExecLuaExtensionContext(t *testing.T) {
	t.Parallel()

	ms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request,
	) {
		if r.RequestURI == "/state" {
			_, _ = w.Write([]byte(stateTestData))

			return
		}

		if r.RequestURI == "/stats.json" {
			_, _ = w.Write([]byte(statsTestData))

			return
		}

		if strings.HasPrefix(r.RequestURI, "/extension/lua/test") {
			if r.URL.Query().Get("test") == "1" {
				_, _ = w.Write([]byte(`{"test":1}`))

				return
			}
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
	res, err := sc.ExecLuaExtension("test",
		[]ExtParam{{Name: "test", Value: "1"}}, node)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 1 {
		t.Fatalf("Expected length: 1, got: %v", len(res))
	}

	v, ok := res["test"].(float64)
	if !ok {
		t.Fatal("Unexpected result test value")
	}

	if v != 1.0 {
		t.Errorf("Expected value: 1, got: %v", v)
	}
}
