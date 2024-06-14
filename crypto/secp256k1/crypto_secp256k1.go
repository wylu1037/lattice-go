package secp256k1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"io"
)

func New() *Api {
	return &Api{}
}

type Api struct {
}

func (i *Api) GenerateKeyPair() (*ecdsa.PrivateKey, error) {
	sk, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return sk, nil
}

func (i *Api) SKToBytes(sk *ecdsa.PrivateKey) ([]byte, error) {
	return nil, nil
}

// SKToHexString 将私钥转为hex string
func (i *Api) SKToHexString(sk *ecdsa.PrivateKey) (string, error) {
	return "", nil
}

func (i *Api) HexToSK(skHex string) (*ecdsa.PrivateKey, error) {
	return nil, nil
}

// PKToBytes 将公钥转为[]byte
func (i *Api) PKToBytes(pk *ecdsa.PublicKey) ([]byte, error) {
	return nil, nil
}

// PKToHexString 将公钥转为hex string
func (i *Api) PKToHexString(pk *ecdsa.PublicKey) (string, error) {
	return "", nil
}

func (i *Api) HexToPK(skHex string) (*ecdsa.PublicKey, error) {
	return nil, nil
}

func (i *Api) PKToAddress(pk *ecdsa.PublicKey) (common.Address, error) {
	return common.Address{}, nil
}

func (i *Api) Sign(hash []byte, sk *ecdsa.PrivateKey) (signature []byte, err error) {
	return nil, nil
}

// SignatureToPK 从签名恢复公钥
func (i *Api) SignatureToPK(hash, signature []byte) (*ecdsa.PublicKey, error) {
	return nil, nil
}

// Verify 验证签名
func (i *Api) Verify(hash []byte, signature []byte, pk *ecdsa.PublicKey) bool {
	return true
}

// CompressPK 压缩公钥
func (i *Api) CompressPK(pk *ecdsa.PublicKey) []byte {
	return nil
}

// DecompressPK 解压缩公钥
func (i *Api) DecompressPK(pk []byte) (*ecdsa.PublicKey, error) {
	return nil, nil
}

// GetCurve 获取椭圆曲线
func (i *Api) GetCurve() elliptic.Curve {
	return secp256k1.S256()
}

func (i *Api) EncodeHash(func(io.Writer)) common.Hash {
	return common.Hash{}
}
