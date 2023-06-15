package controller

import (
	"github.com/Fantom-foundation/Norma/load/app"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
)

func runGeneratorLoop(generator app.TransactionGenerator, trigger <-chan struct{}, output chan<- *types.Transaction) {
	for range trigger {
		tx, err := generator.GenerateTx()
		if err != nil {
			log.Printf("failed to generate tx; %v", err)
		}
		output <- tx
	}
}
