package builtin

import "github.com/wylu1037/lattice-go/abi"

func NewContractManagementContract() {}

type ContractManagementContract interface{}

type contractManagementContract struct {
	abi abi.LatticeAbi
}

var ContractManagementBuiltinContract = Contract{
	Description: "合约内部管理合约",
	AbiString: `[
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "Address",
					"type": "address"
				},
				{
					"internalType": "string",
					"name": "operation",
					"type": "string"
				}
			],
			"name": "launch",
			"outputs": [
				{
					"internalType": "string",
					"name": "",
					"type": "string"
				}
			],
			"stateMutability": "pure",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "Address",
					"type": "address"
				},
				{
					"components": [
						{
							"internalType": "uint8",
							"name": "permissionMode",
							"type": "uint8"
						},
						{
							"internalType": "uint64",
							"name": "threshold",
							"type": "uint64"
						},
						{
							"internalType": "address[]",
							"name": "blackList",
							"type": "address[]"
						},
						{
							"internalType": "address[]",
							"name": "whiteList",
							"type": "address[]"
						},
						{
							"components": [
								{
									"internalType": "address",
									"name": "Address",
									"type": "address"
								},
								{
									"internalType": "uint8",
									"name": "weight",
									"type": "uint8"
								}
							],
							"internalType": "struct chainbychain.Manager[]",
							"name": "managerList",
							"type": "tuple[]"
						}
					],
					"internalType": "struct chainbychain.Args",
					"name": "permissionList",
					"type": "tuple"
				}
			],
			"name": "init",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		}
	]`,
	Address: "zltc_ZDdPo8P72X7dtMNTxBeKU8pT7bDXb7NtV",
}
