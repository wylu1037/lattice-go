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

// NodeInfo 节点
//
//   - ID 示例：16Uiu2HAmQ7Da6iuScYSYs8XGJs95hiKdS6tgmbqUUuKC62Xh3s4V
//   - Name 示例：ZLTC2_1
//   - Version
//   - INode 节点连接信息，示例：/ip4/192.168.1.185/tcp/13801/p2p/16Uiu2HAmQ7Da6iuScYSYs8XGJs95hiKdS6tgmbqUUuKC62Xh3s4V
//   - Inr
//   - IP
//   - Ports
//   - ListenAddress
type NodeInfo struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	INode   string `json:"inode,omitempty"`
	Inr     string `json:"inr,omitempty"`
	IP      string `json:"ip,omitempty"`
	Ports   struct {
		P2PPort       uint16 `json:"discovery,omitempty"`
		WebsocketPort uint16 `json:"listener,omitempty"`
		HTTPPort      uint16 `json:"httpPort,omitempty"`
	}
	ListenAddress string `json:"listenAddr,omitempty"`
}
