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

package app

import (
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestMnemonic(t *testing.T) {
	keyGen, err := NewKeyGenerator(Mnemonic, 0, 0)
	if err != nil {
		t.Fatalf("failed to create key generator; %v", err)
	}

	privateKey1, err := keyGen.GeneratePrivateKey(0)
	if err != nil {
		t.Fatalf("failed to create key 0; %v", err)
	}

	address1 := crypto.PubkeyToAddress(privateKey1.PublicKey)
	if address1.String() != "0x333314e70012Fed4bfA14FbEd2D5F0075db00652" {
		t.Fatalf("address of key 0 does not match: %s != %s", address1, "0x333314e70012Fed4bfA14FbEd2D5F0075db00652")
	}
}
