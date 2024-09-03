package builtin

import "github.com/wylu1037/lattice-go/abi"

func NewBlockPeekabooContract() {}

type BlockPeekabooContract interface{}

type blockPeekabooContract struct {
	abi abi.LatticeAbi
}

var BlockPeekabooBuiltinContract = Contract{
	Description: "区块隐藏合约",
	Address:     "zltc_a8Nx2gcs2XHye7MKVWykdanumqDkWXqRH",
	AbiString: `[
		{
			"constant": false,
			"inputs": [
				{
					"internalType": "bytes32",
					"name": "_hash",
					"type": "bytes32"
				}
			],
			"name": "addPayload",
			"outputs": [],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
				{
					"internalType": "bytes32",
					"name": "_hash",
					"type": "bytes32"
				}
			],
			"name": "delPayload",
			"outputs": [],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
				{
					"internalType": "bytes32",
					"name": "_hash",
					"type": "bytes32"
				}
			],
			"name": "addHash",
			"outputs": [],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
				{
					"internalType": "bytes32",
					"name": "_hash",
					"type": "bytes32"
				}
			],
			"name": "addCode",
			"outputs": [],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
				{
					"internalType": "bytes32",
					"name": "_hash",
					"type": "bytes32"
				}
			],
			"name": "delHash",
			"outputs": [],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
				{
					"internalType": "bytes32",
					"name": "_hash",
					"type": "bytes32"
				}
			],
			"name": "delCode",
			"outputs": [],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		}
	]`,
}
