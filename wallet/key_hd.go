package wallet

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// EntropySize the length of entropy(bip39)
//   - EntropySize128 generate 12 words
//   - EntropySize160 generate 15 words
//   - EntropySize192 generate 18 words
//   - EntropySize256 generate 24 words
type EntropySize int

const (
	EntropySize128 EntropySize = 128
	EntropySize160 EntropySize = 160
	EntropySize192 EntropySize = 192
	EntropySize256 EntropySize = 256
)

// GenerateEntropy 1
func GenerateEntropy(entropySize EntropySize) ([]byte, error) {
	entropy, err := bip39.NewEntropy(int(entropySize))
	if err != nil {
		return nil, err
	}
	return entropy, nil
}

// GenerateSeed Generate a Bip32 HD wallet for the mnemonic and a user supplied password
func GenerateSeed(mnemonic, passphrase string) []byte {
	// check
	return bip39.NewSeed(mnemonic, passphrase)
}

// NewMasterKey 生成主账户密钥
func NewMasterKey(seed []byte) string {
	masterKey, _ := bip32.NewMasterKey(seed)
	publicKey := masterKey.PublicKey()

	// 以太坊的币种类型是60
	// FirstHardenedChild = uint32(0x80000000) 是一个常量
	// 以路径（path: "m/44'/60'/0'/0/0"）为例
	key, _ := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)  // 强化派生 对应 purpose'
	key, _ = key.NewChildKey(bip32.FirstHardenedChild + uint32(60)) // 强化派生 对应 coin_type'
	key, _ = key.NewChildKey(bip32.FirstHardenedChild + uint32(0))  // 强化派生 对应 account'
	key, _ = key.NewChildKey(uint32(0))                             // 常规派生 对应 change
	key, _ = key.NewChildKey(uint32(0))                             // 常规派生 对应 address_index

	// 生成地址
	pubKey, _ := crypto.DecompressPubkey(key.PublicKey().Key)
	address := crypto.PubkeyToAddress(*pubKey).Hex()

	// Display mnemonic and keys
	fmt.Println("Master private key: ", masterKey)
	fmt.Println("Master public key: ", publicKey)
	fmt.Println("Master address: ", address)
	return ""
}
