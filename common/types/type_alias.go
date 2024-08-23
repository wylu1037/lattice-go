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
	VotingRuleNO        VotingRule = 0
	VotingRuleLEADER    VotingRule = 1
	VotingRuleCONSENSUS VotingRule = 2
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
	NodeTypeGENESIS   NodeType = 0
	NodeTypeCONSENSUS NodeType = 1
	NodeTypeWITNESS   NodeType = 2
	NodeTypeUNKNOWN   NodeType = 3
)
