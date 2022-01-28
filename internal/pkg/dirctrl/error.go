package dirctrl

const (
	DirUnexpectedError = iota + 1
	DirUnprivilegedError
)

// DirectoryError represents an error that occurs while accessing the filesystem.
type DirectoryError struct {
	Err       error
	ErrorCode int
}

// Error returns the message for this error.
func (d *DirectoryError) Error() string {
	return d.Err.Error()
}
