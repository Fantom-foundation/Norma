// Copyright 2024 Fantom Foundation
// This file is part of Norma System Testing Infrastructure for Sonic.
//
// Norma is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Norma is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Norma. If not, see <http://www.gnu.org/licenses/>.

package controller

import (
	"log"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/load/app"
)

func runGeneratorLoop(user app.User, trigger <-chan struct{}, network driver.Network) {
	rpcClient, err := network.DialRandomRpc()
	if err != nil {
		log.Fatalf("generator loop; failed to dial random rpc; %v", err)
		return
	}
	defer rpcClient.Close()
	for range trigger {
		if err := user.SendTransaction(rpcClient); err != nil {
			log.Printf("generator loop: failed to send tx; %v", err)
			continue
		}
	}
}
