package util

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"git.oxl.at/open-bot-list/pkg/downloader/config"
)

// ParseCSVFile reads a CSV file from the runtime path, skips the header, and returns the data records.
func ParseCSVFile(filename string) ([][]string, error) {
	targetPath := filepath.Join(config.PATH_RUNTIME, strings.ReplaceAll(filename, "/", "_"))

	f, err := os.Open(targetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open manifest file %s: %w", targetPath, err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	_, err = reader.Read()
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read header from %s: %w", targetPath, err)
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV records from %s: %w", targetPath, err)
	}

	return records, nil
}
