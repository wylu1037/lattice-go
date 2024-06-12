package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZltcToAddress(t *testing.T) {
	zltcAddr := "zltc_dhdfbm9JEoyDvYoCDVsABiZj52TAo9Ei6"
	ethAddr := "0x9293c604c644BfAc34F498998cC3402F203d4D6B"
	addr, err := ZltcToAddress(zltcAddr)
	assert.Nil(t, err)
	assert.Equal(t, ethAddr, addr.Hex())
}
