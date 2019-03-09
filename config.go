package gosnowth

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

// Config values represent configuration information SnowthClient values.
type Config struct {
	// DialTimeout sets the initial connection timeout duration for attempt to
	// connect to IRONdb. The default value is 500 milliseconds.
	DialTimeout time.Duration

	// Discover sets whether the client should attempt to discover other IRONdb
	// nodes in the same cluster as the provided node servers.
	Discover bool

	// Servers is a list of IRONdb node servers to be used by a SnowthClient.
	Servers []string

	// Timeout sets the timeout duration for HTTP requests to IRONdb. The
	// default value is 10 seconds.
	Timeout time.Duration

	// WatchInterval sets the frequency at which a SnowthClient will check for
	// updates to the active status of its nodes if WatchAndUpdate() is called.
	WatchInterval time.Duration
}

// NewConfig creates and initializes a new SnowthClient configuration value.
func NewConfig(servers ...string) *Config {
	return &Config{
		DialTimeout:   time.Duration(500 * time.Millisecond),
		Discover:      false,
		Servers:       servers,
		Timeout:       time.Duration(10 * time.Second),
		WatchInterval: time.Duration(30 * time.Second),
	}
}

// MarshalJSON encodes a Config value into a JSON format byte slice.
func (c *Config) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}
	if c.DialTimeout != 0 {
		m["dial_timeout"] = c.DialTimeout.String()
	}

	m["discover"] = c.Discover
	if c.Timeout != 0 {
		m["timeout"] = c.Timeout.String()
	}

	if c.WatchInterval != 0 {
		m["watch_interval"] = c.WatchInterval.String()
	}

	if len(c.Servers) > 0 {
		m["servers"] = c.Servers
	}

	return json.Marshal(m)
}

// UnmarshalJSON decodes a JSON format byte slice into the Config value.
func (c *Config) UnmarshalJSON(b []byte) error {
	m := map[string]interface{}{}
	var err error
	if err = json.Unmarshal(b, &m); err != nil {
		errors.Wrap(err, "unable to unmarshal JSON data")
		return err
	}

	if v, ok := m["dial_timeout"].(string); ok {
		c.DialTimeout, err = time.ParseDuration(v)
		if err != nil {
			return errors.Wrap(err, "unable to parse dial timeout")
		}
	}

	if v, ok := m["discover"].(bool); ok {
		c.Discover = v
	}

	if v, ok := m["timeout"].(string); ok {
		c.Timeout, err = time.ParseDuration(v)
		if err != nil {
			return errors.Wrap(err, "unable to parse timeout")
		}
	}

	if v, ok := m["watch_interval"].(string); ok {
		c.WatchInterval, err = time.ParseDuration(v)
		if err != nil {
			return errors.Wrap(err, "unable to parse watch interval")
		}
	}

	if v, ok := m["servers"].([]interface{}); ok {
		for _, vv := range v {
			if vs, ok := vv.(string); ok {
				c.Servers = append(c.Servers, vs)
			}
		}
	}

	return nil
}

// WithDialTimeout sets a new dial timeout and returns a pointer to the updated
// configuration value.
func (c *Config) WithDialTimeout(t time.Duration) *Config {
	c.DialTimeout = t
	return c
}

// WithDiscover sets a new discover value and returns a pointer to the updated
// configuration value.
func (c *Config) WithDiscover(d bool) *Config {
	c.Discover = d
	return c
}

// WithTimeout sets a new timeout duration and returns a pointer to the updated
// configuration value.
func (c *Config) WithTimeout(t time.Duration) *Config {
	c.Timeout = t
	return c
}

// WithServers assigns a new list of servers and returns a pointer to the
// updated configuration value.
func (c *Config) WithServers(addrs ...string) *Config {
	c.Servers = addrs
	return c
}

// WithWatchInterval sets a new interval for watch and update functionality and
// returns a pointer to the updated configuration value.
func (c *Config) WithWatchInterval(i time.Duration) *Config {
	c.WatchInterval = i
	return c
}
