package iplist

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/theory/jsonpath"

	"git.oxl.at/open-bot-list/downloader/internal/config"
)

// run json-query on the json-data and return it as list of strings (should always be in this format)
func jsonpathFind(data interface{}, path string) ([]string, error) {
	p, err := jsonpath.Parse(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "   - WARNING: Failed to parse JSON-query\n")
		return nil, err
	}
	resultRaw := p.Select(data)

	var result []string
	for _, item := range resultRaw {
		s, ok := item.(string)
		if !ok {
			fmt.Fprintf(os.Stderr, "   - WARNING: Got non-string JSON-query result\n")
			return nil, fmt.Errorf("got non-string value")
		}
		result = append(result, s)
	}
	return result, nil
}

// nlsv (New Line Separated Values) parser
func parseIPListNLsv(filePath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var ips []string
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Ignore lines starting with non-IP/CIDR characters (commented-out or invalid)
		// Heuristic: If it starts with non-digit, non-bracket, non-dot, it's a comment/non-IP.
		if len(line) > 0 && (line[0] < '0' || line[0] > '9') && line[0] != '[' && line[0] != '.' {
			continue
		}

		// If there is a hash comment, ignore anything after it
		if hashIdx := strings.Index(line, "#"); hashIdx != -1 {
			line = line[:hashIdx]
		}

		// If there is a space after the IP/CIDR - ignore anything after it (inline comments)
		if spaceIdx := strings.Index(line, " "); spaceIdx != -1 {
			line = line[:spaceIdx]
		}

		line = strings.TrimSpace(line)
		if line != "" {
			ips = append(ips, line)
		}
	}
	return ips, nil
}

// csv parser
func parseIPListCsv(filePath string, csvField string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)

	header, err := reader.Read()
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Determine column index. Default to 0 if not found.
	fieldIndex := -1
	for i, col := range header {
		if col == csvField {
			fieldIndex = i
			break
		}
	}
	if fieldIndex == -1 {
		fieldIndex = 0
	}

	if fieldIndex >= len(header) {
		return nil, fmt.Errorf("CSV field '%s' not found or index out of bounds", csvField)
	}

	var ips []string

	// Check if the actual header value is a valid IP/CIDR. If so, put it back into the list.
	// This handles the "ignore first line if not IP (headers)" requirement by checking if the content is an IP.
	headerIP := strings.TrimSpace(header[fieldIndex])
	if _, ipErr := ParseIPNet(headerIP); ipErr == nil {
		ips = append(ips, headerIP)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV record: %w", err)
		}

		if fieldIndex < len(record) {
			ips = append(ips, record[fieldIndex])
		}
	}

	return ips, nil
}

// json parser (uses rfc9535 jsonpath)
func parseIPListJson(filePath string, jsonPath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		fmt.Fprintf(os.Stderr, "   - WARNING: Not valid JSON\n")
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	results, err := jsonpathFind(jsonData, jsonPath)
	if err != nil {
		return nil, fmt.Errorf("JSONPath execution failed with path '%s': %w", jsonPath, err)
	}
	return results, nil
}

// plain parser
func parseIPListPlain(plainValue string) ([]string, error) {
	ips := strings.Split(plainValue, config.VALUE_MULTI_DELIMITER)
	var processedIPs []string
	for _, ip := range ips {
		ip = strings.TrimSpace(ip)
		if ip != "" {
			processedIPs = append(processedIPs, ip)
		}
	}
	return processedIPs, nil
}

// html-json parser
func parseIPListHtmlJson(filePath string, regexStr string, jsonPath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(regexStr)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex '%s': %w", regexStr, err)
	}

	// Find the first match. The captured content (group 1) is expected to be JSON.
	matches := re.FindSubmatch(data)
	if len(matches) < 2 {
		return nil, fmt.Errorf("regex did not match or did not capture content (need at least one capture group)")
	}

	// The captured content is expected to be JSON
	jsonContent := matches[1]

	var jsonData interface{}
	if err := json.Unmarshal(jsonContent, &jsonData); err != nil {
		fmt.Fprintf(os.Stderr, "   - WARNING: Regex result was not valid JSON: %s\n", string(bytes.TrimSpace(jsonContent)))
		return nil, fmt.Errorf("failed to unmarshal JSON from regex result: %w", err)
	}

	results, err := jsonpathFind(jsonData, jsonPath)
	if err != nil {
		return nil, fmt.Errorf("JSONPath execution failed with path '%s': %w", jsonPath, err)
	}
	return results, nil
}
