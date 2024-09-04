package types

import (
	"encoding/hex"
	"encoding/json"
	"github.com/wylu1037/lattice-go/common/constant"
	"math/big"
)

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
//   - ID                              链ID
//   - Name                            名称
//   - Desc                            描述
//   - LatcGodAddr                     守护链的地址
//   - LatcSaints                      共识节点列表
//   - Consensus                       共识
//   - Epoch                           重置投票和检查点的纪元长度
//   - Tokenless                       false:有通证 true:无通证
//   - Period                          出块间隔
//   - EnableNoTxDelayedMining         是否不允许无交易时快速出空块，无交易时延迟出块
//   - NoTxDelayedMiningPeriodMultiple 无交易时的延迟出块间隔倍数
//   - IsGM                            是否使用了Sm2p256v1曲线
//   - RootPublicKey                   中心化CA根证书公钥
//   - EnableContractLifecycle         是否开启合约生命周期
//   - EnableVotingDictatorship        是否开启投票(合约生命周期)时盟主一票制度
//   - ContractDeploymentVotingRule    合约部署的投票规则
//   - EnableContractManagement        是否开启合约管理
//   - ChainByChainVotingRule          以链建链投票规则
//   - ProposalExpirationDays          提案的过期天数，默认7天
//   - ConfigurationModifyVotingRule   配置修改的投票规则
type Subchain struct {
	ID                              uint64     `json:"latcId,omitempty"`
	Name                            string     `json:"name,omitempty"`
	Desc                            string     `json:"desc,omitempty"`
	LatcGodAddr                     string     `json:"latcGod,omitempty"`
	LatcSaints                      []string   `json:"latcSaints,omitempty"`
	Consensus                       string     `json:"consensus,omitempty"`
	Epoch                           uint       `json:"epoch,omitempty"`
	Tokenless                       bool       `json:"tokenless,omitempty"`
	Period                          uint       `json:"period,omitempty"`
	EnableNoTxDelayedMining         bool       `json:"noEmptyAnchor,omitempty"`
	NoTxDelayedMiningPeriodMultiple uint32     `json:"emptyAnchorPeriodMul,omitempty"`
	IsGM                            bool       `json:"GM,omitempty"`
	RootPublicKey                   string     `json:"rootPublicKey,omitempty"`
	EnableContractLifecycle         bool       `json:"isContractVote,omitempty"`
	EnableVotingDictatorship        bool       `json:"isDictatorship,omitempty"`
	ContractDeploymentVotingRule    VotingRule `json:"deployRule,omitempty"`
	EnableContractManagement        bool       `json:"contractPermission,omitempty"`
	ChainByChainVotingRule          VotingRule `json:"chainByChainVote,omitempty"`
	ProposalExpirationDays          uint       `json:"ProposalExpireTime,omitempty"`
	ConfigurationModifyVotingRule   VotingRule `json:"configModifyRule,omitempty"`
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

// NodeConfig 节点配置信息
type NodeConfig struct {
	Lattice struct {
		NetworkIDGroup []int `json:"networkIDGroup"`
	} `json:"latc"`
	Node struct {
		Name                                    string   `json:"name"`
		DataDir                                 string   `json:"dataDir"`
		SecondaryDir                            string   `json:"secondaryDir"`
		Host                                    string   `json:"Host"`
		HTTPPort                                int      `json:"HTTPPort"`
		WSPort                                  int      `json:"WSPort"`
		P2PPort                                 int      `json:"P2PPort"`
		GinHTTPPort                             int      `json:"GinHTTPPort"`
		JWTEnable                               bool     `json:"JWTEnable"`
		JWTSecret                               string   `json:"JWTSecret"`
		Bootstrap                               []string `json:"bootstrap"`
		MultilingualContractBaseDir             string   `json:"basedir"`
		MultilingualContractPattern             string   `json:"pattern"`    // 生成合约执行路径pattern
		MultilingualContractRedundantDeployment bool     `json:"redundancy"` // 是否支持冗余部署，即对一个合约文件部署多次，生成多个合约地址
	}
}

// ContractInformation 合约信息
//
//   - ContractAddress 合约地址
//   - Owner           合约的部署者地址
//   - State           合约的状态
//   - Version         合约的版本
//   - ProposalId      合约的提案ID，包括 部署、升级、吊销
//   - CreatedAt       合约的部署时间
//   - UpdatedAt       合约的修改时间
type ContractInformation struct {
	ContractAddress string        `json:"address"`
	Owner           string        `json:"deploymentAddress"`
	State           ContractState `json:"state"`
	Version         uint8         `json:"version"`
	ProposalId      string        `json:"votingProposalId,omitempty"`
	CreatedAt       uint64        `json:"createAt"`
	UpdatedAt       uint64        `json:"modifiedAt"`
}

// ContractManagement 合约管理信息
//
//   - Mode           合约管理模式，白名单 or 黑名单
//   - Threshold      投票通过的阈值，大于10则按照权重加和，小于等于10则按照百分比
//   - Whitelist	  合约白名单
//   - Blacklist	  合约黑名单
//   - Administrators 合约管理员：`{"zltc_Z1pnS94bP4hQSYLs4aP4UwBP9pH8bEvhi": 10}`
type ContractManagement struct {
	Mode           ContractManagementMode `json:"permissionMode"`
	Threshold      uint64                 `json:"threshold"`
	Whitelist      []string               `json:"whiteList"`
	Blacklist      []string               `json:"blackList"`
	Administrators map[string]uint8       `json:"managerList"`
}

// DeployMultilingualContractCode 部署多语言智能合约的代码
//   - FileName 上传到链上的合约文件名
type DeployMultilingualContractCode struct {
	FileName string `json:"contractName,omitempty"`
}

// UpgradeMultilingualContractCode 升级多语言智能合约的代码
//   - FileName 上传到链上的合约文件名
type UpgradeMultilingualContractCode struct {
	FileName string `json:"contractName,omitempty"`
}

// CallMultilingualContractCode 调用多语言智能合约的代码
//   - Method	 调用的合约方法名，示例：`double`
//   - Arguments 调用的合约方法参数，示例：`{"number":[56,50,55]}`
type CallMultilingualContractCode struct {
	Method    string            `json:"methodName,omitempty"`
	Arguments map[string][]byte `json:"methodArgs,omitempty"`
}

func (c *DeployMultilingualContractCode) Encode() string {
	if bytes, err := json.Marshal(c); err != nil {
		return constant.HexPrefix
	} else {
		return constant.HexPrefix + hex.EncodeToString(bytes)
	}
}

func (c *UpgradeMultilingualContractCode) Encode() string {
	if bytes, err := json.Marshal(c); err != nil {
		return constant.HexPrefix
	} else {
		return constant.HexPrefix + hex.EncodeToString(bytes)
	}
}

func (c *CallMultilingualContractCode) Encode() string {
	if bytes, err := json.Marshal(c); err != nil {
		return constant.HexPrefix
	} else {
		return constant.HexPrefix + hex.EncodeToString(bytes)
	}
}

// NodeProtocol 节点的网络协议信息
type NodeProtocol struct {
	Genesis        string              `json:"genesis"`
	NetWorkIdGroup []int               `json:"netWorkIdGroup"`
	Config         *NodeProtocolConfig `json:"config"`
}

// NodeProtocolConfig
// - LatcID                          链ID
// - Name                            名称
// - Desc                            描述
// - LatcGodAddr                     守护链的地址
// - LatcSaints                      共识节点列表
// - Consensus                       共识
// - Epoch                           重置投票和检查点的纪元长度
// - Tokenless                       false:有通证 true:无通证
// - Period                          出块间隔
// - EnableNoTxDelayedMining         是否不允许无交易时快速出空块，无交易时延迟出块
// - NoTxDelayedMiningPeriodMultiple 无交易时的延迟出块间隔倍数
// - IsGM                            是否使用了Sm2p256v1曲线
// - RootPublicKey                   中心化CA根证书公钥
// - EnableContractLifecycle         是否开启合约生命周期
// - EnableVotingDictatorship        是否开启投票(合约生命周期)时盟主一票制度
// - ContractDeploymentVotingRule    合约部署的投票规则
// - EnableContractManagement        是否开启合约管理
// - ChainByChainVotingRule          以链建链投票规则
// - ProposalExpirationDays          提案的过期天数，默认7天
// - ConfigurationModifyVotingRule   配置修改的投票规则
type NodeProtocolConfig struct {
	LatcID                          *big.Int   `json:"latcId,omitempty"`
	Name                            string     `json:"Name,omitempty"`
	Desc                            string     `json:"Desc,omitempty"`
	LatcGodAddr                     string     `json:"latcGod,omitempty"`
	LatcSaints                      []string   `json:"latcSaints,omitempty"`
	Consensus                       string     `json:"consensus,omitempty"`
	Epoch                           uint       `json:"epoch,omitempty"`
	Tokenless                       bool       `json:"tokenless,omitempty"`
	Period                          uint       `json:"period,omitempty"`
	EnableNoTxDelayedMining         bool       `json:"noEmptyAnchor,omitempty"`
	NoTxDelayedMiningPeriodMultiple uint64     `json:"emptyAnchorPeriodMul,omitempty"`
	IsGM                            bool       `json:"GM,omitempty"`
	RootPublicKey                   string     `json:"rootPublicKey,omitempty"`
	EnableContractLifecycle         bool       `json:"isContractVote,omitempty"`
	EnableVotingDictatorship        bool       `json:"isDictatorship,omitempty"`
	ContractDeploymentVotingRule    VotingRule `json:"deployRule,omitempty"`
	EnableContractManagement        bool       `json:"contractPermission,omitempty"`
	ChainByChainVotingRule          VotingRule `json:"chainByChainVote,omitempty"`
	ProposalExpirationDays          uint       `json:"ProposalExpireTime,omitempty"`
	ConfigurationModifyVotingRule   VotingRule `json:"configModifyRule,omitempty"`
}

// Evidences 留痕信息
type Evidences struct {
	Total uint64                 `json:"total"`
	Data  map[string]interface{} `json:"data"`
}
