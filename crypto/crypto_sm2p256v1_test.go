package crypto

import "testing"

func TestSm2p256v1Api_GenerateKeyPair(t *testing.T) {
	sk, err := NewSm2p256v1Api().GenerateKeyPair()
	if err != nil {
		t.Error(err)
	}
	t.Log(sk)
	t.Log(sk.PublicKey)
}
