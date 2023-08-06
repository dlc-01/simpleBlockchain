package mining

import (
	"context"
	"fmt"
	"github.con/dlc-01/simpleBloclchain/internal/block"
	"github.con/dlc-01/simpleBloclchain/internal/txs"
	"strings"
	"time"
)

func miner(ctx context.Context, res chan block.Block, b block.Block) {
	nulls := block.GetNNUlls()
	start := time.Now()
	for !strings.HasPrefix(b.HashNow, strings.Repeat("0", nulls)) || nulls == 0 {
		b.CalculateHash()
		if nulls == 0 {
			break
		}
	}
	b.TimeGen = time.Now().Sub(start).Seconds()
	select {
	case <-ctx.Done():
		break
	default:
		res <- b
	}
}

func Mining(ctx context.Context, blocks []block.Block) {
	for i := 1; i < 5; i++ {
		tx := txs.GetTxsData()

		newCtx, cancel := context.WithCancel(ctx)

		res := make(chan block.Block)
		b := block.Block{
			ID:           i + 1,
			Timestamp:    time.Now().UnixNano(),
			HashPrevious: blocks[i-1].HashNow,
		}
		for j := 0; j < 10; j++ {
			b.MinerID = j + 1
			go miner(newCtx, res, b)
		}
		for v := range res {
			cancel()
			blocks[i] = v
			break
		}

		coinbase := txs.Transaction{FromU: "Blockchain", Amount: 100, ToU: fmt.Sprintf("miner%v", blocks[i].MinerID)}
		coinbase.ConfirmTx()
		blocks[i].Transactions = append([]txs.Transaction{coinbase}, tx...)

		blocks[i].PrintBlock()

	}
}
