package app

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

const Mnemonic = "wide noise until sentence come nothing they diagram miracle universe recall fringe"

func NewKeyGenerator(mnemonic string, param1, param2 uint32) (*KeyGenerator, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic")
	}

	seed := bip39.NewSeed(mnemonic, "")
	key, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to create master key; %v", err)
	}

	// m/44'/60'/param1'/param2/...
	path := []uint32{bip32.FirstHardenedChild + 44, bip32.FirstHardenedChild + 60, bip32.FirstHardenedChild + param1, param2}
	for _, ix := range path {
		key, err = key.NewChildKey(ix)
		if err != nil {
			return nil, fmt.Errorf("failed to create child key; %v", err)
		}
	}
	return (*KeyGenerator)(key), nil
}

type KeyGenerator bip32.Key

func (g *KeyGenerator) GeneratePrivateKey(i uint32) (*ecdsa.PrivateKey, error) {
	childKey, err := (*bip32.Key)(g).NewChildKey(i)
	if err != nil {
		return nil, fmt.Errorf("failed to create terminating child key; %v", err)
	}
	return crypto.ToECDSA(childKey.Key)
}
