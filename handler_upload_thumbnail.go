package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

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

	mediaType := header.Header.Get("Content-Type")
	imageData, err := io.ReadAll(file)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not read from uploaded file", err)
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

	newThumbnail := thumbnail{
		data: imageData,
		mediaType: mediaType,
	}
	imgDataEncode := base64.StdEncoding.EncodeToString([]byte(newThumbnail.data))
	thumbnailDataURL := fmt.Sprintf("data:%s;base64,%s", newThumbnail.mediaType, imgDataEncode)
	videoMetaData.ThumbnailURL = &thumbnailDataURL
	cfg.db.UpdateVideo(videoMetaData)

	respondWithJSON(w, http.StatusOK, videoMetaData)
}
