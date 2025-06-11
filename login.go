package main

import (
	"time"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
	"github.com/google/uuid"
)

func createJWTAndRefreshToken(cfg *apiConfig, userID uuid.UUID) (string, string, error) {
	accessToken, err := createJWT(cfg, userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := createRefreshToken(cfg, userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func createJWT(cfg *apiConfig, userID uuid.UUID) (string, error) {
	accessToken, err := auth.MakeJWT(
		userID,
		cfg.jwtSecret,
		time.Hour * 24 * 30,
	)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func createRefreshToken(cfg *apiConfig, userID uuid.UUID) (string, error) {
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		return "", err
	}

	_, err = cfg.db.CreateRefreshToken(database.CreateRefreshTokenParams{
		UserID: userID,
		Token: refreshToken,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 + 60),
	})
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}