package secp256k1

import "testing"

func TestSecp256k1Api_GenerateKeyPair(t *testing.T) {
	sk, err := New().GenerateKeyPair()
	if err != nil {
		t.Error(err)
	}
	t.Log(sk)
}
