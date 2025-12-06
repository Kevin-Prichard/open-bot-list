package config

type FormatManifestMatchOverall struct {
	File  string `csv:"file"`
	Kind  string `csv:"kind"`
	Match string `csv:"match"`
}

type FormatManifestUserAgent struct {
	Organization string   `csv:"organization"`
	MatchName    string   `csv:"match_name"`
	UserAgentSub []string `csv:"user_agent_sub"`
	DocsUrl      string   `csv:"docs_url"`
}

type FormatManifestIPNet struct {
	Organization    string   `csv:"organization"`
	MatchName       string   `csv:"match_name"`
	Url             []string `csv:"url"`
	Format          string   `csv:"format"`
	JsonPathRFC9535 string   `csv:"jsonpath_rfc9535"`
	JsonPathJQ      string   `csv:"jsonpath_jq"`
	Regex           string   `csv:"regex"`
	CsvField        string   `csv:"csv_field"`
	Plain           string   `csv:"plain"`
	DocsUrl         string   `csv:"docs_url"`
}

type FormatManifestFingerprint struct {
	Kind        string   `csv:"kind"`
	MatchName   string   `csv:"match_name"`
	Fingerprint []string `csv:"fingerprint"`
	Clients     []string `csv:"clients"`
}
