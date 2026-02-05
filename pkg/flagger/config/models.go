package config

// LogEntry represents the relevant fields from one line of the input CSV.
type LogEntry struct {
	ClientIP       string
	FingerprintJA4 string
	UserAgent      string
	AllFields      []string
}

// EnrichedLog represents one line of the final output CSV.
// It includes the core fields from LogEntry and the serialized flags string.
type EnrichedLog struct {
	ClientIP       string
	FingerprintJA4 string
	UserAgent      string
	Flags          string
	AllFields      []string
}
