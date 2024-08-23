package types

// Fee 费用
//
// 最小单位为 wei
// 1eth = 1e18wei
// 1eth = 1e9gwei
type Fee string

// VotingRule 投票规则类型
//   - VotingRuleNO 不需要投票
//   - VotingRuleLEADER 盟主一票制
//   - VotingRuleCONSENSUS 共识投票
type VotingRule uint8

const (
	VotingRuleNO        VotingRule = 0
	VotingRuleLEADER    VotingRule = 1
	VotingRuleCONSENSUS VotingRule = 2
)

// Consensus 共识类型
//   - ConsensusPOA poa共识
//   - ConsensusPBFT pbft共识
//   - ConsensusRAFT raft共识
type Consensus string

const (
	ConsensusPOA  Consensus = "POA"
	ConsensusPBFT Consensus = "PBFT"
	ConsensusRAFT Consensus = "RAFT"
)