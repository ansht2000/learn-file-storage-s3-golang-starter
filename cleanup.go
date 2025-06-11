package main

import (
	"log"
	"os"
	"path"
	"strings"
)

func getFilenameFromURL(URL string) string {
	parts := strings.Split(URL, "/")
	return parts[len(parts) - 1]
}

func (cfg *apiConfig) cleanupPreviousThumbnail(thumbnailURL *string) {
	// check for both cases where thumbnail is an empty string
	// or thumbnail is not set
	if thumbnailURL != nil && *thumbnailURL != "" {
		filename := getFilenameFromURL(*thumbnailURL)
		filepath := path.Join(cfg.assetsRoot, filename)
		if err := os.Remove(filepath); err != nil {
			log.Printf("failed to remove old thumbail %s: %v", filepath, err)
		} else {
			log.Printf("successfully removed old thumbnail: %s", filepath)
		}
	}
}