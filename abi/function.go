package abi

import (
	"encoding/json"
	"fmt"
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

type LatticeFunction interface {
}

type latticeFunction struct {
	method *abi.Method
}

func (f *latticeFunction) Encode() string {
	return f.method.Name
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
			return nil, fmt.Errorf("unsupported argument type: %T", param)
		}
	// Example: true or "true"
	case abi.BoolTy:
		if b, ok := param.(bool); ok {
			return b, nil
		} else if s, ok := param.(string); ok {
			val, err := strconv.ParseBool(s)
			if err != nil {
				return nil, fmt.Errorf("failed to parse bool %q: %v", s, err)
			}
			return val, nil
		}
	// Example: "School"
	case abi.StringTy:
		if s, ok := param.(string); ok {
			return s, nil
		} else {
			return nil, fmt.Errorf("arg need string type, but now is %T", param)
		}
	// Example: ["apple", "banana"] | [1, 2, 3]
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

		for i, elem := range inputArray {
			fmt.Println(i, elem)
			return nil, nil
		}
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
			return nil, fmt.Errorf("address need string type, but now is: %T", param)
		}
	// "0x5f2be9a02b43f748ee460bf36eed24fafa109920"
	case abi.BytesTy:
		if s, ok := param.(string); ok {
			if strings.HasPrefix(s, "0x") {
				bytes, err := hexutil.Decode(s)
				if err != nil {
					return nil, fmt.Errorf("failed to decode bytes: %v", err)
				}
				return common.BytesToAddress(bytes), nil
			} else if strings.HasPrefix(s, "zltc") {

			} else {
				return nil, fmt.Errorf("invalid bytes type: %s", s)
			}
		}
	// "0x5f2be9a02b43f748ee460bf36eed24fafa109920", return common.Hash
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
		}

	case abi.FixedBytesTy:
		return nil, fmt.Errorf("arg need string type, but now is %T", param)
	// struct and tuple
	case abi.TupleTy:
		return nil, fmt.Errorf("unsupported input type %v", abiType)
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
