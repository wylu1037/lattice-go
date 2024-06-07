package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
)

type CryptographyApi interface {
	// GenerateKeyPair 生成密钥对
	GenerateKeyPair() (*ecdsa.PrivateKey, error)
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
