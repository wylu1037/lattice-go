package abi

import (
	"encoding/json"
	"errors"
	"fmt"
	abi2 "github.com/defiweb/go-eth/abi"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"lattice-go/common/constant"
	"lattice-go/common/convert"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
)

func NewLatticeFunction(
	abiString string,
	abi *abi.ABI,
	methodName string,
	args []interface{},
	method *abi.Method) LatticeFunction {
	return &latticeFunction{
		abiString:  abiString,
		abi:        abi,
		methodName: methodName,
		args:       args,
		method:     method,
	}
}

type LatticeFunction interface {
	Encode() (string, error)
}

type latticeFunction struct {
	abiString  string
	abi        *abi.ABI
	methodName string
	args       []interface{}
	method     *abi.Method
}

func (f *latticeFunction) Encode() (string, error) {
	var err error
	convertedArgs, err := f.ConvertArguments(f.method.Inputs, f.args)
	if err != nil {
		return "", err
	}

	var data []byte
	if f.inputsContainsTuple() || f.inputsContainsSlice() {
		contract, err := abi2.ParseJSON([]byte(f.abiString))
		if err != nil {
			return "", err
		}
		m, ok := contract.Methods[f.methodName]
		if !ok {
			return "", errors.New(fmt.Sprintf("no such method: %s", f.methodName))
		}
		if data, err = m.EncodeArgs(convertedArgs...); err != nil {
			return "", err
		}
	} else {
		if data, err = f.abi.Pack(f.methodName, convertedArgs...); err != nil {
			return "", err
		}
	}

	return hexutil.Encode(data), nil
}

func (f *latticeFunction) Decode() []interface{} {
	return []interface{}{}
}

func (f *latticeFunction) Inputs() abi.Arguments {
	return f.method.Inputs
}

func (f *latticeFunction) ConvertArguments(args abi.Arguments, params []interface{}) ([]interface{}, error) {
	if len(args) != len(params) {
		return nil, fmt.Errorf("mismatched argument (%d) and parameter (%d) counts", len(args), len(params))
	}
	var convertedParams []interface{}
	for i, input := range args {
		param, err := f.ConvertArgument(input.Type, params[i])
		if err != nil {
			return nil, err
		}
		convertedParams = append(convertedParams, param)
	}
	return convertedParams, nil
}

func (f *latticeFunction) ConvertArgument(abiType abi.Type, param interface{}) (interface{}, error) {
	size := abiType.Size
	switch abiType.T {
	// Input example: "100"
	case abi.IntTy, abi.UintTy:
		if j, ok := param.(json.Number); ok {
			param = string(j)
		}
		if s, ok := param.(string); ok {
			val, ok := new(big.Int).SetString(s, 0)
			if !ok {
				return nil, fmt.Errorf("failed to parse big.Int: %s", s)
			}
			return ConvertInt(abiType.T == abi.IntTy, size, val)
		} else if i, ok := param.(*big.Int); ok {
			return ConvertInt(abiType.T == abi.IntTy, size, i)
		}
		v := reflect.ValueOf(param)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i := new(big.Int).SetInt64(v.Int())
			return ConvertInt(abiType.T == abi.IntTy, size, i)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			i := new(big.Int).SetUint64(v.Uint())
			return ConvertInt(abiType.T == abi.IntTy, size, i)
		case reflect.Float64, reflect.Float32:
			return nil, fmt.Errorf("floating point numbers are not valid in web3 - please use an integer or string instead (including big.Int and json.Number)")
		default:
			return nil, fmt.Errorf("unsupported argument type: %T, int type or uint type expect string number value", param)
		}
	// Input example: true or "true"
	case abi.BoolTy:
		if b, ok := param.(bool); ok {
			return b, nil
		} else if s, ok := param.(string); ok {
			val, err := strconv.ParseBool(s)
			if err != nil {
				return nil, fmt.Errorf("failed to parse bool %q: %v", s, err)
			}
			return val, nil
		} else {
			return nil, fmt.Errorf("unsupported argument type: %T, bool type expect string or bool value", param)
		}
	// Input example: "School"
	case abi.StringTy:
		if s, ok := param.(string); ok {
			return s, nil
		} else {
			return nil, fmt.Errorf("unsupported argument type: %T, string type expect string value", param)
		}
	// Input example: ["apple", "banana"] | [1, 2, 3]
	case abi.SliceTy, abi.ArrayTy:
		r := reflect.ValueOf(param)
		inputArray := make([]interface{}, 0)
		if r.Kind() == reflect.Array || r.Kind() == reflect.Slice {
			for i := 0; i < r.Len(); i++ {
				inputArray = append(inputArray, r.Index(i).Interface())
			}

		} else {
			s, ok := param.(string)
			if !ok {
				return nil, fmt.Errorf("invalid array: %s", s)
			}
			s = strings.TrimPrefix(s, "[")
			s = strings.TrimSuffix(s, "]")
			strArr := strings.Split(s, ",")
			for _, str := range strArr {
				inputArray = append(inputArray, str)
			}
		}

		convertedArgs := make([]interface{}, len(inputArray))
		for i, input := range inputArray {
			convertedArg, err := f.ConvertArgument(*abiType.Elem, input)
			if err != nil {
				return nil, err
			}
			convertedArgs[i] = convertedArg
		}
		return convertedArgs, nil
	// return string, Example: input`zltc_Z1pnS94bP4hQSYLs4aP4UwBP9pH8bEvhi`, output`0x5f2be9a02b43f748ee460bf36eed24fafa109920`
	case abi.AddressTy:
		if s, ok := param.(string); ok {
			if strings.HasPrefix(s, constant.HexPrefix) {
				return s, nil
			} else if strings.HasPrefix(s, constant.AddressTitle) {
				address, err := convert.ZltcToAddress(s)
				if err != nil {
					return nil, fmt.Errorf("invalid base58 address: %s", s)
				}
				return address.Hex(), nil
			}
		} else {
			return nil, fmt.Errorf("unsupported argument type: %T, address type expect hex string(42) or zltc address(38) value", param)
		}
	// Input example: "0x5f2be9a02b43f748ee460bf36eed24fafa109920"
	case abi.BytesTy:
		if s, ok := param.(string); ok {
			if strings.HasPrefix(s, "0x") {
				bytes, err := hexutil.Decode(s)
				if err != nil {
					return nil, fmt.Errorf("failed to decode bytes: %v", err)
				}
				return bytes, nil
			} else if strings.HasPrefix(s, "zltc") {
				addr, err := convert.ZltcToAddress(s)
				if err != nil {
					return nil, fmt.Errorf("failed to convert %s to common.Address,  %v", s, err)
				}
				return addr.Bytes(), nil
			} else {
				return nil, fmt.Errorf("invalid bytes type: %s, you can input hex string or zltc addr", s)
			}
		} else {
			return nil, fmt.Errorf("unsupported argument type: %T, bytes type expect hex string or zltc address(38) value", param)
		}
	// Input example: "0x5f2be9a02b43f748ee460bf36eed24fafa109920", return common.Hash
	case abi.HashTy:
		if s, ok := param.(string); ok {
			bytes, err := hexutil.Decode(s)
			if err != nil {
				return nil, fmt.Errorf("failed to parse hash %q: %v", s, err)
			}
			if len(bytes) != common.HashLength {
				return nil, fmt.Errorf("invalid hash length %d:hash must be 32 bytes", len(bytes))
			}
			return common.BytesToHash(bytes), nil
		} else {
			return nil, fmt.Errorf("unsupported argument type: %T, hash type expect hex string value", param)
		}
	case abi.FixedBytesTy:
		return nil, fmt.Errorf("arg need string type, but now is %T", param)
	// struct and tuple, Return type: Map, Return example: {"name": "Jack", "age": 28}
	case abi.TupleTy:
		// See: https://github.com/defiweb/go-eth
		// contract := MustParseJSON([]byte(`[{"inputs":[],"name":"get","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"components":[{"internalType":"uint256","name":"id","type":"uint256"},{"internalType":"string","name":"name","type":"string"},{"internalType":"bool","name":"isMan","type":"bool"},{"internalType":"string[]","name":"tags","type":"string[]"}],"internalType":"struct Test.User","name":"user","type":"tuple"}],"name":"set","outputs":[],"stateMutability":"nonpayable","type":"function"}]`))
		// method := contract.Methods["set"]
		// encodedData, err := method.EncodeArgs(map[string]interface{}{
		//		"id":    big.NewInt(18),
		//		"name":  "Jack",
		//		"isMan": true,
		//		"tags":  []string{"Hello, world!"},
		//	})
		encodedArgsMap := make(map[string]interface{})
		keys := abiType.TupleRawNames

		paramsArr := param.([]interface{})
		for i, elem := range abiType.TupleElems {
			convertedArg, err := f.ConvertArgument(*elem, paramsArr[i])
			if err != nil {
				return nil, err
			}
			encodedArgsMap[keys[i]] = convertedArg
		}

		return encodedArgsMap, nil
	// 固定精度的小数类型
	case abi.FixedPointTy:
		return nil, fmt.Errorf("unsupported input type %v", abiType)
	// 函数类型
	case abi.FunctionTy:
		return nil, fmt.Errorf("unsupported input type %v", abiType)
	default:
		return nil, fmt.Errorf("unsupported input type %v", abiType)
	}
	return param, nil
}

// ConvertInt converts a big.Int in to the provided type.
func ConvertInt(signed bool, size int, i *big.Int) (interface{}, error) {
	if signed {
		return convertSignedInt(size, i)
	}
	return convertUnsignedInt(size, i)
}

func convertSignedInt(size int, i *big.Int) (interface{}, error) {
	if size > 64 {
		return i, nil
	}
	int64Val := i.Int64()
	switch {
	case size > 32:
		if int64Val > math.MaxInt64 || int64Val < math.MinInt64 {
			return nil, fmt.Errorf("integer overflows int64: %s", i)
		}
		return int64Val, nil
	case size > 16:
		if int64Val > math.MaxInt32 || int64Val < math.MinInt32 {
			return nil, fmt.Errorf("integer overflows int32: %s", i)
		}
		return int32(int64Val), nil
	case size > 8:
		if int64Val > math.MaxInt16 || int64Val < math.MinInt16 {
			return nil, fmt.Errorf("integer overflows int16: %s", i)
		}
		return int16(int64Val), nil
	default:
		if int64Val > math.MaxInt8 || int64Val < math.MinInt8 {
			return nil, fmt.Errorf("integer overflows int8: %s", i)
		}
		return int8(int64Val), nil
	}
}

func convertUnsignedInt(size int, i *big.Int) (interface{}, error) {
	if i.Sign() == -1 {
		return nil, fmt.Errorf("negative value in unsigned field: %s", i)
	}
	uint64Val := i.Uint64()
	switch {
	case size > 64:
		return i, nil
	case size > 32:
		if uint64Val > math.MaxUint64 {
			return nil, fmt.Errorf("integer overflows uint64: %s", i)
		}
		return uint64Val, nil
	case size > 16:
		if uint64Val > math.MaxUint32 {
			return nil, fmt.Errorf("integer overflows uint32: %s", i)
		}
		return uint32(uint64Val), nil
	case size > 8:
		if uint64Val > math.MaxUint16 {
			return nil, fmt.Errorf("integer overflows uint16: %s", i)
		}
		return uint16(uint64Val), nil
	default:
		if uint64Val > math.MaxUint8 {
			return nil, fmt.Errorf("integer overflows uint8: %s", i)
		}
		return uint8(uint64Val), nil
	}
}

// 检查合约的入参是否包含元组或者结构体
// Parameters
//
// Returns
//   - bool: false-不包含，true-包含
func (f *latticeFunction) inputsContainsTuple() bool {
	contains := false
	for _, arg := range f.method.Inputs {
		if len(arg.Type.TupleElems) > 0 {
			contains = true
			break
		}
	}
	return contains
}

// 检查合约的入参是否包含切片和数组
// Parameters
//
// Returns
//   - bool: false-不包含，true-包含
func (f *latticeFunction) inputsContainsSlice() bool {
	contains := false
	for _, arg := range f.method.Inputs {
		if arg.Type.T == abi.SliceTy || arg.Type.T == abi.ArrayTy {
			contains = true
			break
		}
	}
	return contains
}
