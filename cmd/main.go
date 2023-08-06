package main

import (
	"context"
	"github.con/dlc-01/simpleBloclchain/internal/block"
	"github.con/dlc-01/simpleBloclchain/internal/mining"
)

func main() {
	blockchain := make([]block.Block, 5)
	ctx := context.Background()
	blockchain[0] = block.GenerateGenius()

	mining.Mining(ctx, blockchain)

}
