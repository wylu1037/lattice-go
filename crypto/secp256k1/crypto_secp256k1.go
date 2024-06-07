package secp256k1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"lattice-go/crypto"
	"lattice-go/crypto/constant"
)

func New() crypto.CryptographyApi {
	return &secp256k1Api{}
}

type secp256k1Api struct {
	curve constant.Curve
}

func (i *secp256k1Api) GenerateKeyPair() (*ecdsa.PrivateKey, error) {
	sk, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return sk, nil
}

func (i *secp256k1Api) Sign(hash []byte, sk *ecdsa.PrivateKey) (signature []byte, err error) {
	return nil, nil
}

// SignatureToPK 从签名恢复公钥
func (i *secp256k1Api) SignatureToPK(hash, signature []byte) (*ecdsa.PublicKey, error) {
	return nil, nil
}

// Verify 验证签名
func (i *secp256k1Api) Verify(hash []byte, signature []byte, pk *ecdsa.PublicKey) error {
	return nil
}

// CompressPK 压缩公钥
func (i *secp256k1Api) CompressPK(pk *ecdsa.PublicKey) []byte {
	return nil
}

// DecompressPK 解压缩公钥
func (i *secp256k1Api) DecompressPK(pk []byte) (*ecdsa.PublicKey, error) {
	return nil, nil
}

// GetCurve 获取椭圆曲线
func (i *secp256k1Api) GetCurve() elliptic.Curve {
	return nil
}
