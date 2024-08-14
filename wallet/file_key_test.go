package wallet

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/wylu1037/lattice-go/crypto"
	"testing"
)

func TestGenerateFileKey(t *testing.T) {
	privateKey := "0xbd7ea728f7e6240507b321cb4a937a8d34ecfd39c275dbacf31ddb4793691dcc"
	passphrase := "Aa123456"
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

func TestFileKey_Decrypt(t *testing.T) {
	fileKeyString := `{"uuid":"bb889ee6-5d1d-474e-9514-5bbf412a42ec","address":"zltc_iEUCcfMhVYy3zcpp8zLjoaTAeN6PZfMBL","cipher":{"aes":{"cipher":"aes-128-ctr","iv":"23b4ddcd8cfea7e37b3c69bbb600934f"},"kdf":{"kdf":"scrypt","kdfParams":{"DKLen":32,"n":262144,"p":1,"r":8,"salt":"87cf307be225ce2eaf255d602233852200195d838b5d98c4078ceb6235ec46e4"}},"cipherText":"672b3de4784fc0d17941ae257908672dd4984a43c616147366a42bc2e9ef2d8a","mac":"bd6ac051c41f4d0238464a66df004de357baf2f3f03ced8ccba0a497e14044bd"},"isGM":true}`
	passphrase := "Root1234"
	sk, err := NewFileKey(fileKeyString).Decrypt(passphrase)
	assert.Nil(t, err)
	skString, err := crypto.NewCrypto(crypto.Sm2p256v1).SKToHexString(sk)
	assert.Nil(t, err)
	assert.Equal(t, skString, "0x23d5b2a2eb0a9c8b86d62cbc3955cfd1fb26ec576ecc379f402d0f5d2b27a7bb")
}
