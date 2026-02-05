package internal

import (
	"fmt"
	"git.oxl.at/open-bot-list/pkg/downloader/config"
	"os"

	"git.oxl.at/open-bot-list/pkg/downloader/iplist"
	"git.oxl.at/open-bot-list/pkg/downloader/manifest"
)

func Run() {
	fmt.Println("Starting Bot List Data Acquisition...")
	if err := manifest.Download(); err != nil {
		fmt.Printf("\nFATAL ERROR: Data acquisition failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nAll manifest files downloaded successfully to %s\n\nStarting Manifest Processing...\n", config.PATH_RUNTIME)

	if err := iplist.ProcessIPLists(); err != nil {
		fmt.Printf("\nFATAL ERROR: IP List processing failed: %v\n", err)
		os.Exit(1)
	}

	if err := iplist.ProcessASNLists(); err != nil {
		fmt.Printf("\nFATAL ERROR: ASN List processing failed: %v\n", err)
		os.Exit(1)
	}

	if err := manifest.ProcessFingerprintManifests(); err != nil {
		fmt.Printf("\nFATAL ERROR: Fingerprint map processing failed: %v\n", err)
		os.Exit(1)
	}

	if err := manifest.ProcessOverallFingerprintManifests(); err != nil {
		fmt.Printf("\nFATAL ERROR: Overall Fingerprint list processing failed: %v\n", err)
		os.Exit(1)
	}

	if err := manifest.ProcessOverallGenericManifests(config.UserAgentCategories, "user_agent", 4, 0, 2, "_sub"); err != nil {
		fmt.Printf("\nFATAL ERROR: Overall User-Agent list processing failed: %v\n", err)
		os.Exit(1)
	}

	if err := manifest.ProcessGenericManifests(config.UserAgentCategories, 4, 1, 2, "http_user_agent", "_sub"); err != nil {
		fmt.Printf("\nFATAL ERROR: User-Agent map processing failed: %v\n", err)
		os.Exit(1)
	}

	if err := manifest.ProcessOverallGenericManifests(config.PtrCategories, "ptr", 4, 0, 2, "_end"); err != nil {
		fmt.Printf("\nFATAL ERROR: Overall PTR list processing failed: %v\n", err)
		os.Exit(1)
	}

	if err := manifest.ProcessGenericManifests(config.PtrCategories, 4, 1, 2, "ptr", "_end"); err != nil {
		fmt.Printf("\nFATAL ERROR: PTR map processing failed: %v\n", err)
		os.Exit(1)
	}

	if err := manifest.ProcessOverallGenericManifests(config.ASNCategories, "asn", 4, 2, 0, ""); err != nil {
		fmt.Printf("\nFATAL ERROR: Overall PTR list processing failed: %v\n", err)
		os.Exit(1)
	}

	if err := manifest.ProcessGenericManifests(config.ASNCategories, 4, 1, 0, "src_asn", ""); err != nil {
		fmt.Printf("\nFATAL ERROR: ASN map processing failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nAll manifest files processed successfully and output written to %s\n", config.PATH_OUTPUT)
}
