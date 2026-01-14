package iplist

import (
	"fmt"
	"net/netip"
	"sort"
	"testing"
)

// mustParsePrefix is a test helper for creating netip.Prefixes easily in tests.
func mustParsePrefix(s string) netip.Prefix {
	p, err := netip.ParsePrefix(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse prefix %q: %v", s, err))
	}
	return p
}

// TestParseIPNet tests the ParseIPNet function for converting raw strings to prefixes.
func TestParseIPNet(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    netip.Prefix
		wantErr bool
	}{
		{
			name:  "IPv4 CIDR",
			input: "192.168.1.0/24",
			want:  mustParsePrefix("192.168.1.0/24"),
		},
		{
			name:  "IPv4 Address to /32",
			input: "8.8.8.8",
			want:  mustParsePrefix("8.8.8.8/32"),
		},
		{
			name:  "IPv6 CIDR",
			input: "2001:db8::/32",
			want:  mustParsePrefix("2001:db8::/32"),
		},
		{
			name:  "IPv6 Address to /128",
			input: "2001:4860:4860::8888",
			want:  mustParsePrefix("2001:4860:4860::8888/128"),
		},
		{
			name:    "Invalid IP",
			input:   "256.256.256.256",
			wantErr: true,
		},
		{
			name:    "Invalid CIDR format",
			input:   "1.1.1.1/33",
			wantErr: true,
		},
		{
			name:    "Empty String",
			input:   "",
			wantErr: true,
		},
		{
			name:    "Whitespace String",
			input:   "  \t",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIPNet(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseIPNet() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseIPNet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsBogon tests the bogon filtering logic against known ranges.
func TestIsBogon(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{name: "Non-Bogon IPv4", ip: "1.1.1.1/32", want: false},
		{name: "Non-Bogon IPv6", ip: "2606:4700::/32", want: false},
		{name: "Bogon IPv4 Private A", ip: "10.0.0.1/32", want: true},
		{name: "Bogon IPv4 Private B", ip: "172.20.0.1/32", want: true},
		{name: "Bogon IPv4 Private C", ip: "192.168.1.1/32", want: true},
		{name: "Bogon IPv4 Loopback", ip: "127.0.0.1/32", want: true},
		{name: "Bogon IPv4 Link-local", ip: "169.254.1.1/32", want: true},
		{name: "Bogon IPv6 ULA", ip: "fc00::1/128", want: true},
		{name: "Bogon IPv6 Link-local", ip: "fe80::1/128", want: true},
		{name: "Bogon IPv6 Docs", ip: "2001:db8::1/128", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prefix := mustParsePrefix(tt.ip)
			if got := isBogon(prefix); got != tt.want {
				t.Errorf("isBogon(%s) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}

// TestPrefixSliceSort tests the custom sort order for netip.Prefixes.
func TestPrefixSliceSort(t *testing.T) {
	input := []netip.Prefix{
		mustParsePrefix("192.168.1.0/24"),  // 3. C
		mustParsePrefix("10.0.0.0/8"),      // 1. A
		mustParsePrefix("192.168.1.0/28"),  // 4. D (same address, longer mask)
		mustParsePrefix("172.16.0.0/12"),   // 2. B
		mustParsePrefix("1.1.1.1/32"),      // 0. (Smallest address)
		mustParsePrefix("2001:db8::/32"),   // 5. IPv6 start
		mustParsePrefix("2001:db8::1/128"), // 6. IPv6 end
	}

	want := []netip.Prefix{
		mustParsePrefix("1.1.1.1/32"),
		mustParsePrefix("10.0.0.0/8"),
		mustParsePrefix("172.16.0.0/12"),
		mustParsePrefix("192.168.1.0/24"),
		mustParsePrefix("192.168.1.0/28"),
		mustParsePrefix("2001:db8::/32"),
		mustParsePrefix("2001:db8::1/128"),
	}

	sort.Sort(prefixSlice(input))

	if len(input) != len(want) {
		t.Fatalf("Length mismatch after sort: got %d, want %d", len(input), len(want))
	}

	for i, got := range input {
		if got != want[i] {
			t.Errorf("Sort mismatch at index %d: got %v, want %v", i, got, want[i])
		}
	}
}

// TestAggregateIPNets tests the aggregation (sorting) function.
func TestAggregateIPNets(t *testing.T) {
	input := []netip.Prefix{
		mustParsePrefix("1.1.1.1/32"),
		mustParsePrefix("1.1.0.0/24"),
		mustParsePrefix("2.2.2.2/32"),
		mustParsePrefix("1.1.1.0/24"),
	}

	want := []netip.Prefix{
		mustParsePrefix("1.1.0.0/23"),
		mustParsePrefix("2.2.2.2/32"),
	}

	got := aggregateIPNets(input)

	if len(got) != len(want) {
		t.Fatalf("Length mismatch: got %d, want %d", len(got), len(want))
	}
	for i, p := range got {
		if p != want[i] {
			t.Errorf("Aggregation/Sort mismatch at index %d: got %v, want %v", i, p, want[i])
		}
	}
}

// TestAggregateIPNets tests the aggregation (summarization and redundancy removal) function.
func TestAggregateIPNets2(t *testing.T) {
	// User requested test case: checks for proper summarization and redundancy removal.
	t.Run("Summarization and Redundancy", func(t *testing.T) {
		input := []netip.Prefix{
			mustParsePrefix("1.1.1.1/32"), // Redundant (contained in 1.1.1.0/24)
			mustParsePrefix("2.2.2.2/32"), // Not mergeable
			mustParsePrefix("1.1.0.0/24"), // Sibling 1
			mustParsePrefix("1.1.1.0/24"), // Sibling 2
		}

		want := []netip.Prefix{
			mustParsePrefix("1.1.0.0/23"), // Merged from 1.1.0.0/24 and 1.1.1.0/24
			mustParsePrefix("2.2.2.2/32"), // Unchanged
		}

		got := aggregateIPNets(input)

		if len(got) != len(want) {
			t.Fatalf("Summarization length mismatch. Got %d, want %d. Got: %v", len(got), len(want), got)
		}
		for i, p := range got {
			if p != want[i] {
				t.Errorf("Summarization mismatch at index %d: got %v, want %v", i, p, want[i])
			}
		}
	})

	t.Run("Multi-Pass Merging", func(t *testing.T) {
		input := []netip.Prefix{
			mustParsePrefix("10.0.0.0/24"),
			mustParsePrefix("10.0.1.0/24"),
			mustParsePrefix("10.0.2.0/24"),
			mustParsePrefix("10.0.3.0/24"),
		}

		// Expected: (10.0.0.0/24 + 10.0.1.0/24) -> 10.0.0.0/23
		// (10.0.2.0/24 + 10.0.3.0/24) -> 10.0.2.0/23
		// (10.0.0.0/23 + 10.0.2.0/23) -> 10.0.0.0/22

		want := []netip.Prefix{
			mustParsePrefix("10.0.0.0/22"),
		}

		got := aggregateIPNets(input)

		if len(got) != len(want) {
			t.Fatalf("Multi-Pass Merging length mismatch. Got %d, want %d. Got: %v", len(got), len(want), got)
		}
		for i, p := range got {
			if p != want[i] {
				t.Errorf("Multi-Pass Merging mismatch at index %d: got %v, want %v", i, p, want[i])
			}
		}
	})
}

// TestParseAndValidate tests the entire pipeline: parse, validate, unique, filter bogon, and sort.
func TestParseAndValidate(t *testing.T) {
	rawIPs := []string{
		"1.1.1.1",          // Valid IPv4, non-bogon, converts to /32
		"10.0.0.1/32",      // Bogon, should be filtered
		"2.2.2.2/32",       // Valid IPv4, non-bogon
		"1.1.1.1/32",       // Duplicate, should be filtered
		"invalid",          // Invalid, should be ignored
		"2001:db8::1",      // Bogon IPv6, should be filtered
		"2606:4700::/32",   // Valid IPv6, non-bogon
		"2606:4700::1/128", // Valid IPv6, non-bogon
	}

	// Expected results (sorted, unique, non-bogon)
	wantV4 := []netip.Prefix{
		mustParsePrefix("1.1.1.1/32"),
		mustParsePrefix("2.2.2.2/32"),
	}
	wantV6 := []netip.Prefix{
		mustParsePrefix("2606:4700::/32"),
	}

	collection := parseAndValidateIPList(rawIPs)

	if len(collection.IPv4) != len(wantV4) {
		t.Fatalf("IPv4 length mismatch. Got %d, want %d. Got: %v", len(collection.IPv4), len(wantV4), collection.IPv4)
	}
	for i, p := range collection.IPv4 {
		if p != wantV4[i] {
			t.Errorf("IPv4 mismatch at %d: got %v, want %v", i, p, wantV4[i])
		}
	}

	if len(collection.IPv6) != len(wantV6) {
		t.Fatalf("IPv6 length mismatch. Got %d, want %d. Got: %v", len(collection.IPv6), len(wantV6), collection.IPv6)
	}
	for i, p := range collection.IPv6 {
		if p != wantV6[i] {
			t.Errorf("IPv6 mismatch at %d: got %v, want %v", i, p, wantV6[i])
		}
	}
}
