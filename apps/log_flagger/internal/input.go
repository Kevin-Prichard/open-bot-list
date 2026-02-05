package internal

import (
	"encoding/csv"
	config2 "git.oxl.at/open-bot-list/pkg/flagger/config"
	"os"
)

// OpenInputFile opens a CSV file and returns a pointer to the os.File and a CSV reader.
func OpenInputFile(filePath string) (*os.File, *csv.Reader, []string, error) { // Updated function signature
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, nil, err
	}

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 0

	header, err := reader.Read() // Capture the header
	if err != nil {
		file.Close()
		return nil, nil, nil, err
	}

	return file, reader, header, nil
}

// ParseRecord takes a CSV record (slice of strings) and maps it to a LogEntry struct.
func ParseRecord(record []string) config2.LogEntry {
	return config2.LogEntry{
		ClientIP:       record[config2.CSV_FIELD_CLIENT_IP],
		FingerprintJA4: record[config2.CSV_FIELD_FP_JA4],
		UserAgent:      record[config2.CSV_FIELD_UA],
		AllFields:      record,
	}
}
