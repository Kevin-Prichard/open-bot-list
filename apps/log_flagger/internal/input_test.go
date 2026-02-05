package internal

import (
	"reflect"
	"testing"

	"git.oxl.at/open-bot-list/pkg/flagger/config"
)

func TestParseRecord(t *testing.T) {
	// Temporarily set the indices to non-default values to ensure the logic works regardless of config.
	originalIPIndex := config.CSV_FIELD_CLIENT_IP
	originalFPIndex := config.CSV_FIELD_FP_JA4
	originalUAIndex := config.CSV_FIELD_UA

	config.CSV_FIELD_CLIENT_IP = 2
	config.CSV_FIELD_FP_JA4 = 4
	config.CSV_FIELD_UA = 5

	t.Cleanup(func() {
		// Restore original config values
		config.CSV_FIELD_CLIENT_IP = originalIPIndex
		config.CSV_FIELD_FP_JA4 = originalFPIndex
		config.CSV_FIELD_UA = originalUAIndex
	})

	record := []string{
		"2025-11-12T13:52:11.144000+01:00",     // 0
		"7c8f2.test.oxl.app",                   // 1
		"213.55.221.232",                       // 2 -> ClientIP
		"200",                                  // 3
		"q13d0311h3_55b375c5d22e_f2a83c8e78ae", // 4 -> FingerprintJA4
		"Python requests someversion",          // 5 -> UserAgent
		"extra_field",                          // 6
	}

	want := config.LogEntry{
		ClientIP:       "213.55.221.232",
		FingerprintJA4: "q13d0311h3_55b375c5d22e_f2a83c8e78ae",
		UserAgent:      "Python requests someversion",
		AllFields:      record,
	}

	got := ParseRecord(record)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseRecord() got = %+v, want %+v", got, want)
	}
}
