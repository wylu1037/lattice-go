package builtin

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/wylu1037/lattice-go/abi"
	"math/big"
)

const (
	SubChainWitnessMember   = iota // 见证身份的链成员
	SubChainConsensusMember        // 共识身份的链成员
)

// NewSubChainRequest 创建一个子链的请求结构体
type NewSubChainRequest struct {
	Consensus            uint8            `json:"consensus"`            // 0:继承主链1: poa 2:pbft 3:raft 默认
	Tokenless            bool             `json:"tokenless"`            // 是否有通证
	GodAmount            *big.Int         `json:"godAmount"`            // 盟主初始余额
	Period               uint64           `json:"period"`               // 出块间隔
	NoEmptyAnchor        bool             `json:"noEmptyAnchor"`        // 不允许快速出空块
	EmptyAnchorPeriodMul uint64           `json:"emptyAnchorPeriodMul"` // 空块等待次数
	IsContractVote       bool             `json:"isContractVote"`       // 开启合约生命周期
	IsDictatorship       bool             `json:"isDictatorship"`       // 开启盟主独裁
	DeployRule           uint8            `json:"deployRule"`           // 合约部署规则
	ChannelName          string           `json:"name"`                 // 链名称
	ChannelId            *big.Int         `json:"chainId"`              // 链id
	Preacher             string           `json:"preacher"`             // 创世节点地址
	BootStrap            string           `json:"bootStrap"`            // 创世节点Inode
	ChainMemberGroup     []SubChainMember `json:"chainMemberGroup"`     // 链成员
	ContractPermission   bool             `json:"contractPermission"`   // 合约内部管理开关
	ChainByChainVote     uint8            `json:"chainByChainVote"`     // 以链建链投票开关
	ProposalExpireTime   uint             `json:"proposalExpireTime"`   // 提案过期时间（天）
	ChannelDesc          string           `json:"desc"`                 // 链描述
	Extra                []byte           `json:"extra"`                // 暂时不用的字段
}

type SubChainMember struct {
	Type   uint8          `json:"memberType"` // 成员类型，0-见证、1-共识, SubChainWitnessMember or SubChainConsensusMember
	Member common.Address `json:"member"`     // 节点ZLTC地址
}

func (req *NewSubChainRequest) ToCallContractParams() (string, error) {
	bytes, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// JoinSubChainRequest 加入子链请求
type JoinSubChainRequest struct {
	ChannelId     *big.Int         `json:"chainId"`       // 待加入的链ID
	NetworkId     uint64           `json:"networkId"`     // 待加入的链的所在的网络ID
	NodeInfo      string           `json:"nodeInfo"`      // 指定一个已经加入该链的节点地址
	AccessMembers []common.Address `json:"accessMembers"` // 指定哪些节点加入该链
}

func NewChainBuildsChainContract() ChainBuildsChainContract {
	return &chainBuildsChainContract{
		abi: abi.NewAbi(ChainBuildsChainBuiltinContract.AbiString),
	}
}

type ChainBuildsChainContract interface {

	// ContractAddress 获取以链建链的合约地址
	//
	// Returns:
	//   - string: 合约地址，zltc_ZDfqCd4ZbBi4WA7uG4cGpFWRyTFqzyHUn
	ContractAddress() string

	// NewSubChain 创建子链
	//
	// Parameters:
	//   - req *NewSubChainRequest
	//
	// Returns:
	//   - data string
	//   - err error
	NewSubChain(req *NewSubChainRequest) (data string, err error)

	// DeleteSubChain 删除子链
	//
	// Parameters:
	//   - SubChainId string: 子链id
	//
	// Returns:
	//   - data string
	//   - err error
	DeleteSubChain(SubChainId string) (data string, err error)

	// JoinSubChain 加入子链
	//
	// Parameters:
	//   - req *JoinSubChainRequest
	//
	// Returns:
	//   - data string
	//   - err error
	JoinSubChain(req *JoinSubChainRequest) (data string, err error)

	// StartSubChain 启动子链
	//
	// Parameters:
	//   - SubChainId string: 子链id
	//
	// Returns:
	//   - data string
	//   - err error
	StartSubChain(SubChainId string) (data string, err error)

	// StopSubChain 停止子链
	//
	// Parameters:
	//   - SubChainId string: 子链id
	//
	// Returns:
	//   - data string
	//   - err error
	StopSubChain(SubChainId string) (data string, err error)
}

type chainBuildsChainContract struct {
	abi abi.LatticeAbi
}

func (c *chainBuildsChainContract) ContractAddress() string {
	return ChainBuildsChainBuiltinContract.Address
}

func (c *chainBuildsChainContract) NewSubChain(req *NewSubChainRequest) (data string, err error) {
	args, err := req.ToCallContractParams()
	if err != nil {
		return "", err
	}
	code, err := c.abi.MyAbi().Pack("newChain", args)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(code), nil
}

func (c *chainBuildsChainContract) DeleteSubChain(ChannelId string) (data string, err error) {
	fn, err := c.abi.GetLatticeFunction("delChain", ChannelId)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *chainBuildsChainContract) JoinSubChain(req *JoinSubChainRequest) (data string, err error) {
	code, err := c.abi.MyAbi().Pack("oldChain", req.ChannelId, req.NetworkId, req.NodeInfo, req.AccessMembers)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(code), nil
}

func (c *chainBuildsChainContract) StartSubChain(ChannelId string) (data string, err error) {
	fn, err := c.abi.GetLatticeFunction("startChain", ChannelId)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *chainBuildsChainContract) StopSubChain(ChannelId string) (data string, err error) {
	fn, err := c.abi.GetLatticeFunction("stopChain", ChannelId)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

var ChainBuildsChainBuiltinContract = Contract{
	Description: "以链建链合约",
	Address:     "zltc_ZDfqCd4ZbBi4WA7uG4cGpFWRyTFqzyHUn",
	AbiString: `[
		{
			"inputs": [
				{
					"internalType": "uint256",
					"name": "chainId",
					"type": "uint256"
				}
			],
			"name": "delChain",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "string",
					"name": "jsonMap",
					"type": "string"
				}
			],
			"name": "newChain",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "uint256",
					"name": "chainId",
					"type": "uint256"
				},
				{
					"internalType": "uint64",
					"name": "networkId",
					"type": "uint64"
				},
				{
					"internalType": "string",
					"name": "nodeInfo",
					"type": "string"
				},
				{
					"internalType": "address[]",
					"name": "accessMembers",
					"type": "address[]"
				}
			],
			"name": "oldChain",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "uint256",
					"name": "chainId",
					"type": "uint256"
				}
			],
			"name": "stopChain",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "uint256",
					"name": "chainId",
					"type": "uint256"
				}
			],
			"name": "startChain",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		}
	]`,
}
