package types

// Fee 费用
//
// 最小单位为 wei
// 1eth = 1e18wei
// 1eth = 1e9gwei
type Fee string

// VotingRule 投票规则类型
//   - VotingRuleNO        不需要投票
//   - VotingRuleLEADER    盟主一票制
//   - VotingRuleCONSENSUS 共识投票
type VotingRule uint8

const (
	VotingRuleNO VotingRule = iota
	VotingRuleLEADER
	VotingRuleCONSENSUS
)

// Consensus 共识类型
//   - ConsensusPOA  poa共识
//   - ConsensusPBFT pbft共识
//   - ConsensusRAFT raft共识
type Consensus string

const (
	ConsensusPOA  Consensus = "POA"
	ConsensusPBFT Consensus = "PBFT"
	ConsensusRAFT Consensus = "RAFT"
)

// NodeType 节点类型
//   - NodeTypeGENESIS   创世节点
//   - NodeTypeCONSENSUS 共识节点
//   - NodeTypeWITNESS   见证节点
//   - NodeTypeUNKNOWN   未知节点
type NodeType uint8

const (
	NodeTypeGENESIS NodeType = iota
	NodeTypeCONSENSUS
	NodeTypeWITNESS
	NodeTypeUNKNOWN
)

// ContractState 合约状态
//   - ContractStatePROHIBITED  禁止执行合约
//   - ContractStateALLOWABLE   允许执行合约
//   - ContractStateUNAVAILABLE 合约不可调用，处于冻结状态
type ContractState uint8

const (
	ContractStatePROHIBITED ContractState = iota
	ContractStateALLOWABLE
	ContractStateUNAVAILABLE
)

// ContractManagementMode 合约管理模式
//   - ContractManagementModeWHITELIST 白名单模式
//   - ContractManagementModeBLACKLIST 黑名单模式
type ContractManagementMode uint8

const (
	ContractManagementModeWHITELIST ContractManagementMode = iota
	ContractManagementModeBLACKLIST
)

// ContractLang 合约语言
//
//   - ContractLangGo	Go
//   - ContractLangJava Java.
type ContractLang string

const (
	ContractLangGo   ContractLang = "Go"
	ContractLangJava ContractLang = "Java"
)
