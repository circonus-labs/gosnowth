package gosnowth

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	if _, err := NewConfig(":invalid"); err == nil {
		t.Fatal("Expected invalid server error")
	}

	cfg, err := NewConfig("test1", "test2")
	if err != nil {
		t.Fatal(err)
	}

	if err := cfg.SetDialTimeout(time.Second); err != nil {
		t.Fatal(err)
	}

	cfg.SetDiscover(true)
	if err := cfg.SetServers("test1", "test2"); err != nil {
		t.Fatal(err)
	}

	if err := cfg.SetTimeout(time.Second); err != nil {
		t.Fatal(err)
	}

	cfg.SetRetries(1)
	cfg.SetConnectRetries(1)
	if err := cfg.SetWatchInterval(time.Second); err != nil {
		t.Fatal(err)
	}

	if cfg.DialTimeout() != time.Second {
		t.Errorf("Expected dial timeout: %v, got: %v",
			time.Second, cfg.DialTimeout())
	}

	if cfg.Discover() != true {
		t.Errorf("Expected discover value: %v, got: %v",
			true, cfg.Discover())
	}

	if len(cfg.Servers()) != 2 {
		t.Fatalf("Expected servers length: %v, got: %v",
			2, len(cfg.Servers()))
	}

	if cfg.Servers()[0] != "test1" {
		t.Errorf("Expected server value: %v, got: %v",
			"test1", cfg.Servers()[0])
	}

	if cfg.Servers()[1] != "test2" {
		t.Errorf("Expected server value: %v, got: %v",
			"test2", cfg.Servers()[1])
	}

	if cfg.Timeout() != time.Second {
		t.Errorf("Expected timeout: %v, got: %v",
			time.Second, cfg.Timeout())
	}

	if cfg.Retries() != 1 {
		t.Errorf("Expected retries: 1, got: %v", cfg.Retries())
	}

	if cfg.ConnectRetries() != 1 {
		t.Errorf("Expected connect retries: 1, got: %v", cfg.ConnectRetries())
	}

	if cfg.WatchInterval() != time.Second {
		t.Errorf("Expected watch interval: %v, got: %v",
			time.Second, cfg.WatchInterval())
	}
}

func TestConfigMarshalJSON(t *testing.T) {
	t.Parallel()

	s := `{"dial_timeout":"100ms","discover":true,"timeout":"1s",` +
		`"watch_interval":"5s","connect_retries":-1,` +
		`"servers":["localhost:8112"]}`
	c, err := NewConfig()
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal([]byte(s), c)
	if err != nil {
		t.Fatal(err)
	}

	r, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}

	if string(r) != s {
		t.Errorf("Expected JSON string: %v, got: %v", s, string(r))
	}

	s = `{"dial_timeout":"100s","discover":true,` +
		`"servers":["localhost:8112"],"timeout":"1s","watch_interval":"5s"}`
	err = json.Unmarshal([]byte(s), c)
	if err == nil || !strings.Contains(err.Error(),
		"invalid dial timeout value") {
		t.Error("Expected error not returned.")
	}

	s = `{"dial_timeout":"aa","discover":true,` +
		`"servers":["localhost:8112"],"timeout":"1s","watch_interval":"5s"}`
	err = json.Unmarshal([]byte(s), c)
	if err == nil || !strings.Contains(err.Error(),
		"unable to parse dial timeout") {
		t.Error("Expected error not returned.")
	}

	s = `{"dial_timeout":"100ms","discover":true,` +
		`"servers":["localhost:8112"],"timeout":"38h","watch_interval":"5s"}`
	err = json.Unmarshal([]byte(s), c)
	if err == nil || !strings.Contains(err.Error(),
		"invalid timeout value") {
		t.Error("Expected error not returned.")
	}

	s = `{"dial_timeout":"100ms","discover":true,` +
		`"servers":["localhost:8112"],"timeout":"aa","watch_interval":"5s"}`
	err = json.Unmarshal([]byte(s), c)
	if err == nil || !strings.Contains(err.Error(),
		"unable to parse timeout") {
		t.Error("Expected error not returned.")
	}

	s = `{"dial_timeout":"100ms","discover":true,` +
		`"servers":["localhost:8112"],"timeout":"10s","watch_interval":"500h"}`
	err = json.Unmarshal([]byte(s), c)
	if err == nil || !strings.Contains(err.Error(),
		"invalid watch interval value") {
		t.Error("Expected error not returned.")
	}

	s = `{"dial_timeout":"100ms","discover":true,` +
		`"servers":["localhost:8112"],"timeout":"10s","watch_interval":"aa"}`
	err = json.Unmarshal([]byte(s), c)
	if err == nil || !strings.Contains(err.Error(),
		"unable to parse watch interval") {
		t.Error("Expected error not returned.")
	}

	s = `{"dial_timeout":"100ms","discover":true,` +
		`"connect_retries":-1,"servers":[":invalid"],"timeout":"10s",` +
		`"watch_interval":"30s"}`
	err = json.Unmarshal([]byte(s), c)
	if err == nil || !strings.Contains(err.Error(),
		"invalid server address") {
		t.Error("Expected error not returned.")
	}

	s = `{"dial_timeout":"100ms","discover":true,"retries":"invalid",` +
		`"connect_retries":"invalid","servers":[":invalid"],"timeout":"10s",` +
		`"watch_interval":"30s"}`
	err = json.Unmarshal([]byte(s), c)
	if err == nil || !strings.Contains(err.Error(),
		"unmarshal string into Go struct field .retries of type") {
		t.Errorf("Expected error not returned, got: %v", err)
	}

	s = `{"dial_timeout":"100ms","discover":true,` +
		`"connect_retries":"invalid","servers":[":invalid"],"timeout":"10s",` +
		`"watch_interval":"30s"}`
	err = json.Unmarshal([]byte(s), c)
	if err == nil || !strings.Contains(err.Error(),
		"unmarshal string into Go struct field .connect_retries of type") {
		t.Errorf("Expected error not returned, got: %v", err)
	}

	s = `{$$$}`
	err = c.UnmarshalJSON([]byte(s))
	if err == nil || !strings.Contains(err.Error(),
		"unable to unmarshal JSON data") {
		t.Error("Expected error not returned.")
	}
}
