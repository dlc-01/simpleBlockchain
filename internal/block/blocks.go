package block

import (
	"crypto/sha256"
	"fmt"
	"github.con/dlc-01/simpleBloclchain/internal/txs"
	"math/rand"
	"time"
)

type Block struct {
	ID           int
	Timestamp    int64
	HashPrevious string
	HashNow      string
	MagicNumber  int32
	MinerID      int
	TimeGen      float64
	ChangN       string
	Transactions []txs.Transaction
}

var nNulls int

func GetNNUlls() int {
	return nNulls
}

func (b *Block) CalculateHash() {
	b.MagicNumber = rand.Int31()
	blockData := fmt.Sprintf("%d%d%d%s", b.ID, b.Timestamp, b.MagicNumber, b.HashPrevious)
	sha256Hash := sha256.New()
	sha256Hash.Write([]byte(blockData))
	b.HashNow = fmt.Sprintf("%x", sha256Hash.Sum(nil))
}

func (b *Block) PrintBlock() {
	b.selectDif()
	var tBlock string

	if b.ID == 1 {
		tBlock = "Genesis Block: \n"

	} else {
		tBlock = fmt.Sprintf("Block: \n"+
			"Created by miner%v\n", b.MinerID)

	}
	fmt.Printf(tBlock+
		"Id: %v\n"+
		"Timestamp: %v\n"+
		"Magic number: %v\n"+
		"Hash of the previous block:\n"+
		"%s\n"+
		"Hash of the block:\n"+
		"%s\n"+
		"Block data:\n",

		b.ID, b.Timestamp, b.MagicNumber, b.HashPrevious, b.HashNow)
	printTransaction(b.Transactions)

	fmt.Printf("Block was generating for %.f seconds\n"+
		b.ChangN+
		"\n", b.TimeGen)
}

func printTransaction(txs []txs.Transaction) {
	if len(txs) == 0 {
		fmt.Println("No transactions")
		return
	}
	fLine := fmt.Sprintf("Transaction #1 (Coinbase):\n")
	id := 2

	for i := 0; i < len(txs); i++ {
		datal := fmt.Sprintf("Transaction ID: %s\n"+
			"Public Key: %v", txs[i].ID, txs[i].Public)
		dataP := fmt.Sprintf("")
		if i != 0 {
			fLine = fmt.Sprintf("Transaction #%v:\n", id)
			id++
		}
		if txs[i].Sign != "" {
			dataP = fmt.Sprintf("\nSignature: %s", txs[i].Sign)
		}
		FrTo := fmt.Sprintf("%s sent %v VC to %s\n", txs[i].FromU, txs[i].Amount, txs[i].ToU)
		fmt.Println(fLine + FrTo + datal + dataP)

	}
}

func (b *Block) selectDif() {
	switch {
	case b.TimeGen < 5:
		nNulls++
		b.ChangN = fmt.Sprintf("N was increased to %v\n", nNulls)
	case b.TimeGen > 10:
		nNulls--
		b.ChangN = fmt.Sprintf("N was decreased to %v\n", nNulls)
	default:
		b.ChangN = "N stays the same\n"
	}
}

func GenerateGenius() Block {
	start := time.Now()
	b := Block{
		ID:           1,
		Timestamp:    time.Now().UnixNano(),
		MagicNumber:  rand.Int31(),
		HashPrevious: "0",
	}
	b.CalculateHash()
	b.TimeGen = time.Now().Sub(start).Seconds()
	b.PrintBlock()
	return b
}
