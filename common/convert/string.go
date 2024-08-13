package convert

import "math/big"

// StringToBigInt 将字符串数子转为big.Int类型
//
// Parameters:
//   - data string
//
// Returns:
//   - *big.Int
func StringToBigInt(data string) *big.Int {
	num := new(big.Int)
	num.SetString(data, 10)
	return num
}
