package abi

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLatticeFunction_ConvertArguments(t *testing.T) {
	abiString := `[{"inputs":[],"name":"get","outputs":[{"internalType":"string[]","name":"","type":"string[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string[]","name":"arr","type":"string[]"}],"name":"set","outputs":[{"internalType":"string[]","name":"","type":"string[]"}],"stateMutability":"nonpayable","type":"function"}]`
	latticeAbi := NewAbi(abiString)
	fn, err := latticeAbi.GetLatticeFunction("set", []string{"1", "jack"})
	assert.Nil(t, err)
	code, err := fn.Encode()
	assert.Nil(t, err)
	t.Log(code)
}
