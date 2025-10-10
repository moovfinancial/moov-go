package moov

import (
	"context"

	"golang.org/x/time/rate"
)

func WithRateLimit(rps int) ClientConfigurable {
	return func(c *Client) error {
		c.rateLimiter = rate.NewLimiter(rate.Limit(rps), 1)
	}
}

func (c *Client) waitForSlot(ctx context.Context) {
	if c == nil || c.rateLimiter == nil {
		return
	}

	c.rateLimiter.Wait(ctx)
}
