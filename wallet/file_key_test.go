package wallet

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"lattice-go/crypto"
	"testing"
)

func TestGenerateFileKey(t *testing.T) {
	privateKey := "0x23d5b2a2eb0a9c8b86d62cbc3955cfd1fb26ec576ecc379f402d0f5d2b27a7bb"
	passphrase := "Root1234"
	fileKey, err := GenerateFileKey(privateKey, passphrase, crypto.Sm2p256v1)
	assert.Nil(t, err)
	bytes, err := json.Marshal(fileKey)
	assert.Nil(t, err)
	fmt.Println(string(bytes))
}

func TestGenCipher(t *testing.T) {
	privateKey := "0x23d5b2a2eb0a9c8b86d62cbc3955cfd1fb26ec576ecc379f402d0f5d2b27a7bb"
	cipher, err := GenCipher(privateKey, "Root1234", crypto.Sm2p256v1)
	assert.Nil(t, err)
	bytes, err := json.Marshal(cipher)
	assert.Nil(t, err)
	fmt.Println(string(bytes))
}
