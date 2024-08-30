package errs

import (
	"fmt"
	"testing"
)

func TestError_Error(t *testing.T) {
	fmt.Println(ErrInvalidAddressFormat.Error())
}
