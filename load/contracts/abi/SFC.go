package abi

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// SFCMetaData contains all meta data concerning the SFC contract.
var SFCMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[],\"name\":\"currentSealedEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"getEpochSnapshot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epochFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalBaseRewardWeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalTxRewardWeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_baseRewardPerSecond\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSupply\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"getLockupInfo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"lockedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fromEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"getStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"getStashedLockupRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"lockupExtraReward\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockupBaseReward\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unlockedReward\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"getValidator\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"status\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivatedTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivatedEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"receivedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdTime\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"auth\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"getValidatorID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"getValidatorPubkey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"getWithdrawalRequest\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"lastValidatorID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minGasPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"slashingRefundRatio\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"stakeTokenizerAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"stashedRewardsUntilEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalActiveStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSlashedStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"treasuryAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"bytes3\",\"name\":\"\",\"type\":\"bytes3\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"v\",\"type\":\"address\"}],\"name\":\"updateConstsAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"constsAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"getEpochValidatorIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"}],\"name\":\"getEpochReceivedStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"}],\"name\":\"getEpochAccumulatedRewardPerToken\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"}],\"name\":\"getEpochAccumulatedUptime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"}],\"name\":\"getEpochAccumulatedOriginatedTxsFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"}],\"name\":\"getEpochOfflineTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"}],\"name\":\"getEpochOfflineBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"}],\"name\":\"rewardsStash\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"toValidatorID\",\"type\":\"uint256\"}],\"name\":\"getLockedStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"}],\"name\":\"createValidator\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"}],\"name\":\"getSelfStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"toValidatorID\",\"type\":\"uint256\"}],\"name\":\"delegate\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"toValidatorID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"wrID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"undelegate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"}],\"name\":\"isSlashed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"toValidatorID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"wrID\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"status\",\"type\":\"uint256\"}],\"name\":\"deactivateValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"toValidatorID\",\"type\":\"uint256\"}],\"name\":\"pendingRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"toValidatorID\",\"type\":\"uint256\"}],\"name\":\"stashRewards\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"toValidatorID\",\"type\":\"uint256\"}],\"name\":\"claimRewards\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"toValidatorID\",\"type\":\"uint256\"}],\"name\":\"restakeRewards\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"updateBaseRewardPerSecond\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blocksNum\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"}],\"name\":\"updateOfflinePenaltyThreshold\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"refundRatio\",\"type\":\"uint256\"}],\"name\":\"updateSlashingRefundRatio\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"updateStakeTokenizerAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"v\",\"type\":\"address\"}],\"name\":\"updateTreasuryAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burnFTM\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"offlineTime\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"offlineBlocks\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"uptimes\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"originatedTxsFee\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"epochGas\",\"type\":\"uint256\"}],\"name\":\"sealEpoch\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"nextValidatorIDs\",\"type\":\"uint256[]\"}],\"name\":\"sealEpochValidators\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"toValidatorID\",\"type\":\"uint256\"}],\"name\":\"isLockedUp\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"toValidatorID\",\"type\":\"uint256\"}],\"name\":\"getUnlockedStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"toValidatorID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockupDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"lockStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"toValidatorID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockupDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"relockStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"toValidatorID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"unlockStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"sealedEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_totalSupply\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"nodeDriver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"lib\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"consts\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"auth\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"status\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivatedEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivatedTime\",\"type\":\"uint256\"}],\"name\":\"setGenesisValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"toValidatorID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockupFromEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockupEndTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockupDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"earlyUnlockPenalty\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rewards\",\"type\":\"uint256\"}],\"name\":\"setGenesisDelegation\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"v\",\"type\":\"address\"}],\"name\":\"updateVoteBookAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"voteBookAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"toValidatorID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"liquidateSFTM\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"v\",\"type\":\"address\"}],\"name\":\"updateSFTMFinalizer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052600436106103765760003560e01c80638b0e9f3f116101d1578063c7be95de11610102578063d96ed505116100a0578063e261641a1161006f578063e261641a146110a2578063e6f45adf146110d2578063e9a505a714611105578063f2fde38b1461111a57610376565b8063d96ed50514610fb0578063dc31e1af14610fc5578063df00c92214610ff5578063e08d7e661461102557610376565b8063cfd47663116100dc578063cfd4766314610ef6578063cfdbb7cd14610f2f578063d46fa51814610f68578063d725e91f14610f7d57610376565b8063c7be95de14610e74578063cc17278414610e89578063cc8343aa14610ec457610376565b8063a86a056f1161016f578063b810e41111610149578063b810e41114610d82578063b88a37e214610dbb578063c5f956af14610e35578063c65ee0e114610e4a57610376565b8063a86a056f14610cab578063b0ef386c14610ce4578063b5d8962714610d1757610376565b806396c7ee46116101ab57806396c7ee4614610b91578063a198d22914610bf0578063a2f6e6bc14610c20578063a45e115414610c5357610376565b80638b0e9f3f14610b3e5780638da5cb5b14610b535780638f32d59b14610b6857610376565b80635fab23a8116102ab5780637667180811610249578063854873e111610223578063854873e1146109da578063860c275014610a79578063873571d214610aac578063893675c614610b2957610376565b8063766718081461097d5780637cacb1d614610992578063841e4561146109a757610376565b80636a225ec2116102855780636a225ec2146108cc5780636f498663146108fc578063715018a614610935578063736de9ae1461094a57610376565b80635fab23a81461084e57806361e53fcc14610863578063670322f81461089357610376565b806339b80c001161031857806354fd4d50116102f257806354fd4d501461062e578063550359a01461067857806358f95b80146106ab578063592fe0c0146106db57610376565b806339b80c0014610566578063468f35ee146105c857806346f1ca35146105fb57610376565b806318160ddd1161035457806318160ddd146104ac5780631f270152146104c157806328f731481461051e5780632ce719601461053357610376565b80630135b1db146103df5780630e559d821461042457806310e51e1414610455575b366103c8576040805162461bcd60e51b815260206004820152601560248201527f7472616e7366657273206e6f7420616c6c6f7765640000000000000000000000604482015290519081900360640190fd5b6080546103dd906001600160a01b031661114d565b005b3480156103eb57600080fd5b506104126004803603602081101561040257600080fd5b50356001600160a01b0316611176565b60408051918252519081900360200190f35b34801561043057600080fd5b50610439611188565b604080516001600160a01b039092168252519081900360200190f35b34801561046157600080fd5b506103dd600480360360c081101561047857600080fd5b508035906020810135906001600160a01b0360408201358116916060810135821691608082013581169160a0013516611197565b3480156104b857600080fd5b50610412611324565b3480156104cd57600080fd5b50610500600480360360608110156104e457600080fd5b506001600160a01b03813516906020810135906040013561132a565b60408051938452602084019290925282820152519081900360600190f35b34801561052a57600080fd5b5061041261135c565b34801561053f57600080fd5b506103dd6004803603602081101561055657600080fd5b50356001600160a01b0316611362565b34801561057257600080fd5b506105906004803603602081101561058957600080fd5b50356113f5565b604080519788526020880196909652868601949094526060860192909252608085015260a084015260c0830152519081900360e00190f35b3480156105d457600080fd5b50610439600480360360208110156105eb57600080fd5b50356001600160a01b0316611437565b34801561060757600080fd5b506103dd6004803603602081101561061e57600080fd5b50356001600160a01b0316611452565b34801561063a57600080fd5b5061064361148b565b604080517fffffff00000000000000000000000000000000000000000000000000000000009092168252519081900360200190f35b34801561068457600080fd5b506103dd6004803603602081101561069b57600080fd5b50356001600160a01b03166114b0565b3480156106b757600080fd5b50610412600480360360408110156106ce57600080fd5b5080359060200135611543565b3480156106e757600080fd5b506103dd600480360360a08110156106fe57600080fd5b81019060208101813564010000000081111561071957600080fd5b82018360208201111561072b57600080fd5b8035906020019184602083028401116401000000008311171561074d57600080fd5b91939092909160208101903564010000000081111561076b57600080fd5b82018360208201111561077d57600080fd5b8035906020019184602083028401116401000000008311171561079f57600080fd5b9193909290916020810190356401000000008111156107bd57600080fd5b8201836020820111156107cf57600080fd5b803590602001918460208302840111640100000000831117156107f157600080fd5b91939092909160208101903564010000000081111561080f57600080fd5b82018360208201111561082157600080fd5b8035906020019184602083028401116401000000008311171561084357600080fd5b919350915035611564565b34801561085a57600080fd5b5061041261181f565b34801561086f57600080fd5b506104126004803603604081101561088657600080fd5b5080359060200135611825565b34801561089f57600080fd5b50610412600480360360408110156108b657600080fd5b506001600160a01b038135169060200135611846565b3480156108d857600080fd5b506103dd600480360360408110156108ef57600080fd5b5080359060200135611887565b34801561090857600080fd5b506104126004803603604081101561091f57600080fd5b506001600160a01b038135169060200135611a10565b34801561094157600080fd5b506103dd611a8e565b34801561095657600080fd5b506104396004803603602081101561096d57600080fd5b50356001600160a01b0316611b49565b34801561098957600080fd5b50610412611b64565b34801561099e57600080fd5b50610412611b6d565b3480156109b357600080fd5b506103dd600480360360208110156109ca57600080fd5b50356001600160a01b0316611b73565b3480156109e657600080fd5b50610a04600480360360208110156109fd57600080fd5b5035611c06565b6040805160208082528351818301528351919283929083019185019080838360005b83811015610a3e578181015183820152602001610a26565b50505050905090810190601f168015610a6b5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b348015610a8557600080fd5b506103dd60048036036020811015610a9c57600080fd5b50356001600160a01b0316611cbf565b348015610ab857600080fd5b506103dd60048036036020811015610acf57600080fd5b810190602081018135640100000000811115610aea57600080fd5b820183602082011115610afc57600080fd5b80359060200191846001830284011164010000000083111715610b1e57600080fd5b509092509050611d52565b348015610b3557600080fd5b5061043961218a565b348015610b4a57600080fd5b50610412612199565b348015610b5f57600080fd5b5061043961219f565b348015610b7457600080fd5b50610b7d6121ae565b604080519115158252519081900360200190f35b348015610b9d57600080fd5b50610bca60048036036040811015610bb457600080fd5b506001600160a01b0381351690602001356121bf565b604080519485526020850193909352838301919091526060830152519081900360800190f35b348015610bfc57600080fd5b5061041260048036036040811015610c1357600080fd5b50803590602001356121f1565b348015610c2c57600080fd5b506103dd60048036036020811015610c4357600080fd5b50356001600160a01b0316612212565b348015610c5f57600080fd5b50610c9260048036036060811015610c7657600080fd5b506001600160a01b0381351690602081013590604001356122a5565b6040805192835260208301919091528051918290030190f35b348015610cb757600080fd5b5061041260048036036040811015610cce57600080fd5b506001600160a01b0381351690602001356122ec565b348015610cf057600080fd5b506103dd60048036036020811015610d0757600080fd5b50356001600160a01b0316612309565b348015610d2357600080fd5b50610d4160048036036020811015610d3a57600080fd5b5035612401565b604080519788526020880196909652868601949094526060860192909252608085015260a08401526001600160a01b031660c0830152519081900360e00190f35b348015610d8e57600080fd5b5061050060048036036040811015610da557600080fd5b506001600160a01b038135169060200135612447565b348015610dc757600080fd5b50610de560048036036020811015610dde57600080fd5b5035612473565b60408051602080825283518183015283519192839290830191858101910280838360005b83811015610e21578181015183820152602001610e09565b505050509050019250505060405180910390f35b348015610e4157600080fd5b506104396124d9565b348015610e5657600080fd5b5061041260048036036020811015610e6d57600080fd5b50356124e8565b348015610e8057600080fd5b506104126124fa565b348015610e9557600080fd5b506103dd60048036036040811015610eac57600080fd5b506001600160a01b0381358116916020013516612500565b348015610ed057600080fd5b506103dd60048036036040811015610ee757600080fd5b5080359060200135151561267e565b348015610f0257600080fd5b5061041260048036036040811015610f1957600080fd5b506001600160a01b0381351690602001356128ad565b348015610f3b57600080fd5b50610b7d60048036036040811015610f5257600080fd5b506001600160a01b0381351690602001356128ca565b348015610f7457600080fd5b50610439612961565b348015610f8957600080fd5b506103dd60048036036020811015610fa057600080fd5b50356001600160a01b0316612970565b348015610fbc57600080fd5b50610412612a98565b348015610fd157600080fd5b5061041260048036036040811015610fe857600080fd5b5080359060200135612a9e565b34801561100157600080fd5b506104126004803603604081101561101857600080fd5b5080359060200135612abf565b34801561103157600080fd5b506103dd6004803603602081101561104857600080fd5b81019060208101813564010000000081111561106357600080fd5b82018360208201111561107557600080fd5b8035906020019184602083028401116401000000008311171561109757600080fd5b509092509050612ae0565b3480156110ae57600080fd5b50610412600480360360408110156110c557600080fd5b5080359060200135612c24565b3480156110de57600080fd5b506103dd600480360360208110156110f557600080fd5b50356001600160a01b0316612c45565b34801561111157600080fd5b50610439612cd8565b34801561112657600080fd5b506103dd6004803603602081101561113d57600080fd5b50356001600160a01b0316612ce7565b3660008037600080366000845af43d6000803e80801561116c573d6000f35b3d6000fd5b505050565b60696020526000908152604090205481565b607b546001600160a01b031681565b600054610100900460ff16806111b057506111b0612d4c565b806111be575060005460ff16155b6111f95760405162461bcd60e51b815260040180806020018281038252602e81526020018061458f602e913960400191505060405180910390fd5b600054610100900460ff1615801561125f57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff909116610100171660011790555b61126882612d52565b6067879055606680546001600160a01b038088167fffffffffffffffffffffffff00000000000000000000000000000000000000009283161790925560808054878416908316179055608180549286169290911691909117905560768690556112cf612eb4565b607e556112da612ebd565b600088815260776020526040902060070155801561131b57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1690555b50505050505050565b60765481565b607160209081526000938452604080852082529284528284209052825290208054600182015460029092015490919083565b606d5481565b61136a6121ae565b6113bb576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b608380547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b607760205280600052604060002060009150905080600701549080600801549080600901549080600a01549080600b01549080600c01549080600d0154905087565b6088602052600090815260409020546001600160a01b031681565b6040516001600160a01b0382169033907f857125196131cfcd709c738c6d1fd2701ce70f2a03785aeadae6f4b47fe73c1d90600090a350565b7f33303500000000000000000000000000000000000000000000000000000000005b90565b6114b86121ae565b611509576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b608280547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b60008281526077602090815260408083208484529091529020545b92915050565b61156d33612ec1565b6115a85760405162461bcd60e51b81526004018080602001828103825260298152602001806145456029913960400191505060405180910390fd5b6000607760006115b6611b64565b8152602001908152602001600020905060608160060180548060200260200160405190810160405280929190818152602001828054801561161657602002820191906000526020600020905b815481526020019060010190808311611602575b5050505050905061169d82828d8d80806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f820116905080830192505050505050508c8c80806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250612ed592505050565b606754600090815260776020526040902060078101546001906116be612ebd565b11156116d55781600701546116d1612ebd565b0390505b611757818584868d8d80806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f820116905080830192505050505050508c8c808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152506130dc92505050565b61176181866138cc565b505061176b611b64565b606755611776612ebd565b6007830155608154604080517fd9a7c1f900000000000000000000000000000000000000000000000000000000815290516001600160a01b039092169163d9a7c1f991600480820192602092909190829003018186803b1580156117d957600080fd5b505afa1580156117ed573d6000803e3d6000fd5b505050506040513d602081101561180357600080fd5b5051600b83015550607654600d90910155505050505050505050565b606e5481565b60009182526077602090815260408084209284526001909201905290205490565b600061185283836128ca565b61185e5750600061155e565b506001600160a01b03919091166000908152607360209081526040808320938352929052205490565b815b81811015611171576000818152606a602090815260409182902080548351601f60027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6101006001861615020190931692909204918201849004840281018401909452808452606093928301828280156119445780601f1061191957610100808354040283529160200191611944565b820191906000526020600020905b81548152906001019060200180831161192757829003601f168201915b505050505090506000815111801561197657508160866000838051906020012081526020019081526020016000205414155b15611a07576086600082805190602001208152602001908152602001600020546000146119ea576040805162461bcd60e51b815260206004820152600e60248201527f616c726561647920657869737473000000000000000000000000000000000000604482015290519081900360640190fd5b805160208083019190912060009081526086909152604090208290555b50600101611889565b6000611a1a6143e2565b506001600160a01b0383166000908152606f6020908152604080832085845282529182902082516060810184528154808252600183015493820184905260029092015493810184905292611a86929091611a7a919063ffffffff613a4516565b9063ffffffff613a4516565b949350505050565b611a966121ae565b611ae7576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b6033546040516000916001600160a01b0316907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3603380547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055565b6089602052600090815260409020546001600160a01b031681565b60675460010190565b60675481565b611b7b6121ae565b611bcc576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b607f80547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b606a6020908152600091825260409182902080548351601f60027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff61010060018616150201909316929092049182018490048402810184019094528084529091830182828015611cb75780601f10611c8c57610100808354040283529160200191611cb7565b820191906000526020600020905b815481529060010190602001808311611c9a57829003601f168201915b505050505081565b611cc76121ae565b611d18576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b608180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b600160005260686020527f82eaf0fca2207f91f5027fcf68136c84edb7e928c081c42aa5bbc2a771c7d37c546001600160a01b031673541e408443a592c38e01bed0cb31f9de8c1322d014611dee576040805162461bcd60e51b815260206004820152600b60248201527f6e6f74206d61696e6e6574000000000000000000000000000000000000000000604482015290519081900360640190fd5b604281148015611e39575081816000818110611e0657fe5b9050013560f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191660c060f81b145b611e8a576040805162461bcd60e51b815260206004820152601060248201527f6d616c666f726d6564207075626b657900000000000000000000000000000000604482015290519081900360640190fd5b33600090815260696020526040902054603b811115611ef0576040805162461bcd60e51b815260206004820152601460248201527f6e6f74206c65676163792076616c696461746f72000000000000000000000000604482015290519081900360640190fd5b611ef981613a9f565b611f4a576040805162461bcd60e51b815260206004820152601760248201527f76616c696461746f7220646f65736e2774206578697374000000000000000000604482015290519081900360640190fd5b606a60008281526020019081526020016000206040518082805460018160011615610100020316600290048015611fb85780601f10611f96576101008083540402835291820191611fb8565b820191906000526020600020905b815481529060010190602001808311611fa4575b50509150506040518091039020838360405180838380828437808301925050509250505060405180910390201415612037576040805162461bcd60e51b815260206004820152600b60248201527f73616d65207075626b6579000000000000000000000000000000000000000000604482015290519081900360640190fd5b608660008484604051808383808284376040805191909301819003902086525060208501959095525050500160002054156120b9576040805162461bcd60e51b815260206004820152600c60248201527f616c726561647920757365640000000000000000000000000000000000000000604482015290519081900360640190fd5b6000818152608560205260409020541561211a576040805162461bcd60e51b815260206004820152601160248201527f616c6c6f776564206f6e6c79206f6e6365000000000000000000000000000000604482015290519081900360640190fd5b60008181526085602052604080822080546001019055518291608691869086908083838082843760408051919093018190039020865250602080860196909652938401600090812096909655505050838352606a909152902061217e908484614403565b5061117181600161267e565b6082546001600160a01b031681565b606c5481565b6033546001600160a01b031690565b6033546001600160a01b0316331490565b607360209081526000928352604080842090915290825290208054600182015460028301546003909301549192909184565b60009182526077602090815260408084209284526005909201905290205490565b61221a6121ae565b61226b576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b607b80547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b608460205282600052604060002060205281600052604060002081815481106122ca57fe5b6000918252602090912060029091020180546001909101549093509150839050565b607060209081526000928352604080842090915290825290205481565b6123116121ae565b612362576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b6087546001600160a01b03828116911614156123c7576040805162461bcd60e51b8152602060048083019190915260248201527f73616d6500000000000000000000000000000000000000000000000000000000604482015290519081900360640190fd5b608780547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b606860205260009081526040902080546001820154600283015460038401546004850154600586015460069096015494959394929391929091906001600160a01b031687565b607460209081526000928352604080842090915290825290208054600182015460029092015490919083565b6000818152607760209081526040918290206006018054835181840281018401909452808452606093928301828280156124cc57602002820191906000526020600020905b8154815260200190600101908083116124b8575b505050505090505b919050565b607f546001600160a01b031681565b607a6020526000908152604090205481565b606b5481565b6087546001600160a01b0316331461255f576040805162461bcd60e51b815260206004820152600e60248201527f6e6f7420617574686f72697a6564000000000000000000000000000000000000604482015290519081900360640190fd5b6001600160a01b03828116600090815260896020526040902054811690821614156125d1576040805162461bcd60e51b815260206004820152601060248201527f616c726561647920636f6d706c65746500000000000000000000000000000000604482015290519081900360640190fd5b806001600160a01b0316826001600160a01b03161415612638576040805162461bcd60e51b815260206004820152600c60248201527f73616d6520616464726573730000000000000000000000000000000000000000604482015290519081900360640190fd5b6001600160a01b03918216600090815260886020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001691909216179055565b61268782613a9f565b6126d8576040805162461bcd60e51b815260206004820152601760248201527f76616c696461746f7220646f65736e2774206578697374000000000000000000604482015290519081900360640190fd5b600082815260686020526040902060038101549054156126f6575060005b606654604080517fa4066fbe000000000000000000000000000000000000000000000000000000008152600481018690526024810184905290516001600160a01b039092169163a4066fbe9160448082019260009290919082900301818387803b15801561276357600080fd5b505af1158015612777573d6000803e3d6000fd5b5050505081801561278757508015155b15611171576066546000848152606a60205260409081902081517f242a6e3f0000000000000000000000000000000000000000000000000000000081526004810187815260248201938452825460027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6001831615610100020190911604604483018190526001600160a01b039095169463242a6e3f948994939091606490910190849080156128785780601f1061284d57610100808354040283529160200191612878565b820191906000526020600020905b81548152906001019060200180831161285b57829003601f168201915b50509350505050600060405180830381600087803b15801561289957600080fd5b505af115801561131b573d6000803e3d6000fd5b607260209081526000928352604080842090915290825290205481565b6001600160a01b03821660009081526073602090815260408083208484529091528120600201541580159061292157506001600160a01b038316600090815260736020908152604080832085845290915290205415155b801561295a57506001600160a01b0383166000908152607360209081526040808320858452909152902060020154612957612ebd565b11155b9392505050565b6081546001600160a01b031690565b336001600160a01b0382166129cc576040805162461bcd60e51b815260206004820152600c60248201527f7a65726f20616464726573730000000000000000000000000000000000000000604482015290519081900360640190fd5b6001600160a01b03818116600090815260886020526040902054811690831614612a3d576040805162461bcd60e51b815260206004820152600a60248201527f6e6f207265717565737400000000000000000000000000000000000000000000604482015290519081900360640190fd5b6001600160a01b0390811660009081526089602090815260408083208054949095167fffffffffffffffffffffffff000000000000000000000000000000000000000094851617909455608890529190912080549091169055565b607e5481565b60009182526077602090815260408084209284526003909201905290205490565b60009182526077602090815260408084209284526002909201905290205490565b612ae933612ec1565b612b245760405162461bcd60e51b81526004018080602001828103825260298152602001806145456029913960400191505060405180910390fd5b600060776000612b32611b64565b8152602001908152602001600020905060008090505b82811015612bab576000848483818110612b5e57fe5b60209081029290920135600081815260688452604080822060030154948890529020839055600c860154909350612b9c91508263ffffffff613a4516565b600c8501555050600101612b48565b50612bba60068201848461449b565b50606654607e54604080517f07aaf3440000000000000000000000000000000000000000000000000000000081526004810192909252516001600160a01b03909216916307aaf3449160248082019260009290919082900301818387803b15801561289957600080fd5b60009182526077602090815260408084209284526004909201905290205490565b612c4d6121ae565b612c9e576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b608080547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b6087546001600160a01b031681565b612cef6121ae565b612d40576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b612d4981613ab6565b50565b303b1590565b600054610100900460ff1680612d6b5750612d6b612d4c565b80612d79575060005460ff16155b612db45760405162461bcd60e51b815260040180806020018281038252602e81526020018061458f602e913960400191505060405180910390fd5b600054610100900460ff16158015612e1a57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff909116610100171660011790555b603380547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0384811691909117918290556040519116906000907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a38015612eb057600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1690555b5050565b64174876e80090565b4290565b6066546001600160a01b0390811691161490565b60005b83518110156130d557608160009054906101000a90046001600160a01b03166001600160a01b0316635a68f01a6040518163ffffffff1660e01b815260040160206040518083038186803b158015612f2f57600080fd5b505afa158015612f43573d6000803e3d6000fd5b505050506040513d6020811015612f5957600080fd5b50518251839083908110612f6957fe5b602002602001015111801561300b5750608160009054906101000a90046001600160a01b03166001600160a01b031662cc7f836040518163ffffffff1660e01b815260040160206040518083038186803b158015612fc657600080fd5b505afa158015612fda573d6000803e3d6000fd5b505050506040513d6020811015612ff057600080fd5b5051835184908390811061300057fe5b602002602001015110155b1561304c5761302e84828151811061301f57fe5b60200260200101516008613b6f565b61304c84828151811061303d57fe5b6020026020010151600061267e565b82818151811061305857fe5b602002602001015185600401600086848151811061307257fe5b602002602001015181526020019081526020016000208190555081818151811061309857fe5b60200260200101518560050160008684815181106130b257fe5b602090810291909101810151825281019190915260400160002055600101612ed8565b5050505050565b6130e46144d5565b6040518060a00160405280855160405190808252806020026020018201604052801561311a578160200160208202803883390190505b508152602001600081526020018551604051908082528060200260200182016040528015613152578160200160208202803883390190505b508152602001600081526020016000815250905060008090505b845181101561326d57600086600301600087848151811061318957fe5b602002602001015181526020019081526020016000205490506000809050818584815181106131b457fe5b602002602001015111156131db57818584815181106131cf57fe5b60200260200101510390505b898684815181106131e857fe5b60200260200101518202816131f957fe5b048460400151848151811061320a57fe5b6020026020010181815250506132448460400151848151811061322957fe5b60200260200101518560600151613a4590919063ffffffff16565b6060850152608084015161325e908263ffffffff613a4516565b6080850152505060010161316c565b5060005b8451811015613336578784828151811061328757fe5b60200260200101518986848151811061329c57fe5b60200260200101518a60000160008a87815181106132b657fe5b602002602001015181526020019081526020016000205402816132d557fe5b0402816132de57fe5b04826000015182815181106132ef57fe5b6020026020010181815250506133298260000151828151811061330e57fe5b60200260200101518360200151613a4590919063ffffffff16565b6020830152600101613271565b5060005b84518110156137755760006133e389608160009054906101000a90046001600160a01b03166001600160a01b031663d9a7c1f96040518163ffffffff1660e01b815260040160206040518083038186803b15801561339757600080fd5b505afa1580156133ab573d6000803e3d6000fd5b505050506040513d60208110156133c157600080fd5b505185518051869081106133d157fe5b60200260200101518660200151613c99565b905061341f61341284608001518560400151858151811061340057fe5b60200260200101518660600151613ce6565b829063ffffffff613a4516565b9050600086838151811061342f57fe5b60209081029190910181015160008181526068835260408082206006015460815482517fa778651500000000000000000000000000000000000000000000000000000000815292519496506001600160a01b039182169593946134e5948994929093169263a77865159260048082019391829003018186803b1580156134b457600080fd5b505afa1580156134c8573d6000803e3d6000fd5b505050506040513d60208110156134de57600080fd5b5051613e4f565b6001600160a01b0383166000908152607260209081526040808320878452909152902054909150801561368c5760008161351f8587611846565b84028161352857fe5b0490508083036135366143e2565b6001600160a01b03861660009081526073602090815260408083208a8452909152902060030154613568908490613e6c565b90506135726143e2565b61357d836000613e6c565b6001600160a01b0388166000908152606f602090815260408083208c845282529182902082516060810184528154815260018201549281019290925260020154918101919091529091506135d290838361402e565b6001600160a01b0388166000818152606f602090815260408083208d84528252808320855181558583015160018083019190915595820151600291820155938352607482528083208d84528252918290208251606081018452815481529481015491850191909152909101549082015261364d90838361402e565b6001600160a01b03881660009081526074602090815260408083208c845282529182902083518155908301516001820155910151600290910155505050505b6000848152606860205260408120600301548387039181156136be57816136b1614049565b8402816136ba57fe5b0490505b808e600101600089815260200190815260200160002054018f6001016000898152602001908152602001600020819055508a89815181106136fb57fe5b60200260200101518f6003016000898152602001908152602001600020819055508b898151811061372857fe5b60200260200101518e600201600089815260200190815260200160002054018f6002016000898152602001908152602001600020819055505050505050505050808060010191505061333a565b50608081015160088701819055602082015160098801556060820151600a88015560765411156137b3576008860154607680549190910390556137b9565b60006076555b607f546001600160a01b03161561131b5760006137d4614049565b608160009054906101000a90046001600160a01b03166001600160a01b03166394c3e9146040518163ffffffff1660e01b815260040160206040518083038186803b15801561382257600080fd5b505afa158015613836573d6000803e3d6000fd5b505050506040513d602081101561384c57600080fd5b50516080840151028161385b57fe5b04905061386781614055565b607f546040516001600160a01b0390911690620f42409083906000818181858888f193505050503d80600081146138ba576040519150601f19603f3d011682016040523d82523d6000602084013e6138bf565b606091505b5050505050505050505050565b608154604080517f3a3ef66c00000000000000000000000000000000000000000000000000000000815290516000926001600160a01b031691633a3ef66c916004808301926020929190829003018186803b15801561392a57600080fd5b505afa15801561393e573d6000803e3d6000fd5b505050506040513d602081101561395457600080fd5b505183026001019050600081613968614049565b84028161397157fe5b0490506000608160009054906101000a90046001600160a01b03166001600160a01b0316632c8c36a56040518163ffffffff1660e01b815260040160206040518083038186803b1580156139c457600080fd5b505afa1580156139d8573d6000803e3d6000fd5b505050506040513d60208110156139ee57600080fd5b505190508481016139fd614049565b82028387020181613a0a57fe5b049150613a16826140f3565b91506000613a22614049565b83607e540281613a2e57fe5b049050613a3a81614161565b607e55505050505050565b60008282018381101561295a576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b600090815260686020526040902060050154151590565b6001600160a01b038116613afb5760405162461bcd60e51b815260040180806020018281038252602681526020018061451f6026913960400191505060405180910390fd5b6033546040516001600160a01b038084169216907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3603380547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b600082815260686020526040902054158015613b8a57508015155b15613bb757600082815260686020526040902060030154606d54613bb39163ffffffff61419716565b606d555b600082815260686020526040902054811115612eb057600082815260686020526040902081815560020154613c5f57613bee611b64565b600083815260686020526040902060020155613c08612ebd565b6000838152606860209081526040918290206001810184905560020154825190815290810192909252805184927fac4801c32a6067ff757446524ee4e7a373797278ac3c883eac5c693b4ad72e4792908290030190a25b60408051828152905183917fcd35267e7654194727477d6c78b541a553483cff7f92a055d17868d3da6e953e919081900360200190a25050565b600082613ca857506000611a86565b6000613cba868663ffffffff6141d916565b9050613cdc83613cd0838763ffffffff6141d916565b9063ffffffff61423216565b9695505050505050565b600082613cf55750600061295a565b6000613d0b83613cd0878763ffffffff6141d916565b9050613e46613d18614049565b608154604080517f94c3e9140000000000000000000000000000000000000000000000000000000081529051613cd0926001600160a01b0316916394c3e914916004808301926020929190829003018186803b158015613d7757600080fd5b505afa158015613d8b573d6000803e3d6000fd5b505050506040513d6020811015613da157600080fd5b5051608154604080517fc74dd62100000000000000000000000000000000000000000000000000000000815290516001600160a01b039092169163c74dd62191600480820192602092909190829003018186803b158015613e0157600080fd5b505afa158015613e15573d6000803e3d6000fd5b505050506040513d6020811015613e2b57600080fd5b5051613e35614049565b0303846141d990919063ffffffff16565b95945050505050565b600061295a613e5c614049565b613cd0858563ffffffff6141d916565b613e746143e2565b60405180606001604052806000815260200160008152602001600081525090506000608160009054906101000a90046001600160a01b03166001600160a01b0316635e2308d26040518163ffffffff1660e01b815260040160206040518083038186803b158015613ee457600080fd5b505afa158015613ef8573d6000803e3d6000fd5b505050506040513d6020811015613f0e57600080fd5b50519050821561400657600081613f23614049565b0390506000613fb5608160009054906101000a90046001600160a01b03166001600160a01b0316630d4955e36040518163ffffffff1660e01b815260040160206040518083038186803b158015613f7957600080fd5b505afa158015613f8d573d6000803e3d6000fd5b505050506040513d6020811015613fa357600080fd5b5051613cd0848863ffffffff6141d916565b90506000613fd6613fc4614049565b613cd08987860163ffffffff6141d916565b9050613ff3613fe3614049565b613cd0898763ffffffff6141d916565b6020860181905290038452506140279050565b614021614011614049565b613cd0868463ffffffff6141d916565b60408301525b5092915050565b6140366143e2565b611a866140438585614274565b83614274565b670de0b6b3a764000090565b606654604080517f66e7ea0f0000000000000000000000000000000000000000000000000000000081523060048201526024810184905290516001600160a01b03909216916366e7ea0f9160448082019260009290919082900301818387803b1580156140c157600080fd5b505af11580156140d5573d6000803e3d6000fd5b50506076546140ed925090508263ffffffff613a4516565b60765550565b600060646140ff614049565b6069028161410957fe5b0482111561412d57606461411b614049565b6069028161412557fe5b0490506124d4565b6064614137614049565b605f028161414157fe5b0482101561415d576064614153614049565b605f028161412557fe5b5090565b600066038d7ea4c68000821115614180575066038d7ea4c680006124d4565b633b9aca0082101561415d5750633b9aca006124d4565b600061295a83836040518060400160405280601e81526020017f536166654d6174683a207375627472616374696f6e206f766572666c6f7700008152506142e6565b6000826141e85750600061155e565b828202828482816141f557fe5b041461295a5760405162461bcd60e51b815260040180806020018281038252602181526020018061456e6021913960400191505060405180910390fd5b600061295a83836040518060400160405280601a81526020017f536166654d6174683a206469766973696f6e206279207a65726f00000000000081525061437d565b61427c6143e2565b604080516060810190915282518451829161429d919063ffffffff613a4516565b81526020016142bd84602001518660200151613a4590919063ffffffff16565b81526020016142dd84604001518660400151613a4590919063ffffffff16565b90529392505050565b600081848411156143755760405162461bcd60e51b81526004018080602001828103825283818151815260200191508051906020019080838360005b8381101561433a578181015183820152602001614322565b50505050905090810190601f1680156143675780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b505050900390565b600081836143cc5760405162461bcd60e51b815260206004820181815283516024840152835190928392604490910191908501908083836000831561433a578181015183820152602001614322565b5060008385816143d857fe5b0495945050505050565b60405180606001604052806000815260200160008152602001600081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10614462578280017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082351617855561448f565b8280016001018555821561448f579182015b8281111561448f578235825591602001919060010190614474565b5061415d929150614504565b82805482825590600052602060002090810192821561448f579160200282018281111561448f578235825591602001919060010190614474565b6040518060a0016040528060608152602001600081526020016060815260200160008152602001600081525090565b6114ad91905b8082111561415d576000815560010161450a56fe4f776e61626c653a206e6577206f776e657220697320746865207a65726f206164647265737363616c6c6572206973206e6f7420746865204e6f64654472697665724175746820636f6e7472616374536166654d6174683a206d756c7469706c69636174696f6e206f766572666c6f77436f6e747261637420696e7374616e63652068617320616c7265616479206265656e20696e697469616c697a6564a265627a7a72315820f930ac1e4aa030cc95832bcc388afd7a33ffda6c05bcf821a74088d69fdd4c2764736f6c63430005110032",
}

// SFCABI is the input ABI used to generate the binding from.
// Deprecated: Use SFCMetaData.ABI instead.
var SFCABI = SFCMetaData.ABI

// SFC is an auto generated Go binding around an Ethereum contract.
type SFC struct {
	SFCCaller     // Read-only binding to the contract
	SFCTransactor // Write-only binding to the contract
	SFCFilterer   // Log filterer for contract events
}

// SFCCaller is an auto generated read-only Go binding around an Ethereum contract.
type SFCCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SFCTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SFCTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SFCFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SFCFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SFCSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SFCSession struct {
	Contract     *SFC              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SFCCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SFCCallerSession struct {
	Contract *SFCCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// SFCTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SFCTransactorSession struct {
	Contract     *SFCTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SFCRaw is an auto generated low-level Go binding around an Ethereum contract.
type SFCRaw struct {
	Contract *SFC // Generic contract binding to access the raw methods on
}

// SFCCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SFCCallerRaw struct {
	Contract *SFCCaller // Generic read-only contract binding to access the raw methods on
}

// SFCTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SFCTransactorRaw struct {
	Contract *SFCTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSFC creates a new instance of SFC, bound to a specific deployed contract.
func NewSFC(address common.Address, backend bind.ContractBackend) (*SFC, error) {
	contract, err := bindSFC(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SFC{SFCCaller: SFCCaller{contract: contract}, SFCTransactor: SFCTransactor{contract: contract}, SFCFilterer: SFCFilterer{contract: contract}}, nil
}

// NewSFCCaller creates a new read-only instance of SFC, bound to a specific deployed contract.
func NewSFCCaller(address common.Address, caller bind.ContractCaller) (*SFCCaller, error) {
	contract, err := bindSFC(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SFCCaller{contract: contract}, nil
}

// NewSFCTransactor creates a new write-only instance of SFC, bound to a specific deployed contract.
func NewSFCTransactor(address common.Address, transactor bind.ContractTransactor) (*SFCTransactor, error) {
	contract, err := bindSFC(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SFCTransactor{contract: contract}, nil
}

// NewSFCFilterer creates a new log filterer instance of SFC, bound to a specific deployed contract.
func NewSFCFilterer(address common.Address, filterer bind.ContractFilterer) (*SFCFilterer, error) {
	contract, err := bindSFC(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SFCFilterer{contract: contract}, nil
}

// bindSFC binds a generic wrapper to an already deployed contract.
func bindSFC(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SFCABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SFC *SFCRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SFC.Contract.SFCCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SFC *SFCRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SFC.Contract.SFCTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SFC *SFCRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SFC.Contract.SFCTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SFC *SFCCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SFC.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SFC *SFCTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SFC.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SFC *SFCTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SFC.Contract.contract.Transact(opts, method, params...)
}

// LastValidatorID is a free data retrieval call binding the contract method 0xc7be95de.
//
// Solidity: function lastValidatorID() view returns(uint256)
func (_SFC *SFCCaller) LastValidatorID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SFC.contract.Call(opts, &out, "lastValidatorID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetValidator is a free data retrieval call binding the contract method 0xb5d89627.
//
// Solidity: function getValidator(uint256 ) view returns(uint256 status, uint256 deactivatedTime, uint256 deactivatedEpoch, uint256 receivedStake, uint256 createdEpoch, uint256 createdTime, address auth)
func (_SFC *SFCCaller) GetValidator(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Status           *big.Int
	DeactivatedTime  *big.Int
	DeactivatedEpoch *big.Int
	ReceivedStake    *big.Int
	CreatedEpoch     *big.Int
	CreatedTime      *big.Int
	Auth             common.Address
}, error) {
	var out []interface{}
	err := _SFC.contract.Call(opts, &out, "getValidator", arg0)

	outstruct := new(struct {
		Status           *big.Int
		DeactivatedTime  *big.Int
		DeactivatedEpoch *big.Int
		ReceivedStake    *big.Int
		CreatedEpoch     *big.Int
		CreatedTime      *big.Int
		Auth             common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Status = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.DeactivatedTime = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.DeactivatedEpoch = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.ReceivedStake = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.CreatedEpoch = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.CreatedTime = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.Auth = *abi.ConvertType(out[6], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// CreateValidator is a paid mutator transaction binding the contract method 0xa5a470ad.
//
// Solidity: function createValidator(bytes pubkey) payable returns()
func (_SFC *SFCCaller) CreateValidator(opts *bind.TransactOpts, pubkey []byte) (*types.Transaction, error) {
	return _SFC.contract.Transact(opts, "createValidator", pubkey)
}

// CreateValidator is a paid mutator transaction binding the contract method 0xa5a470ad.
//
// Solidity: function createValidator(bytes pubkey) payable returns()
func (_SFC *SFCSession) CreateValidator(pubkey []byte) (*types.Transaction, error) {
	return _SFC.Contract.CreateValidator(&_SFC.TransactOpts, pubkey)
}

// SFCApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the SFC contract.
type SFCApprovalIterator struct {
	Event *SFCApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SFCApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SFCApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SFCApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SFCApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SFCApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SFCApproval represents a Approval event raised by the SFC contract.
type SFCApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_SFC *SFCFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*SFCApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _SFC.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &SFCApprovalIterator{contract: _SFC.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_SFC *SFCFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *SFCApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _SFC.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SFCApproval)
				if err := _SFC.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_SFC *SFCFilterer) ParseApproval(log types.Log) (*SFCApproval, error) {
	event := new(SFCApproval)
	if err := _SFC.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SFCTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the SFC contract.
type SFCTransferIterator struct {
	Event *SFCTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SFCTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SFCTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SFCTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SFCTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SFCTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SFCTransfer represents a Transfer event raised by the SFC contract.
type SFCTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_SFC *SFCFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SFCTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SFC.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &SFCTransferIterator{contract: _SFC.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_SFC *SFCFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *SFCTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SFC.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SFCTransfer)
				if err := _SFC.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_SFC *SFCFilterer) ParseTransfer(log types.Log) (*SFCTransfer, error) {
	event := new(SFCTransfer)
	if err := _SFC.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
