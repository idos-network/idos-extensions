// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package registry

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

// FractalRegistryGrant is an auto generated low-level Go binding around an user-defined struct.
type FractalRegistryGrant struct {
	Owner       common.Address
	Grantee     common.Address
	DataId      string
	LockedUntil *big.Int
}

// RegistryMetaData contains all meta data concerning the Registry contract.
var RegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_grantee\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_dataId\",\"type\":\"string\"}],\"name\":\"delete_grant\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"grants\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"grantee\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"dataId\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"lockedUntil\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_grantee\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_dataId\",\"type\":\"string\"}],\"name\":\"grants_for\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"grantee\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"dataId\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"lockedUntil\",\"type\":\"uint256\"}],\"internalType\":\"structFractalRegistry.Grant[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_grantee\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_dataId\",\"type\":\"string\"}],\"name\":\"insert_grant\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// RegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use RegistryMetaData.ABI instead.
var RegistryABI = RegistryMetaData.ABI

// Registry is an auto generated Go binding around an Ethereum contract.
type Registry struct {
	RegistryCaller     // Read-only binding to the contract
	RegistryTransactor // Write-only binding to the contract
	RegistryFilterer   // Log filterer for contract events
}

// RegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type RegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RegistrySession struct {
	Contract     *Registry         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RegistryCallerSession struct {
	Contract *RegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// RegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RegistryTransactorSession struct {
	Contract     *RegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// RegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type RegistryRaw struct {
	Contract *Registry // Generic contract binding to access the raw methods on
}

// RegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RegistryCallerRaw struct {
	Contract *RegistryCaller // Generic read-only contract binding to access the raw methods on
}

// RegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RegistryTransactorRaw struct {
	Contract *RegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRegistry creates a new instance of Registry, bound to a specific deployed contract.
func NewRegistry(address common.Address, backend bind.ContractBackend) (*Registry, error) {
	contract, err := bindRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Registry{RegistryCaller: RegistryCaller{contract: contract}, RegistryTransactor: RegistryTransactor{contract: contract}, RegistryFilterer: RegistryFilterer{contract: contract}}, nil
}

// NewRegistryCaller creates a new read-only instance of Registry, bound to a specific deployed contract.
func NewRegistryCaller(address common.Address, caller bind.ContractCaller) (*RegistryCaller, error) {
	contract, err := bindRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryCaller{contract: contract}, nil
}

// NewRegistryTransactor creates a new write-only instance of Registry, bound to a specific deployed contract.
func NewRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*RegistryTransactor, error) {
	contract, err := bindRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryTransactor{contract: contract}, nil
}

// NewRegistryFilterer creates a new log filterer instance of Registry, bound to a specific deployed contract.
func NewRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*RegistryFilterer, error) {
	contract, err := bindRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RegistryFilterer{contract: contract}, nil
}

// bindRegistry binds a generic wrapper to an already deployed contract.
func bindRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registry *RegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Registry.Contract.RegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registry *RegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registry.Contract.RegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registry *RegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registry.Contract.RegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registry *RegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Registry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registry *RegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registry *RegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registry.Contract.contract.Transact(opts, method, params...)
}

// Grants is a free data retrieval call binding the contract method 0x8b956d28.
//
// Solidity: function grants(address , string , uint256 ) view returns(address owner, address grantee, string dataId, uint256 lockedUntil)
func (_Registry *RegistryCaller) Grants(opts *bind.CallOpts, arg0 common.Address, arg1 string, arg2 *big.Int) (struct {
	Owner       common.Address
	Grantee     common.Address
	DataId      string
	LockedUntil *big.Int
}, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "grants", arg0, arg1, arg2)

	outstruct := new(struct {
		Owner       common.Address
		Grantee     common.Address
		DataId      string
		LockedUntil *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Owner = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Grantee = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.DataId = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.LockedUntil = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Grants is a free data retrieval call binding the contract method 0x8b956d28.
//
// Solidity: function grants(address , string , uint256 ) view returns(address owner, address grantee, string dataId, uint256 lockedUntil)
func (_Registry *RegistrySession) Grants(arg0 common.Address, arg1 string, arg2 *big.Int) (struct {
	Owner       common.Address
	Grantee     common.Address
	DataId      string
	LockedUntil *big.Int
}, error) {
	return _Registry.Contract.Grants(&_Registry.CallOpts, arg0, arg1, arg2)
}

// Grants is a free data retrieval call binding the contract method 0x8b956d28.
//
// Solidity: function grants(address , string , uint256 ) view returns(address owner, address grantee, string dataId, uint256 lockedUntil)
func (_Registry *RegistryCallerSession) Grants(arg0 common.Address, arg1 string, arg2 *big.Int) (struct {
	Owner       common.Address
	Grantee     common.Address
	DataId      string
	LockedUntil *big.Int
}, error) {
	return _Registry.Contract.Grants(&_Registry.CallOpts, arg0, arg1, arg2)
}

// GrantsFor is a free data retrieval call binding the contract method 0xb59cd130.
//
// Solidity: function grants_for(address _grantee, string _dataId) view returns((address,address,string,uint256)[])
func (_Registry *RegistryCaller) GrantsFor(opts *bind.CallOpts, _grantee common.Address, _dataId string) ([]FractalRegistryGrant, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "grants_for", _grantee, _dataId)

	if err != nil {
		return *new([]FractalRegistryGrant), err
	}

	out0 := *abi.ConvertType(out[0], new([]FractalRegistryGrant)).(*[]FractalRegistryGrant)

	return out0, err

}

// GrantsFor is a free data retrieval call binding the contract method 0xb59cd130.
//
// Solidity: function grants_for(address _grantee, string _dataId) view returns((address,address,string,uint256)[])
func (_Registry *RegistrySession) GrantsFor(_grantee common.Address, _dataId string) ([]FractalRegistryGrant, error) {
	return _Registry.Contract.GrantsFor(&_Registry.CallOpts, _grantee, _dataId)
}

// GrantsFor is a free data retrieval call binding the contract method 0xb59cd130.
//
// Solidity: function grants_for(address _grantee, string _dataId) view returns((address,address,string,uint256)[])
func (_Registry *RegistryCallerSession) GrantsFor(_grantee common.Address, _dataId string) ([]FractalRegistryGrant, error) {
	return _Registry.Contract.GrantsFor(&_Registry.CallOpts, _grantee, _dataId)
}

// DeleteGrant is a paid mutator transaction binding the contract method 0xb0c2418d.
//
// Solidity: function delete_grant(address _grantee, string _dataId) returns()
func (_Registry *RegistryTransactor) DeleteGrant(opts *bind.TransactOpts, _grantee common.Address, _dataId string) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "delete_grant", _grantee, _dataId)
}

// DeleteGrant is a paid mutator transaction binding the contract method 0xb0c2418d.
//
// Solidity: function delete_grant(address _grantee, string _dataId) returns()
func (_Registry *RegistrySession) DeleteGrant(_grantee common.Address, _dataId string) (*types.Transaction, error) {
	return _Registry.Contract.DeleteGrant(&_Registry.TransactOpts, _grantee, _dataId)
}

// DeleteGrant is a paid mutator transaction binding the contract method 0xb0c2418d.
//
// Solidity: function delete_grant(address _grantee, string _dataId) returns()
func (_Registry *RegistryTransactorSession) DeleteGrant(_grantee common.Address, _dataId string) (*types.Transaction, error) {
	return _Registry.Contract.DeleteGrant(&_Registry.TransactOpts, _grantee, _dataId)
}

// InsertGrant is a paid mutator transaction binding the contract method 0xbd8f9b17.
//
// Solidity: function insert_grant(address _grantee, string _dataId) returns()
func (_Registry *RegistryTransactor) InsertGrant(opts *bind.TransactOpts, _grantee common.Address, _dataId string) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "insert_grant", _grantee, _dataId)
}

// InsertGrant is a paid mutator transaction binding the contract method 0xbd8f9b17.
//
// Solidity: function insert_grant(address _grantee, string _dataId) returns()
func (_Registry *RegistrySession) InsertGrant(_grantee common.Address, _dataId string) (*types.Transaction, error) {
	return _Registry.Contract.InsertGrant(&_Registry.TransactOpts, _grantee, _dataId)
}

// InsertGrant is a paid mutator transaction binding the contract method 0xbd8f9b17.
//
// Solidity: function insert_grant(address _grantee, string _dataId) returns()
func (_Registry *RegistryTransactorSession) InsertGrant(_grantee common.Address, _dataId string) (*types.Transaction, error) {
	return _Registry.Contract.InsertGrant(&_Registry.TransactOpts, _grantee, _dataId)
}
