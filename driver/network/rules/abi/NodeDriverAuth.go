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

// NodeDriverAuthMetaData contains all meta data concerning the NodeDriverAuth contract.
var NodeDriverAuthMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_sfc\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_driver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newDriverAuth\",\"type\":\"address\"}],\"name\":\"migrateTo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"executable\",\"type\":\"address\"}],\"name\":\"execute\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"executable\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"selfCodeHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"driverCodeHash\",\"type\":\"bytes32\"}],\"name\":\"mutExecute\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"acc\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"diff\",\"type\":\"uint256\"}],\"name\":\"incBalance\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"acc\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"upgradeCode\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"acc\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"copyCode\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"acc\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"diff\",\"type\":\"uint256\"}],\"name\":\"incNonce\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"diff\",\"type\":\"bytes\"}],\"name\":\"updateNetworkRules\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minGasPrice\",\"type\":\"uint256\"}],\"name\":\"updateMinGasPrice\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"version\",\"type\":\"uint256\"}],\"name\":\"updateNetworkVersion\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"name\":\"advanceEpochs\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"updateValidatorWeight\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"}],\"name\":\"updateValidatorPubkey\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_auth\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"status\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivatedEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivatedTime\",\"type\":\"uint256\"}],\"name\":\"setGenesisValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"toValidatorID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockupFromEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockupEndTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockupDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"earlyUnlockPenalty\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rewards\",\"type\":\"uint256\"}],\"name\":\"setGenesisDelegation\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"status\",\"type\":\"uint256\"}],\"name\":\"deactivateValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"nextValidatorIDs\",\"type\":\"uint256[]\"}],\"name\":\"sealEpochValidators\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"offlineTimes\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"offlineBlocks\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"uptimes\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"originatedTxsFee\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"usedGas\",\"type\":\"uint256\"}],\"name\":\"sealEpoch\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// NodeDriverAuthABI is the input ABI used to generate the binding from.
// Deprecated: Use NodeDriverAuthMetaData.ABI instead.
var NodeDriverAuthABI = NodeDriverAuthMetaData.ABI

// NodeDriverAuth is an auto generated Go binding around an Ethereum contract.
type NodeDriverAuth struct {
	NodeDriverAuthCaller     // Read-only binding to the contract
	NodeDriverAuthTransactor // Write-only binding to the contract
	NodeDriverAuthFilterer   // Log filterer for contract events
}

// NodeDriverAuthCaller is an auto generated read-only Go binding around an Ethereum contract.
type NodeDriverAuthCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeDriverAuthTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NodeDriverAuthTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeDriverAuthFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NodeDriverAuthFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeDriverAuthSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NodeDriverAuthSession struct {
	Contract     *NodeDriverAuth   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NodeDriverAuthCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NodeDriverAuthCallerSession struct {
	Contract *NodeDriverAuthCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// NodeDriverAuthTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NodeDriverAuthTransactorSession struct {
	Contract     *NodeDriverAuthTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// NodeDriverAuthRaw is an auto generated low-level Go binding around an Ethereum contract.
type NodeDriverAuthRaw struct {
	Contract *NodeDriverAuth // Generic contract binding to access the raw methods on
}

// NodeDriverAuthCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NodeDriverAuthCallerRaw struct {
	Contract *NodeDriverAuthCaller // Generic read-only contract binding to access the raw methods on
}

// NodeDriverAuthTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NodeDriverAuthTransactorRaw struct {
	Contract *NodeDriverAuthTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNodeDriverAuth creates a new instance of NodeDriverAuth, bound to a specific deployed contract.
func NewNodeDriverAuth(address common.Address, backend bind.ContractBackend) (*NodeDriverAuth, error) {
	contract, err := bindNodeDriverAuth(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NodeDriverAuth{NodeDriverAuthCaller: NodeDriverAuthCaller{contract: contract}, NodeDriverAuthTransactor: NodeDriverAuthTransactor{contract: contract}, NodeDriverAuthFilterer: NodeDriverAuthFilterer{contract: contract}}, nil
}

// NewNodeDriverAuthCaller creates a new read-only instance of NodeDriverAuth, bound to a specific deployed contract.
func NewNodeDriverAuthCaller(address common.Address, caller bind.ContractCaller) (*NodeDriverAuthCaller, error) {
	contract, err := bindNodeDriverAuth(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NodeDriverAuthCaller{contract: contract}, nil
}

// NewNodeDriverAuthTransactor creates a new write-only instance of NodeDriverAuth, bound to a specific deployed contract.
func NewNodeDriverAuthTransactor(address common.Address, transactor bind.ContractTransactor) (*NodeDriverAuthTransactor, error) {
	contract, err := bindNodeDriverAuth(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NodeDriverAuthTransactor{contract: contract}, nil
}

// NewNodeDriverAuthFilterer creates a new log filterer instance of NodeDriverAuth, bound to a specific deployed contract.
func NewNodeDriverAuthFilterer(address common.Address, filterer bind.ContractFilterer) (*NodeDriverAuthFilterer, error) {
	contract, err := bindNodeDriverAuth(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NodeDriverAuthFilterer{contract: contract}, nil
}

// bindNodeDriverAuth binds a generic wrapper to an already deployed contract.
func bindNodeDriverAuth(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NodeDriverAuthABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeDriverAuth *NodeDriverAuthRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeDriverAuth.Contract.NodeDriverAuthCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeDriverAuth *NodeDriverAuthRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.NodeDriverAuthTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeDriverAuth *NodeDriverAuthRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.NodeDriverAuthTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeDriverAuth *NodeDriverAuthCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeDriverAuth.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeDriverAuth *NodeDriverAuthTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeDriverAuth *NodeDriverAuthTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.contract.Transact(opts, method, params...)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_NodeDriverAuth *NodeDriverAuthCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _NodeDriverAuth.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_NodeDriverAuth *NodeDriverAuthSession) IsOwner() (bool, error) {
	return _NodeDriverAuth.Contract.IsOwner(&_NodeDriverAuth.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_NodeDriverAuth *NodeDriverAuthCallerSession) IsOwner() (bool, error) {
	return _NodeDriverAuth.Contract.IsOwner(&_NodeDriverAuth.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NodeDriverAuth *NodeDriverAuthCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NodeDriverAuth.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NodeDriverAuth *NodeDriverAuthSession) Owner() (common.Address, error) {
	return _NodeDriverAuth.Contract.Owner(&_NodeDriverAuth.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NodeDriverAuth *NodeDriverAuthCallerSession) Owner() (common.Address, error) {
	return _NodeDriverAuth.Contract.Owner(&_NodeDriverAuth.CallOpts)
}

// AdvanceEpochs is a paid mutator transaction binding the contract method 0x0aeeca00.
//
// Solidity: function advanceEpochs(uint256 num) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactor) AdvanceEpochs(opts *bind.TransactOpts, num *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.contract.Transact(opts, "advanceEpochs", num)
}

// AdvanceEpochs is a paid mutator transaction binding the contract method 0x0aeeca00.
//
// Solidity: function advanceEpochs(uint256 num) returns()
func (_NodeDriverAuth *NodeDriverAuthSession) AdvanceEpochs(num *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.AdvanceEpochs(&_NodeDriverAuth.TransactOpts, num)
}

// AdvanceEpochs is a paid mutator transaction binding the contract method 0x0aeeca00.
//
// Solidity: function advanceEpochs(uint256 num) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactorSession) AdvanceEpochs(num *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.AdvanceEpochs(&_NodeDriverAuth.TransactOpts, num)
}

// CopyCode is a paid mutator transaction binding the contract method 0xd6a0c7af.
//
// Solidity: function copyCode(address acc, address from) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactor) CopyCode(opts *bind.TransactOpts, acc common.Address, from common.Address) (*types.Transaction, error) {
	return _NodeDriverAuth.contract.Transact(opts, "copyCode", acc, from)
}

// CopyCode is a paid mutator transaction binding the contract method 0xd6a0c7af.
//
// Solidity: function copyCode(address acc, address from) returns()
func (_NodeDriverAuth *NodeDriverAuthSession) CopyCode(acc common.Address, from common.Address) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.CopyCode(&_NodeDriverAuth.TransactOpts, acc, from)
}

// CopyCode is a paid mutator transaction binding the contract method 0xd6a0c7af.
//
// Solidity: function copyCode(address acc, address from) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactorSession) CopyCode(acc common.Address, from common.Address) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.CopyCode(&_NodeDriverAuth.TransactOpts, acc, from)
}

// DeactivateValidator is a paid mutator transaction binding the contract method 0x1e702f83.
//
// Solidity: function deactivateValidator(uint256 validatorID, uint256 status) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactor) DeactivateValidator(opts *bind.TransactOpts, validatorID *big.Int, status *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.contract.Transact(opts, "deactivateValidator", validatorID, status)
}

// DeactivateValidator is a paid mutator transaction binding the contract method 0x1e702f83.
//
// Solidity: function deactivateValidator(uint256 validatorID, uint256 status) returns()
func (_NodeDriverAuth *NodeDriverAuthSession) DeactivateValidator(validatorID *big.Int, status *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.DeactivateValidator(&_NodeDriverAuth.TransactOpts, validatorID, status)
}

// DeactivateValidator is a paid mutator transaction binding the contract method 0x1e702f83.
//
// Solidity: function deactivateValidator(uint256 validatorID, uint256 status) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactorSession) DeactivateValidator(validatorID *big.Int, status *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.DeactivateValidator(&_NodeDriverAuth.TransactOpts, validatorID, status)
}

// Execute is a paid mutator transaction binding the contract method 0x4b64e492.
//
// Solidity: function execute(address executable) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactor) Execute(opts *bind.TransactOpts, executable common.Address) (*types.Transaction, error) {
	return _NodeDriverAuth.contract.Transact(opts, "execute", executable)
}

// Execute is a paid mutator transaction binding the contract method 0x4b64e492.
//
// Solidity: function execute(address executable) returns()
func (_NodeDriverAuth *NodeDriverAuthSession) Execute(executable common.Address) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.Execute(&_NodeDriverAuth.TransactOpts, executable)
}

// Execute is a paid mutator transaction binding the contract method 0x4b64e492.
//
// Solidity: function execute(address executable) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactorSession) Execute(executable common.Address) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.Execute(&_NodeDriverAuth.TransactOpts, executable)
}

// IncBalance is a paid mutator transaction binding the contract method 0x66e7ea0f.
//
// Solidity: function incBalance(address acc, uint256 diff) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactor) IncBalance(opts *bind.TransactOpts, acc common.Address, diff *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.contract.Transact(opts, "incBalance", acc, diff)
}

// IncBalance is a paid mutator transaction binding the contract method 0x66e7ea0f.
//
// Solidity: function incBalance(address acc, uint256 diff) returns()
func (_NodeDriverAuth *NodeDriverAuthSession) IncBalance(acc common.Address, diff *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.IncBalance(&_NodeDriverAuth.TransactOpts, acc, diff)
}

// IncBalance is a paid mutator transaction binding the contract method 0x66e7ea0f.
//
// Solidity: function incBalance(address acc, uint256 diff) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactorSession) IncBalance(acc common.Address, diff *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.IncBalance(&_NodeDriverAuth.TransactOpts, acc, diff)
}

// IncNonce is a paid mutator transaction binding the contract method 0x79bead38.
//
// Solidity: function incNonce(address acc, uint256 diff) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactor) IncNonce(opts *bind.TransactOpts, acc common.Address, diff *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.contract.Transact(opts, "incNonce", acc, diff)
}

// IncNonce is a paid mutator transaction binding the contract method 0x79bead38.
//
// Solidity: function incNonce(address acc, uint256 diff) returns()
func (_NodeDriverAuth *NodeDriverAuthSession) IncNonce(acc common.Address, diff *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.IncNonce(&_NodeDriverAuth.TransactOpts, acc, diff)
}

// IncNonce is a paid mutator transaction binding the contract method 0x79bead38.
//
// Solidity: function incNonce(address acc, uint256 diff) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactorSession) IncNonce(acc common.Address, diff *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.IncNonce(&_NodeDriverAuth.TransactOpts, acc, diff)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address _sfc, address _driver, address _owner) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactor) Initialize(opts *bind.TransactOpts, _sfc common.Address, _driver common.Address, _owner common.Address) (*types.Transaction, error) {
	return _NodeDriverAuth.contract.Transact(opts, "initialize", _sfc, _driver, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address _sfc, address _driver, address _owner) returns()
func (_NodeDriverAuth *NodeDriverAuthSession) Initialize(_sfc common.Address, _driver common.Address, _owner common.Address) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.Initialize(&_NodeDriverAuth.TransactOpts, _sfc, _driver, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address _sfc, address _driver, address _owner) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactorSession) Initialize(_sfc common.Address, _driver common.Address, _owner common.Address) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.Initialize(&_NodeDriverAuth.TransactOpts, _sfc, _driver, _owner)
}

// MigrateTo is a paid mutator transaction binding the contract method 0x4ddaf8f2.
//
// Solidity: function migrateTo(address newDriverAuth) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactor) MigrateTo(opts *bind.TransactOpts, newDriverAuth common.Address) (*types.Transaction, error) {
	return _NodeDriverAuth.contract.Transact(opts, "migrateTo", newDriverAuth)
}

// MigrateTo is a paid mutator transaction binding the contract method 0x4ddaf8f2.
//
// Solidity: function migrateTo(address newDriverAuth) returns()
func (_NodeDriverAuth *NodeDriverAuthSession) MigrateTo(newDriverAuth common.Address) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.MigrateTo(&_NodeDriverAuth.TransactOpts, newDriverAuth)
}

// MigrateTo is a paid mutator transaction binding the contract method 0x4ddaf8f2.
//
// Solidity: function migrateTo(address newDriverAuth) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactorSession) MigrateTo(newDriverAuth common.Address) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.MigrateTo(&_NodeDriverAuth.TransactOpts, newDriverAuth)
}

// MutExecute is a paid mutator transaction binding the contract method 0x1cef4fab.
//
// Solidity: function mutExecute(address executable, address newOwner, bytes32 selfCodeHash, bytes32 driverCodeHash) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactor) MutExecute(opts *bind.TransactOpts, executable common.Address, newOwner common.Address, selfCodeHash [32]byte, driverCodeHash [32]byte) (*types.Transaction, error) {
	return _NodeDriverAuth.contract.Transact(opts, "mutExecute", executable, newOwner, selfCodeHash, driverCodeHash)
}

// MutExecute is a paid mutator transaction binding the contract method 0x1cef4fab.
//
// Solidity: function mutExecute(address executable, address newOwner, bytes32 selfCodeHash, bytes32 driverCodeHash) returns()
func (_NodeDriverAuth *NodeDriverAuthSession) MutExecute(executable common.Address, newOwner common.Address, selfCodeHash [32]byte, driverCodeHash [32]byte) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.MutExecute(&_NodeDriverAuth.TransactOpts, executable, newOwner, selfCodeHash, driverCodeHash)
}

// MutExecute is a paid mutator transaction binding the contract method 0x1cef4fab.
//
// Solidity: function mutExecute(address executable, address newOwner, bytes32 selfCodeHash, bytes32 driverCodeHash) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactorSession) MutExecute(executable common.Address, newOwner common.Address, selfCodeHash [32]byte, driverCodeHash [32]byte) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.MutExecute(&_NodeDriverAuth.TransactOpts, executable, newOwner, selfCodeHash, driverCodeHash)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NodeDriverAuth *NodeDriverAuthTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeDriverAuth.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NodeDriverAuth *NodeDriverAuthSession) RenounceOwnership() (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.RenounceOwnership(&_NodeDriverAuth.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NodeDriverAuth *NodeDriverAuthTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.RenounceOwnership(&_NodeDriverAuth.TransactOpts)
}

// SealEpoch is a paid mutator transaction binding the contract method 0x592fe0c0.
//
// Solidity: function sealEpoch(uint256[] offlineTimes, uint256[] offlineBlocks, uint256[] uptimes, uint256[] originatedTxsFee, uint256 usedGas) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactor) SealEpoch(opts *bind.TransactOpts, offlineTimes []*big.Int, offlineBlocks []*big.Int, uptimes []*big.Int, originatedTxsFee []*big.Int, usedGas *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.contract.Transact(opts, "sealEpoch", offlineTimes, offlineBlocks, uptimes, originatedTxsFee, usedGas)
}

// SealEpoch is a paid mutator transaction binding the contract method 0x592fe0c0.
//
// Solidity: function sealEpoch(uint256[] offlineTimes, uint256[] offlineBlocks, uint256[] uptimes, uint256[] originatedTxsFee, uint256 usedGas) returns()
func (_NodeDriverAuth *NodeDriverAuthSession) SealEpoch(offlineTimes []*big.Int, offlineBlocks []*big.Int, uptimes []*big.Int, originatedTxsFee []*big.Int, usedGas *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.SealEpoch(&_NodeDriverAuth.TransactOpts, offlineTimes, offlineBlocks, uptimes, originatedTxsFee, usedGas)
}

// SealEpoch is a paid mutator transaction binding the contract method 0x592fe0c0.
//
// Solidity: function sealEpoch(uint256[] offlineTimes, uint256[] offlineBlocks, uint256[] uptimes, uint256[] originatedTxsFee, uint256 usedGas) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactorSession) SealEpoch(offlineTimes []*big.Int, offlineBlocks []*big.Int, uptimes []*big.Int, originatedTxsFee []*big.Int, usedGas *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.SealEpoch(&_NodeDriverAuth.TransactOpts, offlineTimes, offlineBlocks, uptimes, originatedTxsFee, usedGas)
}

// SealEpochValidators is a paid mutator transaction binding the contract method 0xe08d7e66.
//
// Solidity: function sealEpochValidators(uint256[] nextValidatorIDs) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactor) SealEpochValidators(opts *bind.TransactOpts, nextValidatorIDs []*big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.contract.Transact(opts, "sealEpochValidators", nextValidatorIDs)
}

// SealEpochValidators is a paid mutator transaction binding the contract method 0xe08d7e66.
//
// Solidity: function sealEpochValidators(uint256[] nextValidatorIDs) returns()
func (_NodeDriverAuth *NodeDriverAuthSession) SealEpochValidators(nextValidatorIDs []*big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.SealEpochValidators(&_NodeDriverAuth.TransactOpts, nextValidatorIDs)
}

// SealEpochValidators is a paid mutator transaction binding the contract method 0xe08d7e66.
//
// Solidity: function sealEpochValidators(uint256[] nextValidatorIDs) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactorSession) SealEpochValidators(nextValidatorIDs []*big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.SealEpochValidators(&_NodeDriverAuth.TransactOpts, nextValidatorIDs)
}

// SetGenesisDelegation is a paid mutator transaction binding the contract method 0x18f628d4.
//
// Solidity: function setGenesisDelegation(address delegator, uint256 toValidatorID, uint256 stake, uint256 lockedStake, uint256 lockupFromEpoch, uint256 lockupEndTime, uint256 lockupDuration, uint256 earlyUnlockPenalty, uint256 rewards) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactor) SetGenesisDelegation(opts *bind.TransactOpts, delegator common.Address, toValidatorID *big.Int, stake *big.Int, lockedStake *big.Int, lockupFromEpoch *big.Int, lockupEndTime *big.Int, lockupDuration *big.Int, earlyUnlockPenalty *big.Int, rewards *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.contract.Transact(opts, "setGenesisDelegation", delegator, toValidatorID, stake, lockedStake, lockupFromEpoch, lockupEndTime, lockupDuration, earlyUnlockPenalty, rewards)
}

// SetGenesisDelegation is a paid mutator transaction binding the contract method 0x18f628d4.
//
// Solidity: function setGenesisDelegation(address delegator, uint256 toValidatorID, uint256 stake, uint256 lockedStake, uint256 lockupFromEpoch, uint256 lockupEndTime, uint256 lockupDuration, uint256 earlyUnlockPenalty, uint256 rewards) returns()
func (_NodeDriverAuth *NodeDriverAuthSession) SetGenesisDelegation(delegator common.Address, toValidatorID *big.Int, stake *big.Int, lockedStake *big.Int, lockupFromEpoch *big.Int, lockupEndTime *big.Int, lockupDuration *big.Int, earlyUnlockPenalty *big.Int, rewards *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.SetGenesisDelegation(&_NodeDriverAuth.TransactOpts, delegator, toValidatorID, stake, lockedStake, lockupFromEpoch, lockupEndTime, lockupDuration, earlyUnlockPenalty, rewards)
}

// SetGenesisDelegation is a paid mutator transaction binding the contract method 0x18f628d4.
//
// Solidity: function setGenesisDelegation(address delegator, uint256 toValidatorID, uint256 stake, uint256 lockedStake, uint256 lockupFromEpoch, uint256 lockupEndTime, uint256 lockupDuration, uint256 earlyUnlockPenalty, uint256 rewards) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactorSession) SetGenesisDelegation(delegator common.Address, toValidatorID *big.Int, stake *big.Int, lockedStake *big.Int, lockupFromEpoch *big.Int, lockupEndTime *big.Int, lockupDuration *big.Int, earlyUnlockPenalty *big.Int, rewards *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.SetGenesisDelegation(&_NodeDriverAuth.TransactOpts, delegator, toValidatorID, stake, lockedStake, lockupFromEpoch, lockupEndTime, lockupDuration, earlyUnlockPenalty, rewards)
}

// SetGenesisValidator is a paid mutator transaction binding the contract method 0x4feb92f3.
//
// Solidity: function setGenesisValidator(address _auth, uint256 validatorID, bytes pubkey, uint256 status, uint256 createdEpoch, uint256 createdTime, uint256 deactivatedEpoch, uint256 deactivatedTime) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactor) SetGenesisValidator(opts *bind.TransactOpts, _auth common.Address, validatorID *big.Int, pubkey []byte, status *big.Int, createdEpoch *big.Int, createdTime *big.Int, deactivatedEpoch *big.Int, deactivatedTime *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.contract.Transact(opts, "setGenesisValidator", _auth, validatorID, pubkey, status, createdEpoch, createdTime, deactivatedEpoch, deactivatedTime)
}

// SetGenesisValidator is a paid mutator transaction binding the contract method 0x4feb92f3.
//
// Solidity: function setGenesisValidator(address _auth, uint256 validatorID, bytes pubkey, uint256 status, uint256 createdEpoch, uint256 createdTime, uint256 deactivatedEpoch, uint256 deactivatedTime) returns()
func (_NodeDriverAuth *NodeDriverAuthSession) SetGenesisValidator(_auth common.Address, validatorID *big.Int, pubkey []byte, status *big.Int, createdEpoch *big.Int, createdTime *big.Int, deactivatedEpoch *big.Int, deactivatedTime *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.SetGenesisValidator(&_NodeDriverAuth.TransactOpts, _auth, validatorID, pubkey, status, createdEpoch, createdTime, deactivatedEpoch, deactivatedTime)
}

// SetGenesisValidator is a paid mutator transaction binding the contract method 0x4feb92f3.
//
// Solidity: function setGenesisValidator(address _auth, uint256 validatorID, bytes pubkey, uint256 status, uint256 createdEpoch, uint256 createdTime, uint256 deactivatedEpoch, uint256 deactivatedTime) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactorSession) SetGenesisValidator(_auth common.Address, validatorID *big.Int, pubkey []byte, status *big.Int, createdEpoch *big.Int, createdTime *big.Int, deactivatedEpoch *big.Int, deactivatedTime *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.SetGenesisValidator(&_NodeDriverAuth.TransactOpts, _auth, validatorID, pubkey, status, createdEpoch, createdTime, deactivatedEpoch, deactivatedTime)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _NodeDriverAuth.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NodeDriverAuth *NodeDriverAuthSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.TransferOwnership(&_NodeDriverAuth.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.TransferOwnership(&_NodeDriverAuth.TransactOpts, newOwner)
}

// UpdateMinGasPrice is a paid mutator transaction binding the contract method 0x07aaf344.
//
// Solidity: function updateMinGasPrice(uint256 minGasPrice) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactor) UpdateMinGasPrice(opts *bind.TransactOpts, minGasPrice *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.contract.Transact(opts, "updateMinGasPrice", minGasPrice)
}

// UpdateMinGasPrice is a paid mutator transaction binding the contract method 0x07aaf344.
//
// Solidity: function updateMinGasPrice(uint256 minGasPrice) returns()
func (_NodeDriverAuth *NodeDriverAuthSession) UpdateMinGasPrice(minGasPrice *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.UpdateMinGasPrice(&_NodeDriverAuth.TransactOpts, minGasPrice)
}

// UpdateMinGasPrice is a paid mutator transaction binding the contract method 0x07aaf344.
//
// Solidity: function updateMinGasPrice(uint256 minGasPrice) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactorSession) UpdateMinGasPrice(minGasPrice *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.UpdateMinGasPrice(&_NodeDriverAuth.TransactOpts, minGasPrice)
}

// UpdateNetworkRules is a paid mutator transaction binding the contract method 0xb9cc6b1c.
//
// Solidity: function updateNetworkRules(bytes diff) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactor) UpdateNetworkRules(opts *bind.TransactOpts, diff []byte) (*types.Transaction, error) {
	return _NodeDriverAuth.contract.Transact(opts, "updateNetworkRules", diff)
}

// UpdateNetworkRules is a paid mutator transaction binding the contract method 0xb9cc6b1c.
//
// Solidity: function updateNetworkRules(bytes diff) returns()
func (_NodeDriverAuth *NodeDriverAuthSession) UpdateNetworkRules(diff []byte) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.UpdateNetworkRules(&_NodeDriverAuth.TransactOpts, diff)
}

// UpdateNetworkRules is a paid mutator transaction binding the contract method 0xb9cc6b1c.
//
// Solidity: function updateNetworkRules(bytes diff) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactorSession) UpdateNetworkRules(diff []byte) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.UpdateNetworkRules(&_NodeDriverAuth.TransactOpts, diff)
}

// UpdateNetworkVersion is a paid mutator transaction binding the contract method 0x267ab446.
//
// Solidity: function updateNetworkVersion(uint256 version) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactor) UpdateNetworkVersion(opts *bind.TransactOpts, version *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.contract.Transact(opts, "updateNetworkVersion", version)
}

// UpdateNetworkVersion is a paid mutator transaction binding the contract method 0x267ab446.
//
// Solidity: function updateNetworkVersion(uint256 version) returns()
func (_NodeDriverAuth *NodeDriverAuthSession) UpdateNetworkVersion(version *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.UpdateNetworkVersion(&_NodeDriverAuth.TransactOpts, version)
}

// UpdateNetworkVersion is a paid mutator transaction binding the contract method 0x267ab446.
//
// Solidity: function updateNetworkVersion(uint256 version) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactorSession) UpdateNetworkVersion(version *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.UpdateNetworkVersion(&_NodeDriverAuth.TransactOpts, version)
}

// UpdateValidatorPubkey is a paid mutator transaction binding the contract method 0x242a6e3f.
//
// Solidity: function updateValidatorPubkey(uint256 validatorID, bytes pubkey) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactor) UpdateValidatorPubkey(opts *bind.TransactOpts, validatorID *big.Int, pubkey []byte) (*types.Transaction, error) {
	return _NodeDriverAuth.contract.Transact(opts, "updateValidatorPubkey", validatorID, pubkey)
}

// UpdateValidatorPubkey is a paid mutator transaction binding the contract method 0x242a6e3f.
//
// Solidity: function updateValidatorPubkey(uint256 validatorID, bytes pubkey) returns()
func (_NodeDriverAuth *NodeDriverAuthSession) UpdateValidatorPubkey(validatorID *big.Int, pubkey []byte) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.UpdateValidatorPubkey(&_NodeDriverAuth.TransactOpts, validatorID, pubkey)
}

// UpdateValidatorPubkey is a paid mutator transaction binding the contract method 0x242a6e3f.
//
// Solidity: function updateValidatorPubkey(uint256 validatorID, bytes pubkey) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactorSession) UpdateValidatorPubkey(validatorID *big.Int, pubkey []byte) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.UpdateValidatorPubkey(&_NodeDriverAuth.TransactOpts, validatorID, pubkey)
}

// UpdateValidatorWeight is a paid mutator transaction binding the contract method 0xa4066fbe.
//
// Solidity: function updateValidatorWeight(uint256 validatorID, uint256 value) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactor) UpdateValidatorWeight(opts *bind.TransactOpts, validatorID *big.Int, value *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.contract.Transact(opts, "updateValidatorWeight", validatorID, value)
}

// UpdateValidatorWeight is a paid mutator transaction binding the contract method 0xa4066fbe.
//
// Solidity: function updateValidatorWeight(uint256 validatorID, uint256 value) returns()
func (_NodeDriverAuth *NodeDriverAuthSession) UpdateValidatorWeight(validatorID *big.Int, value *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.UpdateValidatorWeight(&_NodeDriverAuth.TransactOpts, validatorID, value)
}

// UpdateValidatorWeight is a paid mutator transaction binding the contract method 0xa4066fbe.
//
// Solidity: function updateValidatorWeight(uint256 validatorID, uint256 value) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactorSession) UpdateValidatorWeight(validatorID *big.Int, value *big.Int) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.UpdateValidatorWeight(&_NodeDriverAuth.TransactOpts, validatorID, value)
}

// UpgradeCode is a paid mutator transaction binding the contract method 0xfd1b6ec1.
//
// Solidity: function upgradeCode(address acc, address from) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactor) UpgradeCode(opts *bind.TransactOpts, acc common.Address, from common.Address) (*types.Transaction, error) {
	return _NodeDriverAuth.contract.Transact(opts, "upgradeCode", acc, from)
}

// UpgradeCode is a paid mutator transaction binding the contract method 0xfd1b6ec1.
//
// Solidity: function upgradeCode(address acc, address from) returns()
func (_NodeDriverAuth *NodeDriverAuthSession) UpgradeCode(acc common.Address, from common.Address) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.UpgradeCode(&_NodeDriverAuth.TransactOpts, acc, from)
}

// UpgradeCode is a paid mutator transaction binding the contract method 0xfd1b6ec1.
//
// Solidity: function upgradeCode(address acc, address from) returns()
func (_NodeDriverAuth *NodeDriverAuthTransactorSession) UpgradeCode(acc common.Address, from common.Address) (*types.Transaction, error) {
	return _NodeDriverAuth.Contract.UpgradeCode(&_NodeDriverAuth.TransactOpts, acc, from)
}

// NodeDriverAuthOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the NodeDriverAuth contract.
type NodeDriverAuthOwnershipTransferredIterator struct {
	Event *NodeDriverAuthOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *NodeDriverAuthOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeDriverAuthOwnershipTransferred)
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
		it.Event = new(NodeDriverAuthOwnershipTransferred)
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
func (it *NodeDriverAuthOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeDriverAuthOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeDriverAuthOwnershipTransferred represents a OwnershipTransferred event raised by the NodeDriverAuth contract.
type NodeDriverAuthOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NodeDriverAuth *NodeDriverAuthFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*NodeDriverAuthOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NodeDriverAuth.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &NodeDriverAuthOwnershipTransferredIterator{contract: _NodeDriverAuth.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NodeDriverAuth *NodeDriverAuthFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NodeDriverAuthOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NodeDriverAuth.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeDriverAuthOwnershipTransferred)
				if err := _NodeDriverAuth.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NodeDriverAuth *NodeDriverAuthFilterer) ParseOwnershipTransferred(log types.Log) (*NodeDriverAuthOwnershipTransferred, error) {
	event := new(NodeDriverAuthOwnershipTransferred)
	if err := _NodeDriverAuth.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
