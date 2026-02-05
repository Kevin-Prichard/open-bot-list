package iplist

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"git.oxl.at/open-bot-list/pkg/downloader/config"
	"git.oxl.at/open-bot-list/pkg/downloader/util"
)

func processListsGeneric(kind string, categoryConfig map[string]string, filePrefix string, callback func(values []string, matchName string) error) error {
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

	for category, filename := range categoryConfig {
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
			var combinedValues []string
			hasError := false

			// Check if any URL exists or if it's a 'plain' format entry
			isPlain := format == "plain"
			needsProcessing := isPlain || len(urls) > 0

			if !needsProcessing {
				continue
			}

			fmt.Printf("   - Processing entry '%s' (Format: %s)... \n", matchName, format)

			if isPlain {
				var parseErr error
				combinedValues, parseErr = parseListPlain(plainValue)
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

					// Reconstruct the runtime file name used during download: <prefix>_<category>_<match_name>_<index>
					targetFileBase := fmt.Sprintf("%s%s_%s_%d", filePrefix, category, matchName, i)
					targetPath := filepath.Join(config.PATH_RUNTIME, targetFileBase)

					var values []string
					var parseErr error

					switch format {
					case "nlsv":
						values, parseErr = parseListNLsv(targetPath)
					case "csv":
						if kind == "ip" {
							values, parseErr = parseIPListCsv(targetPath, csvField)
						}
					case "json":
						values, parseErr = parseListJson(targetPath, jsonPath)
					case "ndjson":
						values, parseErr = parseListNdjson(targetPath, jsonPath)
					case "regex-json":
						values, parseErr = parseListRegexJson(targetPath, regexStr, jsonPath)
					default:
						parseErr = fmt.Errorf("unsupported format: %s", format)
					}

					if parseErr != nil {
						fmt.Printf("\n     -> Parsing file %d failed: %v\n", i, parseErr)
						hasError = true
						// Continue processing other URLs if possible, but mark entry as failed
						continue
					}
					combinedValues = append(combinedValues, values...)
				}
			}

			if hasError {
				continue
			}

			err = callback(combinedValues, matchName)
			if err != nil {
				fmt.Printf("%v", err)
				continue
			}
			fmt.Println("SUCCESS")
		}
	}

	return nil
}

func processIPListsFinish(values []string, matchName string) error {
	collection := parseAndValidateIPList(values)
	if err := writeIPLists(config.PATH_OUTPUT, matchName, collection); err != nil {
		return fmt.Errorf("FAILED writing output: %v", err)
	}
	return nil
}

// ProcessIPLists reads the IP list manifest files, processes the downloaded IP lists based on format,
// performs validation and aggregation, and writes the final output files.
func ProcessIPLists() error {
	fmt.Println("\nStarting IP-List Processing from Manifests...")
	return processListsGeneric("ip", config.IPListCategories, config.FILE_PREFIX_IPLIST, processIPListsFinish)
}

func processASNListsFinish(values []string, matchName string) error {
	var collection []int

	for _, val := range values {
		asn, err := strconv.Atoi(val)
		if err != nil {
			// fmt.Fprintf(os.Stderr, "   - SKIPPING: Could not parse '%s' as ASN integer\n", val)
			continue
		}
		collection = append(collection, asn)
	}

	if err := writeASNLists(config.PATH_OUTPUT, matchName, collection); err != nil {
		return fmt.Errorf("FAILED writing output: %v", err)
	}
	return nil
}

// ProcessASNLists reads the ASN list manifest files, processes the downloaded ASN lists based on format,
// performs validation and aggregation, and writes the final output files.
func ProcessASNLists() error {
	fmt.Println("\nStarting ASN-List Processing from Manifests...")
	return processListsGeneric("asn", config.ASNListCategories, config.FILE_PREFIX_ASNLIST, processASNListsFinish)
}

// DownloadListsFromManifests reads the IP list manifest files and downloads the IP lists.
// It skips the _overall.csv manifest.
func DownloadListsFromManifests(categoryConfig map[string]string, filePrefix string) error {
	fmt.Println("\nStarting List Download from Manifests...")

	const (
		MatchNameCol = 1 // 0-based index
		URLCol       = 2 // 0-based index
	)

	for category, filename := range categoryConfig {
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
				targetFileBase := fmt.Sprintf("%s%s_%s_%d", filePrefix, category, matchName, i)
				targetPath := filepath.Join(config.PATH_RUNTIME, targetFileBase)

				fmt.Printf("   - Downloading URL %d for %s (%s)... ", i, matchName, url)

				if err := util.DownloadFileWithCache(url, targetPath); err != nil {
					fmt.Printf("FAILED: %v\n", err)
					continue
				}
				fmt.Printf("SUCCESS -> %s\n", targetFileBase)
			}
		}
	}

	return nil
}
