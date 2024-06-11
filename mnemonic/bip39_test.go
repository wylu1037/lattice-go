package mnemonic

import (
	"fmt"
	"testing"
)

func TestGenerateMnemonic(t *testing.T) {
	fmt.Println(GenerateMnemonic())
}

func TestNewMasterKey(t *testing.T) {
	NewMasterKey(GenerateSeed(GenerateMnemonic(), "Root1234"))
}
