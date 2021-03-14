package errs

// HTTPError ...
type HTTPError struct {
	Code int
	Err  error
}

func (e *HTTPError) Error() string { return e.Err.Error() }
