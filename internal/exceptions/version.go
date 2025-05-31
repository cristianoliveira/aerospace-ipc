package exceptions

import (
	"errors"
	"fmt"
)

var ErrVersion = errors.New("server version does not match the minimum client version")

type ErrVersionMismatch struct {
	MinorVersion   int
	MajorVersion   int
	CurrentVersion string
}

func NewErrVersionMismatch(major, minor int, currentVersion string) *ErrVersionMismatch {
	return &ErrVersionMismatch{
		MajorVersion:   major,
		MinorVersion:   minor,
		CurrentVersion: currentVersion,
	}
}

func (e *ErrVersionMismatch) Error() string {
	return fmt.Sprintf(
		"Server version %s does not match the minimum required version %d.%d.x",
		e.CurrentVersion,
		e.MajorVersion,
		e.MinorVersion,
	)
}

// Add Is to allow errors.Is(err, ErrMyCustom)
func (e *ErrVersionMismatch) Is(target error) bool {
	return target == ErrVersion
}
