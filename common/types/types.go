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

// Subchain 子链信息
//   - ID 链ID
//   - Name 名称
//   - Desc 描述
//   - LatcGodAddr 守护链的地址
//   - LatcSaints 共识节点列表
//   - Consensus 共识
//   - Epoch 重置投票和检查点的纪元长度
//   - Tokenless false:有通证 true:无通证
//   - Period 出块间隔
//   - NoEmptyAnchor 是否不允许快速出空块
//   - EmptyAnchorPeriodMul 无交易时出快的时间间隔倍数
//   - IsGM 是否使用了Sm2p256v1曲线
//   - RootPublicKey 中心化CA根证书公钥
//   - EnableContractLifecycle 是否开启合约生命周期
//   - EnableVotingDictatorship 是否开启投票时盟主一票制度
//   - ContractDeploymentVotingRule 合约部署的投票规则
//   - EnableContractManagement 是否开启合约管理
//   - ChainByChainVotingRule 以链建链投票规则
//   - ProposalExpirationDays 提案的过期天数，默认7天
//   - ConfigurationModifyVotingRule 配置修改的投票规则
type Subchain struct {
	ID                            uint64     `json:"latcId,omitempty"`
	Name                          string     `json:"name,omitempty"`
	Desc                          string     `json:"desc,omitempty"`
	LatcGodAddr                   string     `json:"latcGod,omitempty"`
	LatcSaints                    []string   `json:"latcSaints,omitempty"`
	Consensus                     string     `json:"consensus,omitempty"`
	Epoch                         uint       `json:"epoch,omitempty"`
	Tokenless                     bool       `json:"tokenless,omitempty"`
	Period                        uint       `json:"period,omitempty"`
	NoEmptyAnchor                 bool       `json:"noEmptyAnchor,omitempty"`
	EmptyAnchorPeriodMul          uint32     `json:"emptyAnchorPeriodMul,omitempty"`
	IsGM                          bool       `json:"GM,omitempty"`
	RootPublicKey                 string     `json:"rootPublicKey,omitempty"`
	EnableContractLifecycle       bool       `json:"isContractVote,omitempty"`
	EnableVotingDictatorship      bool       `json:"isDictatorship,omitempty"`
	ContractDeploymentVotingRule  VotingRule `json:"deployRule,omitempty"`
	EnableContractManagement      bool       `json:"contractPermission,omitempty"`
	ChainByChainVotingRule        VotingRule `json:"chainByChainVote,omitempty"`
	ProposalExpirationDays        uint       `json:"ProposalExpireTime,omitempty"`
	ConfigurationModifyVotingRule VotingRule `json:"configModifyRule,omitempty"`
}

// ConsensusNodeStatus 共识节点的状态
//   - Address						   节点地址
//   - WitnessedBlockCount             见证区块数量
//   - FailureWitnessedBlockCount 	   见证失败的区块数量
//   - ShouldWitnessedBlockCount	   应当见证的区块数量
type ConsensusNodeStatus struct {
	Address                    string `json:"Addr,omitempty"`
	WitnessedBlockCount        uint64 `json:"SignatureCount,omitempty"`
	FailureWitnessedBlockCount uint64 `json:"SignatureFailCount,omitempty"`
	ShouldWitnessedBlockCount  uint64 `json:"ShouldSignatureCount,omitempty"`
	Online                     bool   `json:"online,omitempty"`
}

// NodePeer peer
type NodePeer struct {
	INode     string            `json:"inode"`
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Caps      []string          `json:"caps"`
	Network   NodePeerNetwork   `json:"network"`
	Protocols NodePeerProtocols `json:"protocols"`
}

type NodePeerNetwork struct {
	LocalAddress  string `json:"localAddress"`
	RemoteAddress string `json:"remoteAddress"`
	Inbound       bool   `json:"inbound"`
	Trusted       bool   `json:"trusted"`
	Static        bool   `json:"static"`
}

type NodePeerProtocols struct {
	Latc string `json:"latc"`
}
