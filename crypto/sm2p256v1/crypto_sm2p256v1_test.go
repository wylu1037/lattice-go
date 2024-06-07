package sm2p256v1

import (
	"testing"
)

func TestSm2p256v1Api_GenerateKeyPair(t *testing.T) {
	sk, err := New().GenerateKeyPair()
	if err != nil {
		t.Error(err)
	}
	t.Log(sk)
	t.Log(sk.PublicKey)
}
