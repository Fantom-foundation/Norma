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
