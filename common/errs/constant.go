package errs

const (
	ErrInvalidAddressCode = 5001
	ErrMarshallStructCode = 5002
)

var (
	ErrInvalidAddressFormat = NewError(ErrInvalidAddressCode, "The address format is invalid", "地址格式不合法")
	ErrMarshallStruct       = NewError(ErrMarshallStructCode, "Serialization structure error", "序列化结构体错误")
)
