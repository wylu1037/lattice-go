package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"github.com/ethereum/go-ethereum/common"
	"io"
	"lattice-go/common/types"
	"lattice-go/crypto/secp256k1"
	"lattice-go/crypto/sm2p256v1"
)

var instance CryptographyApi

func NewCrypto(curve types.Curve) CryptographyApi {
	if instance != nil {
		return instance
	}

	switch curve {
	case Sm2p256v1:
		instance = sm2p256v1.New()
	case Secp256k1:
		instance = secp256k1.New()
	default:
		instance = sm2p256v1.New()
	}
	return instance
}

type CryptographyApi interface {
	// GenerateKeyPair 生成密钥对
	GenerateKeyPair() (*ecdsa.PrivateKey, error)
	// SKToBytes 将私钥转为[]byte
	SKToBytes(sk *ecdsa.PrivateKey) ([]byte, error)
	// SKToHexString 将私钥转为hex string
	SKToHexString(sk *ecdsa.PrivateKey) (string, error)
	// HexToSK 将hex字符串的私钥转为私钥
	HexToSK(skHex string) (*ecdsa.PrivateKey, error)
	// PKToBytes 将公钥转为[]byte
	PKToBytes(pk *ecdsa.PublicKey) ([]byte, error)
	// PKToHexString 将公钥转为hex string
	PKToHexString(pk *ecdsa.PublicKey) (string, error)
	// HexToPK 将hex字符串的公钥转为公钥
	HexToPK(pkHex string) (*ecdsa.PublicKey, error)
	// BytesToPK 将[]byte转为公钥
	BytesToPK(pk []byte) (*ecdsa.PublicKey, error)
	BytesToSK(sk []byte) (*ecdsa.PrivateKey, error)
	// PKToAddress 将公钥转为地址
	PKToAddress(pk *ecdsa.PublicKey) (common.Address, error)
	// Sign 签名
	Sign(hash []byte, sk *ecdsa.PrivateKey) (signature []byte, err error)
	// SignatureToPK 从签名恢复公钥
	SignatureToPK(hash, signature []byte) (*ecdsa.PublicKey, error)
	// Verify 验证签名
	Verify(hash []byte, signature []byte, pk *ecdsa.PublicKey) bool
	// CompressPK 压缩公钥
	CompressPK(pk *ecdsa.PublicKey) []byte
	// DecompressPK 解压缩公钥
	DecompressPK(pk []byte) (*ecdsa.PublicKey, error)
	// GetCurve 获取椭圆曲线
	GetCurve() elliptic.Curve
	// EncodeHash Rlp encode and hash
	EncodeHash(encodeFunc func(io.Writer)) common.Hash
	// Hash 哈希
	Hash(data ...[]byte) common.Hash
	// Encrypt 非对称加密
	Encrypt(data []byte, pk string) ([]byte, error)
	// Decrypt 非对称解密
	Decrypt(cipher []byte, sk string) ([]byte, error)
}
