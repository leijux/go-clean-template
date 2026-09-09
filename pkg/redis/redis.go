// Package redis implements redis connection.
package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

const (
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Redis -.
type Redis struct {
	connAttempts int
	connTimeout  time.Duration

	Client *redis.Client
}

// New -.
func New(addr, password string, db int, opts ...Option) (*Redis, error) {
	r := &Redis{
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(r)
	}

	r.Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	err := redisotel.InstrumentTracing(r.Client)
	if err != nil {
		return nil, fmt.Errorf("redis - New - redisotel.InstrumentTracing: %w", err)
	}

	for r.connAttempts > 0 {
		err = r.Client.Ping(context.Background()).Err()
		if err == nil {
			break
		}

		log.Printf("Redis is trying to connect, attempts left: %d", r.connAttempts)

		time.Sleep(r.connTimeout)

		r.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("redis - New - connAttempts == 0: %w", err)
	}

	return r, nil
}

// Close -.
func (r *Redis) Close() error {
	if r.Client == nil {
		return nil
	}

	return r.Client.Close()
}
