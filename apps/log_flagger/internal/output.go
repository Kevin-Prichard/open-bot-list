package internal

import (
	"encoding/csv"
	"os"
	"slices"
	"strings"

	"git.oxl.at/open-bot-list/log_flagger/internal/config"
)

// SerializeFlags converts the Flags struct into a pipe-separated string of flag names.
func SerializeFlags(flags []string) string {
	out := []string{}
	for _, flag := range flags {
		if strings.HasPrefix(flag, "src_net_") || strings.HasPrefix(flag, "http_user_agent_") {
			continue
		}
		if slices.Contains(out, flag) {
			continue
		}
		out = append(out, flag)
	}

	slices.Sort(out)
	return strings.Join(out, "|")
}

// PrepareOutput creates the final output struct, combining LogEntry and serialized flags.
func PrepareOutput(entry config.LogEntry, flagsString string) config.EnrichedLog {
	return config.EnrichedLog{
		ClientIP:       entry.ClientIP,
		FingerprintJA4: entry.FingerprintJA4,
		UserAgent:      entry.UserAgent,
		Flags:          flagsString,
		AllFields:      entry.AllFields,
	}
}

// OpenOutputFile creates and returns a pointer to the os.File and a CSV writer,
// writing the provided input header plus the "flags" column.
func OpenOutputFile(filePath string, inputHeader []string) (*os.File, *csv.Writer, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, nil, err
	}

	writer := csv.NewWriter(file)
	outputHeader := append(inputHeader, "flags")

	if err := writer.Write(outputHeader); err != nil {
		file.Close()
		return nil, nil, err
	}

	return file, writer, nil
}

// WriteRecord takes an EnrichedLog struct and writes the original fields (AllFields)
// plus the calculated Flags as a new CSV record.
func WriteRecord(writer *csv.Writer, log config.EnrichedLog) error {
	record := append(log.AllFields, log.Flags)
	return writer.Write(record)
}
