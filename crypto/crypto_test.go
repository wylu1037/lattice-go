package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCrypto(t *testing.T) {
	c := NewCrypto(Sm2p256v1)
	assert.NotNil(t, c)
}
