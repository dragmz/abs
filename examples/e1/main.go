package main

import (
	"abs"
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/pkg/errors"
)

type args struct {
	Algod      string
	AlgodToken string

	Round uint64
}

func run(a args) error {
	ac, err := algod.MakeClient(a.Algod, a.AlgodToken)
	if err != nil {
		return errors.Wrap(err, "failed to make algod client")
	}

	blocks, err := abs.MakeBlocks(ac, abs.WithRetry(time.Second))
	if err != nil {
		return errors.Wrap(err, "failed to make blocks")
	}

	ch := make(chan types.Block)

	go func() {
		for b := range ch {
			fmt.Printf("Block: %d at %s\n", b.Round, time.Now())
		}
	}()

	err = blocks.Stream(context.Background(), a.Round, ch)
	if err != nil {
		return errors.Wrap(err, "failed to stream blocks")
	}

	return nil
}

func main() {
	var a args

	flag.StringVar(&a.Algod, "algod", "http://mainnet-api.algonode.network", "algod address")
	flag.StringVar(&a.AlgodToken, "algod-token", "", "algod token")
	flag.Uint64Var(&a.Round, "round", 0, "start round")

	flag.Parse()

	err := run(a)
	if err != nil {
		panic(err)
	}
}
