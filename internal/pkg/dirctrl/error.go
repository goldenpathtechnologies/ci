package dirctrl

const (
	DirUnexpectedError = iota + 1
	DirUnprivilegedError
)

type DirectoryError struct {
	Err       error
	ErrorCode int
}

func (d *DirectoryError) Error() string {
	return d.Err.Error()
}
