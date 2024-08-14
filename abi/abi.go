package abi

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"strings"
)

func NewAbi(abiString string) LatticeAbi {
	return &latticeAbi{
		abiString: abiString,
		abi:       FromJson(abiString),
	}
}

type LatticeAbi interface {
	MyAbi() *abi.ABI
	Function(method string) (*abi.Method, error)
	GetLatticeFunction(methodName string, args ...interface{}) (LatticeFunction, error)
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

func (i *latticeAbi) MyAbi() *abi.ABI {
	return i.abi
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

func (i *latticeAbi) GetLatticeFunction(methodName string, args ...interface{}) (LatticeFunction, error) {
	method, err := i.Function(methodName)
	if err != nil {
		return nil, err
	}
	return NewLatticeFunction(i.abiString, i.abi, methodName, args, method), nil
}

// DecodeReturn 解码合约调用结果
//
// Parameters:
//   - myabi *abi.ABI
//   - functionName string: 方法名
//   - contractReturn string: 合约调用结果
//
// Returns:
//   - string: abi解码后的合约调用结果
//   - error
func DecodeReturn(myabi *abi.ABI, functionName, contractReturn string) (string, error) {
	method, ok := myabi.Methods[functionName]
	if !ok {
		return "", fmt.Errorf("合约方法【%s】不存在", functionName)
	}

	bytes, err := hexutil.Decode(contractReturn)
	if err != nil {
		return "", err
	}
	res, err := method.Outputs.UnpackValues(bytes)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(res[0])
	if err != nil {
		return "", err
	}

	return string(data), nil
}
