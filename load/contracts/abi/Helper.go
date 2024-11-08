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
	_ = abi.ConvertType
)

// HelperMetaData contains all meta data concerning the Helper contract.
var HelperMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"addresspayable[]\",\"name\":\"receivers\",\"type\":\"address[]\"}],\"name\":\"distribute\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600f57600080fd5b506101d48061001f6000396000f3fe60806040526004361061001e5760003560e01c80636138889b14610023575b600080fd5b6100366100313660046100bf565b610038565b005b60006100448234610136565b905060005b828110156100b95783838281811061006357610063610158565b9050602002016020810190610078919061016e565b6001600160a01b03166108fc839081150290604051600060405180830381858888f193505050501580156100b0573d6000803e3d6000fd5b50600101610049565b50505050565b600080602083850312156100d257600080fd5b823567ffffffffffffffff8111156100e957600080fd5b8301601f810185136100fa57600080fd5b803567ffffffffffffffff81111561011157600080fd5b8560208260051b840101111561012657600080fd5b6020919091019590945092505050565b60008261015357634e487b7160e01b600052601260045260246000fd5b500490565b634e487b7160e01b600052603260045260246000fd5b60006020828403121561018057600080fd5b81356001600160a01b038116811461019757600080fd5b939250505056fea26469706673582212205052e12a87dfe9a082e5acf10c6c3ff8ec18e14b73c46862856688c3f4c1b26b64736f6c634300081c0033",
}

// HelperABI is the input ABI used to generate the binding from.
// Deprecated: Use HelperMetaData.ABI instead.
var HelperABI = HelperMetaData.ABI

// HelperBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use HelperMetaData.Bin instead.
var HelperBin = HelperMetaData.Bin

// DeployHelper deploys a new Ethereum contract, binding an instance of Helper to it.
func DeployHelper(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Helper, error) {
	parsed, err := HelperMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(HelperBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Helper{HelperCaller: HelperCaller{contract: contract}, HelperTransactor: HelperTransactor{contract: contract}, HelperFilterer: HelperFilterer{contract: contract}}, nil
}

// Helper is an auto generated Go binding around an Ethereum contract.
type Helper struct {
	HelperCaller     // Read-only binding to the contract
	HelperTransactor // Write-only binding to the contract
	HelperFilterer   // Log filterer for contract events
}

// HelperCaller is an auto generated read-only Go binding around an Ethereum contract.
type HelperCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HelperTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HelperTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HelperFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HelperFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HelperSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HelperSession struct {
	Contract     *Helper           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HelperCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HelperCallerSession struct {
	Contract *HelperCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// HelperTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HelperTransactorSession struct {
	Contract     *HelperTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HelperRaw is an auto generated low-level Go binding around an Ethereum contract.
type HelperRaw struct {
	Contract *Helper // Generic contract binding to access the raw methods on
}

// HelperCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HelperCallerRaw struct {
	Contract *HelperCaller // Generic read-only contract binding to access the raw methods on
}

// HelperTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HelperTransactorRaw struct {
	Contract *HelperTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHelper creates a new instance of Helper, bound to a specific deployed contract.
func NewHelper(address common.Address, backend bind.ContractBackend) (*Helper, error) {
	contract, err := bindHelper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Helper{HelperCaller: HelperCaller{contract: contract}, HelperTransactor: HelperTransactor{contract: contract}, HelperFilterer: HelperFilterer{contract: contract}}, nil
}

// NewHelperCaller creates a new read-only instance of Helper, bound to a specific deployed contract.
func NewHelperCaller(address common.Address, caller bind.ContractCaller) (*HelperCaller, error) {
	contract, err := bindHelper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HelperCaller{contract: contract}, nil
}

// NewHelperTransactor creates a new write-only instance of Helper, bound to a specific deployed contract.
func NewHelperTransactor(address common.Address, transactor bind.ContractTransactor) (*HelperTransactor, error) {
	contract, err := bindHelper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HelperTransactor{contract: contract}, nil
}

// NewHelperFilterer creates a new log filterer instance of Helper, bound to a specific deployed contract.
func NewHelperFilterer(address common.Address, filterer bind.ContractFilterer) (*HelperFilterer, error) {
	contract, err := bindHelper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HelperFilterer{contract: contract}, nil
}

// bindHelper binds a generic wrapper to an already deployed contract.
func bindHelper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := HelperMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Helper *HelperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Helper.Contract.HelperCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Helper *HelperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Helper.Contract.HelperTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Helper *HelperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Helper.Contract.HelperTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Helper *HelperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Helper.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Helper *HelperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Helper.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Helper *HelperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Helper.Contract.contract.Transact(opts, method, params...)
}

// Distribute is a paid mutator transaction binding the contract method 0x6138889b.
//
// Solidity: function distribute(address[] receivers) payable returns()
func (_Helper *HelperTransactor) Distribute(opts *bind.TransactOpts, receivers []common.Address) (*types.Transaction, error) {
	return _Helper.contract.Transact(opts, "distribute", receivers)
}

// Distribute is a paid mutator transaction binding the contract method 0x6138889b.
//
// Solidity: function distribute(address[] receivers) payable returns()
func (_Helper *HelperSession) Distribute(receivers []common.Address) (*types.Transaction, error) {
	return _Helper.Contract.Distribute(&_Helper.TransactOpts, receivers)
}

// Distribute is a paid mutator transaction binding the contract method 0x6138889b.
//
// Solidity: function distribute(address[] receivers) payable returns()
func (_Helper *HelperTransactorSession) Distribute(receivers []common.Address) (*types.Transaction, error) {
	return _Helper.Contract.Distribute(&_Helper.TransactOpts, receivers)
}
