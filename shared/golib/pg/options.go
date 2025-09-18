package pg

import (
	"time"
)

// options holds all configuration for the PostgreSQL client (unexported)
type options struct {
	// Connection settings
	host     string
	port     int
	database string
	username string
	password string
	sslMode  string

	// Pool settings
	maxConns        int
	minConns        int
	maxConnLifetime time.Duration
	maxConnIdleTime time.Duration

	// Timeout settings
	connectTimeout time.Duration
	queryTimeout   time.Duration

	// Retry settings
	enableRetry      bool
	maxRetryAttempts int
	retryDelay       time.Duration
	maxRetryDelay    time.Duration
}

// Option is a function that configures options
type Option func(*options)

// WithHost sets the database host
func WithHost(host string) Option {
	return func(o *options) {
		o.host = host
	}
}

// WithPort sets the database port
func WithPort(port int) Option {
	return func(o *options) {
		o.port = port
	}
}

// WithDatabase sets the database name
func WithDatabase(database string) Option {
	return func(o *options) {
		o.database = database
	}
}

// WithCredentials sets username and password
func WithCredentials(username, password string) Option {
	return func(o *options) {
		o.username = username
		o.password = password
	}
}

// WithSSLMode sets the SSL mode (disable, allow, prefer, require, verify-ca, verify-full)
func WithSSLMode(sslMode string) Option {
	return func(o *options) {
		o.sslMode = sslMode
	}
}

// WithMaxConnections sets the maximum number of connections in the pool
func WithMaxConnections(maxConns int) Option {
	return func(o *options) {
		o.maxConns = maxConns
	}
}

// WithMinConnections sets the minimum number of connections in the pool
func WithMinConnections(minConns int) Option {
	return func(o *options) {
		o.minConns = minConns
	}
}

// WithConnectionLifetime sets the maximum lifetime of a connection
func WithConnectionLifetime(lifetime time.Duration) Option {
	return func(o *options) {
		o.maxConnLifetime = lifetime
	}
}

// WithConnectionIdleTime sets the maximum idle time of a connection
func WithConnectionIdleTime(idleTime time.Duration) Option {
	return func(o *options) {
		o.maxConnIdleTime = idleTime
	}
}

// WithConnectTimeout sets the connection timeout
func WithConnectTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.connectTimeout = timeout
	}
}

// WithQueryTimeout sets the default query timeout
func WithQueryTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.queryTimeout = timeout
	}
}

// WithRetrySettings enables retry with custom settings
func WithRetrySettings(maxAttempts int, initialDelay, maxDelay time.Duration) Option {
	return func(o *options) {
		o.enableRetry = true
		o.maxRetryAttempts = maxAttempts
		o.retryDelay = initialDelay
		o.maxRetryDelay = maxDelay
	}
}

// WithoutRetry disables automatic retry for operations
func WithoutRetry() Option {
	return func(o *options) {
		o.enableRetry = false
	}
}

// validate checks if the options are valid
func (o *options) validate() error {
	if o.host == "" {
		return ErrMissingHost
	}
	if o.database == "" {
		return ErrMissingDatabase
	}
	if o.username == "" {
		return ErrMissingUsername
	}
	if o.password == "" {
		return ErrMissingPassword
	}
	if o.port < 1 || o.port > 65535 {
		return ErrInvalidPort
	}
	if o.maxConns <= 0 {
		return ErrInvalidMaxConns
	}
	if o.minConns < 0 {
		return ErrInvalidMinConns
	}
	if o.minConns > o.maxConns {
		return ErrMinConnsGreaterThanMax
	}

	validSSLModes := map[string]bool{
		"disable": true, "allow": true, "prefer": true,
		"require": true, "verify-ca": true, "verify-full": true,
	}
	if !validSSLModes[o.sslMode] {
		return ErrInvalidSSLMode
	}

	return nil
}

// defaultOptions returns options with sensible defaults
func defaultOptions() *options {
	return &options{
		port:             5432,
		sslMode:          "prefer",
		maxConns:         25,
		minConns:         5,
		maxConnLifetime:  30 * time.Minute,
		maxConnIdleTime:  15 * time.Minute,
		connectTimeout:   10 * time.Second,
		queryTimeout:     30 * time.Second,
		enableRetry:      true,
		maxRetryAttempts: 3,
		retryDelay:       100 * time.Millisecond,
		maxRetryDelay:    5 * time.Second,
	}
}
