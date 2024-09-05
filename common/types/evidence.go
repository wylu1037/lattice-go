package types

import (
	"bytes"
	"encoding/json"
	"github.com/buger/jsonparser"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strconv"
	"strings"
)

// Evidences 留痕信息
type Evidences struct {
	Total uint64                 `json:"total"`
	Data  map[string]interface{} `json:"data"`
}

type ErrorEvidence struct {
	Lvl       int           `json:"lvl"`
	Timestamp uint64        `json:"timestamp"`
	Call      string        `json:"call"`
	Msg       string        `json:"msg"`
	Ctx       []interface{} `json:"ctx"`
}

type TBlockEvidence struct {
	Number     *big.Int       `json:"Number"`
	Hash       common.Hash    `json:"hash"`
	Owner      common.Address `json:"owner"`
	Timestamp  uint64         `json:"timestamp"`
	TBlockType uint8          `json:"tblockType"`
}

type DBlockEvidence struct {
	Number     *big.Int       `json:"Number"`
	Coinbase   common.Address `json:"miner"`
	ParentHash common.Hash    `json:"parentHash"`
	Timestamp  uint64         `json:"timestamp"`
	Hash       common.Hash    `json:"hash"`
}

type VoteEvidence struct {
	Signer     string `json:"signer"`
	Proposal   string `json:"proposal"`
	SwitchAddr string `json:"switchAddr"`
	VoteType   byte   `json:"voteType"`
	Timestamp  uint64 `json:"timestamp"`
}

type SignEvidence struct {
	From      common.Address `json:"from"`
	Owner     common.Address `json:"owner"`
	Hash      common.Hash    `json:"hash"`
	Number    *big.Int       `json:"number"`
	Sign      []byte         `json:"sign"`
	Timestamp uint64         `json:"timestamp"`
}

type ReceiptEvidence struct {
	Success         bool        `json:"success"`
	TblockHash      common.Hash `json:"tblockHash"`
	ContractAddress string      `json:"contractAddress"`
	ContractRet     string      `json:"contractRet"`
	Timestamp       uint64      `json:"timestamp"`
}

type DeployCallEvidence struct {
	Status          uint8          `json:"status"`
	TBHash          common.Hash    `json:"tblockHash"`
	ContractAddress common.Address `json:"contractAddress,omitempty"`
	Ret             []byte         `json:"ret,omitempty"`
	ContractType    uint8          `json:"contractType"`
	Timestamp       uint64         `json:"timestamp"`
}

type AccountEvidence struct {
	Address     string `json:"address"`
	AccountType string `json:"accountType"`
	Success     bool   `json:"success"`
	Timestamp   uint64 `json:"timestamp"`
}

func (d *Evidences) Unmarshall(data []byte) error {
	if err := jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		if bytes.Equal(key, []byte("total")) {
			total, err := strconv.ParseUint(string(value), 10, 64)
			if err != nil {
				return err
			}
			d.Total = total
		} else if bytes.Equal(key, []byte("data")) {
			d.Data = make(map[string]interface{})
			if err := jsonparser.ObjectEach(value, func(key []byte, iterVal []byte, dataType jsonparser.ValueType, offset int) error {
				iterKey := string(key)
				splitKey := strings.Split(iterKey, "_")
				evidenceLevel := EvidenceLevel(splitKey[1])
				evidenceType := EvidenceType(splitKey[2])
				if evidenceLevel == EvidenceLevelERROR || evidenceLevel == EvidenceLevelCRITICAL {
					var evidence ErrorEvidence
					if err := json.Unmarshal(iterVal, &evidence); err != nil {
						return err
					}
					d.Data[iterKey] = evidence
				} else {
					switch evidenceType {
					case EvidenceTypeTBLOCK, EvidenceTypeONCHAIN, EvidenceTypeEXECUTE, EvidenceTypeUPGRADE:
						var evidence TBlockEvidence
						if err := json.Unmarshal(iterVal, &evidence); err != nil {
							return err
						}
						d.Data[iterKey] = evidence
					case EvidenceTypeDBLOCK:
						var evidence DBlockEvidence
						if err := json.Unmarshal(iterVal, &evidence); err != nil {
							return err
						}
						d.Data[iterKey] = evidence
					case EvidenceTypeVOTING:
						var evidence VoteEvidence
						if err := json.Unmarshal(iterVal, &evidence); err != nil {
							return err
						}
						d.Data[iterKey] = evidence
					case EvidenceTypeSIGN:
						var evidence SignEvidence
						if err := json.Unmarshal(iterVal, &evidence); err != nil {
							return err
						}
						d.Data[iterKey] = evidence
					case EvidenceTypePRECALL:
						var evidence ReceiptEvidence
						if err := json.Unmarshal(iterVal, &evidence); err != nil {
							return err
						}
						d.Data[iterKey] = evidence
					case EvidenceTypeDEPLOY, EvidenceTypeCALL, EvidenceTypeUPDATE, EvidenceTypeREVOKE, EvidenceTypeFREEZE, EvidenceTypeUNFREEZE:
						var evidence DeployCallEvidence
						if err := json.Unmarshal(iterVal, &evidence); err != nil {
							return err
						}
						d.Data[iterKey] = evidence
					case EvidenceTypeADDED, EvidenceTypeDELETED, EvidenceTypeUNLOCKED, EvidenceTypeLOCKED:
						var evidence AccountEvidence
						if err := json.Unmarshal(iterVal, &evidence); err != nil {
							return err
						}
						d.Data[iterKey] = evidence
					default:
					}
				}
				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
