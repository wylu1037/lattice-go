package lattice

import (
	"lattice-go/common/types"
	"net/http"
)

func NewLattice(ip string, httpPort, wsPort uint16, options *Options) Lattice {
	return &lattice{
		ip:       ip,
		httpPort: httpPort,
		wsPort:   wsPort,
	}
}

type Lattice interface {
}

type lattice struct {
	chainId    types.Number
	curve      types.Curve
	passphrase string
	ip         string
	httpPort   uint16
	wsPort     uint16
	transport  *http.Transport
}

type Options struct {
	Insecure            bool
	MaxIdleConns        int
	MaxIdleConnsPerHost int
}
