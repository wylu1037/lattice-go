package errs

import "fmt"

const (
	// english
	en = "en"
	// chinese
	zh = "zh-Hans"
)

var local = zh

// NewError create a custom error
// Parameters:
//   - code int: error code
//   - enMsg string: english error message
//   - zhMsg string: chinese error message
//
// Returns:
//   - error
func NewError(code int, enMsg, zhMsg string) error {
	return &Error{
		Code: code,
		Message: map[string]string{
			en: enMsg,
			zh: zhMsg,
		},
	}
}

type Error struct {
	Code    int               `json:"code"`
	Message map[string]string `json:"message"`
}

func (e *Error) Error() string {
	if msg, ok := e.Message[local]; ok {
		return fmt.Sprintf("%d:%s", e.Code, msg)
	}
	return fmt.Sprintf("%d:%s", e.Code, e.Message[en])
}
