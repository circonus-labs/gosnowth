package gosnowth

import (
	"encoding/json"
	"net/url"
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
	c := &Config{
		DialTimeout:   time.Duration(500 * time.Millisecond),
		Discover:      false,
		Servers:       []string{},
		Timeout:       time.Duration(10 * time.Second),
		WatchInterval: time.Duration(30 * time.Second),
	}

	for _, svr := range servers {
		if _, err := url.Parse(svr); err == nil {
			c.Servers = append(c.Servers, svr)
		}
	}

	return c
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
	if err := json.Unmarshal(b, &m); err != nil {
		return errors.Wrap(err, "unable to unmarshal JSON data")
	}

	if v, ok := m["dial_timeout"].(string); ok {
		d, err := time.ParseDuration(v)
		if err != nil {
			return errors.Wrap(err, "unable to parse dial timeout")
		}

		if d < 0 || d > time.Minute {
			return errors.New("invalid dial timeout value")
		}

		c.DialTimeout = d
	}

	if v, ok := m["discover"].(bool); ok {
		c.Discover = v
	}

	if v, ok := m["timeout"].(string); ok {
		d, err := time.ParseDuration(v)
		if err != nil {
			return errors.Wrap(err, "unable to parse timeout")
		}

		if d < 0 || d > (5*time.Minute) {
			return errors.New("invalid timeout value")
		}

		c.Timeout = d
	}

	if v, ok := m["watch_interval"].(string); ok {
		d, err := time.ParseDuration(v)
		if err != nil {
			return errors.Wrap(err, "unable to parse watch interval")
		}

		if d < 0 || d > (24*time.Hour) {
			return errors.New("invalid watch interval value")
		}

		c.WatchInterval = d
	}

	if v, ok := m["servers"].([]interface{}); ok {
		for _, vv := range v {
			if svr, ok := vv.(string); ok {
				if _, err := url.Parse(svr); err == nil {
					c.Servers = append(c.Servers, svr)
				}
			}
		}
	}

	return nil
}

// WithDialTimeout sets a new dial timeout and returns a pointer to the updated
// configuration value.
func (c *Config) WithDialTimeout(t time.Duration) *Config {
	if t >= 0 && t < time.Minute {
		c.DialTimeout = t
	}

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
	if t >= 0 && t < (5*time.Minute) {
		c.Timeout = t
	}

	return c
}

// WithServers assigns a new list of servers and returns a pointer to the
// updated configuration value.
func (c *Config) WithServers(servers ...string) *Config {
	c.Servers = []string{}
	for _, svr := range servers {
		if _, err := url.Parse(svr); err == nil {
			c.Servers = append(c.Servers, svr)
		}
	}

	return c
}

// WithWatchInterval sets a new interval for watch and update functionality and
// returns a pointer to the updated configuration value.
func (c *Config) WithWatchInterval(i time.Duration) *Config {
	if i >= 0 && i <= (time.Hour*24) {
		c.WatchInterval = i
	}

	return c
}
