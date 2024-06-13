package sm2p256v1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm3"
	"io"
	"lattice-go/common/constant"
	"lattice-go/common/convert"
	"lattice-go/crypto"
	cryptoConstant "lattice-go/crypto/constant"
	"math/big"
)

var (
	secp256k1N, _  = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN = new(big.Int).Div(secp256k1N, big.NewInt(2))
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

func (i *sm2p256v1Api) HexToSK(skHex string) (*ecdsa.PrivateKey, error) {
	bytes, err := hexutil.Decode(skHex)
	if err != nil {
		return nil, err
	}

	return i.hexStringToSK(bytes, true)
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

func (i *sm2p256v1Api) HexToPK(skHex string) (*ecdsa.PublicKey, error) {
	bytes, err := hexutil.Decode(skHex)
	if err != nil {
		return nil, err
	}
	x, y := elliptic.Unmarshal(i.GetCurve(), bytes)
	if x == nil {
		return nil, fmt.Errorf("invalid public key")
	}

	return &ecdsa.PublicKey{
		Curve: i.GetCurve(),
		X:     x,
		Y:     y,
	}, nil
}

// PKToAddress 将公钥(取后20位字节)转为地址
// Parameters
//   - pk *ecdsa.PublicKey: 公钥
//
// Returns
//   - common.Address: 地址
//   - error
func (i *sm2p256v1Api) PKToAddress(pk *ecdsa.PublicKey) (common.Address, error) {
	bytes, err := i.PKToBytes(pk)
	if err != nil {
		return common.Address{}, err
	}
	return common.BytesToAddress(bytes), nil
}

// Sign 签名
func (i *sm2p256v1Api) Sign(hash []byte, sk *ecdsa.PrivateKey) (signature []byte, err error) {
	if len(hash) != constant.HashLength {
		return nil, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
	}

	privateKey := convert.EcdsaSKToSm2SK(sk)
	// use default uid: []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38}
	r, s, err := sm2.Sm2Sign(privateKey, hash, nil, rand.Reader)
	if err != nil {
		return nil, err
	}

	signature = make([]byte, 65)
	copy(signature[32-len(r.Bytes()):], r.Bytes())
	copy(signature[64-len(s.Bytes()):], s.Bytes())
	signature[64] = cryptoConstant.Sm2p256v1SignatureRemark

	if len(signature) != 65 {
		return nil, errors.New(fmt.Sprintf("sig length is wrong !!! sig length is %d ", len(signature)))
	}

	// calculate E
	digest, err := privateKey.PublicKey.Sm3Digest(hash, nil)
	if err != nil {
		return nil, err
	}
	e := new(big.Int).SetBytes(digest)

	var pad [32]byte
	buffer := e.Bytes()
	copy(pad[32-len(buffer):], buffer)
	signature = append(signature, pad[:]...)
	return signature, nil
}

// SignatureToPK 从签名恢复公钥
func (i *sm2p256v1Api) SignatureToPK(hash, signature []byte) (*ecdsa.PublicKey, error) {
	_ = new(big.Int).SetBytes(signature[65:])
	return nil, nil
}

// Verify 验证签名
func (i *sm2p256v1Api) Verify(hash []byte, signature []byte, pk *ecdsa.PublicKey) bool {
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:64])
	return sm2.Sm2Verify(convert.EcdsaPKToSm2PK(pk), hash, nil, r, s)
}

// CompressPK 压缩公钥
func (i *sm2p256v1Api) CompressPK(pk *ecdsa.PublicKey) []byte {
	if pk == nil || pk.X == nil || pk.Y == nil {
		return nil
	}
	return sm2.Compress(convert.EcdsaPKToSm2PK(pk))
}

// DecompressPK 解压缩公钥
func (i *sm2p256v1Api) DecompressPK(pk []byte) (*ecdsa.PublicKey, error) {
	if len(pk) != 33 {
		return nil, errors.New(fmt.Sprintf("DecompressPubKey length is wrong !,lenth is %d", len(pk)))
	}
	return convert.Sm2PKToEcdsaPK(sm2.Decompress(pk)), nil
}

// GetCurve 获取椭圆曲线
func (i *sm2p256v1Api) GetCurve() elliptic.Curve {
	return sm2.P256Sm2()
}

func (i *sm2p256v1Api) hexStringToSK(skBytes []byte, strict bool) (*ecdsa.PrivateKey, error) {
	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = i.GetCurve()
	if strict && 8*len(skBytes) != privateKey.Params().BitSize {
		return nil, errors.New("skBytes length is wrong")
	}

	privateKey.D = new(big.Int).SetBytes(skBytes)

	// The priv.D must < N
	if privateKey.D.Cmp(secp256k1N) >= 0 {
		return nil, fmt.Errorf("invalid private key, >=N")
	}
	// The priv.D must not be zero or negative.
	if privateKey.D.Sign() <= 0 {
		return nil, fmt.Errorf("invalid private key, zero or negative")
	}

	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(skBytes)
	if privateKey.PublicKey.X == nil {
		return nil, errors.New("invalid private key")
	}
	return privateKey, nil
}

func (i *sm2p256v1Api) EncodeHash(encodeFunc func(io.Writer)) (h common.Hash) {
	hash := sm3.New()
	encodeFunc(hash)
	hash.Sum(h[:0])
	return h
}
