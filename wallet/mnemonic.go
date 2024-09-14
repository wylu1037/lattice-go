package wallet

import "github.com/tyler-smith/go-bip39"

func GenerateMnemonic() string {
	entropy, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	bip39.GetWordList()
	return mnemonic
}
