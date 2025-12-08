package iplist

import (
	"fmt"
	"net/netip"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// IPNetCollection holds the aggregated and validated IP networks, separated by version.
type IPNetCollection struct {
	IPv4 []netip.Prefix
	IPv6 []netip.Prefix
}

// List of well-known IANA/RFC/Bogon ranges for IPv4 and IPv6 to filter out.
var bogonRanges = []netip.Prefix{
	// IPv4 Bogons (common)
	netip.MustParsePrefix("0.0.0.0/8"),       // "This network"
	netip.MustParsePrefix("10.0.0.0/8"),      // Private use (RFC 1918)
	netip.MustParsePrefix("100.64.0.0/10"),   // Shared Address Space (RFC 6598)
	netip.MustParsePrefix("127.0.0.0/8"),     // Loopback (RFC 6890)
	netip.MustParsePrefix("169.254.0.0/16"),  // Link local (RFC 3927)
	netip.MustParsePrefix("172.16.0.0/12"),   // Private use (RFC 1918)
	netip.MustParsePrefix("192.0.0.0/24"),    // IANA reserved (RFC 6890)
	netip.MustParsePrefix("192.0.2.0/24"),    // TEST-NET-1 (RFC 5737)
	netip.MustParsePrefix("192.88.99.0/24"),  // 6to4 relay anycast (deprecated)
	netip.MustParsePrefix("192.168.0.0/16"),  // Private use (RFC 1918)
	netip.MustParsePrefix("198.18.0.0/15"),   // Benchmarking (RFC 2544)
	netip.MustParsePrefix("198.51.100.0/24"), // TEST-NET-2 (RFC 5737)
	netip.MustParsePrefix("203.0.113.0/24"),  // TEST-NET-3 (RFC 5737)
	netip.MustParsePrefix("224.0.0.0/4"),     // Multicast
	netip.MustParsePrefix("240.0.0.0/4"),     // Reserved
	// IPv6 Bogons (common)
	netip.MustParsePrefix("::/128"),        // Unspecified address
	netip.MustParsePrefix("::1/128"),       // Loopback address
	netip.MustParsePrefix("::ffff:0:0/96"), // IPv4 mapped addresses
	netip.MustParsePrefix("100::/64"),      // Discard prefix (RFC 6666)
	netip.MustParsePrefix("2001:db8::/32"), // Documentation (RFC 3849)
	netip.MustParsePrefix("fc00::/7"),      // Unique local address (ULA)
	netip.MustParsePrefix("fe80::/10"),     // Link-local unicast (RFC 4291)
	netip.MustParsePrefix("ff00::/8"),      // Multicast
}

// isBogon checks if the given netip.Prefix's address is contained within a known bogon range.
func isBogon(p netip.Prefix) bool {
	for _, bogon := range bogonRanges {
		if bogon.Contains(p.Addr()) {
			return true
		}
	}
	return false
}

// prefixSlice implements sort.Interface for []netip.Prefix based on address value.
type prefixSlice []netip.Prefix

func (p prefixSlice) Len() int      { return len(p) }
func (p prefixSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p prefixSlice) Less(i, j int) bool {
	// Sort primarily by address, then by mask length.
	addrI := p[i].Addr()
	addrJ := p[j].Addr()

	if addrI != addrJ {
		return addrI.Less(addrJ)
	}
	return p[i].Bits() < p[j].Bits()
}

// aggregateIPNets aggregates the list of prefixes using an iterative, greedy algorithm.
// It removes redundant (contained) prefixes and merges adjacent prefixes.
func aggregateIPNets(nets []netip.Prefix) []netip.Prefix {
	if len(nets) < 2 {
		return nets
	}

	sort.Sort(prefixSlice(nets))

	for {
		changed := false
		var aggregated []netip.Prefix

		if len(nets) == 0 {
			break
		}

		current := nets[0]

		for i := 1; i < len(nets); i++ {
			next := nets[i]

			if current.Addr().Is4() != next.Addr().Is4() {
				aggregated = append(aggregated, current)
				current = next
				continue
			}

			// 1. Redundancy check: if 'next' is contained in 'current', discard 'next'.
			// Since the list is sorted, this primarily handles cases where a shorter prefix (e.g., /24) precedes a longer one (e.g., /32) that is within its range.
			if current.Contains(next.Addr()) && current.Bits() <= next.Bits() {
				// next is contained in current. Discard next, keep current.
				changed = true
				continue
			}

			// 2. Merge check (same length L, shared parent L-1)
			if current.Bits() == next.Bits() && current.Bits() > 0 {
				newBits := current.Bits() - 1

				// Calculate the common parent prefix of length L-1
				parent := netip.PrefixFrom(current.Addr(), newBits).Masked()

				// Check if 'next' also belongs to 'parent'.
				if parent == netip.PrefixFrom(next.Addr(), newBits).Masked() {
					// Success: merge current and next into 'parent'.
					// Continue the loop with the merged prefix as the new 'current'.
					current = parent
					changed = true
					continue
				}
			}

			// 3. No redundancy/merge possible. Add 'current' to the aggregated list and move to 'next'.
			aggregated = append(aggregated, current)
			current = next
		}

		// Add the last accumulated prefix
		aggregated = append(aggregated, current)

		if !changed {
			return aggregated
		}

		// A change occurred (merge or removal). Update 'nets' and re-sort for the next pass.
		nets = aggregated
		sort.Sort(prefixSlice(nets))
	}
	return nets
}

// parseAndValidate processes a list of raw IP/CIDR strings, validates them,
// filters bogons, and aggregates the resulting prefixes.
func parseAndValidate(rawIPs []string) IPNetCollection {
	var v4Nets []netip.Prefix
	var v6Nets []netip.Prefix
	uniqueNets := make(map[netip.Prefix]struct{})

	for _, rawIP := range rawIPs {
		rawIP = strings.TrimSpace(rawIP)
		if rawIP == "" {
			continue
		}

		p, err := ParseIPNet(rawIP)
		if err != nil {
			fmt.Fprintf(os.Stderr, "   - WARNING: Ignoring invalid IP/CIDR '%s'\n", rawIP)
			continue
		}

		if _, ok := uniqueNets[p]; ok {
			continue
		}

		if isBogon(p) {
			fmt.Printf("   - WARNING: Ignoring bogon network '%s'\n", p.String())
			continue
		}

		uniqueNets[p] = struct{}{}

		if p.Addr().Is4() {
			v4Nets = append(v4Nets, p)
		} else if p.Addr().Is6() {
			v6Nets = append(v6Nets, p)
		}
	}

	return IPNetCollection{
		IPv4: aggregateIPNets(v4Nets),
		IPv6: aggregateIPNets(v6Nets),
	}
}

// writeIPLists writes the aggregated IP lists to the specified output files.
func writeIPLists(outputDir, baseName string, collection IPNetCollection) error {
	v4Count := len(collection.IPv4)
	v6Count := len(collection.IPv6)

	totalCount := v4Count + v6Count
	if totalCount > 10000 {
		fmt.Printf("WARNING: IP list '%s' contains %d total entries (> 10000) ", baseName, totalCount)
	}

	writer := func(filename string, prefixes []netip.Prefix) error {
		outputPath := filepath.Join(outputDir, filename)
		f, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
		}
		defer f.Close()

		for _, p := range prefixes {
			if _, err := fmt.Fprintf(f, "%s\n", p.String()); err != nil {
				return fmt.Errorf("failed to write to file %s: %w", outputPath, err)
			}
		}

		if err := f.Sync(); err != nil {
			return fmt.Errorf("failed to sync output file %s: %w", outputPath, err)
		}
		return nil
	}

	v4FileName := fmt.Sprintf("%s_net4.lst", baseName)
	if err := writer(v4FileName, collection.IPv4); err != nil {
		return err
	}

	v6FileName := fmt.Sprintf("%s_net6.lst", baseName)
	if err := writer(v6FileName, collection.IPv6); err != nil {
		return err
	}

	allList := append(collection.IPv4, collection.IPv6...)
	allFileName := fmt.Sprintf("%s_all.lst", baseName)
	if err := writer(allFileName, allList); err != nil {
		return err
	}

	return nil
}

// ParseIPNet is a helper to parse a raw string into a netip.Prefix, handling single IPs as /32 or /128.
func ParseIPNet(s string) (netip.Prefix, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return netip.Prefix{}, fmt.Errorf("empty string")
	}

	if p, err := netip.ParsePrefix(s); err == nil {
		return p, nil
	}

	if addr, err := netip.ParseAddr(s); err == nil {
		return netip.PrefixFrom(addr, addr.BitLen()), nil
	}

	return netip.Prefix{}, fmt.Errorf("not a valid IP or CIDR")
}
