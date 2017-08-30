package utils

import "os"
import "path"

// EnsureDirectory takes in the full path of a file and ensures that
// that it's containing directory exists
func EnsureDirectory(filepath string) error {
	dir, _ := path.Split(filepath)
	if exists, err := PathExists(dir); err != nil {
		return err
	} else if !exists {
		err := os.MkdirAll(dir, 0700)
		return err
	}
	return nil
}

// PathExists checks if a file path exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
