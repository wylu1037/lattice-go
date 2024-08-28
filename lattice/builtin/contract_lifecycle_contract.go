package builtin

import "github.com/wylu1037/lattice-go/abi"

const (
	ContractLifecycleREVOKE = iota
	ContractLifecycleUNFREEZE
	ContractLifecycleFREEZE
)

func NewContractLifecycleContract() ContractLifecycleContract {
	return &contractLifecycleContract{
		abi: abi.NewAbi(ContractLifecycleBuiltinContract.AbiString),
	}
}

type ContractLifecycleContract interface {

	// ContractAddress 获取以链建链的合约地址
	//
	// Returns:
	//   - string: 合约地址，zltc_ZQJjaw74CKMjqYJFMKdEDaNTDMq5QKi3T
	ContractAddress() string

	// Freeze 冻结合约
	//
	// Parameters:
	//   - contractAddress string: 合约地址
	Freeze(contractAddress string) (string, error)

	// Unfreeze 解冻合约
	//
	// Parameters:
	//   - contractAddress string: 合约地址
	Unfreeze(contractAddress string) (string, error)

	// Revoke 吊销合约
	//
	// Parameters:
	//   - contractAddress string: 合约地址
	Revoke(contractAddress string) (string, error)
}

type contractLifecycleContract struct {
	abi abi.LatticeAbi
}

func (c *contractLifecycleContract) ContractAddress() string {
	return ContractLifecycleBuiltinContract.Address
}

func (c *contractLifecycleContract) Freeze(contractAddress string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("launch", contractAddress, ContractLifecycleFREEZE)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *contractLifecycleContract) Unfreeze(contractAddress string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("launch", contractAddress, ContractLifecycleUNFREEZE)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *contractLifecycleContract) Revoke(contractAddress string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("launch", contractAddress, ContractLifecycleREVOKE)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

var ContractLifecycleBuiltinContract = Contract{
	Description: "合约生命周期合约",
	Address:     "zltc_ZQJjaw74CKMjqYJFMKdEDaNTDMq5QKi3T",
	AbiString: `[
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "Address",
					"type": "address"
				},
				{
					"internalType": "uint8",
					"name": "IsRevoke",
					"type": "uint8"
				}
			],
			"name": "launch",
			"outputs": [
				{
					"internalType": "bytes",
					"name": "",
					"type": "bytes"
				}
			],
			"stateMutability": "pure",
			"type": "function"
		}
	]`,
}
