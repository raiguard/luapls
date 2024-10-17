package util

import "os"

// Ptr returns a pointer to the given value.
func Ptr[T any](value T) *T {
	return &value
}

// FileExists returns whether the given file exists on the filesystem.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
