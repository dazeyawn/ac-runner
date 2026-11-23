package internal

import (
	"fmt"
	"os"
)

func CleanupDir(targetDir string) error {
	if err := os.RemoveAll(targetDir); err != nil {
		return fmt.Errorf("error: failed to remove samples dir (%v)", err)
	}

	return nil
}
