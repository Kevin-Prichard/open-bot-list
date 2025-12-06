package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func downloadCategoryManifests(categories map[string]string) error {
	for category, filename := range categories {
		targetPath := filepath.Join(PATH_RUNTIME, FILE_PREFIX_MATCH+strings.ReplaceAll(filename, "/", "_"))
		url := URL_BASE_MATCHES + "/" + filename

		fmt.Printf("   - Fetching %s %s (%s)... ", category, filename, url)

		if err := DownloadFile(url, targetPath); err != nil { // Uses exported DownloadFile
			fmt.Println("FAILED")
			return fmt.Errorf("failed to download %s: %w", filename, err)
		}
		fmt.Println("SUCCESS")
	}

	return nil
}

// DownloadFile fetches a file from a URL and saves it to a local path.
// It is now exported to be used by iplist_downloader.go.
func DownloadFile(url string, filepath string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", USER_AGENT_STRING)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK HTTP status: %d", resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func DownloadManifests() error {
	if err := downloadCategoryManifests(FingerprintCategories); err != nil {
		return err
	}

	if err := downloadCategoryManifests(IPListCategories); err != nil {
		return err
	}

	if err := downloadCategoryManifests(PtrCategories); err != nil {
		return err
	}

	if err := downloadCategoryManifests(UserAgentCategories); err != nil {
		return err
	}

	// NEW: Download the actual IP lists from the manifests
	if err := DownloadIPListsFromManifests(); err != nil {
		return err
	}

	return nil
}
