package kit

import (
	"context"
	"net"
	"time"
)

type DialOptions struct {
	Timeout         time.Duration
	MaxRetries      int
	RetryBackoff    time.Duration
	MaxRetryBackoff time.Duration
}

const (
	defaultRetryBackoff    = 250 * time.Millisecond
	defaultMaxRetryBackoff = 10 * time.Second
)

func (d DialOptions) Backoff(backoff time.Duration) time.Duration {
	max, step := d.MaxRetryBackoff, d.RetryBackoff
	if step <= 0 {
		step = defaultRetryBackoff
		if max <= 0 {
			max = defaultMaxRetryBackoff
		}
	}
	if max > 0 {
		if backoff == 0 {
			return step
		}
		return backoff * 2
	}
	return backoff + step
}

func (d DialOptions) DialContext(ctx context.Context, network, address string) (conn net.Conn, err error) {
	var (
		backoff time.Duration
		retries = d.MaxRetries
	)
	if ctx == nil {
		ctx = context.Background()
	}
	for {
		if conn, err = net.DialTimeout(network, address, d.Timeout); err == nil || retries == 0 {
			return
		}
		switch e := err.(type) {
		case net.Error:
			switch {
			case e.Temporary():
			case e.Timeout():
			default:
				return
			}
		default:
			return
		}

		retries--
		backoff = d.Backoff(backoff)
		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
