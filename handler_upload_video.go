package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadVideo(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "improperly configured jwt", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not validate jwt", err)
		return
	}

	fmt.Printf("uploading video: %s, by user: %s\n", videoID, userID)

	// limit of 1 GB for video size, 2^30 bytes 
	const maxMemory = 1 << 30

	videoMetadata, err := cfg.db.GetVideo(videoID)
	if videoMetadata.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "unauthorized action", err)
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not retrieve video", err)
		return
	}

	if err = r.ParseMultipartForm(maxMemory); err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not parse uploaded file", err)
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not parse uploaded file", err)
		return
	}
	defer file.Close()

	mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not parse content-type from header", err)
		return
	}
	if mediaType != "video/mp4" {
		respondWithError(w, http.StatusBadRequest, "invalid file type", err)
		return
	}

	tempFile, err := os.CreateTemp(cfg.assetsRoot, "tubely-video-*.mp4")
	if err != nil {
		log.Printf("failed to create temp file for video: %s", videoID)
		respondWithError(w, http.StatusInternalServerError, "failed to upload video", err)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	_, err = io.Copy(tempFile, file)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create video file", err)
		return
	}

	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create video file", err)
		return
	}

	fmt.Println(tempFile.Name())
	aspectRatio, err := getVideoAspectRatio(tempFile.Name())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "an error occurred processing the video", err)
		return
	}
	filepath := getS3VideoAssetPath(mediaType, aspectRatio)
	putObjectParams := s3.PutObjectInput{
		Bucket: &cfg.s3Bucket,
		Key: &filepath,
		Body: tempFile,
		ContentType: &mediaType,
	}
	cfg.s3client.PutObject(r.Context(), &putObjectParams)

	newVideoURL := cfg.getS3AssetURL(filepath)
	videoMetadata.VideoURL = &newVideoURL
	if err = cfg.db.UpdateVideo(videoMetadata); err != nil {
		respondWithError(w, http.StatusInternalServerError, "error updating video url", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, videoMetadata)
}
