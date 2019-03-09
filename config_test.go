package gosnowth

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig().
		WithDialTimeout(time.Second).
		WithDiscover(true).
		WithServers("test1", "test2").
		WithTimeout(time.Second).
		WithWatchInterval(time.Second)
	if cfg.DialTimeout != time.Second {
		t.Errorf("Expected dial timeout: %v, got: %v",
			time.Second, cfg.DialTimeout)
	}

	if cfg.Discover != true {
		t.Errorf("Expected discover value: %v, got: %v",
			true, cfg.Discover)
	}

	if len(cfg.Servers) != 2 {
		t.Fatalf("Expected servers length: %v, got: %v",
			2, len(cfg.Servers))
	}

	if cfg.Servers[0] != "test1" {
		t.Errorf("Expected server value: %v, got: %v",
			"test1", cfg.Servers[0])
	}

	if cfg.Servers[1] != "test2" {
		t.Errorf("Expected server value: %v, got: %v",
			"test2", cfg.Servers[1])
	}

	if cfg.Timeout != time.Second {
		t.Errorf("Expected timeout: %v, got: %v",
			time.Second, cfg.Timeout)
	}

	if cfg.WatchInterval != time.Second {
		t.Errorf("Expected watch interval: %v, got: %v",
			time.Second, cfg.WatchInterval)
	}
}

func TestConfigMarshalJSON(t *testing.T) {
	s := `{"dial_timeout":"100ms","discover":true,` +
		`"servers":["localhost:8112"],"timeout":"1s","watch_interval":"5s"}`
	c := NewConfig()
	err := json.Unmarshal([]byte(s), c)
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
}
