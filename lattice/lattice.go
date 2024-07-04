package lattice

import (
	"lattice-go/common/types"
	"net/http"
)

func NewLattice(chainConfig *ChainConfig, nodeConfig *NodeConfig, identityConfig *IdentityConfig, options *Options) Lattice {
	return &lattice{
		ChainConfig:    chainConfig,
		NodeConfig:     nodeConfig,
		IdentityConfig: identityConfig,
		Options:        options,
	}
}

type Lattice interface {
}

type lattice struct {
	ChainConfig    *ChainConfig
	NodeConfig     *NodeConfig
	IdentityConfig *IdentityConfig
	Options        *Options
}

type ChainConfig struct {
	ChainId uint64
	Curve   types.Curve
}

type NodeConfig struct {
	Insecure      bool
	Ip            string
	HttpPort      uint16
	WebsocketPort uint16
}

type IdentityConfig struct {
	Passphrase string
	PrivateKey string
}

type Options struct {
	Transport *http.Transport
}
