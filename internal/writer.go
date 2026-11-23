package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func WriteFiles(wd string, sampleNum int, sample Sample) error {
	dirName := fmt.Sprintf("sample-%d", sampleNum)
	dirPath := filepath.Join(wd, "samples", dirName)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("error: failed to create dir %s (%v)", dirPath, err)
	}
	if err := os.WriteFile(fmt.Sprintf("%s/in.txt", dirPath), []byte(sample.Input), 0644); err != nil {
		return fmt.Errorf("error: failed to write in.txt (%v)", err)
	}
	if err := os.WriteFile(fmt.Sprintf("%s/out.txt", dirPath), []byte(sample.Output), 0644); err != nil {
		return fmt.Errorf("error: failed to write out.txt (%v)", err)
	}
	return nil
}
