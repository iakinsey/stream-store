package errs

// ResponseError ...
type ResponseError struct {
	Code byte
	Err  error
}

func (e *ResponseError) Error() string { return e.Err.Error() }
