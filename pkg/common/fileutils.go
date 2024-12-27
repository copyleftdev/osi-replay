package common

import (
	"fmt"
	"os"
)

func EnsureFileWritable(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("unable to create file %s: %w", path, err)
		}
		f.Close()
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("unable to remove test file %s: %w", path, err)
		}
	} else {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("file %s not writable: %w", path, err)
		}
		f.Close()
	}
	return nil
}
