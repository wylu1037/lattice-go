package types

type Proposal[T ContractLifecycleProposal | ModifyChainConfigProposal] struct {
	Type    uint8 `json:"proposalType"`
	Content *T    `json:"proposalContent"`
}

// ContractLifecycleProposal 合约生命周期提案
type ContractLifecycleProposal struct {
	Id              string `json:"proposalId"`
	State           uint8  `json:"proposalState"`
	Nonce           uint64 `json:"nonce"`
	ContractAddress string `json:"contractAddress"`
	IsRevoke        uint32 `json:"isRevoke"`
	Period          uint8  `json:"period"`
}

// ModifyChainConfigProposal 修改链配置提案
type ModifyChainConfigProposal struct {
	Id             string   `json:"proposalId"`
	State          uint8    `json:"proposalState"`
	Nonce          uint64   `json:"nonce"`
	Type           uint8    `json:"modifyType"`
	Period         uint32   `json:"period"`
	IsDictatorship bool     `json:"isDictatorship"`
	NoEmptyAnchor  bool     `json:"noEmptyAnchor"`
	DeployRule     uint8    `json:"deployRule"`
	LatcSaint      []string `json:"latcSaint"`
	Consensus      string   `json:"consensus"`
}

// ProposalState 提案状态
type ProposalState uint8

const (
	ProposalStateNone ProposalState = iota
	ProposalStateInitial
	ProposalStateSuccess
	ProposalStateFailed
	ProposalStateExpired
	ProposalStateError
)

// ProposalType 提案类型
//   - ProposalTypeNone
//   - ProposalTypeContractManagement		合约内部管理
//   - ProposalTypeContractLifecycle		合约生命周期
//   - ProposalTypeModifyChainConfiguration 修改链配置
type ProposalType uint8

const (
	ProposalTypeNone ProposalType = iota
	ProposalTypeContractManagement
	ProposalTypeContractLifecycle
	ProposalTypeModifyChainConfiguration
)

// VoteSuggestion 投票建议
//   - VoteSuggestionDISAPPROVE 反对
//   - VoteSuggestionAPPROVE	同意
type VoteSuggestion uint8

const (
	VoteSuggestionDISAPPROVE VoteSuggestion = iota
	VoteSuggestionAPPROVE
)

// VoteDetails 投票详情
//   - VoteId
//   - ProposalId
//   - VoteSuggestion
//   - Address
//   - ProposalType
//   - Nonce
//   - CreatedAt
type VoteDetails struct {
	VoteId         string         `json:"voteId"`
	ProposalId     string         `json:"proposalId"`
	VoteSuggestion VoteSuggestion `json:"voteSuggestion"`
	Address        string         `json:"address"`
	ProposalType   ProposalType   `json:"proposalType"`
	Nonce          uint64         `json:"nonce"`
	CreatedAt      uint64         `json:"createAt"`
}
