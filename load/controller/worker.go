package controller

import (
	"github.com/Fantom-foundation/Norma/load/generator"
	"log"
)

type Worker struct {
	generator generator.TransactionGenerator
	trigger   <-chan struct{}
}

func (w *Worker) Run() {
	for range w.trigger {
		err := w.generator.SendTx()
		if err != nil {
			log.Printf("failed to send generated tx; %v", err)
		}
	}
	err := w.generator.Close()
	if err != nil {
		log.Printf("failed to close generator; %v", err)
	}
}
