package secp256k1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	//"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"io"
	"math/big"
)

var (
	secp256k1N, _  = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN = new(big.Int).Div(secp256k1N, big.NewInt(2))
)

func New() *NistApi {
	return &NistApi{}
}

type NistApi struct {
}

func (i *NistApi) GenerateKeyPair() (*ecdsa.PrivateKey, error) {
	/*sk, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return sk, nil*/
	return nil, nil
}

func (i *NistApi) SKToBytes(sk *ecdsa.PrivateKey) ([]byte, error) {
	if sk == nil {
		return nil, errors.New("sk is nil")
	}

	// calculate the length of sk
	length := sk.Params().BitSize / 8
	// check the actual length of sk
	if sk.D.BitLen()/8 > length {
		return sk.D.Bytes(), errors.New("sk is too big")
	}

	bytes := make([]byte, length)
	// padding zero on the top of arr
	copy(bytes[length-len(sk.D.Bytes()):], sk.D.Bytes())
	return bytes, nil
}

// SKToHexString 将私钥转为hex string
func (i *NistApi) SKToHexString(sk *ecdsa.PrivateKey) (string, error) {
	bytes, err := i.SKToBytes(sk)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(bytes), nil
}

func (i *NistApi) HexToSK(skHex string) (*ecdsa.PrivateKey, error) {
	bytes, err := hexutil.Decode(skHex)
	if err != nil {
		return nil, err
	}

	return i.bytesToSK(bytes, true)
}

// PKToBytes 将公钥转为[]byte
func (i *NistApi) PKToBytes(pk *ecdsa.PublicKey) ([]byte, error) {
	if pk == nil || pk.X == nil || pk.Y == nil {
		return nil, errors.New("pk is invalid")
	}

	return elliptic.Marshal(i.GetCurve(), pk.X, pk.Y), nil
}

// PKToHexString 将公钥转为hex string
func (i *NistApi) PKToHexString(pk *ecdsa.PublicKey) (string, error) {
	bytes, err := i.PKToBytes(pk)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(bytes), nil
}

func (i *NistApi) HexToPK(pkHex string) (*ecdsa.PublicKey, error) {
	bytes, err := hexutil.Decode(pkHex)
	if err != nil {
		return nil, err
	}

	return i.BytesToPK(bytes)
}

func (i *NistApi) BytesToPK(pkBytes []byte) (*ecdsa.PublicKey, error) {
	x, y := elliptic.Unmarshal(i.GetCurve(), pkBytes)
	if x == nil {
		return nil, fmt.Errorf("invalid public key")
	}

	return &ecdsa.PublicKey{
		Curve: i.GetCurve(),
		X:     x,
		Y:     y,
	}, nil
}

func (i *NistApi) BytesToSK(sk []byte) (*ecdsa.PrivateKey, error) {
	return i.bytesToSK(sk, true)
}

func (i *NistApi) PKToAddress(pk *ecdsa.PublicKey) (common.Address, error) {
	bytes, err := i.PKToBytes(pk)
	if err != nil {
		return common.Address{}, err
	}
	return common.BytesToAddress(bytes), nil
}

func (i *NistApi) Sign(hash []byte, sk *ecdsa.PrivateKey) (signature []byte, err error) {
	skBytes, err := i.SKToBytes(sk)
	if err != nil {
		return nil, err
	}

	fmt.Println(skBytes)

	// return secp256k1.Sign(hash, skBytes)

	return nil, nil
}

// SignatureToPK 从签名恢复公钥
func (i *NistApi) SignatureToPK(hash, signature []byte) (*ecdsa.PublicKey, error) {
	if len(signature) == 97 {
		signature = signature[:65]
	}
	/*pkBytes, err := secp256k1.RecoverPubkey(hash, signature)
	if err != nil {
		return nil, err
	}
	return i.BytesToPK(pkBytes)*/
	return nil, nil
}

// Verify 验证签名
func (i *NistApi) Verify(hash []byte, signature []byte, pk *ecdsa.PublicKey) bool {
	pkBytes, err := i.PKToBytes(pk)
	if err != nil {
		return false
	}
	if len(signature) == 65 {
		signature = signature[:64]
	}
	fmt.Println(pkBytes)
	//return secp256k1.VerifySignature(pkBytes, hash, signature)
	return false
}

// CompressPK 压缩公钥
func (i *NistApi) CompressPK(pk *ecdsa.PublicKey) []byte {
	//return secp256k1.CompressPubkey(pk.X, pk.Y)
	return nil
}

// DecompressPK 解压缩公钥
func (i *NistApi) DecompressPK(pk []byte) (*ecdsa.PublicKey, error) {
	// x, y := secp256k1.DecompressPubkey(pk)
	x := big.NewInt(0)
	y := big.NewInt(0)
	if x == nil {
		return nil, errors.New("invalid public key")
	}
	return &ecdsa.PublicKey{
		Curve: i.GetCurve(),
		X:     x,
		Y:     y,
	}, nil
}

// GetCurve 获取椭圆曲线
func (i *NistApi) GetCurve() elliptic.Curve {
	//return secp256k1.S256()
	return i.GetCurve()
}

func (i *NistApi) EncodeHash(encodeFunc func(io.Writer)) (h common.Hash) {
	hash := sha256.New()
	encodeFunc(hash)
	hash.Sum(h[:0])
	return h
}

func (i *NistApi) bytesToSK(skBytes []byte, strict bool) (*ecdsa.PrivateKey, error) {
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

func (i *NistApi) Hash(data ...[]byte) (h common.Hash) {
	hash256 := sha256.New()
	for _, b := range data {
		hash256.Write(b)
	}
	hash256.Sum(h[:0])
	return h
}

func (i *NistApi) Encrypt(data []byte, pk string) ([]byte, error) {
	return nil, nil
}

func (i *NistApi) Decrypt(cipher []byte, sk string) ([]byte, error) {
	return nil, nil
}
