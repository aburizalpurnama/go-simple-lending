package custerror

var _ error = new(Error)

type Error struct {
	HttpStatusCode int
	Message        string
	Err            error
}

func New(httpCode int, msg string, err error) error {
	return &Error{
		HttpStatusCode: httpCode,
		Message:        msg,
		Err:            err,
	}
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Is(target error) bool {
	return e == target
}

func (e *Error) Unwrap() error {
	return e.Err
}
