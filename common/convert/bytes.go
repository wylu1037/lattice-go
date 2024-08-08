package convert

import "encoding/hex"

// BytesToBytes32Arr 将[]byte转为[][32]byte，长度不足时在数组尾部补0
// Parameters:
//   - bytes []byte: Example: [1, 2, 3, 4, 5]
//
// Returns:
//   - [][32]byte: Example: [[1 2 3 4 5 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]]
func BytesToBytes32Arr(bytes []byte) [][32]byte {
	var bytes32Arr [][32]byte
	bytes = PadToMultipleOf32(bytes)
	for i := 0; i < len(bytes); i += 32 {
		var b32 [32]byte
		copy(b32[:], bytes[i:i+32])
		bytes32Arr = append(bytes32Arr, b32)
	}
	return bytes32Arr
}

func BytesToBytes32HexArr(bytes []byte) []string {
	var bytes32Arr []string
	bytes = PadToMultipleOf32(bytes)
	for i := 0; i < len(bytes); i += 32 {
		var b32 [32]byte
		copy(b32[:], bytes[i:i+32])
		bytes32Arr = append(bytes32Arr, hex.EncodeToString(b32[:]))
	}
	return bytes32Arr
}

// PadToMultipleOf32 将[]byte补0至为32的整数倍
// Parameters:
//   - bytes []byte: Example: [1, 2, 3, 4, 5]
//
// Returns:
//   - []byte: Example: [1 2 3 4 5 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
func PadToMultipleOf32(bytes []byte) []byte {
	padding := 32 - len(bytes)%32
	if padding == 32 {
		return bytes
	}

	padded := make([]byte, len(bytes)+padding)
	copy(padded, bytes)
	return padded
}

func BytesToBytes1(bytes []byte) [1]byte {
	var bytes1 [1]byte
	copy(bytes1[:], bytes)
	return bytes1
}

func BytesToBytes2(bytes []byte) [2]byte {
	var bytes2 [2]byte
	copy(bytes2[:], bytes)
	return bytes2
}

func BytesToBytes3(bytes []byte) [3]byte {
	var bytes3 [3]byte
	copy(bytes3[:], bytes)
	return bytes3
}

func BytesToBytes4(bytes []byte) [4]byte {
	var bytes4 [4]byte
	copy(bytes4[:], bytes)
	return bytes4
}
func BytesToBytes5(bytes []byte) [5]byte {
	var bytes5 [5]byte
	copy(bytes5[:], bytes)
	return bytes5
}
func BytesToBytes6(bytes []byte) [6]byte {
	var bytes6 [6]byte
	copy(bytes6[:], bytes)
	return bytes6
}
func BytesToBytes7(bytes []byte) [7]byte {
	var bytes7 [7]byte
	copy(bytes7[:], bytes)
	return bytes7
}
func BytesToBytes8(bytes []byte) [8]byte {
	var bytes8 [8]byte
	copy(bytes8[:], bytes)
	return bytes8
}
func BytesToBytes9(bytes []byte) [9]byte {
	var bytes9 [9]byte
	copy(bytes9[:], bytes)
	return bytes9
}
func BytesToBytes10(bytes []byte) [10]byte {
	var bytes10 [10]byte
	copy(bytes10[:], bytes)
	return bytes10
}
func BytesToBytes11(bytes []byte) [11]byte {
	var bytes11 [11]byte
	copy(bytes11[:], bytes)
	return bytes11
}
func BytesToBytes12(bytes []byte) [12]byte {
	var bytes12 [12]byte
	copy(bytes12[:], bytes)
	return bytes12
}
func BytesToBytes13(bytes []byte) [13]byte {
	var bytes13 [13]byte
	copy(bytes13[:], bytes)
	return bytes13
}
func BytesToBytes14(bytes []byte) [14]byte {
	var bytes14 [14]byte
	copy(bytes14[:], bytes)
	return bytes14
}
func BytesToBytes15(bytes []byte) [15]byte {
	var bytes15 [15]byte
	copy(bytes15[:], bytes)
	return bytes15
}
func BytesToBytes16(bytes []byte) [16]byte {
	var bytes16 [16]byte
	copy(bytes16[:], bytes)
	return bytes16
}
func BytesToBytes17(bytes []byte) [17]byte {
	var bytes17 [17]byte
	copy(bytes17[:], bytes)
	return bytes17
}
func BytesToBytes18(bytes []byte) [18]byte {
	var bytes18 [18]byte
	copy(bytes18[:], bytes)
	return bytes18
}
func BytesToBytes19(bytes []byte) [19]byte {
	var bytes19 [19]byte
	copy(bytes19[:], bytes)
	return bytes19
}
func BytesToBytes20(bytes []byte) [20]byte {
	var bytes20 [20]byte
	copy(bytes20[:], bytes)
	return bytes20
}
func BytesToBytes21(bytes []byte) [21]byte {
	var bytes21 [21]byte
	copy(bytes21[:], bytes)
	return bytes21
}
func BytesToBytes22(bytes []byte) [22]byte {
	var bytes22 [22]byte
	copy(bytes22[:], bytes)
	return bytes22
}
func BytesToBytes23(bytes []byte) [23]byte {
	var bytes23 [23]byte
	copy(bytes23[:], bytes)
	return bytes23
}
func BytesToBytes24(bytes []byte) [24]byte {
	var bytes24 [24]byte
	copy(bytes24[:], bytes)
	return bytes24
}
func BytesToBytes25(bytes []byte) [25]byte {
	var bytes25 [25]byte
	copy(bytes25[:], bytes)
	return bytes25
}
func BytesToBytes26(bytes []byte) [26]byte {
	var bytes26 [26]byte
	copy(bytes26[:], bytes)
	return bytes26
}
func BytesToBytes27(bytes []byte) [27]byte {
	var bytes27 [27]byte
	copy(bytes27[:], bytes)
	return bytes27
}
func BytesToBytes28(bytes []byte) [28]byte {
	var bytes28 [28]byte
	copy(bytes28[:], bytes)
	return bytes28
}
func BytesToBytes29(bytes []byte) [29]byte {
	var bytes29 [29]byte
	copy(bytes29[:], bytes)
	return bytes29
}
func BytesToBytes30(bytes []byte) [30]byte {
	var bytes30 [30]byte
	copy(bytes30[:], bytes)
	return bytes30
}

func BytesToBytes31(bytes []byte) [31]byte {
	var bytes31 [31]byte
	copy(bytes31[:], bytes)
	return bytes31
}

func BytesToBytes32(bytes []byte) [32]byte {
	var bytes32 [32]byte
	copy(bytes32[:], bytes)
	return bytes32
}
