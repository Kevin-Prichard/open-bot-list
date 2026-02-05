package manifest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git.oxl.at/open-bot-list/pkg/downloader/config"
	"git.oxl.at/open-bot-list/pkg/downloader/util"
)

// ProcessFingerprintManifests reads all fingerprint manifests (excluding _overall) and writes the resulting map files.
// Output format: "<fingerprint><space><clients>" in "fingerprint_<category>.map"
func ProcessFingerprintManifests() error {
	fmt.Println("Processing Fingerprint Manifests...")

	const (
		FingerprintCol = 2
		ClientCol      = 3
	)

	for category, filename := range config.FingerprintCategories {
		if category == "_overall" {
			continue
		}
		outputFileName := fmt.Sprintf("fingerprint_%s.map", category)
		outputPath := filepath.Join(config.PATH_OUTPUT, outputFileName)

		records, err := util.ParseCSVFile(config.FILE_PREFIX_MATCH + filename)
		if err != nil {
			fmt.Printf("   - WARNING: Could not read manifest for category '%s' (%s): %v\n", category, filename, err)
			continue
		}

		f, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
		}
		defer f.Close()

		fmt.Printf("   - Writing %s (%d records) -> %s\n", filename, len(records), outputPath)

		for _, record := range records {
			if len(record) < 3 {
				continue
			}

			clientName := record[ClientCol]
			fingerprintStr := record[FingerprintCol]

			fingerprints := strings.Split(fingerprintStr, config.VALUE_MULTI_DELIMITER)
			for _, fp := range fingerprints {
				fp = strings.TrimSpace(fp)
				if fp != "" {
					// Format: "<fingerprint><space><clients>"
					if _, err := fmt.Fprintf(f, "%s %s\n", fp, clientName); err != nil {
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

	overallFilename := config.FingerprintCategories["_overall"]
	overallRecords, err := util.ParseCSVFile(config.FILE_PREFIX_MATCH + overallFilename)
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
		categoryRecords, err := util.ParseCSVFile(config.FILE_PREFIX_MATCH + sourceFilename)
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

			fingerprints := strings.Split(fingerprintStr, config.VALUE_MULTI_DELIMITER)
			for _, fp := range fingerprints {
				fp = strings.TrimSpace(fp)
				if fp != "" {
					uniqueFingerprints[fp] = struct{}{}
				}
			}
		}

		// Write output file: "<match>.lst"
		outputFileName := fmt.Sprintf("%s.lst", outputMatchName)
		outputPath := filepath.Join(config.PATH_OUTPUT, outputFileName)

		f, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
		}
		defer f.Close()

		fmt.Printf("   - Writing list of %d fingerprints for match '%s' -> %s\n", len(uniqueFingerprints), outputMatchName, outputPath)

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

// ProcessOverallPTRManifests reads the overall PTR-manifest to generate
// a list of unique values, filtered by 'kind' for each source file.
// Output format: "<value>" (one per line) in a file named "<match_name><outAppend>.lst"
func ProcessOverallGenericManifests(categoryConfig map[string]string, categoryName string, cols int, kindCol int, valueCol int, outAppend string) error {
	fmt.Println("Processing Overall Generic Manifest...")

	const (
		OverallFileCol       = 0
		OverallKindFilterCol = 1
		OverallMatchNameCol  = 2
		OverallCols          = 3
	)

	overallFilename := categoryConfig["_overall"]
	overallRecords, err := util.ParseCSVFile(config.FILE_PREFIX_MATCH + overallFilename)
	if err != nil {
		return fmt.Errorf("failed to read _overall manifest %s: %w", overallFilename, err)
	}

	for _, overallRecord := range overallRecords {
		if len(overallRecord) < OverallCols {
			continue
		}

		sourceFileBaseName := overallRecord[OverallFileCol]
		kindFilter := overallRecord[OverallKindFilterCol]
		outputMatchName := overallRecord[OverallMatchNameCol]

		sourceFilename := fmt.Sprintf("%s/%s.csv", categoryName, sourceFileBaseName)
		categoryRecords, err := util.ParseCSVFile(config.FILE_PREFIX_MATCH + sourceFilename)
		if err != nil {
			fmt.Printf("   - WARNING: Could not read source manifest '%s' for overall match: %v\n", sourceFilename, err)
			continue
		}

		uniqueValues := make(map[string]struct{})

		for _, record := range categoryRecords {
			if len(record) < cols {
				continue
			}

			recordOrganization := record[kindCol]
			multiValues := record[valueCol]

			if kindFilter != "*" && recordOrganization != kindFilter {
				continue
			}

			for _, value := range strings.Split(multiValues, config.VALUE_MULTI_DELIMITER) {
				value = strings.TrimSpace(value)
				if value != "" {
					uniqueValues[value] = struct{}{}
				}
			}
		}

		// Write output file: "<match_name><outAppend>.lst"
		outputFileName := fmt.Sprintf("%s%s.lst", outputMatchName, outAppend)
		outputPath := filepath.Join(config.PATH_OUTPUT, outputFileName)

		f, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
		}
		defer f.Close()

		fmt.Printf("   - Writing list of %d values for match '%s' -> %s\n", len(uniqueValues), outputMatchName, outputPath)

		// Write unique values, one per line
		for uas := range uniqueValues {
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

// ProcessUserAgentManifests reads all manifests of this category (excluding _overall) and writes the resulting map files.
// Output format: "<match_name><space><value>" in "<outPrefix><category><outAppend>.map"
func ProcessGenericManifests(categoryConfig map[string]string, cols int, colMatchName int, colValue int, outPrefix string, outAppend string) error {
	fmt.Printf("Processing Manifests of %s...\n", outPrefix)

	for category, filename := range categoryConfig {
		if category == "_overall" {
			continue
		}
		outputFileName := fmt.Sprintf("%s_%s%s.map", outPrefix, category, outAppend)
		outputPath := filepath.Join(config.PATH_OUTPUT, outputFileName)

		records, err := util.ParseCSVFile(config.FILE_PREFIX_MATCH + filename)
		if err != nil {
			fmt.Printf("   - WARNING: Could not read manifest for category '%s' (%s): %v\n", category, filename, err)
			continue
		}

		f, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
		}
		defer f.Close()

		fmt.Printf("   - Writing %s (%d records) -> %s\n", filename, len(records), outputPath)

		for _, record := range records {
			if len(record) < cols {
				continue
			}

			matchName := record[colMatchName]
			multiValues := record[colValue]

			for _, value := range strings.Split(multiValues, config.VALUE_MULTI_DELIMITER) {
				value = strings.TrimSpace(value)
				if value != "" {
					// Format: "<match_name><space><value>"
					if _, err := fmt.Fprintf(f, "%s %s\n", matchName, value); err != nil {
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
