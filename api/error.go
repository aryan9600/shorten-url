package api

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

var AlreadyPresentErrorString = "Path already exists."
var NotFoundErrorString = "Path not found."

func AlreadyPresentError(text string) error{
	return &errorString{text}
}

func NotFoundError(text string) error{
	return &errorString{text}
}
