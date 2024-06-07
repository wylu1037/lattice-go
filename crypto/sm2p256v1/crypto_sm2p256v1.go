package sm2p256v1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/tjfoc/gmsm/sm2"
	"lattice-go/common/types"
	"lattice-go/crypto"
)

func New() crypto.CryptographyApi {
	return &sm2p256v1Api{}
}

type sm2p256v1Api struct {
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

func (i *sm2p256v1Api) SKToBytes(sk *ecdsa.PrivateKey) ([]byte, error) {
	if sk == nil {
		return nil, errors.New("sk is nil")
	}
	length := sk.Params().BitSize / 8
	if sk.D.BitLen()/8 > length {
		return sk.D.Bytes(), errors.New("sk is too big")
	}

	bytes := make([]byte, length)
	// padding zero on the top of arr
	copy(bytes[len(bytes)-len(sk.D.Bytes()):], sk.D.Bytes())
	return bytes, nil
}

// SKToHexString 将私钥转为hex string
func (i *sm2p256v1Api) SKToHexString(sk *ecdsa.PrivateKey) (string, error) {
	bytes, err := i.SKToBytes(sk)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("0x%s", hex.EncodeToString(bytes)), nil
}

// PKToBytes 将公钥转为[]byte
func (i *sm2p256v1Api) PKToBytes(pk *ecdsa.PublicKey) ([]byte, error) {
	if pk == nil || pk.X == nil || pk.Y == nil {
		return nil, errors.New("pk is invalid")
	}

	return elliptic.Marshal(sm2.P256Sm2(), pk.X, pk.Y), nil
}

// PKToHexString 将公钥转为hex string
func (i *sm2p256v1Api) PKToHexString(pk *ecdsa.PublicKey) (string, error) {
	bytes, err := i.PKToBytes(pk)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("0x%s", hex.EncodeToString(bytes)), nil
}

func (i *sm2p256v1Api) PKToAddress(pk *ecdsa.PublicKey) (types.Address, error) {
	bytes, err := i.PKToBytes(pk)
	if err != nil {
		return types.Address{}, err
	}
	fmt.Println(bytes)
	return types.Address{}, nil
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
	return sm2.P256Sm2()
}
