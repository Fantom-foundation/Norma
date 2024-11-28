package app

import (
	"fmt"
	contract "github.com/Fantom-foundation/Norma/load/contracts/abi"
	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/Fantom-foundation/go-opera/inter/validatorpk"
	"github.com/Fantom-foundation/go-opera/opera/contracts/sfc"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

// RegisterValidatorNode registers a validator in the SFC contract.
func RegisterValidatorNode(factory RpcClientFactory) (int, error) {
	newValId := 0

	rpcClient, err := factory.DialRandomRpc()
	if err != nil {
		return 0, fmt.Errorf("failed to connect to network: %w", err)
	}

	// get a representation of the deployed contract
	SFCContract, err := contract.NewSFC(sfc.ContractAddress, rpcClient)
	if err != nil {
		return 0, fmt.Errorf("failed to get SFC contract representation; %v", err)
	}

	var lastValId *big.Int
	lastValId, err = SFCContract.LastValidatorID(nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get validator count; %v", err)
	}

	newValId = int(lastValId.Int64()) + 1

	const chainID = 0xfa3
	privateKeyECDSA := evmcore.FakeKey(uint32(newValId))
	txOpts, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, big.NewInt(chainID))
	if err != nil {
		return 0, fmt.Errorf("failed to create txOpts; %v", err)
	}

	txOpts.Value = big.NewInt(0).Mul(big.NewInt(5_000_000), big.NewInt(1_000_000_000_000_000_000)) // 5_000_000 FTM

	validatorPubKey := validatorpk.PubKey{
		Raw:  crypto.FromECDSAPub(&privateKeyECDSA.PublicKey),
		Type: validatorpk.Types.Secp256k1,
	}

	tx, err := SFCContract.CreateValidator(txOpts, validatorPubKey.Bytes())
	if err != nil {
		return 0, fmt.Errorf("failed to create validator; %v", err)
	}

	receipt, err := GetReceipt(tx.Hash(), rpcClient)
	if err != nil {
		return 0, fmt.Errorf("failed to get receipt; %v", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return 0, fmt.Errorf("failed to deploy helper contract: transaction reverted")
	}

	lastValId, err = SFCContract.LastValidatorID(nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get validator count; %v", err)
	}
	if newValId != int(lastValId.Int64()) {
		return 0, fmt.Errorf("failed to create validator %d", newValId)
	}

	return newValId, nil
}
