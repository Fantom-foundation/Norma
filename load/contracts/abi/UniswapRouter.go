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

// UniswapRouterMetaData contains all meta data concerning the UniswapRouter contract.
var UniswapRouterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getCount\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"pairsPath\",\"type\":\"address[]\"}],\"name\":\"swapExactTokensForTokens\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040526000805534801561001457600080fd5b50610ce1806100246000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063a87d942c1461003b578063ddba27a714610051575b600080fd5b6000546040519081526020015b60405180910390f35b61006461005f3660046109cd565b610071565b6040516100489190610a47565b60606100e18686868080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050604080516020808a028281018201909352898252909350899250889182918501908490808284376000920191909152506101e492505050565b9050610157858560008181106100f9576100f9610a8b565b905060200201602081019061010e9190610aa1565b338585600081811061012257610122610a8b565b90506020020160208101906101379190610aa1565b8460008151811061014a5761014a610a8b565b60200260200101516103f3565b6101c78186868080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050604080516020808a02828101820190935289825290935089925088918291850190849080828437600092019190915250339250610523915050565b6000805490806101d683610ae7565b919050555095945050505050565b606060028351101561023d5760405162461bcd60e51b815260206004820152601e60248201527f556e697377617056324c6962726172793a20494e56414c49445f50415448000060448201526064015b60405180910390fd5b815161024a906001610b06565b8351146102a35760405162461bcd60e51b815260206004820152602160248201527f696e76616c6964206c656e677468206f662070616972735061746820706172616044820152606d60f81b6064820152608401610234565b825167ffffffffffffffff8111156102bd576102bd610b19565b6040519080825280602002602001820160405280156102e6578160200160208202803683370190505b50905083816000815181106102fd576102fd610a8b565b60200260200101818152505060005b6001845161031a9190610b2f565b8110156103eb5760008061038686848151811061033957610339610a8b565b60200260200101518785600161034f9190610b06565b8151811061035f5761035f610a8b565b602002602001015187868151811061037957610379610a8b565b60200260200101516106f0565b915091506103ae84848151811061039f5761039f610a8b565b602002602001015183836107a0565b846103ba856001610b06565b815181106103ca576103ca610a8b565b602002602001018181525050505080806103e390610b42565b91505061030c565b509392505050565b604080516001600160a01b0385811660248301528481166044830152606480830185905283518084039091018152608490920183526020820180516001600160e01b03166323b872dd60e01b17905291516000928392908816916104579190610b78565b6000604051808303816000865af19150503d8060008114610494576040519150601f19603f3d011682016040523d82523d6000602084013e610499565b606091505b50915091508180156104c35750805115806104c35750808060200190518101906104c39190610b94565b61051b5760405162461bcd60e51b8152602060048201526024808201527f5472616e7366657248656c7065723a205452414e534645525f46524f4d5f46416044820152631253115160e21b6064820152608401610234565b505050505050565b60005b82518110156106e95760008085838151811061054457610544610a8b565b60200260200101518684600161055a9190610b06565b8151811061056a5761056a610a8b565b6020026020010151915091506000878460016105869190610b06565b8151811061059657610596610a8b565b60200260200101519050600080836001600160a01b0316856001600160a01b0316106105c4578260006105c8565b6000835b91509150600088518760016105dd9190610b06565b106105e8578761060d565b886105f4886001610b06565b8151811061060457610604610a8b565b60200260200101515b905088878151811061062157610621610a8b565b60200260200101516001600160a01b031663022c0d9f848484600067ffffffffffffffff81111561065457610654610b19565b6040519080825280601f01601f19166020018201604052801561067e576020820181803683370190505b506040518563ffffffff1660e01b815260040161069e9493929190610bb6565b600060405180830381600087803b1580156106b857600080fd5b505af11580156106cc573d6000803e3d6000fd5b5050505050505050505080806106e190610b42565b915050610526565b5050505050565b600080600080846001600160a01b0316630902f1ac6040518163ffffffff1660e01b8152600401606060405180830381865afa158015610734573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107589190610c22565b506001600160701b031691506001600160701b03169150856001600160a01b0316876001600160a01b03161061078f578082610792565b81815b909890975095505050505050565b60008084116108055760405162461bcd60e51b815260206004820152602b60248201527f556e697377617056324c6962726172793a20494e53554646494349454e545f4960448201526a1394155517d05353d5539560aa1b6064820152608401610234565b6000831180156108155750600082115b6108725760405162461bcd60e51b815260206004820152602860248201527f556e697377617056324c6962726172793a20494e53554646494349454e545f4c604482015267495155494449545960c01b6064820152608401610234565b6000610880856103e56108bf565b9050600061088e82856108bf565b905060006108a8836108a2886103e86108bf565b9061092c565b90506108b48183610c72565b979650505050505050565b60008115806108e3575082826108d58183610c94565b92506108e19083610c72565b145b6109265760405162461bcd60e51b815260206004820152601460248201527364732d6d6174682d6d756c2d6f766572666c6f7760601b6044820152606401610234565b92915050565b6000826109398382610b06565b91508110156109265760405162461bcd60e51b815260206004820152601460248201527364732d6d6174682d6164642d6f766572666c6f7760601b6044820152606401610234565b60008083601f84011261099357600080fd5b50813567ffffffffffffffff8111156109ab57600080fd5b6020830191508360208260051b85010111156109c657600080fd5b9250929050565b6000806000806000606086880312156109e557600080fd5b85359450602086013567ffffffffffffffff80821115610a0457600080fd5b610a1089838a01610981565b90965094506040880135915080821115610a2957600080fd5b50610a3688828901610981565b969995985093965092949392505050565b6020808252825182820181905260009190848201906040850190845b81811015610a7f57835183529284019291840191600101610a63565b50909695505050505050565b634e487b7160e01b600052603260045260246000fd5b600060208284031215610ab357600080fd5b81356001600160a01b0381168114610aca57600080fd5b9392505050565b634e487b7160e01b600052601160045260246000fd5b60006001600160ff1b018201610aff57610aff610ad1565b5060010190565b8082018082111561092657610926610ad1565b634e487b7160e01b600052604160045260246000fd5b8181038181111561092657610926610ad1565b600060018201610aff57610aff610ad1565b60005b83811015610b6f578181015183820152602001610b57565b50506000910152565b60008251610b8a818460208701610b54565b9190910192915050565b600060208284031215610ba657600080fd5b81518015158114610aca57600080fd5b84815283602082015260018060a01b03831660408201526080606082015260008251806080840152610bef8160a0850160208701610b54565b601f01601f19169190910160a00195945050505050565b80516001600160701b0381168114610c1d57600080fd5b919050565b600080600060608486031215610c3757600080fd5b610c4084610c06565b9250610c4e60208501610c06565b9150604084015163ffffffff81168114610c6757600080fd5b809150509250925092565b600082610c8f57634e487b7160e01b600052601260045260246000fd5b500490565b808202811582820484141761092657610926610ad156fea2646970667358221220d80586d285f22fb21b235f89aeef42c3c3f9781f4ffda35422a76fb63247a9be64736f6c63430008130033",
}

// UniswapRouterABI is the input ABI used to generate the binding from.
// Deprecated: Use UniswapRouterMetaData.ABI instead.
var UniswapRouterABI = UniswapRouterMetaData.ABI

// UniswapRouterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UniswapRouterMetaData.Bin instead.
var UniswapRouterBin = UniswapRouterMetaData.Bin

// DeployUniswapRouter deploys a new Ethereum contract, binding an instance of UniswapRouter to it.
func DeployUniswapRouter(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *UniswapRouter, error) {
	parsed, err := UniswapRouterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UniswapRouterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UniswapRouter{UniswapRouterCaller: UniswapRouterCaller{contract: contract}, UniswapRouterTransactor: UniswapRouterTransactor{contract: contract}, UniswapRouterFilterer: UniswapRouterFilterer{contract: contract}}, nil
}

// UniswapRouter is an auto generated Go binding around an Ethereum contract.
type UniswapRouter struct {
	UniswapRouterCaller     // Read-only binding to the contract
	UniswapRouterTransactor // Write-only binding to the contract
	UniswapRouterFilterer   // Log filterer for contract events
}

// UniswapRouterCaller is an auto generated read-only Go binding around an Ethereum contract.
type UniswapRouterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapRouterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UniswapRouterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapRouterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UniswapRouterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapRouterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UniswapRouterSession struct {
	Contract     *UniswapRouter    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UniswapRouterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UniswapRouterCallerSession struct {
	Contract *UniswapRouterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// UniswapRouterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UniswapRouterTransactorSession struct {
	Contract     *UniswapRouterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// UniswapRouterRaw is an auto generated low-level Go binding around an Ethereum contract.
type UniswapRouterRaw struct {
	Contract *UniswapRouter // Generic contract binding to access the raw methods on
}

// UniswapRouterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UniswapRouterCallerRaw struct {
	Contract *UniswapRouterCaller // Generic read-only contract binding to access the raw methods on
}

// UniswapRouterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UniswapRouterTransactorRaw struct {
	Contract *UniswapRouterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUniswapRouter creates a new instance of UniswapRouter, bound to a specific deployed contract.
func NewUniswapRouter(address common.Address, backend bind.ContractBackend) (*UniswapRouter, error) {
	contract, err := bindUniswapRouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UniswapRouter{UniswapRouterCaller: UniswapRouterCaller{contract: contract}, UniswapRouterTransactor: UniswapRouterTransactor{contract: contract}, UniswapRouterFilterer: UniswapRouterFilterer{contract: contract}}, nil
}

// NewUniswapRouterCaller creates a new read-only instance of UniswapRouter, bound to a specific deployed contract.
func NewUniswapRouterCaller(address common.Address, caller bind.ContractCaller) (*UniswapRouterCaller, error) {
	contract, err := bindUniswapRouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapRouterCaller{contract: contract}, nil
}

// NewUniswapRouterTransactor creates a new write-only instance of UniswapRouter, bound to a specific deployed contract.
func NewUniswapRouterTransactor(address common.Address, transactor bind.ContractTransactor) (*UniswapRouterTransactor, error) {
	contract, err := bindUniswapRouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapRouterTransactor{contract: contract}, nil
}

// NewUniswapRouterFilterer creates a new log filterer instance of UniswapRouter, bound to a specific deployed contract.
func NewUniswapRouterFilterer(address common.Address, filterer bind.ContractFilterer) (*UniswapRouterFilterer, error) {
	contract, err := bindUniswapRouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UniswapRouterFilterer{contract: contract}, nil
}

// bindUniswapRouter binds a generic wrapper to an already deployed contract.
func bindUniswapRouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UniswapRouterABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapRouter *UniswapRouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapRouter.Contract.UniswapRouterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapRouter *UniswapRouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapRouter.Contract.UniswapRouterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapRouter *UniswapRouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapRouter.Contract.UniswapRouterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapRouter *UniswapRouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapRouter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapRouter *UniswapRouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapRouter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapRouter *UniswapRouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapRouter.Contract.contract.Transact(opts, method, params...)
}

// GetCount is a free data retrieval call binding the contract method 0xa87d942c.
//
// Solidity: function getCount() view returns(int256)
func (_UniswapRouter *UniswapRouterCaller) GetCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapRouter.contract.Call(opts, &out, "getCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCount is a free data retrieval call binding the contract method 0xa87d942c.
//
// Solidity: function getCount() view returns(int256)
func (_UniswapRouter *UniswapRouterSession) GetCount() (*big.Int, error) {
	return _UniswapRouter.Contract.GetCount(&_UniswapRouter.CallOpts)
}

// GetCount is a free data retrieval call binding the contract method 0xa87d942c.
//
// Solidity: function getCount() view returns(int256)
func (_UniswapRouter *UniswapRouterCallerSession) GetCount() (*big.Int, error) {
	return _UniswapRouter.Contract.GetCount(&_UniswapRouter.CallOpts)
}

// SwapExactTokensForTokens is a paid mutator transaction binding the contract method 0xddba27a7.
//
// Solidity: function swapExactTokensForTokens(uint256 amountIn, address[] path, address[] pairsPath) returns(uint256[] amounts)
func (_UniswapRouter *UniswapRouterTransactor) SwapExactTokensForTokens(opts *bind.TransactOpts, amountIn *big.Int, path []common.Address, pairsPath []common.Address) (*types.Transaction, error) {
	return _UniswapRouter.contract.Transact(opts, "swapExactTokensForTokens", amountIn, path, pairsPath)
}

// SwapExactTokensForTokens is a paid mutator transaction binding the contract method 0xddba27a7.
//
// Solidity: function swapExactTokensForTokens(uint256 amountIn, address[] path, address[] pairsPath) returns(uint256[] amounts)
func (_UniswapRouter *UniswapRouterSession) SwapExactTokensForTokens(amountIn *big.Int, path []common.Address, pairsPath []common.Address) (*types.Transaction, error) {
	return _UniswapRouter.Contract.SwapExactTokensForTokens(&_UniswapRouter.TransactOpts, amountIn, path, pairsPath)
}

// SwapExactTokensForTokens is a paid mutator transaction binding the contract method 0xddba27a7.
//
// Solidity: function swapExactTokensForTokens(uint256 amountIn, address[] path, address[] pairsPath) returns(uint256[] amounts)
func (_UniswapRouter *UniswapRouterTransactorSession) SwapExactTokensForTokens(amountIn *big.Int, path []common.Address, pairsPath []common.Address) (*types.Transaction, error) {
	return _UniswapRouter.Contract.SwapExactTokensForTokens(&_UniswapRouter.TransactOpts, amountIn, path, pairsPath)
}
