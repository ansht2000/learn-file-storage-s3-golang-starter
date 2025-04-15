package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}


	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	// TODO: implement the upload here
	const maxMemory = 10 << 20

	if err = r.ParseMultipartForm(maxMemory); err != nil {
		respondWithError(w, http.StatusBadRequest, "could not parse uploaded file", err)
		return
	}

	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to parse uploaded file", err)
		return
	}
	defer file.Close()

	mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not parse content-type from header", err)
		return
	}
	if mediaType != "image/jpg" && mediaType != "image/png" {
		respondWithError(w, http.StatusBadRequest, "invalid file type", err)
		return
	}

	filepath := getAssetPath(mediaType)
	assetDiskPath := cfg.getAssetDiskPath(filepath)
	createdFile, err := os.Create(assetDiskPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create asset file", err)
		return
	}
	defer createdFile.Close()
	_, err = io.Copy(createdFile, file)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create asset file", err)
		return
	}

	videoMetaData, err := cfg.db.GetVideo(videoID)
	if videoMetaData.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "unauthorized action", err)
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not retrieve video metadata", err)
		return
	}
	
	newThumbnailURL := cfg.getAssetURL(filepath)
	cfg.cleanupPreviousThumbnail(videoMetaData.ThumbnailURL)
	videoMetaData.ThumbnailURL = &newThumbnailURL
	cfg.db.UpdateVideo(videoMetaData)

	respondWithJSON(w, http.StatusOK, videoMetaData)
}
