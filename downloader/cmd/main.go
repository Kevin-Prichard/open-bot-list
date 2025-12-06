package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"git.oxl.at/open-bot-list/downloader/internal"
)

func main() {
	flag.StringVar(&internal.PATH_RUNTIME, "runtime-dir", filepath.Join(os.TempDir(), "oxl-open-bot-list"), "Runtime directory to store downloaded manifest files.")
	flag.StringVar(&internal.PATH_OUTPUT, "output-dir", "", "Output directory for processed data (required).")

	flag.Parse()

	if internal.PATH_OUTPUT == "" {
		fmt.Println("Error: The -output-dir argument is required.")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("Starting Bot List Data Acquisition...")
	if err := internal.DownloadManifests(); err != nil {
		fmt.Printf("\nFATAL ERROR: Data acquisition failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nAll manifest files downloaded successfully to %s\n", internal.PATH_RUNTIME)

	if err := os.MkdirAll(internal.PATH_OUTPUT, 0755); err != nil {
		log.Fatalln("Failed to create runtime-dir")
	}

	// --- PROCESSING STEP ---

	fmt.Println("\nStarting Manifest Processing...")

	// 1. Process category fingerprints (generates .map files)
	if err := internal.ProcessFingerprintManifests(); err != nil {
		fmt.Printf("\nFATAL ERROR: Fingerprint map processing failed: %v\n", err)
		os.Exit(1)
	}

	// 2. Process overall fingerprint manifest (generates filtered .lst files)
	if err := internal.ProcessOverallFingerprintManifests(); err != nil {
		fmt.Printf("\nFATAL ERROR: Overall Fingerprint list processing failed: %v\n", err)
		os.Exit(1)
	}

	// 3. Process overall user agent manifest (generates filtered .lst files)
	if err := internal.ProcessOverallUserAgentManifests(); err != nil {
		fmt.Printf("\nFATAL ERROR: Overall User-Agent list processing failed: %v\n", err)
		os.Exit(1)
	}

	// 4. Process individual user agent manifests (generates .map files)
	if err := internal.ProcessUserAgentManifests(); err != nil {
		fmt.Printf("\nFATAL ERROR: User-Agent map processing failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nAll manifest files processed successfully and output written to %s\n", internal.PATH_OUTPUT)

}
