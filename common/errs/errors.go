package errs

const (
	// english
	en = "en"
	// chinese
	zh = "zh-Hans"
)

var local = zh

// NewError create a custom error
// Parameters:
//   - code(int): error code
//   - en(string): english error message
//   - zh(string): chinese error message
//
// Returns:
//   - error
func NewError(code int, en, zh string) error {
	return &Error{
		Code: code,
		Message: map[string]string{
			en: en,
			zh: zh,
		},
	}
}

type Error struct {
	Code    int               `json:"code"`
	Message map[string]string `json:"message"`
}

func (e *Error) Error() string {
	if msg, ok := e.Message[local]; ok {
		return msg
	}
	return e.Message[en]
}
