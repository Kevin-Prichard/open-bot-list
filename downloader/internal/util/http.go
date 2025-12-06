package util

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"git.oxl.at/open-bot-list/downloader/internal/config"
)

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
