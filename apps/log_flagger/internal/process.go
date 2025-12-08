package internal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"git.oxl.at/open-bot-list/log_flagger/internal/config"
)

const TEST_LIMIT_COUNT = 250

// LoadAllLists walks the data directory, finds all *.lst files, reads their content,
// and passes them to LoadListContents for parsing.
func LoadAllLists() error {
	fmt.Printf("Searching for list files in: %s\n", config.PATH_OPEN_BOT_DATA)

	if _, err := os.Stat(config.PATH_OPEN_BOT_DATA); os.IsNotExist(err) {
		return fmt.Errorf("list directory not found at: %s", config.PATH_OPEN_BOT_DATA)
	}

	err := filepath.Walk(config.PATH_OPEN_BOT_DATA, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		name := info.Name()
		isListOrMap := strings.HasSuffix(name, ".lst") || strings.HasSuffix(name, ".map")
		isExcludedList := strings.HasSuffix(name, "_all.lst")

		if info.IsDir() || !isListOrMap || isExcludedList {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		if err := LoadListContents(filepath.Base(path), string(content)); err != nil {
			fmt.Printf("Warning: Could not load %s: %v\n", path, err)

		} else if config.DEBUG {
			fmt.Printf("Loaded: %s\n", path)
		}

		return nil
	})

	return err
}

func logDebugLoaded() {
	if !config.DEBUG {
		return
	}

	fmt.Println("\n--- List Loading Complete. Proceeding with application logic. ---")

	fmt.Println("\n--- Loaded IP Lists (partial view) ---")
	countIP := 0
	for key, entries := range LoadedIPLists {
		if len(entries) > 0 {
			if config.DEBUG {
				fmt.Printf("Key: %s, Total entries: %d\n", key, len(entries))
			}
			countIP++
		}
	}
	if countIP == 0 {
		fmt.Println("No IP Lists Loaded.")
	}

	fmt.Println("\n--- Loaded User-Agent Lists (partial view) ---")
	countUA := 0
	for key, entries := range LoadedUserAgentLists {
		if len(entries) > 0 {
			if config.DEBUG {
				fmt.Printf("Key: %s, Total entries: %d\n", key, len(entries))
			}
			countUA++
		}
	}
	if countUA == 0 {
		fmt.Println("No User-Agent Lists Loaded.")
	}

	fmt.Println("\n--- Loaded Fingerprint Lists (partial view) ---")
	countFP := 0
	for key, entries := range LoadedFingerprintLists {
		if len(entries) > 0 {
			if config.DEBUG {
				fmt.Printf("Key: %s, Total entries: %d\n", key, len(entries))
			}
			countFP++
		}
	}
	if countFP == 0 {
		fmt.Println("No Fingerprint Lists Loaded.")
	}
}

// Run is the main entry point for the application logic.
// It orchestrates data loading and processing.
func Run() error {
	if err := LoadAllLists(); err != nil {
		return fmt.Errorf("fatal error loading lists: %w", err)
	}

	logDebugLoaded()

	fmt.Println("Starting log file processing")

	file, reader, header, err := OpenInputFile(config.PATH_INPUT)
	if err != nil {
		return fmt.Errorf("failed to open/read input file header: %w", err)
	}
	defer file.Close()

	outFile, writer, err := OpenOutputFile(config.PATH_OUTPUT, header)
	if err != nil {
		return fmt.Errorf("failed to open/write output file header: %w", err)
	}
	defer func() {
		writer.Flush()
		outFile.Close()
	}()

	count := 0
	for {
		if config.MODE_TEST && count >= TEST_LIMIT_COUNT {
			fmt.Printf("TEST MODE: Stopped processing after %d records.\n", TEST_LIMIT_COUNT)
			break
		}

		record, err := reader.Read()
		if err == io.EOF {
			fmt.Println("Finished processing file.")
			break
		}
		if err != nil {
			fmt.Printf("Warning: Failed to read CSV record %d: %v. Skipping.\n", count+1, err)
			count++
			continue
		}

		logEntry := ParseRecord(record)
		enrichedLog := EnrichLog(logEntry)

		if err := WriteRecord(writer, enrichedLog); err != nil {
			fmt.Printf("Warning: Failed to write CSV record %d: %v. Skipping.\n", count+1, err)
		}
		count++
	}

	// fmt.Println(LookupIP("66.249.77.224"))

	return nil
}
