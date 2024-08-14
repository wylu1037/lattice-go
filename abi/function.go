package abi

import (
	"encoding/json"
	"errors"
	"fmt"
	//abi2 "github.com/defiweb/go-eth/abi"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/wylu1037/lattice-go/common/constant"
	"github.com/wylu1037/lattice-go/common/convert"
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

// Encode abi encode
//
// Returns:
//   - string: 带0x前缀的16进制字符串
//   - error
func (f *latticeFunction) Encode() (string, error) {
	var err error
	convertedArgs, err := f.ConvertArguments(f.method.Inputs, f.args)
	if err != nil {
		return "", err
	}

	var data []byte
	if f.inputsContainsTuple() {
		/*contract, err := abi2.ParseJSON([]byte(f.abiString))
		if err != nil {
			return "", err
		}
		m, ok := contract.Methods[f.methodName]
		if !ok {
			return "", errors.New(fmt.Sprintf("no such method: %s", f.methodName))
		}
		if data, err = m.EncodeArgs(convertedArgs...); err != nil {
			return "", err
		}*/
		return "", errors.New("tuple is not supported")
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

		switch abiType.Elem.T {
		case abi.StringTy:
			stringArr := make([]string, len(convertedArgs))
			for i, converted := range convertedArgs {
				stringArr[i] = converted.(string)
			}
			return stringArr, nil
		case abi.AddressTy:
			addressArr := make([]common.Address, len(convertedArgs))
			for i, converted := range convertedArgs {
				addressArr[i] = converted.(common.Address)
			}
			return addressArr, nil
		case abi.BoolTy:
			boolArr := make([]bool, len(convertedArgs))
			for i, converted := range convertedArgs {
				boolArr[i] = converted.(bool)
			}
			return boolArr, nil
		case abi.BytesTy:
			bytesArr := make([][]byte, len(convertedArgs))
			for i, converted := range convertedArgs {
				bytesArr[i] = converted.([]byte)
			}
		case abi.HashTy:
			hashArr := make([]common.Hash, len(convertedArgs))
			for i, converted := range convertedArgs {
				hashArr[i] = converted.(common.Hash)
			}
			return hashArr, nil
		case abi.TupleTy:
			tupleArr := make([]map[string]interface{}, len(convertedArgs))
			for i, converted := range convertedArgs {
				tupleArr[i] = converted.(map[string]interface{})
			}
			return tupleArr, nil
		case abi.FixedBytesTy: // 1~32
			switch abiType.Elem.Size {
			case 1:
				fixedBytes1Arr := make([][1]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes1Arr[i] = converted.([1]byte)
				}
				return fixedBytes1Arr, nil
			case 2:
				fixedBytes2Arr := make([][2]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes2Arr[i] = converted.([2]byte)
				}
				return fixedBytes2Arr, nil
			case 3:
				fixedBytes3Arr := make([][3]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes3Arr[i] = converted.([3]byte)
				}
				return fixedBytes3Arr, nil
			case 4:
				fixedBytes4Arr := make([][4]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes4Arr[i] = converted.([4]byte)
				}
				return fixedBytes4Arr, nil
			case 5:
				fixedBytes5Arr := make([][5]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes5Arr[i] = converted.([5]byte)
				}
				return fixedBytes5Arr, nil
			case 6:
				fixedBytes6Arr := make([][6]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes6Arr[i] = converted.([6]byte)
				}
				return fixedBytes6Arr, nil
			case 7:
				fixedBytes7Arr := make([][7]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes7Arr[i] = converted.([7]byte)
				}
				return fixedBytes7Arr, nil
			case 8:
				fixedBytes8Arr := make([][8]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes8Arr[i] = converted.([8]byte)
				}
				return fixedBytes8Arr, nil
			case 9:
				fixedBytes9Arr := make([][9]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes9Arr[i] = converted.([9]byte)
				}
				return fixedBytes9Arr, nil
			case 10:
				fixedBytes10Arr := make([][10]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes10Arr[i] = converted.([10]byte)
				}
				return fixedBytes10Arr, nil
			case 11:
				fixedBytes11Arr := make([][11]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes11Arr[i] = converted.([11]byte)
				}
				return fixedBytes11Arr, nil
			case 12:
				fixedBytes12Arr := make([][12]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes12Arr[i] = converted.([12]byte)
				}
				return fixedBytes12Arr, nil
			case 13:
				fixedBytes13Arr := make([][13]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes13Arr[i] = converted.([13]byte)
				}
				return fixedBytes13Arr, nil
			case 14:
				fixedBytes14Arr := make([][14]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes14Arr[i] = converted.([14]byte)
				}
				return fixedBytes14Arr, nil
			case 15:
				fixedBytes15Arr := make([][15]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes15Arr[i] = converted.([15]byte)
				}
				return fixedBytes15Arr, nil
			case 16:
				fixedBytes16Arr := make([][16]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes16Arr[i] = converted.([16]byte)
				}
				return fixedBytes16Arr, nil
			case 17:
				fixedBytes17Arr := make([][17]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes17Arr[i] = converted.([17]byte)
				}
				return fixedBytes17Arr, nil
			case 18:
				fixedBytes18Arr := make([][18]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes18Arr[i] = converted.([18]byte)
				}
				return fixedBytes18Arr, nil
			case 19:
				fixedBytes19Arr := make([][19]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes19Arr[i] = converted.([19]byte)
				}
				return fixedBytes19Arr, nil
			case 20:
				fixedBytes20Arr := make([][20]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes20Arr[i] = converted.([20]byte)
				}
				return fixedBytes20Arr, nil
			case 21:
				fixedBytes21Arr := make([][21]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes21Arr[i] = converted.([21]byte)
				}
				return fixedBytes21Arr, nil
			case 22:
				fixedBytes22Arr := make([][22]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes22Arr[i] = converted.([22]byte)
				}
				return fixedBytes22Arr, nil
			case 23:
				fixedBytes23Arr := make([][23]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes23Arr[i] = converted.([23]byte)
				}
				return fixedBytes23Arr, nil
			case 24:
				fixedBytes24Arr := make([][24]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes24Arr[i] = converted.([24]byte)
				}
				return fixedBytes24Arr, nil
			case 25:
				fixedBytes25Arr := make([][25]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes25Arr[i] = converted.([25]byte)
				}
				return fixedBytes25Arr, nil
			case 26:
				fixedBytes26Arr := make([][26]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes26Arr[i] = converted.([26]byte)
				}
				return fixedBytes26Arr, nil
			case 27:
				fixedBytes27Arr := make([][27]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes27Arr[i] = converted.([27]byte)
				}
				return fixedBytes27Arr, nil
			case 28:
				fixedBytes28Arr := make([][28]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes28Arr[i] = converted.([28]byte)
				}
				return fixedBytes28Arr, nil
			case 29:
				fixedBytes29Arr := make([][29]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes29Arr[i] = converted.([29]byte)
				}
				return fixedBytes29Arr, nil
			case 30:
				fixedBytes30Arr := make([][30]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes30Arr[i] = converted.([30]byte)
				}
				return fixedBytes30Arr, nil
			case 31:
				fixedBytes31Arr := make([][31]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes31Arr[i] = converted.([31]byte)
				}
				return fixedBytes31Arr, nil
			case 32:
				fixedBytes32Arr := make([][32]byte, len(convertedArgs))
				for i, converted := range convertedArgs {
					fixedBytes32Arr[i] = converted.([32]byte)
				}
				return fixedBytes32Arr, nil
			}
		case abi.FixedPointTy, abi.FunctionTy:
		default:
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
				return address, nil
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
	case abi.FixedBytesTy: // 1~32
		switch size {
		case 1:
			if b1, ok := param.([1]byte); ok {
				return b1, nil
			} else if b1Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b1Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b1Str, err)
				}
				if len(bytes) != 1 {
					return nil, fmt.Errorf("expect length is 1, but got %d", len(bytes))
				}
				return convert.BytesToBytes1(bytes), nil
			}
		case 2:
			if b2, ok := param.([2]byte); ok {
				return b2, nil
			} else if b2Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b2Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b2Str, err)
				}
				if len(bytes) != 2 {
					return nil, fmt.Errorf("expect length is 2, but got %d", len(bytes))
				}
				return convert.BytesToBytes2(bytes), nil
			}
		case 3:
			if b3, ok := param.([3]byte); ok {
				return b3, nil
			} else if b3Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b3Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b3Str, err)
				}
				if len(bytes) != 3 {
					return nil, fmt.Errorf("expect length is 3, but got %d", len(bytes))
				}
				return convert.BytesToBytes3(bytes), nil
			}
		case 4:
			if b4, ok := param.([4]byte); ok {
				return b4, nil
			} else if b4Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b4Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b4Str, err)
				}
				if len(bytes) != 4 {
					return nil, fmt.Errorf("expect length is 4, but got %d", len(bytes))
				}
				return convert.BytesToBytes4(bytes), nil
			}
		case 5:
			if b5, ok := param.([5]byte); ok {
				return b5, nil
			} else if b5Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b5Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b5Str, err)
				}
				if len(bytes) != 5 {
					return nil, fmt.Errorf("expect length is 5, but got %d", len(bytes))
				}
				return convert.BytesToBytes5(bytes), nil
			}
		case 6:
			if b6, ok := param.([6]byte); ok {
				return b6, nil
			} else if b6Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b6Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b6Str, err)
				}
				if len(bytes) != 6 {
					return nil, fmt.Errorf("expect length is 6, but got %d", len(bytes))
				}
				return convert.BytesToBytes6(bytes), nil
			}
		case 7:
			if b7, ok := param.([7]byte); ok {
				return b7, nil
			} else if b7Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b7Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b7Str, err)
				}
				if len(bytes) != 7 {
					return nil, fmt.Errorf("expect length is 7, but got %d", len(bytes))
				}
				return convert.BytesToBytes7(bytes), nil
			}
		case 8:
			if b8, ok := param.([8]byte); ok {
				return b8, nil
			} else if b8Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b8Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b8Str, err)
				}
				if len(bytes) != 8 {
					return nil, fmt.Errorf("expect length is 8, but got %d", len(bytes))
				}
				return convert.BytesToBytes8(bytes), nil
			}
		case 9:
			if b9, ok := param.([9]byte); ok {
				return b9, nil
			} else if b9Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b9Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b9Str, err)
				}
				if len(bytes) != 9 {
					return nil, fmt.Errorf("expect length is 9, but got %d", len(bytes))
				}
				return convert.BytesToBytes9(bytes), nil
			}
		case 10:
			if b10, ok := param.([10]byte); ok {
				return b10, nil
			} else if b10Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b10Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b10Str, err)
				}
				if len(bytes) != 10 {
					return nil, fmt.Errorf("expect length is 10, but got %d", len(bytes))
				}
				return convert.BytesToBytes10(bytes), nil
			}
		case 11:
			if b11, ok := param.([11]byte); ok {
				return b11, nil
			} else if b11Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b11Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b11Str, err)
				}
				if len(bytes) != 11 {
					return nil, fmt.Errorf("expect length is 11, but got %d", len(bytes))
				}
				return convert.BytesToBytes11(bytes), nil
			}
		case 12:
			if b12, ok := param.([12]byte); ok {
				return b12, nil
			} else if b12Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b12Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b12Str, err)
				}
				if len(bytes) != 12 {
					return nil, fmt.Errorf("expect length is 12, but got %d", len(bytes))
				}
				return convert.BytesToBytes12(bytes), nil
			}
		case 13:
			if b13, ok := param.([13]byte); ok {
				return b13, nil
			} else if b13Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b13Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b13Str, err)
				}
				if len(bytes) != 13 {
					return nil, fmt.Errorf("expect length is 13, but got %d", len(bytes))
				}
				return convert.BytesToBytes13(bytes), nil
			}
		case 14:
			if b14, ok := param.([14]byte); ok {
				return b14, nil
			} else if b14Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b14Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b14Str, err)
				}
				if len(bytes) != 14 {
					return nil, fmt.Errorf("expect length is 14, but got %d", len(bytes))
				}
				return convert.BytesToBytes14(bytes), nil
			}
		case 15:
			if b15, ok := param.([15]byte); ok {
				return b15, nil
			} else if b15Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b15Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b15Str, err)
				}
				if len(bytes) != 15 {
					return nil, fmt.Errorf("expect length is 15, but got %d", len(bytes))
				}
				return convert.BytesToBytes15(bytes), nil
			}
		case 16:
			if b16, ok := param.([16]byte); ok {
				return b16, nil
			} else if b16Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b16Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b16Str, err)
				}
				if len(bytes) != 16 {
					return nil, fmt.Errorf("expect length is 16, but got %d", len(bytes))
				}
				return convert.BytesToBytes16(bytes), nil
			}
		case 17:
			if b17, ok := param.([17]byte); ok {
				return b17, nil
			} else if b17Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b17Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b17Str, err)
				}
				if len(bytes) != 17 {
					return nil, fmt.Errorf("expect length is 17, but got %d", len(bytes))
				}
				return convert.BytesToBytes17(bytes), nil
			}
		case 18:
			if b18, ok := param.([18]byte); ok {
				return b18, nil
			} else if b18Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b18Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b18Str, err)
				}
				if len(bytes) != 18 {
					return nil, fmt.Errorf("expect length is 18, but got %d", len(bytes))
				}
				return convert.BytesToBytes18(bytes), nil
			}
		case 19:
			if b19, ok := param.([19]byte); ok {
				return b19, nil
			} else if b19Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b19Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b19Str, err)
				}
				if len(bytes) != 19 {
					return nil, fmt.Errorf("expect length is 19, but got %d", len(bytes))
				}
				return convert.BytesToBytes19(bytes), nil
			}
		case 20:
			if b20, ok := param.([20]byte); ok {
				return b20, nil
			} else if b20Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b20Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b20Str, err)
				}
				if len(bytes) != 20 {
					return nil, fmt.Errorf("expect length is 20, but got %d", len(bytes))
				}
				return convert.BytesToBytes20(bytes), nil
			}
		case 21:
			if b21, ok := param.([21]byte); ok {
				return b21, nil
			} else if b21Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b21Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b21Str, err)
				}
				if len(bytes) != 21 {
					return nil, fmt.Errorf("expect length is 21, but got %d", len(bytes))
				}
				return convert.BytesToBytes21(bytes), nil
			}
		case 22:
			if b22, ok := param.([22]byte); ok {
				return b22, nil
			} else if b22Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b22Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b22Str, err)
				}
				if len(bytes) != 22 {
					return nil, fmt.Errorf("expect length is 22, but got %d", len(bytes))
				}
				return convert.BytesToBytes22(bytes), nil
			}
		case 23:
			if b23, ok := param.([23]byte); ok {
				return b23, nil
			} else if b23Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b23Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b23Str, err)
				}
				if len(bytes) != 23 {
					return nil, fmt.Errorf("expect length is 23, but got %d", len(bytes))
				}
				return convert.BytesToBytes23(bytes), nil
			}
		case 24:
			if b24, ok := param.([24]byte); ok {
				return b24, nil
			} else if b24Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b24Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b24Str, err)
				}
				if len(bytes) != 24 {
					return nil, fmt.Errorf("expect length is 24, but got %d", len(bytes))
				}
				return convert.BytesToBytes24(bytes), nil
			}
		case 25:
			if b25, ok := param.([25]byte); ok {
				return b25, nil
			} else if b25Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b25Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b25Str, err)
				}
				if len(bytes) != 25 {
					return nil, fmt.Errorf("expect length is 25, but got %d", len(bytes))
				}
				return convert.BytesToBytes25(bytes), nil
			}
		case 26:
			if b26, ok := param.([26]byte); ok {
				return b26, nil
			} else if b26Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b26Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b26Str, err)
				}
				if len(bytes) != 26 {
					return nil, fmt.Errorf("expect length is 26, but got %d", len(bytes))
				}
				return convert.BytesToBytes29(bytes), nil
			}
		case 27:
			if b27, ok := param.([27]byte); ok {
				return b27, nil
			} else if b27Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b27Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b27Str, err)
				}
				if len(bytes) != 27 {
					return nil, fmt.Errorf("expect length is 27, but got %d", len(bytes))
				}
				return convert.BytesToBytes27(bytes), nil
			}
		case 28:
			if b28, ok := param.([28]byte); ok {
				return b28, nil
			} else if b28Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b28Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b28Str, err)
				}
				if len(bytes) != 28 {
					return nil, fmt.Errorf("expect length is 28, but got %d", len(bytes))
				}
				return convert.BytesToBytes28(bytes), nil
			}
		case 29:
			if b29, ok := param.([29]byte); ok {
				return b29, nil
			} else if b29Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b29Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b29Str, err)
				}
				if len(bytes) != 29 {
					return nil, fmt.Errorf("expect length is 29, but got %d", len(bytes))
				}
				return convert.BytesToBytes29(bytes), nil
			}
		case 30:
			if b30, ok := param.([30]byte); ok {
				return b30, nil
			} else if b30Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b30Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b30Str, err)
				}
				if len(bytes) != 30 {
					return nil, fmt.Errorf("expect length is 30, but got %d", len(bytes))
				}
				return convert.BytesToBytes30(bytes), nil
			}
		case 31:
			if b31, ok := param.([31]byte); ok {
				return b31, nil
			} else if b31Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b31Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b31Str, err)
				}
				if len(bytes) != 31 {
					return nil, fmt.Errorf("expect length is 31, but got %d", len(bytes))
				}
				return convert.BytesToBytes31(bytes), nil
			}
		case 32:
			if b32, ok := param.([32]byte); ok {
				return b32, nil
			} else if b32Str, ok := param.(string); ok {
				bytes, err := hexutil.Decode(b32Str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string: %s，err: %v", b32Str, err)
				}
				if len(bytes) != 32 {
					return nil, fmt.Errorf("expect length is 32, but got %d", len(bytes))
				}
				return convert.BytesToBytes32(bytes), nil
			}
		}
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
