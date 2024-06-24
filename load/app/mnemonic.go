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
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// Mnemonic to use to generate private keys of all accounts used by Norma.
const Mnemonic = "wide noise until sentence come nothing they diagram miracle universe recall fringe"

// NewKeyGenerator creates a KeyGenerator providing a sequence of private keys,
// based on a mnemonic phrase and feederId/appId, which will be used as BIP-32 key path.
func NewKeyGenerator(mnemonic string, feederId, appId uint32) (*KeyGenerator, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic")
	}

	seed := bip39.NewSeed(mnemonic, "")
	key, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to create master key; %v", err)
	}

	// m/44'/60'/feederId'/appId/...
	path := []uint32{bip32.FirstHardenedChild + 44, bip32.FirstHardenedChild + 60, bip32.FirstHardenedChild + feederId, appId}
	for _, ix := range path {
		key, err = key.NewChildKey(ix)
		if err != nil {
			return nil, fmt.Errorf("failed to create child key; %v", err)
		}
	}
	return (*KeyGenerator)(key), nil
}

// KeyGenerator generates a sequence of private keys based on a BIP-32 key.
type KeyGenerator bip32.Key

func (g *KeyGenerator) GeneratePrivateKey(i uint32) (*ecdsa.PrivateKey, error) {
	childKey, err := (*bip32.Key)(g).NewChildKey(i)
	if err != nil {
		return nil, fmt.Errorf("failed to create terminating child key; %v", err)
	}
	return crypto.ToECDSA(childKey.Key)
}
