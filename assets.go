package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func getAssetPath(mediaType string) string {
	randBytes := make([]byte, 32)
	rand.Read(randBytes)
	randString := base64.RawURLEncoding.EncodeToString(randBytes)
	ext := mediaTypeToExt(mediaType)
	return fmt.Sprintf("%s%s", randString, ext)
}

func getS3VideoAssetPath(mediaType string, aspectRatio string) string {
	assetPath := getAssetPath(mediaType)
	var directoryString string
	if aspectRatio == "16:9" {
		directoryString = "landscape"
	} else if aspectRatio == "9:16" {
		directoryString = "portrait"
	} else {
		directoryString = "other"
	}

	return directoryString + "/" + assetPath
}

func (cfg apiConfig) getAssetDiskPath(assetPath string) string {
	return filepath.Join(cfg.assetsRoot, assetPath)
}

func (cfg apiConfig) getAssetURL(assetPath string) string {
	return fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, assetPath)
}

func (cfg apiConfig) getS3AssetURL(assetPath string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", cfg.s3Bucket, cfg.s3Region, assetPath)
}

func mediaTypeToExt(mediaType string) string {
	parts := strings.Split(mediaType, "/")
	if len(parts) != 2 {
		return ".bin"
	}
	return "." + parts[1]
}