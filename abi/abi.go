package abi

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
)

func NewAbi(abiString string) LatticeAbi {
	return &latticeAbi{
		abiString: abiString,
		abi:       FromJson(abiString),
	}
}

type LatticeAbi interface {
	Function(method string) (*abi.Method, error)
	NewLatticeFunction(methodName string, args []interface{}) (LatticeFunction, error)
}

type latticeAbi struct {
	abiString string
	abi       *abi.ABI
}

func FromJson(abiString string) *abi.ABI {
	decoder := json.NewDecoder(strings.NewReader(abiString))

	var myAbi abi.ABI
	if err := decoder.Decode(&myAbi); err != nil {
		return nil
	}
	return &myAbi
}

func (i *latticeAbi) Function(methodName string) (*abi.Method, error) {
	if m, ok := i.abi.Methods[methodName]; ok {
		return &m, nil
	} else {
		return nil, fmt.Errorf("method %s not found", methodName)
	}
}

func (i *latticeAbi) Constructor() *abi.Method {
	return &i.abi.Constructor
}

func (i *latticeAbi) NewLatticeFunction(methodName string, args []interface{}) (LatticeFunction, error) {
	method, err := i.Function(methodName)
	if err != nil {
		return nil, err
	}
	return NewLatticeFunction(i.abiString, i.abi, methodName, args, method), nil
}
