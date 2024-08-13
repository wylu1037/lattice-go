package wallet

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/uuid"
	"golang.org/x/crypto/scrypt"
	"io"
	"lattice-go/common/convert"
	"lattice-go/common/types"
	"lattice-go/crypto"
	"strings"
)

const (
	aes128Ctr    = "aes-128-ctr"
	kdfScrypt    = "scrypt"
	ScryptN      = 1 << 18
	ScryptP      = 1
	ScryptR      = 8
	ScryptKeyLen = 32
)

type FileKey struct {
	Uuid    string  `json:"uuid"`
	Address string  `json:"address"`
	Cipher  *Cipher `json:"cipher"`
	IsGM    bool    `json:"isGM"`
}

type Cipher struct {
	Aes        *Aes   `json:"aes"`
	Kdf        *Kdf   `json:"kdf"`
	CipherText string `json:"cipherText"`
	Mac        string `json:"mac"`
}

type Aes struct {
	Cipher string `json:"cipher"` // 密码算法：aes-128-ctr
	Iv     string `json:"iv"`     // 初始化向量：1ad693b4d8089da0492b9c8c49bc60d3
}

type Kdf struct {
	Kdf       string     `json:"kdf"` // scrypt, PBKDF2, bcrypt, HKDF
	KdfParams *KdfParams `json:"kdfParams"`
}

type KdfParams struct {
	DkLen uint32 `json:"DKLen"` // 生成的密钥长度，单位byte
	N     uint32 `json:"n"`     // CPU/内存成本因子，控制计算和内存的使用量。
	P     uint32 `json:"p"`     // 并行度因子，控制 scrypt 函数的并行度。
	R     uint32 `json:"r"`     // 块大小因子，影响内部工作状态和内存占用。
	Salt  string `json:"salt"`  // 盐值，在密钥派生过程中加入随机性。
}

// NewFileKey 通过FileKey的JSON字符串初始化FileKey
//
// Parameters:
//   - fileKeyJsonString string
//
// Returns:
//   - *FileKey
func NewFileKey(fileKeyJsonString string) *FileKey {
	var fileKey FileKey
	err := json.Unmarshal([]byte(fileKeyJsonString), &fileKey)
	if err != nil {
		return nil
	}
	return &fileKey
}

// GenerateFileKey 生成一个FileKey
//
// Parameters:
//   - privateKey string: 带0x前缀的16进制的私钥
//   - passphrase string: 身份密码
//   - curve types.Curve: 曲线类型，crypto.Sm2p256v1 or crypto.Secp256k1
//
// Returns:
//   - *FileKey
//   - error
func GenerateFileKey(privateKey, passphrase string, curve types.Curve) (*FileKey, error) {
	instance := crypto.NewCrypto(curve)
	secretKey, err := instance.HexToSK(privateKey)
	if err != nil {
		return nil, err
	}

	address, err := instance.PKToAddress(&secretKey.PublicKey)
	if err != nil {
		return nil, err
	}

	ciphertext, err := GenCipher(privateKey, passphrase, curve)
	if err != nil {
		return nil, err
	}

	return &FileKey{
		Uuid:    uuid.New().String(),
		Address: convert.AddressToZltc(address),
		Cipher:  ciphertext,
		IsGM:    curve == crypto.Sm2p256v1,
	}, nil
}

// GenCipher 生成私钥的密文
//
// Parameters:
//   - privateKey string: 带0x前缀的16进制的私钥
//   - passphrase string: 身份密码
//   - curve types.Curve: 曲线类型，crypto.Sm2p256v1 or crypto.Secp256k1
//
// Returns:
//   - *Cipher
//   - error
func GenCipher(privateKey, passphrase string, curve types.Curve) (*Cipher, error) {
	// generate salt
	salt, err := random(32)
	if err != nil {
		return nil, err
	}

	key, err := scryptKey([]byte(passphrase), salt, ScryptN)
	if err != nil {
		return nil, err
	}
	aesKey := key[:aes.BlockSize]
	hashKey := key[aes.BlockSize:32] // compact mac

	ivBytes, err := random(aes.BlockSize) //16 equals aes.BlockSize
	if err != nil {
		return nil, err
	}
	ciphertext, err := aesCtr(aesKey, ivBytes, hexutil.MustDecode(privateKey))
	if err != nil {
		return nil, err
	}

	mac := crypto.NewCrypto(curve).Hash(hashKey, ciphertext)

	return &Cipher{
		Aes: &Aes{
			Cipher: aes128Ctr,
			Iv:     hex.EncodeToString(ivBytes),
		},
		Kdf: &Kdf{
			Kdf: kdfScrypt,
			KdfParams: &KdfParams{
				DkLen: ScryptKeyLen,
				N:     ScryptN,
				P:     ScryptP,
				R:     ScryptR,
				Salt:  hex.EncodeToString(salt),
			},
		},
		CipherText: hex.EncodeToString(ciphertext),
		Mac:        strings.TrimPrefix(mac.Hex(), "0x"),
	}, nil
}

// Decrypt 解密FileKey获取私钥
//
// Parameters:
//   - passphrase string: 身份密码
//
// Returns:
//   - *ecdsa.PrivateKey: 私钥
//   - error
func (e *FileKey) Decrypt(passphrase string) (*ecdsa.PrivateKey, error) {
	salt, err := hex.DecodeString(e.Cipher.Kdf.KdfParams.Salt)
	if err != nil {
		return nil, err
	}
	key, err := scryptKey([]byte(passphrase), salt, ScryptN)
	if err != nil {
		return nil, err
	}

	aesKey := key[:aes.BlockSize]
	ciphertext, err := hex.DecodeString(e.Cipher.CipherText)
	if err != nil {
		return nil, err
	}

	var curve types.Curve
	if e.IsGM {
		curve = crypto.Sm2p256v1
	} else {
		curve = crypto.Secp256k1
	}
	hashKey := key[aes.BlockSize:32] // compact mac
	actualMac := crypto.NewCrypto(curve).Hash(hashKey, ciphertext)
	expectMac, err := hex.DecodeString(e.Cipher.Mac)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(actualMac.Bytes(), expectMac) {
		return nil, fmt.Errorf("根据密码无法解析出私钥，请检查密码")
	}

	iv, err := hex.DecodeString(e.Cipher.Aes.Iv)
	if err != nil {
		return nil, err
	}
	privateKey, err := aesCtr(aesKey, iv, ciphertext)
	if err != nil {
		return nil, err
	}

	return crypto.NewCrypto(curve).BytesToSK(privateKey)
}

// 使用KDF中的script(基于密码的密钥导出算法（Password-Based Key Derivation Function, KDF），其主要作用是通过消耗大量内存和计算资源来增强密码的安全性，防止暴力破解和专用硬件攻击)
//
// Parameters:
//   - passphrase, salt []byte, n int
//
// Returns:
//   - []byte
//   - error
func scryptKey(passphrase, salt []byte, n int) ([]byte, error) {
	return scrypt.Key(passphrase, salt, n, ScryptR, ScryptP, ScryptKeyLen)
}

func aesCtr(key, iv, secretKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)
	ciphertext := make([]byte, len(secretKey))
	stream.XORKeyStream(ciphertext, secretKey)
	return ciphertext, nil
}

// 生成指定长度的随机byte数组
//
// Parameters:
//   - length int
//
// Returns:
//   - []byte
//   - error
func random(length int) ([]byte, error) {
	bs := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, bs); err != nil {
		return nil, fmt.Errorf("reading from crypto/rand failed: %s", err.Error())
	}
	return bs, nil
}
