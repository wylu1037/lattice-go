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
type ProposalType uint8

const (
	ProposalTypeNone ProposalType = iota
	ProposalTypeContractInnerManager
	ProposalTypeContractLifecycle
	ProposalTypeModifyChainConfig
)
