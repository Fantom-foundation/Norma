// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

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

// NodeDriverMetaData contains all meta data concerning the NodeDriver contract.
var NodeDriverMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"name\":\"AdvanceEpochs\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"diff\",\"type\":\"bytes\"}],\"name\":\"UpdateNetworkRules\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"version\",\"type\":\"uint256\"}],\"name\":\"UpdateNetworkVersion\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"}],\"name\":\"UpdateValidatorPubkey\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"weight\",\"type\":\"uint256\"}],\"name\":\"UpdateValidatorWeight\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"backend\",\"type\":\"address\"}],\"name\":\"UpdatedBackend\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_backend\",\"type\":\"address\"}],\"name\":\"setBackend\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_backend\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_evmWriterAddress\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"acc\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setBalance\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"acc\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"copyCode\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"acc\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"with\",\"type\":\"address\"}],\"name\":\"swapCode\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"acc\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"value\",\"type\":\"bytes32\"}],\"name\":\"setStorage\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"acc\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"diff\",\"type\":\"uint256\"}],\"name\":\"incNonce\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"diff\",\"type\":\"bytes\"}],\"name\":\"updateNetworkRules\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"version\",\"type\":\"uint256\"}],\"name\":\"updateNetworkVersion\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"name\":\"advanceEpochs\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"updateValidatorWeight\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"}],\"name\":\"updateValidatorPubkey\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_auth\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"status\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivatedEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivatedTime\",\"type\":\"uint256\"}],\"name\":\"setGenesisValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"toValidatorID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockupFromEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockupEndTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockupDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"earlyUnlockPenalty\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rewards\",\"type\":\"uint256\"}],\"name\":\"setGenesisDelegation\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"status\",\"type\":\"uint256\"}],\"name\":\"deactivateValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"nextValidatorIDs\",\"type\":\"uint256[]\"}],\"name\":\"sealEpochValidators\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"offlineTimes\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"offlineBlocks\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"uptimes\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"originatedTxsFee\",\"type\":\"uint256[]\"}],\"name\":\"sealEpoch\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"offlineTimes\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"offlineBlocks\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"uptimes\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"originatedTxsFee\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"usedGas\",\"type\":\"uint256\"}],\"name\":\"sealEpochV1\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// NodeDriverABI is the input ABI used to generate the binding from.
// Deprecated: Use NodeDriverMetaData.ABI instead.
var NodeDriverABI = NodeDriverMetaData.ABI

// NodeDriver is an auto generated Go binding around an Ethereum contract.
type NodeDriver struct {
	NodeDriverCaller     // Read-only binding to the contract
	NodeDriverTransactor // Write-only binding to the contract
	NodeDriverFilterer   // Log filterer for contract events
}

// NodeDriverCaller is an auto generated read-only Go binding around an Ethereum contract.
type NodeDriverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeDriverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NodeDriverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeDriverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NodeDriverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeDriverSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NodeDriverSession struct {
	Contract     *NodeDriver       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NodeDriverCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NodeDriverCallerSession struct {
	Contract *NodeDriverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// NodeDriverTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NodeDriverTransactorSession struct {
	Contract     *NodeDriverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// NodeDriverRaw is an auto generated low-level Go binding around an Ethereum contract.
type NodeDriverRaw struct {
	Contract *NodeDriver // Generic contract binding to access the raw methods on
}

// NodeDriverCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NodeDriverCallerRaw struct {
	Contract *NodeDriverCaller // Generic read-only contract binding to access the raw methods on
}

// NodeDriverTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NodeDriverTransactorRaw struct {
	Contract *NodeDriverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNodeDriver creates a new instance of NodeDriver, bound to a specific deployed contract.
func NewNodeDriver(address common.Address, backend bind.ContractBackend) (*NodeDriver, error) {
	contract, err := bindNodeDriver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NodeDriver{NodeDriverCaller: NodeDriverCaller{contract: contract}, NodeDriverTransactor: NodeDriverTransactor{contract: contract}, NodeDriverFilterer: NodeDriverFilterer{contract: contract}}, nil
}

// NewNodeDriverCaller creates a new read-only instance of NodeDriver, bound to a specific deployed contract.
func NewNodeDriverCaller(address common.Address, caller bind.ContractCaller) (*NodeDriverCaller, error) {
	contract, err := bindNodeDriver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NodeDriverCaller{contract: contract}, nil
}

// NewNodeDriverTransactor creates a new write-only instance of NodeDriver, bound to a specific deployed contract.
func NewNodeDriverTransactor(address common.Address, transactor bind.ContractTransactor) (*NodeDriverTransactor, error) {
	contract, err := bindNodeDriver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NodeDriverTransactor{contract: contract}, nil
}

// NewNodeDriverFilterer creates a new log filterer instance of NodeDriver, bound to a specific deployed contract.
func NewNodeDriverFilterer(address common.Address, filterer bind.ContractFilterer) (*NodeDriverFilterer, error) {
	contract, err := bindNodeDriver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NodeDriverFilterer{contract: contract}, nil
}

// bindNodeDriver binds a generic wrapper to an already deployed contract.
func bindNodeDriver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NodeDriverABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeDriver *NodeDriverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeDriver.Contract.NodeDriverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeDriver *NodeDriverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeDriver.Contract.NodeDriverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeDriver *NodeDriverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeDriver.Contract.NodeDriverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeDriver *NodeDriverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeDriver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeDriver *NodeDriverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeDriver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeDriver *NodeDriverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeDriver.Contract.contract.Transact(opts, method, params...)
}

// AdvanceEpochs is a paid mutator transaction binding the contract method 0x0aeeca00.
//
// Solidity: function advanceEpochs(uint256 num) returns()
func (_NodeDriver *NodeDriverTransactor) AdvanceEpochs(opts *bind.TransactOpts, num *big.Int) (*types.Transaction, error) {
	return _NodeDriver.contract.Transact(opts, "advanceEpochs", num)
}

// AdvanceEpochs is a paid mutator transaction binding the contract method 0x0aeeca00.
//
// Solidity: function advanceEpochs(uint256 num) returns()
func (_NodeDriver *NodeDriverSession) AdvanceEpochs(num *big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.AdvanceEpochs(&_NodeDriver.TransactOpts, num)
}

// AdvanceEpochs is a paid mutator transaction binding the contract method 0x0aeeca00.
//
// Solidity: function advanceEpochs(uint256 num) returns()
func (_NodeDriver *NodeDriverTransactorSession) AdvanceEpochs(num *big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.AdvanceEpochs(&_NodeDriver.TransactOpts, num)
}

// CopyCode is a paid mutator transaction binding the contract method 0xd6a0c7af.
//
// Solidity: function copyCode(address acc, address from) returns()
func (_NodeDriver *NodeDriverTransactor) CopyCode(opts *bind.TransactOpts, acc common.Address, from common.Address) (*types.Transaction, error) {
	return _NodeDriver.contract.Transact(opts, "copyCode", acc, from)
}

// CopyCode is a paid mutator transaction binding the contract method 0xd6a0c7af.
//
// Solidity: function copyCode(address acc, address from) returns()
func (_NodeDriver *NodeDriverSession) CopyCode(acc common.Address, from common.Address) (*types.Transaction, error) {
	return _NodeDriver.Contract.CopyCode(&_NodeDriver.TransactOpts, acc, from)
}

// CopyCode is a paid mutator transaction binding the contract method 0xd6a0c7af.
//
// Solidity: function copyCode(address acc, address from) returns()
func (_NodeDriver *NodeDriverTransactorSession) CopyCode(acc common.Address, from common.Address) (*types.Transaction, error) {
	return _NodeDriver.Contract.CopyCode(&_NodeDriver.TransactOpts, acc, from)
}

// DeactivateValidator is a paid mutator transaction binding the contract method 0x1e702f83.
//
// Solidity: function deactivateValidator(uint256 validatorID, uint256 status) returns()
func (_NodeDriver *NodeDriverTransactor) DeactivateValidator(opts *bind.TransactOpts, validatorID *big.Int, status *big.Int) (*types.Transaction, error) {
	return _NodeDriver.contract.Transact(opts, "deactivateValidator", validatorID, status)
}

// DeactivateValidator is a paid mutator transaction binding the contract method 0x1e702f83.
//
// Solidity: function deactivateValidator(uint256 validatorID, uint256 status) returns()
func (_NodeDriver *NodeDriverSession) DeactivateValidator(validatorID *big.Int, status *big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.DeactivateValidator(&_NodeDriver.TransactOpts, validatorID, status)
}

// DeactivateValidator is a paid mutator transaction binding the contract method 0x1e702f83.
//
// Solidity: function deactivateValidator(uint256 validatorID, uint256 status) returns()
func (_NodeDriver *NodeDriverTransactorSession) DeactivateValidator(validatorID *big.Int, status *big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.DeactivateValidator(&_NodeDriver.TransactOpts, validatorID, status)
}

// IncNonce is a paid mutator transaction binding the contract method 0x79bead38.
//
// Solidity: function incNonce(address acc, uint256 diff) returns()
func (_NodeDriver *NodeDriverTransactor) IncNonce(opts *bind.TransactOpts, acc common.Address, diff *big.Int) (*types.Transaction, error) {
	return _NodeDriver.contract.Transact(opts, "incNonce", acc, diff)
}

// IncNonce is a paid mutator transaction binding the contract method 0x79bead38.
//
// Solidity: function incNonce(address acc, uint256 diff) returns()
func (_NodeDriver *NodeDriverSession) IncNonce(acc common.Address, diff *big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.IncNonce(&_NodeDriver.TransactOpts, acc, diff)
}

// IncNonce is a paid mutator transaction binding the contract method 0x79bead38.
//
// Solidity: function incNonce(address acc, uint256 diff) returns()
func (_NodeDriver *NodeDriverTransactorSession) IncNonce(acc common.Address, diff *big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.IncNonce(&_NodeDriver.TransactOpts, acc, diff)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _backend, address _evmWriterAddress) returns()
func (_NodeDriver *NodeDriverTransactor) Initialize(opts *bind.TransactOpts, _backend common.Address, _evmWriterAddress common.Address) (*types.Transaction, error) {
	return _NodeDriver.contract.Transact(opts, "initialize", _backend, _evmWriterAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _backend, address _evmWriterAddress) returns()
func (_NodeDriver *NodeDriverSession) Initialize(_backend common.Address, _evmWriterAddress common.Address) (*types.Transaction, error) {
	return _NodeDriver.Contract.Initialize(&_NodeDriver.TransactOpts, _backend, _evmWriterAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _backend, address _evmWriterAddress) returns()
func (_NodeDriver *NodeDriverTransactorSession) Initialize(_backend common.Address, _evmWriterAddress common.Address) (*types.Transaction, error) {
	return _NodeDriver.Contract.Initialize(&_NodeDriver.TransactOpts, _backend, _evmWriterAddress)
}

// SealEpoch is a paid mutator transaction binding the contract method 0xebdf104c.
//
// Solidity: function sealEpoch(uint256[] offlineTimes, uint256[] offlineBlocks, uint256[] uptimes, uint256[] originatedTxsFee) returns()
func (_NodeDriver *NodeDriverTransactor) SealEpoch(opts *bind.TransactOpts, offlineTimes []*big.Int, offlineBlocks []*big.Int, uptimes []*big.Int, originatedTxsFee []*big.Int) (*types.Transaction, error) {
	return _NodeDriver.contract.Transact(opts, "sealEpoch", offlineTimes, offlineBlocks, uptimes, originatedTxsFee)
}

// SealEpoch is a paid mutator transaction binding the contract method 0xebdf104c.
//
// Solidity: function sealEpoch(uint256[] offlineTimes, uint256[] offlineBlocks, uint256[] uptimes, uint256[] originatedTxsFee) returns()
func (_NodeDriver *NodeDriverSession) SealEpoch(offlineTimes []*big.Int, offlineBlocks []*big.Int, uptimes []*big.Int, originatedTxsFee []*big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.SealEpoch(&_NodeDriver.TransactOpts, offlineTimes, offlineBlocks, uptimes, originatedTxsFee)
}

// SealEpoch is a paid mutator transaction binding the contract method 0xebdf104c.
//
// Solidity: function sealEpoch(uint256[] offlineTimes, uint256[] offlineBlocks, uint256[] uptimes, uint256[] originatedTxsFee) returns()
func (_NodeDriver *NodeDriverTransactorSession) SealEpoch(offlineTimes []*big.Int, offlineBlocks []*big.Int, uptimes []*big.Int, originatedTxsFee []*big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.SealEpoch(&_NodeDriver.TransactOpts, offlineTimes, offlineBlocks, uptimes, originatedTxsFee)
}

// SealEpochV1 is a paid mutator transaction binding the contract method 0x7f52e13e.
//
// Solidity: function sealEpochV1(uint256[] offlineTimes, uint256[] offlineBlocks, uint256[] uptimes, uint256[] originatedTxsFee, uint256 usedGas) returns()
func (_NodeDriver *NodeDriverTransactor) SealEpochV1(opts *bind.TransactOpts, offlineTimes []*big.Int, offlineBlocks []*big.Int, uptimes []*big.Int, originatedTxsFee []*big.Int, usedGas *big.Int) (*types.Transaction, error) {
	return _NodeDriver.contract.Transact(opts, "sealEpochV1", offlineTimes, offlineBlocks, uptimes, originatedTxsFee, usedGas)
}

// SealEpochV1 is a paid mutator transaction binding the contract method 0x7f52e13e.
//
// Solidity: function sealEpochV1(uint256[] offlineTimes, uint256[] offlineBlocks, uint256[] uptimes, uint256[] originatedTxsFee, uint256 usedGas) returns()
func (_NodeDriver *NodeDriverSession) SealEpochV1(offlineTimes []*big.Int, offlineBlocks []*big.Int, uptimes []*big.Int, originatedTxsFee []*big.Int, usedGas *big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.SealEpochV1(&_NodeDriver.TransactOpts, offlineTimes, offlineBlocks, uptimes, originatedTxsFee, usedGas)
}

// SealEpochV1 is a paid mutator transaction binding the contract method 0x7f52e13e.
//
// Solidity: function sealEpochV1(uint256[] offlineTimes, uint256[] offlineBlocks, uint256[] uptimes, uint256[] originatedTxsFee, uint256 usedGas) returns()
func (_NodeDriver *NodeDriverTransactorSession) SealEpochV1(offlineTimes []*big.Int, offlineBlocks []*big.Int, uptimes []*big.Int, originatedTxsFee []*big.Int, usedGas *big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.SealEpochV1(&_NodeDriver.TransactOpts, offlineTimes, offlineBlocks, uptimes, originatedTxsFee, usedGas)
}

// SealEpochValidators is a paid mutator transaction binding the contract method 0xe08d7e66.
//
// Solidity: function sealEpochValidators(uint256[] nextValidatorIDs) returns()
func (_NodeDriver *NodeDriverTransactor) SealEpochValidators(opts *bind.TransactOpts, nextValidatorIDs []*big.Int) (*types.Transaction, error) {
	return _NodeDriver.contract.Transact(opts, "sealEpochValidators", nextValidatorIDs)
}

// SealEpochValidators is a paid mutator transaction binding the contract method 0xe08d7e66.
//
// Solidity: function sealEpochValidators(uint256[] nextValidatorIDs) returns()
func (_NodeDriver *NodeDriverSession) SealEpochValidators(nextValidatorIDs []*big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.SealEpochValidators(&_NodeDriver.TransactOpts, nextValidatorIDs)
}

// SealEpochValidators is a paid mutator transaction binding the contract method 0xe08d7e66.
//
// Solidity: function sealEpochValidators(uint256[] nextValidatorIDs) returns()
func (_NodeDriver *NodeDriverTransactorSession) SealEpochValidators(nextValidatorIDs []*big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.SealEpochValidators(&_NodeDriver.TransactOpts, nextValidatorIDs)
}

// SetBackend is a paid mutator transaction binding the contract method 0xda7fc24f.
//
// Solidity: function setBackend(address _backend) returns()
func (_NodeDriver *NodeDriverTransactor) SetBackend(opts *bind.TransactOpts, _backend common.Address) (*types.Transaction, error) {
	return _NodeDriver.contract.Transact(opts, "setBackend", _backend)
}

// SetBackend is a paid mutator transaction binding the contract method 0xda7fc24f.
//
// Solidity: function setBackend(address _backend) returns()
func (_NodeDriver *NodeDriverSession) SetBackend(_backend common.Address) (*types.Transaction, error) {
	return _NodeDriver.Contract.SetBackend(&_NodeDriver.TransactOpts, _backend)
}

// SetBackend is a paid mutator transaction binding the contract method 0xda7fc24f.
//
// Solidity: function setBackend(address _backend) returns()
func (_NodeDriver *NodeDriverTransactorSession) SetBackend(_backend common.Address) (*types.Transaction, error) {
	return _NodeDriver.Contract.SetBackend(&_NodeDriver.TransactOpts, _backend)
}

// SetBalance is a paid mutator transaction binding the contract method 0xe30443bc.
//
// Solidity: function setBalance(address acc, uint256 value) returns()
func (_NodeDriver *NodeDriverTransactor) SetBalance(opts *bind.TransactOpts, acc common.Address, value *big.Int) (*types.Transaction, error) {
	return _NodeDriver.contract.Transact(opts, "setBalance", acc, value)
}

// SetBalance is a paid mutator transaction binding the contract method 0xe30443bc.
//
// Solidity: function setBalance(address acc, uint256 value) returns()
func (_NodeDriver *NodeDriverSession) SetBalance(acc common.Address, value *big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.SetBalance(&_NodeDriver.TransactOpts, acc, value)
}

// SetBalance is a paid mutator transaction binding the contract method 0xe30443bc.
//
// Solidity: function setBalance(address acc, uint256 value) returns()
func (_NodeDriver *NodeDriverTransactorSession) SetBalance(acc common.Address, value *big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.SetBalance(&_NodeDriver.TransactOpts, acc, value)
}

// SetGenesisDelegation is a paid mutator transaction binding the contract method 0x18f628d4.
//
// Solidity: function setGenesisDelegation(address delegator, uint256 toValidatorID, uint256 stake, uint256 lockedStake, uint256 lockupFromEpoch, uint256 lockupEndTime, uint256 lockupDuration, uint256 earlyUnlockPenalty, uint256 rewards) returns()
func (_NodeDriver *NodeDriverTransactor) SetGenesisDelegation(opts *bind.TransactOpts, delegator common.Address, toValidatorID *big.Int, stake *big.Int, lockedStake *big.Int, lockupFromEpoch *big.Int, lockupEndTime *big.Int, lockupDuration *big.Int, earlyUnlockPenalty *big.Int, rewards *big.Int) (*types.Transaction, error) {
	return _NodeDriver.contract.Transact(opts, "setGenesisDelegation", delegator, toValidatorID, stake, lockedStake, lockupFromEpoch, lockupEndTime, lockupDuration, earlyUnlockPenalty, rewards)
}

// SetGenesisDelegation is a paid mutator transaction binding the contract method 0x18f628d4.
//
// Solidity: function setGenesisDelegation(address delegator, uint256 toValidatorID, uint256 stake, uint256 lockedStake, uint256 lockupFromEpoch, uint256 lockupEndTime, uint256 lockupDuration, uint256 earlyUnlockPenalty, uint256 rewards) returns()
func (_NodeDriver *NodeDriverSession) SetGenesisDelegation(delegator common.Address, toValidatorID *big.Int, stake *big.Int, lockedStake *big.Int, lockupFromEpoch *big.Int, lockupEndTime *big.Int, lockupDuration *big.Int, earlyUnlockPenalty *big.Int, rewards *big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.SetGenesisDelegation(&_NodeDriver.TransactOpts, delegator, toValidatorID, stake, lockedStake, lockupFromEpoch, lockupEndTime, lockupDuration, earlyUnlockPenalty, rewards)
}

// SetGenesisDelegation is a paid mutator transaction binding the contract method 0x18f628d4.
//
// Solidity: function setGenesisDelegation(address delegator, uint256 toValidatorID, uint256 stake, uint256 lockedStake, uint256 lockupFromEpoch, uint256 lockupEndTime, uint256 lockupDuration, uint256 earlyUnlockPenalty, uint256 rewards) returns()
func (_NodeDriver *NodeDriverTransactorSession) SetGenesisDelegation(delegator common.Address, toValidatorID *big.Int, stake *big.Int, lockedStake *big.Int, lockupFromEpoch *big.Int, lockupEndTime *big.Int, lockupDuration *big.Int, earlyUnlockPenalty *big.Int, rewards *big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.SetGenesisDelegation(&_NodeDriver.TransactOpts, delegator, toValidatorID, stake, lockedStake, lockupFromEpoch, lockupEndTime, lockupDuration, earlyUnlockPenalty, rewards)
}

// SetGenesisValidator is a paid mutator transaction binding the contract method 0x4feb92f3.
//
// Solidity: function setGenesisValidator(address _auth, uint256 validatorID, bytes pubkey, uint256 status, uint256 createdEpoch, uint256 createdTime, uint256 deactivatedEpoch, uint256 deactivatedTime) returns()
func (_NodeDriver *NodeDriverTransactor) SetGenesisValidator(opts *bind.TransactOpts, _auth common.Address, validatorID *big.Int, pubkey []byte, status *big.Int, createdEpoch *big.Int, createdTime *big.Int, deactivatedEpoch *big.Int, deactivatedTime *big.Int) (*types.Transaction, error) {
	return _NodeDriver.contract.Transact(opts, "setGenesisValidator", _auth, validatorID, pubkey, status, createdEpoch, createdTime, deactivatedEpoch, deactivatedTime)
}

// SetGenesisValidator is a paid mutator transaction binding the contract method 0x4feb92f3.
//
// Solidity: function setGenesisValidator(address _auth, uint256 validatorID, bytes pubkey, uint256 status, uint256 createdEpoch, uint256 createdTime, uint256 deactivatedEpoch, uint256 deactivatedTime) returns()
func (_NodeDriver *NodeDriverSession) SetGenesisValidator(_auth common.Address, validatorID *big.Int, pubkey []byte, status *big.Int, createdEpoch *big.Int, createdTime *big.Int, deactivatedEpoch *big.Int, deactivatedTime *big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.SetGenesisValidator(&_NodeDriver.TransactOpts, _auth, validatorID, pubkey, status, createdEpoch, createdTime, deactivatedEpoch, deactivatedTime)
}

// SetGenesisValidator is a paid mutator transaction binding the contract method 0x4feb92f3.
//
// Solidity: function setGenesisValidator(address _auth, uint256 validatorID, bytes pubkey, uint256 status, uint256 createdEpoch, uint256 createdTime, uint256 deactivatedEpoch, uint256 deactivatedTime) returns()
func (_NodeDriver *NodeDriverTransactorSession) SetGenesisValidator(_auth common.Address, validatorID *big.Int, pubkey []byte, status *big.Int, createdEpoch *big.Int, createdTime *big.Int, deactivatedEpoch *big.Int, deactivatedTime *big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.SetGenesisValidator(&_NodeDriver.TransactOpts, _auth, validatorID, pubkey, status, createdEpoch, createdTime, deactivatedEpoch, deactivatedTime)
}

// SetStorage is a paid mutator transaction binding the contract method 0x39e503ab.
//
// Solidity: function setStorage(address acc, bytes32 key, bytes32 value) returns()
func (_NodeDriver *NodeDriverTransactor) SetStorage(opts *bind.TransactOpts, acc common.Address, key [32]byte, value [32]byte) (*types.Transaction, error) {
	return _NodeDriver.contract.Transact(opts, "setStorage", acc, key, value)
}

// SetStorage is a paid mutator transaction binding the contract method 0x39e503ab.
//
// Solidity: function setStorage(address acc, bytes32 key, bytes32 value) returns()
func (_NodeDriver *NodeDriverSession) SetStorage(acc common.Address, key [32]byte, value [32]byte) (*types.Transaction, error) {
	return _NodeDriver.Contract.SetStorage(&_NodeDriver.TransactOpts, acc, key, value)
}

// SetStorage is a paid mutator transaction binding the contract method 0x39e503ab.
//
// Solidity: function setStorage(address acc, bytes32 key, bytes32 value) returns()
func (_NodeDriver *NodeDriverTransactorSession) SetStorage(acc common.Address, key [32]byte, value [32]byte) (*types.Transaction, error) {
	return _NodeDriver.Contract.SetStorage(&_NodeDriver.TransactOpts, acc, key, value)
}

// SwapCode is a paid mutator transaction binding the contract method 0x07690b2a.
//
// Solidity: function swapCode(address acc, address with) returns()
func (_NodeDriver *NodeDriverTransactor) SwapCode(opts *bind.TransactOpts, acc common.Address, with common.Address) (*types.Transaction, error) {
	return _NodeDriver.contract.Transact(opts, "swapCode", acc, with)
}

// SwapCode is a paid mutator transaction binding the contract method 0x07690b2a.
//
// Solidity: function swapCode(address acc, address with) returns()
func (_NodeDriver *NodeDriverSession) SwapCode(acc common.Address, with common.Address) (*types.Transaction, error) {
	return _NodeDriver.Contract.SwapCode(&_NodeDriver.TransactOpts, acc, with)
}

// SwapCode is a paid mutator transaction binding the contract method 0x07690b2a.
//
// Solidity: function swapCode(address acc, address with) returns()
func (_NodeDriver *NodeDriverTransactorSession) SwapCode(acc common.Address, with common.Address) (*types.Transaction, error) {
	return _NodeDriver.Contract.SwapCode(&_NodeDriver.TransactOpts, acc, with)
}

// UpdateNetworkRules is a paid mutator transaction binding the contract method 0xb9cc6b1c.
//
// Solidity: function updateNetworkRules(bytes diff) returns()
func (_NodeDriver *NodeDriverTransactor) UpdateNetworkRules(opts *bind.TransactOpts, diff []byte) (*types.Transaction, error) {
	return _NodeDriver.contract.Transact(opts, "updateNetworkRules", diff)
}

// UpdateNetworkRules is a paid mutator transaction binding the contract method 0xb9cc6b1c.
//
// Solidity: function updateNetworkRules(bytes diff) returns()
func (_NodeDriver *NodeDriverSession) UpdateNetworkRules(diff []byte) (*types.Transaction, error) {
	return _NodeDriver.Contract.UpdateNetworkRules(&_NodeDriver.TransactOpts, diff)
}

// UpdateNetworkRules is a paid mutator transaction binding the contract method 0xb9cc6b1c.
//
// Solidity: function updateNetworkRules(bytes diff) returns()
func (_NodeDriver *NodeDriverTransactorSession) UpdateNetworkRules(diff []byte) (*types.Transaction, error) {
	return _NodeDriver.Contract.UpdateNetworkRules(&_NodeDriver.TransactOpts, diff)
}

// UpdateNetworkVersion is a paid mutator transaction binding the contract method 0x267ab446.
//
// Solidity: function updateNetworkVersion(uint256 version) returns()
func (_NodeDriver *NodeDriverTransactor) UpdateNetworkVersion(opts *bind.TransactOpts, version *big.Int) (*types.Transaction, error) {
	return _NodeDriver.contract.Transact(opts, "updateNetworkVersion", version)
}

// UpdateNetworkVersion is a paid mutator transaction binding the contract method 0x267ab446.
//
// Solidity: function updateNetworkVersion(uint256 version) returns()
func (_NodeDriver *NodeDriverSession) UpdateNetworkVersion(version *big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.UpdateNetworkVersion(&_NodeDriver.TransactOpts, version)
}

// UpdateNetworkVersion is a paid mutator transaction binding the contract method 0x267ab446.
//
// Solidity: function updateNetworkVersion(uint256 version) returns()
func (_NodeDriver *NodeDriverTransactorSession) UpdateNetworkVersion(version *big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.UpdateNetworkVersion(&_NodeDriver.TransactOpts, version)
}

// UpdateValidatorPubkey is a paid mutator transaction binding the contract method 0x242a6e3f.
//
// Solidity: function updateValidatorPubkey(uint256 validatorID, bytes pubkey) returns()
func (_NodeDriver *NodeDriverTransactor) UpdateValidatorPubkey(opts *bind.TransactOpts, validatorID *big.Int, pubkey []byte) (*types.Transaction, error) {
	return _NodeDriver.contract.Transact(opts, "updateValidatorPubkey", validatorID, pubkey)
}

// UpdateValidatorPubkey is a paid mutator transaction binding the contract method 0x242a6e3f.
//
// Solidity: function updateValidatorPubkey(uint256 validatorID, bytes pubkey) returns()
func (_NodeDriver *NodeDriverSession) UpdateValidatorPubkey(validatorID *big.Int, pubkey []byte) (*types.Transaction, error) {
	return _NodeDriver.Contract.UpdateValidatorPubkey(&_NodeDriver.TransactOpts, validatorID, pubkey)
}

// UpdateValidatorPubkey is a paid mutator transaction binding the contract method 0x242a6e3f.
//
// Solidity: function updateValidatorPubkey(uint256 validatorID, bytes pubkey) returns()
func (_NodeDriver *NodeDriverTransactorSession) UpdateValidatorPubkey(validatorID *big.Int, pubkey []byte) (*types.Transaction, error) {
	return _NodeDriver.Contract.UpdateValidatorPubkey(&_NodeDriver.TransactOpts, validatorID, pubkey)
}

// UpdateValidatorWeight is a paid mutator transaction binding the contract method 0xa4066fbe.
//
// Solidity: function updateValidatorWeight(uint256 validatorID, uint256 value) returns()
func (_NodeDriver *NodeDriverTransactor) UpdateValidatorWeight(opts *bind.TransactOpts, validatorID *big.Int, value *big.Int) (*types.Transaction, error) {
	return _NodeDriver.contract.Transact(opts, "updateValidatorWeight", validatorID, value)
}

// UpdateValidatorWeight is a paid mutator transaction binding the contract method 0xa4066fbe.
//
// Solidity: function updateValidatorWeight(uint256 validatorID, uint256 value) returns()
func (_NodeDriver *NodeDriverSession) UpdateValidatorWeight(validatorID *big.Int, value *big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.UpdateValidatorWeight(&_NodeDriver.TransactOpts, validatorID, value)
}

// UpdateValidatorWeight is a paid mutator transaction binding the contract method 0xa4066fbe.
//
// Solidity: function updateValidatorWeight(uint256 validatorID, uint256 value) returns()
func (_NodeDriver *NodeDriverTransactorSession) UpdateValidatorWeight(validatorID *big.Int, value *big.Int) (*types.Transaction, error) {
	return _NodeDriver.Contract.UpdateValidatorWeight(&_NodeDriver.TransactOpts, validatorID, value)
}

// NodeDriverAdvanceEpochsIterator is returned from FilterAdvanceEpochs and is used to iterate over the raw logs and unpacked data for AdvanceEpochs events raised by the NodeDriver contract.
type NodeDriverAdvanceEpochsIterator struct {
	Event *NodeDriverAdvanceEpochs // Event containing the contract specifics and raw log

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
func (it *NodeDriverAdvanceEpochsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeDriverAdvanceEpochs)
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
		it.Event = new(NodeDriverAdvanceEpochs)
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
func (it *NodeDriverAdvanceEpochsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeDriverAdvanceEpochsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeDriverAdvanceEpochs represents a AdvanceEpochs event raised by the NodeDriver contract.
type NodeDriverAdvanceEpochs struct {
	Num *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterAdvanceEpochs is a free log retrieval operation binding the contract event 0x0151256d62457b809bbc891b1f81c6dd0b9987552c70ce915b519750cd434dd1.
//
// Solidity: event AdvanceEpochs(uint256 num)
func (_NodeDriver *NodeDriverFilterer) FilterAdvanceEpochs(opts *bind.FilterOpts) (*NodeDriverAdvanceEpochsIterator, error) {

	logs, sub, err := _NodeDriver.contract.FilterLogs(opts, "AdvanceEpochs")
	if err != nil {
		return nil, err
	}
	return &NodeDriverAdvanceEpochsIterator{contract: _NodeDriver.contract, event: "AdvanceEpochs", logs: logs, sub: sub}, nil
}

// WatchAdvanceEpochs is a free log subscription operation binding the contract event 0x0151256d62457b809bbc891b1f81c6dd0b9987552c70ce915b519750cd434dd1.
//
// Solidity: event AdvanceEpochs(uint256 num)
func (_NodeDriver *NodeDriverFilterer) WatchAdvanceEpochs(opts *bind.WatchOpts, sink chan<- *NodeDriverAdvanceEpochs) (event.Subscription, error) {

	logs, sub, err := _NodeDriver.contract.WatchLogs(opts, "AdvanceEpochs")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeDriverAdvanceEpochs)
				if err := _NodeDriver.contract.UnpackLog(event, "AdvanceEpochs", log); err != nil {
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

// ParseAdvanceEpochs is a log parse operation binding the contract event 0x0151256d62457b809bbc891b1f81c6dd0b9987552c70ce915b519750cd434dd1.
//
// Solidity: event AdvanceEpochs(uint256 num)
func (_NodeDriver *NodeDriverFilterer) ParseAdvanceEpochs(log types.Log) (*NodeDriverAdvanceEpochs, error) {
	event := new(NodeDriverAdvanceEpochs)
	if err := _NodeDriver.contract.UnpackLog(event, "AdvanceEpochs", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeDriverUpdateNetworkRulesIterator is returned from FilterUpdateNetworkRules and is used to iterate over the raw logs and unpacked data for UpdateNetworkRules events raised by the NodeDriver contract.
type NodeDriverUpdateNetworkRulesIterator struct {
	Event *NodeDriverUpdateNetworkRules // Event containing the contract specifics and raw log

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
func (it *NodeDriverUpdateNetworkRulesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeDriverUpdateNetworkRules)
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
		it.Event = new(NodeDriverUpdateNetworkRules)
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
func (it *NodeDriverUpdateNetworkRulesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeDriverUpdateNetworkRulesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeDriverUpdateNetworkRules represents a UpdateNetworkRules event raised by the NodeDriver contract.
type NodeDriverUpdateNetworkRules struct {
	Diff []byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterUpdateNetworkRules is a free log retrieval operation binding the contract event 0x47d10eed096a44e3d0abc586c7e3a5d6cb5358cc90e7d437cd0627f7e765fb99.
//
// Solidity: event UpdateNetworkRules(bytes diff)
func (_NodeDriver *NodeDriverFilterer) FilterUpdateNetworkRules(opts *bind.FilterOpts) (*NodeDriverUpdateNetworkRulesIterator, error) {

	logs, sub, err := _NodeDriver.contract.FilterLogs(opts, "UpdateNetworkRules")
	if err != nil {
		return nil, err
	}
	return &NodeDriverUpdateNetworkRulesIterator{contract: _NodeDriver.contract, event: "UpdateNetworkRules", logs: logs, sub: sub}, nil
}

// WatchUpdateNetworkRules is a free log subscription operation binding the contract event 0x47d10eed096a44e3d0abc586c7e3a5d6cb5358cc90e7d437cd0627f7e765fb99.
//
// Solidity: event UpdateNetworkRules(bytes diff)
func (_NodeDriver *NodeDriverFilterer) WatchUpdateNetworkRules(opts *bind.WatchOpts, sink chan<- *NodeDriverUpdateNetworkRules) (event.Subscription, error) {

	logs, sub, err := _NodeDriver.contract.WatchLogs(opts, "UpdateNetworkRules")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeDriverUpdateNetworkRules)
				if err := _NodeDriver.contract.UnpackLog(event, "UpdateNetworkRules", log); err != nil {
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

// ParseUpdateNetworkRules is a log parse operation binding the contract event 0x47d10eed096a44e3d0abc586c7e3a5d6cb5358cc90e7d437cd0627f7e765fb99.
//
// Solidity: event UpdateNetworkRules(bytes diff)
func (_NodeDriver *NodeDriverFilterer) ParseUpdateNetworkRules(log types.Log) (*NodeDriverUpdateNetworkRules, error) {
	event := new(NodeDriverUpdateNetworkRules)
	if err := _NodeDriver.contract.UnpackLog(event, "UpdateNetworkRules", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeDriverUpdateNetworkVersionIterator is returned from FilterUpdateNetworkVersion and is used to iterate over the raw logs and unpacked data for UpdateNetworkVersion events raised by the NodeDriver contract.
type NodeDriverUpdateNetworkVersionIterator struct {
	Event *NodeDriverUpdateNetworkVersion // Event containing the contract specifics and raw log

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
func (it *NodeDriverUpdateNetworkVersionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeDriverUpdateNetworkVersion)
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
		it.Event = new(NodeDriverUpdateNetworkVersion)
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
func (it *NodeDriverUpdateNetworkVersionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeDriverUpdateNetworkVersionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeDriverUpdateNetworkVersion represents a UpdateNetworkVersion event raised by the NodeDriver contract.
type NodeDriverUpdateNetworkVersion struct {
	Version *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUpdateNetworkVersion is a free log retrieval operation binding the contract event 0x2ccdfd47cf0c1f1069d949f1789bb79b2f12821f021634fc835af1de66ea2feb.
//
// Solidity: event UpdateNetworkVersion(uint256 version)
func (_NodeDriver *NodeDriverFilterer) FilterUpdateNetworkVersion(opts *bind.FilterOpts) (*NodeDriverUpdateNetworkVersionIterator, error) {

	logs, sub, err := _NodeDriver.contract.FilterLogs(opts, "UpdateNetworkVersion")
	if err != nil {
		return nil, err
	}
	return &NodeDriverUpdateNetworkVersionIterator{contract: _NodeDriver.contract, event: "UpdateNetworkVersion", logs: logs, sub: sub}, nil
}

// WatchUpdateNetworkVersion is a free log subscription operation binding the contract event 0x2ccdfd47cf0c1f1069d949f1789bb79b2f12821f021634fc835af1de66ea2feb.
//
// Solidity: event UpdateNetworkVersion(uint256 version)
func (_NodeDriver *NodeDriverFilterer) WatchUpdateNetworkVersion(opts *bind.WatchOpts, sink chan<- *NodeDriverUpdateNetworkVersion) (event.Subscription, error) {

	logs, sub, err := _NodeDriver.contract.WatchLogs(opts, "UpdateNetworkVersion")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeDriverUpdateNetworkVersion)
				if err := _NodeDriver.contract.UnpackLog(event, "UpdateNetworkVersion", log); err != nil {
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

// ParseUpdateNetworkVersion is a log parse operation binding the contract event 0x2ccdfd47cf0c1f1069d949f1789bb79b2f12821f021634fc835af1de66ea2feb.
//
// Solidity: event UpdateNetworkVersion(uint256 version)
func (_NodeDriver *NodeDriverFilterer) ParseUpdateNetworkVersion(log types.Log) (*NodeDriverUpdateNetworkVersion, error) {
	event := new(NodeDriverUpdateNetworkVersion)
	if err := _NodeDriver.contract.UnpackLog(event, "UpdateNetworkVersion", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeDriverUpdateValidatorPubkeyIterator is returned from FilterUpdateValidatorPubkey and is used to iterate over the raw logs and unpacked data for UpdateValidatorPubkey events raised by the NodeDriver contract.
type NodeDriverUpdateValidatorPubkeyIterator struct {
	Event *NodeDriverUpdateValidatorPubkey // Event containing the contract specifics and raw log

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
func (it *NodeDriverUpdateValidatorPubkeyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeDriverUpdateValidatorPubkey)
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
		it.Event = new(NodeDriverUpdateValidatorPubkey)
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
func (it *NodeDriverUpdateValidatorPubkeyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeDriverUpdateValidatorPubkeyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeDriverUpdateValidatorPubkey represents a UpdateValidatorPubkey event raised by the NodeDriver contract.
type NodeDriverUpdateValidatorPubkey struct {
	ValidatorID *big.Int
	Pubkey      []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUpdateValidatorPubkey is a free log retrieval operation binding the contract event 0x0f0ef1ab97439def0a9d2c6d9dc166207f1b13b99e62b442b2993d6153c63a6e.
//
// Solidity: event UpdateValidatorPubkey(uint256 indexed validatorID, bytes pubkey)
func (_NodeDriver *NodeDriverFilterer) FilterUpdateValidatorPubkey(opts *bind.FilterOpts, validatorID []*big.Int) (*NodeDriverUpdateValidatorPubkeyIterator, error) {

	var validatorIDRule []interface{}
	for _, validatorIDItem := range validatorID {
		validatorIDRule = append(validatorIDRule, validatorIDItem)
	}

	logs, sub, err := _NodeDriver.contract.FilterLogs(opts, "UpdateValidatorPubkey", validatorIDRule)
	if err != nil {
		return nil, err
	}
	return &NodeDriverUpdateValidatorPubkeyIterator{contract: _NodeDriver.contract, event: "UpdateValidatorPubkey", logs: logs, sub: sub}, nil
}

// WatchUpdateValidatorPubkey is a free log subscription operation binding the contract event 0x0f0ef1ab97439def0a9d2c6d9dc166207f1b13b99e62b442b2993d6153c63a6e.
//
// Solidity: event UpdateValidatorPubkey(uint256 indexed validatorID, bytes pubkey)
func (_NodeDriver *NodeDriverFilterer) WatchUpdateValidatorPubkey(opts *bind.WatchOpts, sink chan<- *NodeDriverUpdateValidatorPubkey, validatorID []*big.Int) (event.Subscription, error) {

	var validatorIDRule []interface{}
	for _, validatorIDItem := range validatorID {
		validatorIDRule = append(validatorIDRule, validatorIDItem)
	}

	logs, sub, err := _NodeDriver.contract.WatchLogs(opts, "UpdateValidatorPubkey", validatorIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeDriverUpdateValidatorPubkey)
				if err := _NodeDriver.contract.UnpackLog(event, "UpdateValidatorPubkey", log); err != nil {
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

// ParseUpdateValidatorPubkey is a log parse operation binding the contract event 0x0f0ef1ab97439def0a9d2c6d9dc166207f1b13b99e62b442b2993d6153c63a6e.
//
// Solidity: event UpdateValidatorPubkey(uint256 indexed validatorID, bytes pubkey)
func (_NodeDriver *NodeDriverFilterer) ParseUpdateValidatorPubkey(log types.Log) (*NodeDriverUpdateValidatorPubkey, error) {
	event := new(NodeDriverUpdateValidatorPubkey)
	if err := _NodeDriver.contract.UnpackLog(event, "UpdateValidatorPubkey", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeDriverUpdateValidatorWeightIterator is returned from FilterUpdateValidatorWeight and is used to iterate over the raw logs and unpacked data for UpdateValidatorWeight events raised by the NodeDriver contract.
type NodeDriverUpdateValidatorWeightIterator struct {
	Event *NodeDriverUpdateValidatorWeight // Event containing the contract specifics and raw log

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
func (it *NodeDriverUpdateValidatorWeightIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeDriverUpdateValidatorWeight)
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
		it.Event = new(NodeDriverUpdateValidatorWeight)
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
func (it *NodeDriverUpdateValidatorWeightIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeDriverUpdateValidatorWeightIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeDriverUpdateValidatorWeight represents a UpdateValidatorWeight event raised by the NodeDriver contract.
type NodeDriverUpdateValidatorWeight struct {
	ValidatorID *big.Int
	Weight      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUpdateValidatorWeight is a free log retrieval operation binding the contract event 0xb975807576e3b1461be7de07ebf7d20e4790ed802d7a0c4fdd0a1a13df72a935.
//
// Solidity: event UpdateValidatorWeight(uint256 indexed validatorID, uint256 weight)
func (_NodeDriver *NodeDriverFilterer) FilterUpdateValidatorWeight(opts *bind.FilterOpts, validatorID []*big.Int) (*NodeDriverUpdateValidatorWeightIterator, error) {

	var validatorIDRule []interface{}
	for _, validatorIDItem := range validatorID {
		validatorIDRule = append(validatorIDRule, validatorIDItem)
	}

	logs, sub, err := _NodeDriver.contract.FilterLogs(opts, "UpdateValidatorWeight", validatorIDRule)
	if err != nil {
		return nil, err
	}
	return &NodeDriverUpdateValidatorWeightIterator{contract: _NodeDriver.contract, event: "UpdateValidatorWeight", logs: logs, sub: sub}, nil
}

// WatchUpdateValidatorWeight is a free log subscription operation binding the contract event 0xb975807576e3b1461be7de07ebf7d20e4790ed802d7a0c4fdd0a1a13df72a935.
//
// Solidity: event UpdateValidatorWeight(uint256 indexed validatorID, uint256 weight)
func (_NodeDriver *NodeDriverFilterer) WatchUpdateValidatorWeight(opts *bind.WatchOpts, sink chan<- *NodeDriverUpdateValidatorWeight, validatorID []*big.Int) (event.Subscription, error) {

	var validatorIDRule []interface{}
	for _, validatorIDItem := range validatorID {
		validatorIDRule = append(validatorIDRule, validatorIDItem)
	}

	logs, sub, err := _NodeDriver.contract.WatchLogs(opts, "UpdateValidatorWeight", validatorIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeDriverUpdateValidatorWeight)
				if err := _NodeDriver.contract.UnpackLog(event, "UpdateValidatorWeight", log); err != nil {
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

// ParseUpdateValidatorWeight is a log parse operation binding the contract event 0xb975807576e3b1461be7de07ebf7d20e4790ed802d7a0c4fdd0a1a13df72a935.
//
// Solidity: event UpdateValidatorWeight(uint256 indexed validatorID, uint256 weight)
func (_NodeDriver *NodeDriverFilterer) ParseUpdateValidatorWeight(log types.Log) (*NodeDriverUpdateValidatorWeight, error) {
	event := new(NodeDriverUpdateValidatorWeight)
	if err := _NodeDriver.contract.UnpackLog(event, "UpdateValidatorWeight", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeDriverUpdatedBackendIterator is returned from FilterUpdatedBackend and is used to iterate over the raw logs and unpacked data for UpdatedBackend events raised by the NodeDriver contract.
type NodeDriverUpdatedBackendIterator struct {
	Event *NodeDriverUpdatedBackend // Event containing the contract specifics and raw log

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
func (it *NodeDriverUpdatedBackendIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeDriverUpdatedBackend)
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
		it.Event = new(NodeDriverUpdatedBackend)
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
func (it *NodeDriverUpdatedBackendIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeDriverUpdatedBackendIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeDriverUpdatedBackend represents a UpdatedBackend event raised by the NodeDriver contract.
type NodeDriverUpdatedBackend struct {
	Backend common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUpdatedBackend is a free log retrieval operation binding the contract event 0x64ee8f7bfc37fc205d7194ee3d64947ab7b57e663cd0d1abd3ef245035830693.
//
// Solidity: event UpdatedBackend(address indexed backend)
func (_NodeDriver *NodeDriverFilterer) FilterUpdatedBackend(opts *bind.FilterOpts, backend []common.Address) (*NodeDriverUpdatedBackendIterator, error) {

	var backendRule []interface{}
	for _, backendItem := range backend {
		backendRule = append(backendRule, backendItem)
	}

	logs, sub, err := _NodeDriver.contract.FilterLogs(opts, "UpdatedBackend", backendRule)
	if err != nil {
		return nil, err
	}
	return &NodeDriverUpdatedBackendIterator{contract: _NodeDriver.contract, event: "UpdatedBackend", logs: logs, sub: sub}, nil
}

// WatchUpdatedBackend is a free log subscription operation binding the contract event 0x64ee8f7bfc37fc205d7194ee3d64947ab7b57e663cd0d1abd3ef245035830693.
//
// Solidity: event UpdatedBackend(address indexed backend)
func (_NodeDriver *NodeDriverFilterer) WatchUpdatedBackend(opts *bind.WatchOpts, sink chan<- *NodeDriverUpdatedBackend, backend []common.Address) (event.Subscription, error) {

	var backendRule []interface{}
	for _, backendItem := range backend {
		backendRule = append(backendRule, backendItem)
	}

	logs, sub, err := _NodeDriver.contract.WatchLogs(opts, "UpdatedBackend", backendRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeDriverUpdatedBackend)
				if err := _NodeDriver.contract.UnpackLog(event, "UpdatedBackend", log); err != nil {
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

// ParseUpdatedBackend is a log parse operation binding the contract event 0x64ee8f7bfc37fc205d7194ee3d64947ab7b57e663cd0d1abd3ef245035830693.
//
// Solidity: event UpdatedBackend(address indexed backend)
func (_NodeDriver *NodeDriverFilterer) ParseUpdatedBackend(log types.Log) (*NodeDriverUpdatedBackend, error) {
	event := new(NodeDriverUpdatedBackend)
	if err := _NodeDriver.contract.UnpackLog(event, "UpdatedBackend", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
