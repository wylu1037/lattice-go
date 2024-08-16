package types

import "math/big"

const (
	AddressVersion = 1
	AddressLength  = 20 // 20 byte
	AddressTitle   = "zltc"
	HashLength     = 32 // 32 byte
)

// Curve Elliptic curve
type Curve string

type Number string

func (n Number) MustToBigInt() *big.Int {
	num := new(big.Int)
	num.SetString(string(n), 10)
	return num
}

// UploadFileResponse 文件上传到链上的返回结果
//
//   - CID 文件唯一标识，示例：GPK3PveRaWoK6S2b53D3ZeJTm4nBvv2vSVjRStRLQcyX
//   - FilePath 文件存储地址，示例：JG-DFS/tempFileDir/20240816/1723793768943848748_avatar.svg"
//   - Message 返回信息，示例：success
//   - OccupiedStorageByte 文件占用的存储字节数，单位为byte，示例：255686
//   - StorageAddress 需要冗余存储文件的节点地址，示例：DFS_beforeSign||zltc_Z1pnS94bP4hQSYLs4aP4UwBP9pH8bEvhi;zltc_nBGgKoo1rzd4thjfauEN6ULj7jp1zXhxE;
type UploadFileResponse struct {
	CID                 string `json:"cid,omitempty"`
	FilePath            string `json:"filePath,omitempty"`
	Message             string `json:"message,omitempty"`
	OccupiedStorageByte int64  `json:"needStorageSize,omitempty"`
	StorageAddress      string `json:"storageAddress,omitempty"`
}
