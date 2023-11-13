package utilities

import (
	"errors"
	"os"
	"strings"
)

func EnsureDirectoryExists(directory string) error {
	_, err := os.Stat(directory)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = os.Mkdir(directory, os.ModeDir)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func BuildPath(basePath string, parts ...string) string {
	result := basePath

	for _, part := range parts {
		if !strings.HasSuffix(result, "/") {
			result += "/"
		}

		result += part
	}

	return result
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func DeleteFile(path string) (bool, error) {
	err := os.Remove(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
