package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"lattice-go/common/types"
)

type CryptographyApi interface {
	// GenerateKeyPair 生成密钥对
	GenerateKeyPair() (*ecdsa.PrivateKey, error)
	// SKToBytes 将私钥转为[]byte
	SKToBytes(sk *ecdsa.PrivateKey) ([]byte, error)
	// SKToHexString 将私钥转为hex string
	SKToHexString(sk *ecdsa.PrivateKey) (string, error)
	// PKToBytes 将公钥转为[]byte
	PKToBytes(pk *ecdsa.PublicKey) ([]byte, error)
	// PKToHexString 将公钥转为hex string
	PKToHexString(pk *ecdsa.PublicKey) (string, error)
	// PKToAddress 将公钥转为地址
	PKToAddress(pk *ecdsa.PublicKey) (types.Address, error)
	// Sign 签名
	Sign(hash []byte, sk *ecdsa.PrivateKey) (signature []byte, err error)
	// SignatureToPK 从签名恢复公钥
	SignatureToPK(hash, signature []byte) (*ecdsa.PublicKey, error)
	// Verify 验证签名
	Verify(hash []byte, signature []byte, pk *ecdsa.PublicKey) error
	// CompressPK 压缩公钥
	CompressPK(pk *ecdsa.PublicKey) []byte
	// DecompressPK 解压缩公钥
	DecompressPK(pk []byte) (*ecdsa.PublicKey, error)
	// GetCurve 获取椭圆曲线
	GetCurve() elliptic.Curve
}
