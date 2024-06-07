package sm2p256v1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"lattice-go/crypto"
	"lattice-go/crypto/constant"

	"github.com/tjfoc/gmsm/sm2"
)

func New() crypto.CryptographyApi {
	return &sm2p256v1Api{}
}

type sm2p256v1Api struct {
	curve constant.Curve
}

// GenerateKeyPair 生成密钥对
func (i *sm2p256v1Api) GenerateKeyPair() (*ecdsa.PrivateKey, error) {
	sk, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	pk := sk.PublicKey

	// assemble public key
	publicKey := new(ecdsa.PublicKey)
	publicKey.Curve = sk.Curve
	publicKey.X = pk.X
	publicKey.Y = pk.Y

	// assemble private key
	privateKey := new(ecdsa.PrivateKey)
	privateKey.Curve = sk.Curve
	privateKey.D = sk.D
	privateKey.PublicKey = *publicKey

	return privateKey, nil
}

// Sign 签名
func (i *sm2p256v1Api) Sign(hash []byte, sk *ecdsa.PrivateKey) (signature []byte, err error) {
	return nil, nil
}

// SignatureToPK 从签名恢复公钥
func (i *sm2p256v1Api) SignatureToPK(hash, signature []byte) (*ecdsa.PublicKey, error) {
	return nil, nil
}

// Verify 验证签名
func (i *sm2p256v1Api) Verify(hash []byte, signature []byte, pk *ecdsa.PublicKey) error {
	return nil
}

// CompressPK 压缩公钥
func (i *sm2p256v1Api) CompressPK(pk *ecdsa.PublicKey) []byte {
	return nil
}

// DecompressPK 解压缩公钥
func (i *sm2p256v1Api) DecompressPK(pk []byte) (*ecdsa.PublicKey, error) {
	return nil, nil
}

// GetCurve 获取椭圆曲线
func (i *sm2p256v1Api) GetCurve() elliptic.Curve {
	return nil
}
