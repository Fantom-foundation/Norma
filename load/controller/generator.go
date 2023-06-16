package controller

import (
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/load/app"
	"log"
)

func runGeneratorLoop(generator app.TransactionGenerator, trigger <-chan struct{}, network driver.Network) {
	for range trigger {
		tx, err := generator.GenerateTx()
		if err != nil {
			log.Printf("failed to generate tx; %v", err)
		} else {
			network.SendTransaction(tx)
		}
	}
}
