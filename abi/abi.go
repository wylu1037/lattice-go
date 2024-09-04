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

	// RawAbi 获取RawABI
	RawAbi() *abi.ABI

	// Constructor 获取构造函数
	//
	// Returns:
	//   - *abi.Method
	Constructor() *abi.Method

	// Function 获取方法
	//
	// Parameters:
	//   - method string: 方法名称
	//
	// Returns:
	//   - *abi.Method
	//   - error: 方法不存在时，抛出错误
	Function(method string) (*abi.Method, error)

	GetConstructor(args ...interface{}) LatticeFunction

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

func (i *latticeAbi) RawAbi() *abi.ABI {
	return i.abi
}

func (i *latticeAbi) Function(methodName string) (*abi.Method, error) {
	if m, ok := i.abi.Methods[methodName]; ok {
		return &m, nil
	} else {
		return nil, fmt.Errorf("合约方法【%s】不存在", methodName)
	}
}

func (i *latticeAbi) Constructor() *abi.Method {
	return &i.abi.Constructor
}

func (i *latticeAbi) GetConstructor(args ...interface{}) LatticeFunction {
	return NewLatticeFunction(i.abiString, i.abi, "", args, i.Constructor())
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

// DecodeCall 解码合约的入参
func DecodeCall(myabi *abi.ABI, functionName, code string) (string, error) {
	method, ok := myabi.Methods[functionName]
	if !ok {
		return "", fmt.Errorf("合约方法【%s】不存在", functionName)
	}

	inputBytes, err := hexutil.Decode(code)
	if err != nil {
		return "", err
	}
	args := method.Inputs
	res, err := args.UnpackValues(inputBytes[4:])
	if err != nil {
		return "", err
	}
	finalRes := make(map[string]interface{}, len(res))
	for i, v := range res {
		finalRes[args[i].Name] = v
	}
	data, err := json.Marshal(finalRes)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
