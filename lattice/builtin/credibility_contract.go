package builtin

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"lattice-go/abi"
	"lattice-go/common/convert"
)

// CreateBusinessContractAddress 创建存证业务的业务合约地址
const CreateBusinessContractAddress = "zltc_QLbz7JHiBTspS9WTWJUrbNsB5wbENMweQ"

func NewCredibilityContract() CredibilityContract {
	return &credibilityContract{
		abi: abi.NewAbi(CredibilityBuiltinContract.AbiString),
	}
}

type WriteLedgerRequest struct {
	Uri                     uint64 `json:"protocolUri"`
	DataId                  string `json:"hash"`
	Data                    []byte `json:"data"`
	BusinessContractAddress string `json:"address"`
}

type CredibilityContract interface {
	// CreateBusiness 创建业务合约地址
	//
	// Returns:
	//   - data string
	//   - err error
	CreateBusiness() (data string, err error)

	// CreateProtocol 创建协议
	//
	// Parameters:
	//   - tradeNumber uint64
	//   - message []byte
	//
	// Returns:
	//   - data string
	//   - err error
	CreateProtocol(tradeNumber uint64, message []byte) (data string, err error)

	// ReadProtocol 读取协议
	//
	// Parameters:
	//   - uri uint64
	//
	// Returns:
	//   - data string
	//   - err error
	ReadProtocol(uri uint64) (data string, err error)

	// UpdateProtocol 更新协议
	//
	// Parameters:
	//   - uri int64
	//   - message []byte
	//
	// Returns:
	//   - data string
	//   - err error
	UpdateProtocol(uri int64, message []byte) (data string, err error)

	// Write 写入存证数据
	//
	// Parameters:
	//   - request WriteLedgerRequest
	//
	// Returns:
	//   - data string
	//   - err error
	Write(request WriteLedgerRequest) (data string, err error)

	// Read 读取存证数据
	//
	// Parameters:
	//   - dataId string
	//   - businessContractAddress string
	//
	// Returns:
	//   - data string
	//   - err error
	Read(dataId, businessContractAddress string) (data string, err error)
}

type credibilityContract struct {
	abi abi.LatticeAbi
}

func (c *credibilityContract) CreateBusiness() (data string, err error) {
	return hexutil.Encode([]byte{49}), nil
}

func (c *credibilityContract) CreateProtocol(tradeNumber uint64, message []byte) (data string, err error) {
	fn, err := c.abi.GetLatticeFunction("addProtocol", tradeNumber, convert.BytesToBytes32Arr(message))
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *credibilityContract) ReadProtocol(uri uint64) (data string, err error) {
	fn, err := c.abi.GetLatticeFunction("getAddress", uri)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *credibilityContract) UpdateProtocol(uri int64, message []byte) (data string, err error) {
	fn, err := c.abi.GetLatticeFunction("updateProtocol", uri, convert.BytesToBytes32Arr(message))
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *credibilityContract) Write(request WriteLedgerRequest) (data string, err error) {
	fn, err := c.abi.GetLatticeFunction("writeTraceability", request.Uri, request.DataId, convert.BytesToBytes32Arr(request.Data), request.BusinessContractAddress)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *credibilityContract) Read(dataId, businessContractAddress string) (data string, err error) {
	fn, err := c.abi.GetLatticeFunction("getTraceability", dataId, businessContractAddress)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

var CredibilityBuiltinContract = Contract{
	Description: "存证溯源合约",
	Address:     "zltc_QLbz7JHiBTspUvTPzLHy5biDS9mu53mmv",
	AbiString: `[
		{
			"inputs": [
				{
					"internalType": "uint64",
					"name": "protocolSuite",
					"type": "uint64"
				},
				{
					"internalType": "bytes32[]",
					"name": "data",
					"type": "bytes32[]"
				}
			],
			"name": "addProtocol",
			"outputs": [
				{
					"internalType": "uint64",
					"name": "protocolUri",
					"type": "uint64"
				}
			],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "uint64",
					"name": "protocolUri",
					"type": "uint64"
				}
			],
			"name": "getAddress",
			"outputs": [
				{
					"components": [
						{
							"internalType": "address",
							"name": "updater",
							"type": "address"
						},
						{
							"internalType": "bytes32[]",
							"name": "data",
							"type": "bytes32[]"
						}
					],
					"internalType": "struct credibilidity.Protocol[]",
					"name": "protocol",
					"type": "tuple[]"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "uint64",
					"name": "protocolUri",
					"type": "uint64"
				},
				{
					"internalType": "bytes32[]",
					"name": "data",
					"type": "bytes32[]"
				}
			],
			"name": "updateProtocol",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "string",
					"name": "hash",
					"type": "string"
				},
				{
					"internalType": "address",
					"name": "address",
					"type": "address"
				}
			],
			"name": "getTraceability",
			"outputs": [
				{
					"components": [
						{
							"internalType": "uint64",
							"name": "number",
							"type": "uint64"
						},
						{
							"internalType": "uint64",
							"name": "protocol",
							"type": "uint64"
						},
						{
							"internalType": "address",
							"name": "updater",
							"type": "address"
						},
						{
							"internalType": "bytes32[]",
							"name": "data",
							"type": "bytes32[]"
						}
					],
					"internalType": "struct credibilidity.Evidence[]",
					"name": "evi",
					"type": "tuple[]"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "uint64",
					"name": "protocolUri",
					"type": "uint64"
				},
				{
					"internalType": "string",
					"name": "hash",
					"type": "string"
				},
				{
					"internalType": "bytes32[]",
					"name": "data",
					"type": "bytes32[]"
				},
				{
					"internalType": "address",
					"name": "address",
					"type": "address"
				}
			],
			"name": "writeTraceability",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"components": [
						{
							"internalType": "uint64",
							"name": "protocolUri",
							"type": "uint64"
						},
						{
							"internalType": "string",
							"name": "hash",
							"type": "string"
						},
						{
							"internalType": "bytes32[]",
							"name": "data",
							"type": "bytes32[]"
						},
						{
							"internalType": "address",
							"name": "address",
							"type": "address"
						}
					],
					"internalType": "struct Business.batch[]",
					"name": "bt",
					"type": "tuple[]"
				}
			],
			"name": "writeTraceabilityBatch",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		}
	]`,
}
