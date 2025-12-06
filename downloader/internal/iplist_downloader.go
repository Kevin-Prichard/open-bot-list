package internal

import (
	"fmt"
	"path/filepath"
	"strings"
)

// DownloadIPListsFromManifests reads the IP list manifest files and downloads the IP lists.
// It skips the _overall.csv manifest.
// The output file name is: "<category>_<match_name>_<index>" (in the runtime dir).
func DownloadIPListsFromManifests() error {
	fmt.Println("\nStarting IP List Download from Manifests...")

	const (
		MatchNameCol = 1
		URLCol       = 2
	)

	for category, filename := range IPListCategories {
		if category == "_overall" {
			continue
		}

		fmt.Printf(" - Processing manifest: %s\n", filename)

		records, err := ParseCSVFile(FILE_PREFIX_MATCH + filename)
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
			urls := strings.Split(urlStr, VALUE_MULTI_DELIMITER)

			for i, url := range urls {
				url = strings.TrimSpace(url)
				if url == "" {
					continue
				}

				// Example: "ip_net/ai.csv" -> runtime_dir/iplist_ai_matchname_0
				targetFileBase := fmt.Sprintf("%s%s_%s_%d", FILE_PREFIX_IPLIST, category, matchName, i)
				targetPath := filepath.Join(PATH_RUNTIME, targetFileBase)

				fmt.Printf("   - Downloading URL %d for %s: %s ... ", i, matchName, url)

				if err := DownloadFile(url, targetPath); err != nil {
					fmt.Printf("FAILED: %v\n", err)
					continue
				}
				fmt.Printf("SUCCESS -> %s\n", targetFileBase)
			}
		}
	}

	return nil
}
