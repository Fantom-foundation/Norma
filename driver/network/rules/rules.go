package rules

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Fantom-foundation/Norma/common/transact"
	"github.com/Fantom-foundation/Norma/driver/network/rules/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var nodeDriverAuthAddress = common.HexToAddress("0xd100ae0000000000000000000000000000000000")

type NetworkRules struct {
	Name      string        `json:",omitempty"`
	NetworkID uint64        `json:",omitempty"`
	Dag       *DagRules     `json:",omitempty"`
	Epochs    *EpochsRules  `json:",omitempty"`
	Blocks    *BlocksRules  `json:",omitempty"`
	Economy   *EconomyRules `json:",omitempty"`
	Upgrades  *Upgrades     `json:",omitempty"`
}

type DagRules struct {
	MaxParents     uint32 `json:",omitempty"`
	MaxFreeParents uint32 `json:",omitempty"` // maximum number of parents with no gas cost
	MaxExtraData   uint32 `json:",omitempty"`
}

type EpochsRules struct {
	MaxEpochGas      uint64 `json:",omitempty"`
	MaxEpochDuration uint64 `json:",omitempty"`
}

type BlocksRules struct {
	MaxBlockGas             uint64 `json:",omitempty"` // technical hard limit, gas is mostly governed by gas power allocation
	MaxEmptyBlockSkipPeriod uint64 `json:",omitempty"`
}

type EconomyRules struct {
	BlockMissedSlack uint64         `json:",omitempty"`
	Gas              *GasRulesRLPV1 `json:",omitempty"`
	MinGasPrice      *big.Int       `json:",omitempty"` // updated by Opera in runtime
	ShortGasPower    *GasPowerRules `json:",omitempty"`
	LongGasPower     *GasPowerRules `json:",omitempty"`
}

type GasRulesRLPV1 struct {
	MaxEventGas  uint64 `json:",omitempty"`
	EventGas     uint64 `json:",omitempty"`
	ParentGas    uint64 `json:",omitempty"`
	ExtraDataGas uint64 `json:",omitempty"`
	// Post-LLR fields
	BlockVotesBaseGas    uint64 `json:",omitempty"`
	BlockVoteGas         uint64 `json:",omitempty"`
	EpochVoteGas         uint64 `json:",omitempty"`
	MisbehaviourProofGas uint64 `json:",omitempty"`
}

type GasPowerRules struct {
	AllocPerSec        uint64 `json:",omitempty"`
	MaxAllocPeriod     uint64 `json:",omitempty"`
	StartupAllocPeriod uint64 `json:",omitempty"`
	MinStartupGas      uint64 `json:",omitempty"`
}

type Upgrades struct {
	Berlin bool
	London bool
	Llr    bool
}

func SetNetworkRules(rpcClient transact.RpcClient, ownerAccount *transact.Account, rulesDiff NetworkRules) error {
	rulesDiffJson, err := json.Marshal(rulesDiff)
	if err != nil {
		return fmt.Errorf("failed to marshal NetworkRules into JSON; %v", err)
	}

	// get price of gas from the network
	regularGasPrice, err := rpcClient.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to suggest gas price; %v", err)
	}

	authContract, err := abi.NewNodeDriverAuth(nodeDriverAuthAddress, rpcClient)
	if err != nil {
		return fmt.Errorf("failed to get NodeDriverAuth contract representation; %v", err)
	}
	txOpts, err := bind.NewKeyedTransactorWithChainID(ownerAccount.PrivateKey, ownerAccount.ChainID)
	if err != nil {
		return fmt.Errorf("failed to create txOpts; %v", err)
	}
	txOpts.GasPrice = transact.GetPriorityGasPrice(regularGasPrice)
	txOpts.Nonce = big.NewInt(int64(ownerAccount.GetNextNonce()))

	_, err = authContract.UpdateNetworkRules(txOpts, rulesDiffJson)
	if err != nil {
		return fmt.Errorf("failed to update network rules; %v", err)
	}
	return nil
}
