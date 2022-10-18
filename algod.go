package abs

import (
	"context"
	"time"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
)

type algodClient struct {
	ac    *algod.Client
	delay time.Duration
}

func (c *algodClient) Status() models.NodeStatus {
	return retry(func() (models.NodeStatus, error) {
		return c.ac.Status().Do(context.Background())
	}, c.delay)
}

func (c *algodClient) Block(round uint64) types.Block {
	return retry(func() (types.Block, error) {
		return c.ac.Block(round).Do(context.Background())
	}, c.delay)
}

func (c *algodClient) StatusAfterBlock(round uint64) models.NodeStatus {
	return retry(func() (models.NodeStatus, error) {
		return c.ac.StatusAfterBlock(round).Do(context.Background())
	}, c.delay)
}
