package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
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
	DkLen uint32 `json:"DkLen"` // 生成的密钥长度，单位byte
	N     uint32 `json:"n"`     // CPU/内存成本因子，控制计算和内存的使用量。
	P     uint32 `json:"p"`     // 并行度因子，控制 scrypt 函数的并行度。
	R     uint32 `json:"r"`     // 块大小因子，影响内部工作状态和内存占用。
	Salt  string `json:"salt"`  // 盐值，在密钥派生过程中加入随机性。
}

func NewFileKey() *FileKey {
	return nil
}

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

func GenCipher(privateKey string, passphrase string, curve types.Curve) (*Cipher, error) {
	// generate salt
	salt, err := random(32)
	if err != nil {
		return nil, err
	}

	key, err := scryptKey([]byte(passphrase), salt, ScryptN)
	if err != nil {
		return nil, err
	}
	aesKey := key[:16]
	hashKey := key[16:32] // compact mac

	ivBytes, err := random(aes.BlockSize) //16 equals aes.BlockSize
	if err != nil {
		return nil, err
	}
	ciphertext, err := aesEncrypt(aesKey, ivBytes, hexutil.MustDecode(privateKey))
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

func (e *FileKey) Decrypt(passphrase string) {

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

func aesEncrypt(key, iv, secretKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)
	ciphertext := make([]byte, len(secretKey))
	stream.XORKeyStream(ciphertext, secretKey)
	return ciphertext, nil
}

func random(length int) ([]byte, error) {
	bytes := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return nil, fmt.Errorf("reading from crypto/rand failed: %s", err.Error())
	}
	return bytes, nil
}
