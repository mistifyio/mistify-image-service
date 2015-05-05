package images

import "errors"

// ErrNotFound is an error for not finding an image
var ErrNotFound = errors.New("image not found")

// IsErrNotFound checks if the error is ErrNotFound
func IsErrNotFound(err error) bool {
	return err == ErrNotFound
}
