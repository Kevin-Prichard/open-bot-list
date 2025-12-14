package iplist

import (
	"fmt"
	"path/filepath"
	"strings"

	"git.oxl.at/open-bot-list/downloader/internal/config"
	"git.oxl.at/open-bot-list/downloader/internal/util"
)

// ProcessIPLists reads the IP list manifest files, processes the downloaded IP lists based on format,
// performs validation and aggregation, and writes the final output files.
func ProcessIPLists() error {
	fmt.Println("\nStarting IP List Processing from Manifests...")

	// Column indices based on FormatManifestIPNet struct
	const (
		MatchNameCol       = 1
		URLCol             = 2
		FormatCol          = 3
		JsonPathRFC9535Col = 4
		RegexCol           = 6
		CsvFieldCol        = 7
		PlainCol           = 8
	)

	for category, filename := range config.IPListCategories {
		if category == "_overall" {
			continue
		}

		fmt.Printf(" - Processing manifest: %s\n", filename)

		// Read the manifest for this category
		records, err := util.ParseCSVFile(config.FILE_PREFIX_MATCH + filename)
		if err != nil {
			fmt.Printf("   - WARNING: Could not read manifest for category '%s' (%s): %v\n", category, filename, err)
			continue
		}

		for _, record := range records {
			if len(record) < 9 {
				continue // Skip rows that don't have the required columns
			}

			matchName := record[MatchNameCol]
			urlStr := record[URLCol]
			format := record[FormatCol]
			jsonPath := record[JsonPathRFC9535Col]
			regexStr := record[RegexCol]
			csvField := record[CsvFieldCol]
			plainValue := record[PlainCol]

			urls := strings.Split(urlStr, config.VALUE_MULTI_DELIMITER)
			var combinedRawIPs []string
			hasError := false

			// Check if any URL exists or if it's a 'plain' format entry
			isPlain := format == "plain"
			needsProcessing := isPlain || len(urls) > 0

			if !needsProcessing {
				continue
			}

			fmt.Printf("   - Processing entry '%s' (Format: %s)... ", matchName, format)

			if isPlain {
				var parseErr error
				combinedRawIPs, parseErr = parseIPListPlain(plainValue)
				if parseErr != nil {
					fmt.Printf("FAILED plain parsing: %v\n", parseErr)
					hasError = true
				}
			} else {
				// Iterate over all downloaded files for this single manifest entry
				for i, url := range urls {
					url = strings.TrimSpace(url)
					if url == "" {
						continue
					}

					// Reconstruct the runtime file name used during download: iplist_<category>_<match_name>_<index>
					targetFileBase := fmt.Sprintf("%s%s_%s_%d", config.FILE_PREFIX_IPLIST, category, matchName, i)
					targetPath := filepath.Join(config.PATH_RUNTIME, targetFileBase)

					var rawIPs []string
					var parseErr error

					// 1. Parse the downloaded file based on format
					switch format {
					case "nlsv":
						rawIPs, parseErr = parseIPListNLsv(targetPath)
					case "csv":
						rawIPs, parseErr = parseIPListCsv(targetPath, csvField)
					case "json":
						rawIPs, parseErr = parseIPListJson(targetPath, jsonPath)
					case "regex-json":
						rawIPs, parseErr = parseIPListRegexJson(targetPath, regexStr, jsonPath)
					default:
						parseErr = fmt.Errorf("unsupported format: %s", format)
					}

					if parseErr != nil {
						fmt.Printf("\n     -> Parsing file %d failed: %v\n", i, parseErr)
						hasError = true
						// Continue processing other URLs if possible, but mark entry as failed
						continue
					}
					combinedRawIPs = append(combinedRawIPs, rawIPs...)
				}
			}

			if hasError {
				continue
			}

			// 2. Post-processing/Validation (done once on combined list)
			collection := parseAndValidate(combinedRawIPs)

			// 3. Write output files (using only matchName, removing category prefix and index)
			outputBaseName := matchName
			if err := writeIPLists(config.PATH_OUTPUT, outputBaseName, collection); err != nil {
				fmt.Printf("FAILED writing output: %v\n", err)
				continue
			}

			fmt.Println("SUCCESS")
		}
	}

	return nil
}

// DownloadIPListsFromManifests reads the IP list manifest files and downloads the IP lists.
// It skips the _overall.csv manifest.
func DownloadIPListsFromManifests() error {
	fmt.Println("\nStarting IP List Download from Manifests...")

	const (
		MatchNameCol = 1 // 0-based index
		URLCol       = 2 // 0-based index
	)

	for category, filename := range config.IPListCategories {
		if category == "_overall" {
			continue
		}

		fmt.Printf(" - Processing manifest: %s\n", filename)

		records, err := util.ParseCSVFile(config.FILE_PREFIX_MATCH + filename)
		if err != nil {
			fmt.Printf("   - WARNING: Could not read manifest for category '%s' (%s): %v\n", category, filename, err)
			continue
		}

		for _, record := range records {
			if len(record) < 3 {
				continue
			}

			matchName := record[MatchNameCol]
			urlStr := record[URLCol]
			urls := strings.Split(urlStr, config.VALUE_MULTI_DELIMITER)

			for i, url := range urls {
				url = strings.TrimSpace(url)
				if url == "" {
					continue
				}

				// The file path includes: category, match_name, and array-index
				// Example: "ip_net/ai.csv" -> runtime_dir/iplist_ai_matchname_0
				targetFileBase := fmt.Sprintf("%s%s_%s_%d", config.FILE_PREFIX_IPLIST, category, matchName, i)
				targetPath := filepath.Join(config.PATH_RUNTIME, targetFileBase)

				fmt.Printf("   - Downloading URL %d for %s (%s)... ", i, matchName, url)

				if err := util.DownloadFile(url, targetPath); err != nil {
					fmt.Printf("FAILED: %v\n", err)
					continue
				}
				fmt.Printf("SUCCESS -> %s\n", targetFileBase)
			}
		}
	}

	return nil
}
