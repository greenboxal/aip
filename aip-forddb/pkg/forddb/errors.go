package forddb

import "errors"

var ErrVersionMismatch = errors.New("version mismatch")
var ErrNotFound = errors.New("not found")

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}
