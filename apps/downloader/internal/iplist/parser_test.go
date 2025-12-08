package iplist

import (
	"encoding/json"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// Helper function to create a temporary file with the given content.
func createTempFile(t *testing.T, content string) string {
	f, err := os.CreateTemp("", "ip-parser-test-")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		os.Remove(f.Name())
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	return f.Name()
}

// TestParseJsonPath tests the jsonpath logic.
func TestParseJsonPath(t *testing.T) {
	content := `
{"addresses":["23.235.32.0/20","43.249.72.0/22","103.244.50.0/24"],"ipv6_addresses":["2a04:4e40::/32","2a04:4e42::/32"]}
`
	jsonPath := "$..['addresses','ipv6_addresses'].*"
	expectedIPs := []string{
		"23.235.32.0/20", "43.249.72.0/22", "103.244.50.0/24", "2a04:4e40::/32", "2a04:4e42::/32",
	}

	var jsonData interface{}
	if err := json.Unmarshal([]byte(content), &jsonData); err != nil {
		t.Errorf("test-content is not valid JSON: %v\n", err)
	}

	res, err := jsonpathFind(jsonData, jsonPath)

	if err != nil {
		t.Errorf("got error from jsonpath: %v\n", err)
	}
	if len(res) != len(expectedIPs) {
		t.Errorf("got unexpected jsonpath result: %d != %d\n", len(res), len(expectedIPs))
	}
	for _, expectedIp := range expectedIPs {
		if !slices.Contains(res, expectedIp) {
			t.Errorf("missing IP in jsonpath result: %s not in %v", expectedIp, res)
		}
	}
}

// TestParseIPListNLsv tests the New Line Separated Values parser logic.
func TestParseIPListNLsv(t *testing.T) {
	content := `
89.187.188.227
89.187.188.228
139.180.134.196
# comment1 
89.187.162.249
89.187.162.242


185.102.217.65 # test (inline comment)
185.93.1.243
; comment 2 (should be skipped by heuristic)
[2a04:4e42::]/32
156.146.40.49
185.59.220.199
185.59.220.198
`
	expected := []string{
		"89.187.188.227",
		"89.187.188.228",
		"139.180.134.196",
		"89.187.162.249",
		"89.187.162.242",
		"185.102.217.65",
		"185.93.1.243",
		"[2a04:4e42::]/32",
		"156.146.40.49",
		"185.59.220.199",
		"185.59.220.198",
	}

	path := createTempFile(t, content)
	defer os.Remove(path)

	ips, err := parseIPListNLsv(path)
	if err != nil {
		t.Fatalf("parseIPListNLsv failed: %v", err)
	}

	if !reflect.DeepEqual(ips, expected) {
		t.Errorf("parseIPListNLsv mismatch.\nGot:\n%s\n\nWant:\n%s", strings.Join(ips, "\n"), strings.Join(expected, "\n"))
	}
}

// TestParseIPListCsv tests the CSV parser logic for extracting IPs from a specific column.
func TestParseIPListCsv(t *testing.T) {
	content := `
172.224.226.0/27,GB,GB-EN,London,
172.224.226.32/31,GB,GB-SC,Aberdeen,
172.224.226.34/31,GB,GB-EN,Oxford,
172.224.226.36/31,GB,GB-EN,Luton,
172.224.226.38/31,GB,GB-NI,Belfast,
172.224.226.40/31,GB,GB-SC,Dundee,
172.224.226.42/31,GB,GB-EN,Brighton,
172.224.226.44/31,GB,GB-EN,Leicester,
172.224.226.46/31,GB,GB-EN,Liverpool,
172.224.226.48/31,GB,GB-SC,Edinburgh,
`
	// Since the first line is a valid IP, the parser logic should include it as data.
	expected := []string{
		"172.224.226.0/27",
		"172.224.226.32/31",
		"172.224.226.34/31",
		"172.224.226.36/31",
		"172.224.226.38/31",
		"172.224.226.40/31",
		"172.224.226.42/31",
		"172.224.226.44/31",
		"172.224.226.46/31",
		"172.224.226.48/31",
	}
	csvField := "0" // The index is derived from this string: 0-based

	path := createTempFile(t, content)
	defer os.Remove(path)

	ips, err := parseIPListCsv(path, csvField)
	if err != nil {
		t.Fatalf("parseIPListCsv failed: %v", err)
	}

	if !reflect.DeepEqual(ips, expected) {
		t.Errorf("parseIPListCsv mismatch.\nGot:\n%s\n\nWant:\n%s", strings.Join(ips, "\n"), strings.Join(expected, "\n"))
	}
}

func TestParseIPListCsv2(t *testing.T) {
	content := `
cidr,country,locale,city,anything
172.224.226.0/27,GB,GB-EN,London,
172.224.226.32/31,GB,GB-SC,Aberdeen,
172.224.226.34/31,GB,GB-EN,Oxford,
172.224.226.36/31,GB,GB-EN,Luton,
172.224.226.38/31,GB,GB-NI,Belfast,
172.224.226.40/31,GB,GB-SC,Dundee,
172.224.226.42/31,GB,GB-EN,Brighton,
172.224.226.44/31,GB,GB-EN,Leicester,
172.224.226.46/31,GB,GB-EN,Liverpool,
172.224.226.48/31,GB,GB-SC,Edinburgh,
`
	expected := []string{
		"172.224.226.0/27",
		"172.224.226.32/31",
		"172.224.226.34/31",
		"172.224.226.36/31",
		"172.224.226.38/31",
		"172.224.226.40/31",
		"172.224.226.42/31",
		"172.224.226.44/31",
		"172.224.226.46/31",
		"172.224.226.48/31",
	}
	csvField := "0" // The index is derived from this string: 0-based

	path := createTempFile(t, content)
	defer os.Remove(path)

	ips, err := parseIPListCsv(path, csvField)
	if err != nil {
		t.Fatalf("parseIPListCsv failed: %v", err)
	}

	if !reflect.DeepEqual(ips, expected) {
		t.Errorf("parseIPListCsv mismatch.\nGot:\n%s\n\nWant:\n%s", strings.Join(ips, "\n"), strings.Join(expected, "\n"))
	}
}

// TestParseIPListJson tests JSON unmarshalling.
// NOTE: Since the real JSONPath implementation is stubbed, this test primarily
// verifies that the file is read and the content is valid JSON.
// The resulting list will be empty unless the stub is manually modified or mocked.
func TestParseIPListJson(t *testing.T) {
	content := `
{"creationTime":"2025-10-30T11:00:00.000000","prefixes":[{"ipv4Prefix":"132.196.86.0/24"},{"ipv4Prefix":"172.182.202.0/25"},{"ipv4Prefix":"172.182.204.0/24"},{"ipv4Prefix":"172.182.207.0/25"},{"ipv4Prefix":"172.182.214.0/24"},{"ipv4Prefix":"172.182.215.0/24"},{"ipv4Prefix":"20.125.66.80/28"},{"ipv4Prefix":"20.171.206.0/24"},{"ipv4Prefix":"20.171.207.0/24"},{"ipv4Prefix":"4.227.36.0/25"},{"ipv4Prefix":"52.230.152.0/24"},{"ipv4Prefix":"74.7.175.128/25"},{"ipv4Prefix":"74.7.227.0/25"},{"ipv4Prefix":"74.7.227.128/25"},{"ipv4Prefix":"74.7.228.0/25"},{"ipv4Prefix":"74.7.230.0/25"},{"ipv4Prefix":"74.7.241.0/25"},{"ipv4Prefix":"74.7.241.128/25"},{"ipv4Prefix":"74.7.242.0/25"},{"ipv4Prefix":"74.7.243.128/25"},{"ipv4Prefix":"74.7.244.0/25"}]}
`
	jsonPath := "$.prefixes.*['ipv4Prefix','ipv6Prefix']"
	expectedIPs := []string{
		"132.196.86.0/24",
		"172.182.202.0/25",
		"172.182.204.0/24",
		"172.182.207.0/25",
		"172.182.214.0/24",
		"172.182.215.0/24",
		"20.125.66.80/28",
		"20.171.206.0/24",
		"20.171.207.0/24",
		"4.227.36.0/25",
		"52.230.152.0/24",
		"74.7.175.128/25",
		"74.7.227.0/25",
		"74.7.227.128/25",
		"74.7.228.0/25",
		"74.7.230.0/25",
		"74.7.241.0/25",
		"74.7.241.128/25",
		"74.7.242.0/25",
		"74.7.243.128/25",
		"74.7.244.0/25",
	}

	t.Run("Valid JSON (Stubbed JSONPath)", func(t *testing.T) {
		path := createTempFile(t, content)
		defer os.Remove(path)

		ips, err := parseIPListJson(path, jsonPath)
		if err != nil {
			t.Errorf("parseIPListJson unexpectedly failed: %v", err)
		}

		if len(ips) != len(expectedIPs) {
			t.Errorf("got unexpected jsonpath result: %d != %d", len(ips), len(expectedIPs))
		}
		for _, expectedIp := range expectedIPs {
			if !slices.Contains(ips, expectedIp) {
				t.Errorf("missing IP in jsonpath result: %s not in %v", expectedIp, ips)
			}
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		path := createTempFile(t, `{"addresses": "no-array"`)
		defer os.Remove(path)

		_, err := parseIPListJson(path, jsonPath)
		if err == nil || !strings.Contains(err.Error(), "failed to unmarshal JSON") {
			t.Errorf("parseIPListJson should have failed on invalid JSON, but got: %v", err)
		}
	})
}
