package redis

import "time"

// Option -.
type Option func(*Redis)

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(c *Redis) {
		c.connAttempts = attempts
	}
}

// ConnTimeout -.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *Redis) {
		c.connTimeout = timeout
	}
}
