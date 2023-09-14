package abs

import (
	"context"
	"time"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

type Blocks struct {
	c *algodClient

	delay time.Duration
}

type BlocksOption func(b *Blocks) error

func WithRetry(delay time.Duration) BlocksOption {
	return func(b *Blocks) error {
		b.delay = delay
		return nil
	}
}

func MakeBlocks(ac *algod.Client, opts ...BlocksOption) (*Blocks, error) {
	b := &Blocks{}

	for _, opt := range opts {
		err := opt(b)
		if err != nil {
			return nil, err
		}
	}

	b.c = &algodClient{ac: ac, delay: b.delay}

	return b, nil
}

func (b *Blocks) Stream(ctx context.Context, current uint64, blocks chan<- types.Block) error {
	status := b.c.Status()

	last := status.LastRound
	if current == 0 {
		current = last
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		for current <= last {
			block := b.c.Block(current)
			blocks <- block

			current++
		}

		status = b.c.StatusAfterBlock(last)
		last = status.LastRound
	}
}
