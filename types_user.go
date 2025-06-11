package main

import "github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"

type userParameters struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type userResponse struct {
	database.User
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}