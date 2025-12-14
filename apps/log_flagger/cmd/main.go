package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"git.oxl.at/open-bot-list/log_flagger/internal"
	"git.oxl.at/open-bot-list/log_flagger/internal/config"
)

func main() {
	fmt.Printf("\nOXL Open-Bot-List Log-Flagger v%v\n", config.VERSION)
	fmt.Println("> © OXL IT Services / Rath Pascal")
	fmt.Println("> git.oxl.at/open-bot-list")
	fmt.Printf("> License: GPLv3\n\n")

	flag.StringVar(
		&config.PATH_OPEN_BOT_DATA, "data-dir", "",
		"Path to the directory containing the 'open-bot-list' data (required). See: https://github.com/O-X-L/open-bot-list/tree/latest?tab=readme-ov-file#downloader-application",
	)
	flag.StringVar(&config.PATH_OUTPUT, "output-file", "", "Path to the output CSV-file to be written (required).")
	flag.StringVar(&config.PATH_INPUT, "input-file", "", "Path to the input CSV-file to process (required).")
	flag.IntVar(&config.CSV_FIELD_CLIENT_IP, "csv-field-client-ip", 0, "Index of the CSV-Field containing the timestamp. (default 0)")
	flag.IntVar(&config.CSV_FIELD_FP_JA4, "csv-field-fp-ja4", 1, "Index of the CSV-Field containing the JA4 client-fingerprint.")
	flag.IntVar(&config.CSV_FIELD_UA, "csv-field-user-agent", 2, "Index of the CSV-Field containing the User-Agent.")
	flag.StringVar(&config.GEOIP_PROVIDER, "geoip-provider", "OXL", "GeoIP-ASN provider to use - supported: IPInfo-Lite, IPInfo, MaxMind, OXL (https://github.com/O-X-L/geoip-asn)")
	flag.StringVar(&config.PATH_GEOIP_ASN_DB, "geoip-asn-file", "", "Path to GeoIP-ASN MMDB-file (if not supplied - ASN flagging will be skipped).")
	flag.BoolVar(&config.LOOKUP_PTR, "lookup-ptr", false, "Enable flagging via PTR-lookups - makes processing much slower.")
	flag.BoolVar(&config.DEBUG, "debug", false, "Enable debug output.")
	flag.StringVar(&config.DEBUG_UA, "debug-user-agent", "___", "Optional User-Agent substring to show debug-output of.")
	flag.Parse()

	if config.PATH_OPEN_BOT_DATA == "" {
		fmt.Println("Error: The -data-dir argument is required.")
		flag.Usage()
		os.Exit(1)
	}

	if config.PATH_OUTPUT == "" {
		fmt.Println("Error: The -output-dir argument is required.")
		flag.Usage()
		os.Exit(1)
	}

	if config.PATH_INPUT == "" {
		fmt.Println("Error: The -input-file argument is required.")
		flag.Usage()
		os.Exit(1)
	}

	info, err := os.Stat(config.PATH_OPEN_BOT_DATA)
	if os.IsNotExist(err) {
		fmt.Printf("Error: specified data directory does not exist: %s\n", config.PATH_OPEN_BOT_DATA)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Printf("Error: specified path is not a directory: %s\n", config.PATH_OPEN_BOT_DATA)
		os.Exit(1)
	}

	config.GEOIP_PROVIDER = strings.ToLower(config.GEOIP_PROVIDER)
	if !slices.Contains(config.SUPPORTED_GEOIP_PROVIDERS, config.GEOIP_PROVIDER) {
		fmt.Printf("Error: specified GeoIP-Provider is not supported: %s\n", config.GEOIP_PROVIDER)
		os.Exit(1)
	}

	if config.PATH_GEOIP_ASN_DB != "" {
		_, err := os.Stat(config.PATH_GEOIP_ASN_DB)
		if os.IsNotExist(err) {
			fmt.Printf("Error: specified GeoIP-ASN-DB file does not exist: %s\n", config.PATH_GEOIP_ASN_DB)
			os.Exit(1)
		}
	}

	if err := internal.Run(); err != nil {
		fmt.Printf("Application run failed: %v\n", err)
		os.Exit(1)
	}
}
