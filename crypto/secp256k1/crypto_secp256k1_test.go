package secp256k1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSecp256k1Api_GenerateKeyPair(t *testing.T) {
	sk, err := New().GenerateKeyPair()
	if err != nil {
		t.Error(err)
	}
	t.Log(sk)
}

func TestApi_SKToHexString(t *testing.T) {
	c := New()
	sk, err := c.GenerateKeyPair()
	assert.Nil(t, err)
	skHex, err := c.SKToHexString(sk)
	assert.Nil(t, err)
	t.Log(skHex)
}

func TestApi_Sign(t *testing.T) {
	c := New()
	sk, err := c.GenerateKeyPair()
	assert.Nil(t, err)
	hash := []byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}
	signature, err := c.Sign(hash, sk)
	assert.Nil(t, err)
	passed := c.Verify(hash, signature, &sk.PublicKey)
	assert.True(t, passed)
}
