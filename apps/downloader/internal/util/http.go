package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"git.oxl.at/open-bot-list/downloader/internal/config"
)

// DownloadFile fetches a file from a URL and saves it to a local path.
// It is now exported to be used by iplist_downloader.go.
func DownloadFileWithCache(url string, filepath string) error {
	info, err := os.Stat(filepath)
	if err == nil {
		if time.Since(info.ModTime()) < config.MIN_LIST_LIFETIME {
			fmt.Printf("SKIP (fresh < %dh) ", config.MIN_LIST_LIFETIME_HOUR)
			return nil
		}
	}
	return DownloadFile(url, filepath)
}

// DownloadFile fetches a file from a URL and saves it to a local path.
// It is now exported to be used by iplist_downloader.go.
func DownloadFile(url string, filepath string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", config.USER_AGENT_STRING)

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
