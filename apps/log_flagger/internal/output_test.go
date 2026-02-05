package internal

import (
	"git.oxl.at/open-bot-list/pkg/flagger/config"
	"reflect"
	"testing"
)

func TestSerializeFlags(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{
			name:  "Basic sorting and joining",
			input: []string{"bot_crawler", "org_google", "bot"},
			want:  "bot|bot_crawler|org_google",
		},
		{
			name:  "Remove duplicates",
			input: []string{"bot", "bot_crawler", "bot", "bot_crawler"},
			want:  "bot|bot_crawler",
		},
		{
			name:  "Remove internal source net flags",
			input: []string{"bot", "src_net_crawler_google_common", "org_google", "http_user_agent_crawler"},
			want:  "bot|org_google",
		},
		{
			name:  "Empty input",
			input: []string{},
			want:  "",
		},
		{
			name:  "Only internal flags (should result in empty string)",
			input: []string{"src_net_crawler_google_common", "http_user_agent_crawler"},
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SerializeFlags(tt.input)
			if got != tt.want {
				t.Errorf("SerializeFlags() got = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPrepareOutput(t *testing.T) {
	logEntry := config.LogEntry{
		ClientIP:       "1.2.3.4",
		FingerprintJA4: "t13d...",
		UserAgent:      "Test UA",
		AllFields:      []string{"time", "host", "1.2.3.4", "200", "t13d...", "Test UA"},
	}
	flagsStr := "bot|bot_crawler|org_google"

	want := config.EnrichedLog{
		ClientIP:       logEntry.ClientIP,
		FingerprintJA4: logEntry.FingerprintJA4,
		UserAgent:      logEntry.UserAgent,
		Flags:          flagsStr,
		AllFields:      logEntry.AllFields,
	}

	got := PrepareOutput(logEntry, flagsStr)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("PrepareOutput() got = %+v, want %+v", got, want)
	}
}
