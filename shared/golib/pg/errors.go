package pg

import "errors"

// Configuration errors
var (
	ErrMissingHost            = errors.New("host is required")
	ErrMissingDatabase        = errors.New("database is required")
	ErrMissingUsername        = errors.New("username is required")
	ErrMissingPassword        = errors.New("password is required")
	ErrInvalidPort            = errors.New("port must be between 1 and 65535")
	ErrInvalidMaxConns        = errors.New("max connections must be positive")
	ErrInvalidMinConns        = errors.New("min connections cannot be negative")
	ErrMinConnsGreaterThanMax = errors.New("min connections cannot be greater than max connections")
	ErrInvalidSSLMode         = errors.New("invalid SSL mode")
)

// Operation errors
var (
	ErrInvalidTransaction = errors.New("invalid transaction state")
	ErrRetryExhausted     = errors.New("retry attempts exhausted")
)
