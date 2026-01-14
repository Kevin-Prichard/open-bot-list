package manifest

import (
	"fmt"
	"path/filepath"
	"strings"

	"git.oxl.at/open-bot-list/downloader/internal/config"
	"git.oxl.at/open-bot-list/downloader/internal/iplist"
	"git.oxl.at/open-bot-list/downloader/internal/util"
)

func downloadCategoryManifests(categories map[string]string) error {
	for _, filename := range categories {
		targetPath := filepath.Join(config.PATH_RUNTIME, config.FILE_PREFIX_MATCH+strings.ReplaceAll(filename, "/", "_"))
		url := config.URL_BASE_MATCHES + "/" + filename

		fmt.Printf("   - Fetching %s (%s)... ", filename, url)

		if err := util.DownloadFile(url, targetPath); err != nil {
			fmt.Println("FAILED")
			return fmt.Errorf("failed to download %s: %w", filename, err)
		}
		fmt.Println("SUCCESS")
	}

	return nil
}

func Download() error {
	if err := downloadCategoryManifests(config.FingerprintCategories); err != nil {
		return err
	}

	if err := downloadCategoryManifests(config.IPListCategories); err != nil {
		return err
	}

	if err := downloadCategoryManifests(config.ASNListCategories); err != nil {
		return err
	}

	if err := downloadCategoryManifests(config.PtrCategories); err != nil {
		return err
	}

	if err := downloadCategoryManifests(config.UserAgentCategories); err != nil {
		return err
	}

	if err := downloadCategoryManifests(config.ASNCategories); err != nil {
		return err
	}

	if err := iplist.DownloadListsFromManifests(config.IPListCategories, config.FILE_PREFIX_IPLIST); err != nil {
		return err
	}

	if err := iplist.DownloadListsFromManifests(config.ASNListCategories, config.FILE_PREFIX_ASNLIST); err != nil {
		return err
	}

	return nil
}
