package internal

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ParseCSVFile reads a CSV file from the runtime path, skips the header, and returns the data records.
func ParseCSVFile(filename string) ([][]string, error) {
	targetPath := filepath.Join(PATH_RUNTIME, strings.ReplaceAll(filename, "/", "_"))

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

// ProcessFingerprintManifests reads all fingerprint manifests (excluding _overall) and writes the resulting map files.
// Output format: "<match_name><space><fingerprint>" in "fingerprint_<category>.map"
func ProcessFingerprintManifests() error {
	fmt.Println("Processing Fingerprint Manifests...")

	const (
		MatchNameCol   = 1
		FingerprintCol = 2
	)

	for category, filename := range FingerprintCategories {
		if category == "_overall" {
			continue
		}
		outputFileName := fmt.Sprintf("fingerprint_%s.map", category)
		outputPath := filepath.Join(PATH_OUTPUT, outputFileName)

		records, err := ParseCSVFile(FILE_PREFIX_MATCH + filename)
		if err != nil {
			fmt.Printf("   - WARNING: Could not read manifest for category '%s' (%s): %v\n", category, filename, err)
			continue
		}

		f, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
		}
		defer f.Close()

		fmt.Printf("   - Writing %s (%d records) -> %s\n", filename, len(records), outputFileName)

		for _, record := range records {
			if len(record) < 3 {
				continue
			}

			matchName := record[MatchNameCol]
			fingerprintStr := record[FingerprintCol]

			fingerprints := strings.Split(fingerprintStr, VALUE_MULTI_DELIMITER)
			for _, fp := range fingerprints {
				fp = strings.TrimSpace(fp)
				if fp != "" {
					// Format: "<match_name><space><fingerprint>"
					if _, err := fmt.Fprintf(f, "%s %s\n", matchName, fp); err != nil {
						return fmt.Errorf("failed to write to file %s: %w", outputPath, err)
					}
				}
			}
		}

		if err := f.Sync(); err != nil {
			return fmt.Errorf("failed to sync output file %s: %w", outputPath, err)
		}
	}

	return nil
}

// ProcessOverallFingerprintManifests reads the overall fingerprint manifest to generate
// a list of unique fingerprints, filtered by 'kind', for each output match name.
// Output format: "<fingerprint>" (one per line) in a file named "<match>.lst"
func ProcessOverallFingerprintManifests() error {
	fmt.Println("Processing Overall Fingerprint Manifest...")

	const (
		OverallFileCol       = 0
		OverallKindFilterCol = 1
		OverallMatchNameCol  = 2
	)

	overallFilename := FingerprintCategories["_overall"]
	overallRecords, err := ParseCSVFile(FILE_PREFIX_MATCH + overallFilename)
	if err != nil {
		return fmt.Errorf("failed to read _overall manifest %s: %w", overallFilename, err)
	}

	const (
		CategoryKindCol        = 0
		CategoryFingerprintCol = 2
	)

	for _, overallRecord := range overallRecords {
		if len(overallRecord) < 3 {
			continue
		}

		sourceFileBaseName := overallRecord[OverallFileCol]
		kindFilter := overallRecord[OverallKindFilterCol]
		outputMatchName := overallRecord[OverallMatchNameCol]

		sourceFilename := fmt.Sprintf("fingerprint/%s.csv", sourceFileBaseName)
		categoryRecords, err := ParseCSVFile(FILE_PREFIX_MATCH + sourceFilename)
		if err != nil {
			fmt.Printf("   - WARNING: Could not read source manifest '%s' for overall match '%s': %v\n", sourceFilename, outputMatchName, err)
			continue
		}

		uniqueFingerprints := make(map[string]struct{})

		for _, record := range categoryRecords {
			if len(record) < 3 {
				continue
			}

			recordKind := record[CategoryKindCol]
			fingerprintStr := record[CategoryFingerprintCol]

			if kindFilter != "*" && recordKind != kindFilter {
				continue
			}

			fingerprints := strings.Split(fingerprintStr, VALUE_MULTI_DELIMITER)
			for _, fp := range fingerprints {
				fp = strings.TrimSpace(fp)
				if fp != "" {
					uniqueFingerprints[fp] = struct{}{}
				}
			}
		}

		// Write output file: "<match>.lst"
		outputFileName := fmt.Sprintf("%s.lst", outputMatchName)
		outputPath := filepath.Join(PATH_OUTPUT, outputFileName)

		f, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
		}
		defer f.Close()

		fmt.Printf("   - Writing list of %d fingerprints for match '%s' -> %s\n", len(uniqueFingerprints), outputMatchName, outputFileName)

		// Write unique fingerprints, one per line
		for fp := range uniqueFingerprints {
			if _, err := fmt.Fprintf(f, "%s\n", fp); err != nil {
				return fmt.Errorf("failed to write to file %s: %w", outputPath, err)
			}
		}

		if err := f.Sync(); err != nil {
			return fmt.Errorf("failed to sync output file %s: %w", outputPath, err)
		}
	}

	return nil
}

// ProcessOverallUserAgentManifests reads the overall user agent manifest to generate
// a list of unique user agent subs, filtered by 'kind' (Organization), for each source file.
// Output format: "<user_agent_sub>" (one per line) in a file named "<match_name>_sub.lst"
func ProcessOverallUserAgentManifests() error {
	fmt.Println("Processing Overall User Agent Manifest...")

	const (
		OverallFileCol       = 0
		OverallKindFilterCol = 1
		OverallMatchNameCol  = 2
	)

	overallFilename := UserAgentCategories["_overall"]
	overallRecords, err := ParseCSVFile(FILE_PREFIX_MATCH + overallFilename)
	if err != nil {
		return fmt.Errorf("failed to read _overall manifest %s: %w", overallFilename, err)
	}

	const (
		CategoryOrganizationCol = 0
		CategoryUserAgentSubCol = 2
	)

	for _, overallRecord := range overallRecords {
		if len(overallRecord) < 3 {
			continue
		}

		sourceFileBaseName := overallRecord[OverallFileCol]
		kindFilter := overallRecord[OverallKindFilterCol]
		outputMatchName := overallRecord[OverallMatchNameCol]

		sourceFilename := fmt.Sprintf("user_agent/%s.csv", sourceFileBaseName)
		categoryRecords, err := ParseCSVFile(FILE_PREFIX_MATCH + sourceFilename)
		if err != nil {
			fmt.Printf("   - WARNING: Could not read source manifest '%s' for overall match: %v\n", sourceFilename, err)
			continue
		}

		uniqueUserAgentSubs := make(map[string]struct{})

		for _, record := range categoryRecords {
			if len(record) < 3 {
				continue
			}

			recordOrganization := record[CategoryOrganizationCol]
			userAgentSubStr := record[CategoryUserAgentSubCol]

			if kindFilter != "*" && recordOrganization != kindFilter {
				continue
			}

			userAgentSubs := strings.Split(userAgentSubStr, VALUE_MULTI_DELIMITER)
			for _, uas := range userAgentSubs {
				uas = strings.TrimSpace(uas)
				if uas != "" {
					uniqueUserAgentSubs[uas] = struct{}{}
				}
			}
		}

		// Write output file: "<match_name>_sub.lst"
		outputFileName := fmt.Sprintf("%s_sub.lst", outputMatchName)
		outputPath := filepath.Join(PATH_OUTPUT, outputFileName)

		f, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
		}
		defer f.Close()

		fmt.Printf("   - Writing list of %d user-agent subs for match '%s' -> %s\n", len(uniqueUserAgentSubs), outputMatchName, outputFileName)

		// Write unique user agent subs, one per line
		for uas := range uniqueUserAgentSubs {
			if _, err := fmt.Fprintf(f, "%s\n", uas); err != nil {
				return fmt.Errorf("failed to write to file %s: %w", outputPath, err)
			}
		}

		if err := f.Sync(); err != nil {
			return fmt.Errorf("failed to sync output file %s: %w", outputPath, err)
		}
	}

	return nil
}

// ProcessUserAgentManifests reads all user agent manifests (excluding _overall) and writes the resulting map files.
// Output format: "<match_name><space><user_agent_sub>" in "http_user_agent_<category>_sub.map"
func ProcessUserAgentManifests() error {
	fmt.Println("Processing User Agent Manifests...")

	const (
		MatchNameCol    = 1
		UserAgentSubCol = 2
	)

	for category, filename := range UserAgentCategories {
		if category == "_overall" {
			continue
		}
		outputFileName := fmt.Sprintf("http_user_agent_%s_sub.map", category)
		outputPath := filepath.Join(PATH_OUTPUT, outputFileName)

		records, err := ParseCSVFile(FILE_PREFIX_MATCH + filename)
		if err != nil {
			fmt.Printf("   - WARNING: Could not read manifest for category '%s' (%s): %v\n", category, filename, err)
			continue
		}

		f, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
		}
		defer f.Close()

		fmt.Printf("   - Writing %s (%d records) -> %s\n", filename, len(records), outputFileName)

		for _, record := range records {
			if len(record) < 3 {
				continue
			}

			matchName := record[MatchNameCol]
			userAgentSubStr := record[UserAgentSubCol]

			userAgentSubs := strings.Split(userAgentSubStr, VALUE_MULTI_DELIMITER)
			for _, uas := range userAgentSubs {
				uas = strings.TrimSpace(uas)
				if uas != "" {
					// Format: "<match_name><space><user_agent_sub>"
					if _, err := fmt.Fprintf(f, "%s %s\n", matchName, uas); err != nil {
						return fmt.Errorf("failed to write to file %s: %w", outputPath, err)
					}
				}
			}
		}

		if err := f.Sync(); err != nil {
			return fmt.Errorf("failed to sync output file %s: %w", outputPath, err)
		}
	}

	return nil
}
