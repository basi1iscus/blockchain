package main

import (
	"blockchain_demo/pkg/blockchain"
	"blockchain_demo/pkg/sign"
	"blockchain_demo/pkg/transaction"
	"blockchain_demo/pkg/transaction/coin_transfer"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	source := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(source)

	var addresses = [3]string{"1234567890abcdef1234567890abcdef12345678", "abcdef1234567890abcdef1234567890abcdef12", "2345678901abcdef2345678901abcdef23456789"}
	var signatures = [3]*sign.Signature{}
	for i := 0; i < 3; i++ {
		signatures[i], _ = sign.GenerateKeyPair()
	}
	creator := "ad23947398423423cd234fe34345345323423423"
	bc, _ := blockchain.NewBlockchain(50, 8, creator)

	var blockCount = rnd.Intn(10)
	for range blockCount {
		var count = rnd.Intn(10)
		for range count {
			senderInd := rnd.Intn(len(addresses))
			reciverInd := rnd.Intn(len(addresses))
			tx, _ := transaction.CreateTransaction(coin_transfer.CoinTransfer, addresses[senderInd], rnd.Int63n(1000), rnd.Int63n(10), map[string]any{
				"recipient": addresses[reciverInd],
			})
			tx.AddSing(signatures[senderInd])
			bc.AddTransactionToPool(tx)
		}

		bc.MineBlockFromPool(creator)

		err := bc.Verify(4)
		if err != nil {
			fmt.Printf("Blockchain verification failed: %v\n", err)
		}
	}

	fmt.Println(bc)
}
