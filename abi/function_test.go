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

func TestLatticeFunction_EncodeConstructor(t *testing.T) {
	abi := `[{"inputs":[{"internalType":"string","name":"_name","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"get","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"index","type":"uint256"}],"name":"getPeople","outputs":[{"components":[{"internalType":"string","name":"name","type":"string"},{"internalType":"uint256","name":"age","type":"uint256"}],"internalType":"struct HelloWorld.Person","name":"","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getPeoples","outputs":[{"components":[{"internalType":"string","name":"name","type":"string"},{"internalType":"uint256","name":"age","type":"uint256"}],"internalType":"struct HelloWorld.Person[]","name":"","type":"tuple[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"people","outputs":[{"internalType":"string","name":"name","type":"string"},{"internalType":"uint256","name":"age","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"_name","type":"string"}],"name":"set","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"components":[{"internalType":"string","name":"name","type":"string"},{"internalType":"uint256","name":"age","type":"uint256"}],"internalType":"struct HelloWorld.Person","name":"_person","type":"tuple"}],"name":"setPeople","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	data, _ := NewAbi(abi).GetConstructor("jack").Encode()
	expectData := "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000046a61636b00000000000000000000000000000000000000000000000000000000"
	assert.Equal(t, expectData, data)
}
