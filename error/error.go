package error

type Error struct {
	Context string
	Err     error
}

func (e Error) Error() string {
	return e.Err.Error()
}
