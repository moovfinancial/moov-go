package moov

import (
	"context"
	"fmt"

	"golang.org/x/time/rate"
)

// WithRateLimit returns a ClientConfigurable that sets up a rate limiter on the client
// with the specified requests per second (RPS) limit. The limiter allows a burst of 1.
func WithRateLimit(rps int) ClientConfigurable {
	return func(c *Client) error {
		if rps <= 0 {
			return fmt.Errorf("rate limit (rps) must be positive, but was %d", rps)
		}
		c.rateLimiter = rate.NewLimiter(rate.Limit(rps), 1)
		return nil
	}
}

// waitForSlot blocks until a rate limit slot is available, if a rate limiter is configured
// on the client. If the client or rate limiter is nil, it returns immediately without blocking.
func (c *Client) waitForSlot(ctx context.Context) {
	if c == nil || c.rateLimiter == nil {
		return
	}

	c.rateLimiter.Wait(ctx)
}
