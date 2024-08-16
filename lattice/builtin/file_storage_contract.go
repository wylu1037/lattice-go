package builtin

import "github.com/wylu1037/lattice-go/abi"

func NewFileStorageContract() FileStorageContract {
	return &fileStorageContract{
		abi: abi.NewAbi(FileStorageBuiltinContract.AbiString),
	}
}

type FileStorageContract interface {

	// UploadFile 上传文件
	//
	// Parameters:
	//   - accountAddress string
	//   - filePath  string
	//   - nodeAddress string
	//   - storageAddress string
	//   - OccupiedStorageByte int64
	//   - cid string
	//
	// Returns:
	//   - string
	//   - error
	UploadFile(accountAddress, filePath, nodeAddress, storageAddress string, OccupiedStorageByte int64, cid string) (string, error)

	// UpdatePermission 分配文件存储的磁盘空间
	//
	// Parameters:
	//   - allocatorAddress string: 分配者账户地址
	//   - allotteeAddress string: 被分配者账户地址
	//   - totalStorageByte int64: 分配的总磁盘空间（字节）
	//
	// Returns:
	//   - string
	//   - error
	UpdatePermission(allocatorAddress, allotteeAddress string, totalStorageByte int64) (string, error)

	// DownloadFile 下载文件
	//
	// Parameters:
	//   - accountAddress string
	//   - cid string
	//
	// Returns:
	//   - string
	//   - error
	DownloadFile(accountAddress, cid string) (string, error)
}

type fileStorageContract struct {
	abi abi.LatticeAbi
}

func (c *fileStorageContract) UploadFile(accountAddress, filePath, nodeAddress, storageAddress string, OccupiedStorageByte int64, cid string) (string, error) {
	function, err := c.abi.GetLatticeFunction("UploadFile", accountAddress, filePath, nodeAddress, storageAddress, OccupiedStorageByte, cid)
	if err != nil {
		return "", err
	}
	return function.Encode()
}

func (c *fileStorageContract) UpdatePermission(allocatorAddress, allotteeAddress string, totalStorageByte int64) (string, error) {
	function, err := c.abi.GetLatticeFunction("UpdatePermission", allocatorAddress, allotteeAddress, totalStorageByte)
	if err != nil {
		return "", err
	}
	return function.Encode()
}

func (c *fileStorageContract) DownloadFile(accountAddress, cid string) (string, error) {
	function, err := c.abi.GetLatticeFunction("DownloadFile", accountAddress, cid)
	if err != nil {
		return "", err
	}
	return function.Encode()
}

var FileStorageBuiltinContract = Contract{
	Description: "文件上链合约",
	Address:     "zltc_ZwptHk17UU4wojKDwywJ3hfB9ihvUhjAq",
	AbiString: `[
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "_address",
					"type": "address"
				},
				{
					"internalType": "string",
					"name": "_filePath",
					"type": "string"
				},
				{
					"internalType": "address",
					"name": "_callSaintAddress",
					"type": "address"
				},
				{
					"internalType": "string",
					"name": "_storageAddress",
					"type": "string"
				},
				{
					"internalType": "int64",
					"name": "_needStorageSize",
					"type": "int64"
				},
				{
					"internalType": "string",
					"name": "_cid",
					"type": "string"
				}
			],
			"name": "UploadFile",
			"outputs": [
				{
					"internalType": "address",
					"name": "",
					"type": "address"
				}
			],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "_address",
					"type": "address"
				},
				{
					"internalType": "address",
					"name": "_permAddress",
					"type": "address"
				},
				{
					"internalType": "int64",
					"name": "_totalStorageSize",
					"type": "int64"
				}
			],
			"name": "UpdatePermission",
			"outputs": [
				{
					"internalType": "address",
					"name": "",
					"type": "address"
				}
			],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "_address",
					"type": "address"
				},
				{
					"internalType": "string",
					"name": "_cid",
					"type": "string"
				}
			],
			"name": "DownloadFile",
			"outputs": [
				{
					"internalType": "address",
					"name": "",
					"type": "address"
				}
			],
			"stateMutability": "nonpayable",
			"type": "function"
		}
	]`,
}
