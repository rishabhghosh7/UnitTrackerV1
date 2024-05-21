package fsutils

import (
	"os"
)

// FileExists checks if a file exists at the specified path.
// It returns true if the file exists, and false otherwise.
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
