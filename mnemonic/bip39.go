package mnemonic

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func GenerateMnemonic() string {
	entropy, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	bip39.GetWordList()
	return mnemonic
}

// GenerateSeed Generate a Bip32 HD wallet for the mnemonic and a user supplied password
func GenerateSeed(mnemonic, passphrase string) []byte {
	// check
	return bip39.NewSeed(mnemonic, passphrase)
}

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
