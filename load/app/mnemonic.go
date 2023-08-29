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
