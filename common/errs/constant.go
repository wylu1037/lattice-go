package errs

const (
	ErrAddressCode = 5001
)

var (
	ErrAddressFormat = NewError(ErrAddressCode, "address format is invalid", "地址格式不合法")
)
