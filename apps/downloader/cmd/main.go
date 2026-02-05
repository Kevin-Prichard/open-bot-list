package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"git.oxl.at/open-bot-list/apps/downloader/internal"
	"git.oxl.at/open-bot-list/pkg/downloader/config"
)

func main() {
	fmt.Printf("\nOXL Open-Bot-List Downloader v%v\n", config.VERSION)
	fmt.Println("> © OXL IT Services / Rath Pascal")
	fmt.Println("> git.oxl.at/open-bot-list")
	fmt.Printf("> License: GPLv3\n\n")

	flag.StringVar(&config.PATH_RUNTIME, "runtime-dir", filepath.Join(os.TempDir(), "oxl-open-bot-list"), "Runtime directory to store downloaded manifest files.")
	flag.StringVar(&config.PATH_OUTPUT, "output-dir", "", "Output directory for processed data (required).")
	flag.Parse()

	if config.PATH_OUTPUT == "" {
		fmt.Println("Error: The -output-dir argument is required.")
		flag.Usage()
		os.Exit(1)
	}

	if err := os.MkdirAll(config.PATH_RUNTIME, 0750); err != nil {
		log.Fatalln("Failed to create runtime-dir")
	}

	if err := os.MkdirAll(config.PATH_OUTPUT, 0750); err != nil {
		log.Fatalln("Failed to create output-dir")
	}
	internal.Run()
}
