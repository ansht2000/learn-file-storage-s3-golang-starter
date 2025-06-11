package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"math"
	"os/exec"
)

var ErrFailedUnmarshallingVideoMetadata = errors.New("error getting video metadata")
var ErrNoStreamsFound = errors.New("no video streams found")
const sixteenToNineAspectRatio = 16.0/9.0
const nineToSixteenAspectRatio = 9.0/16.0
const threshold = 1.0/1000.0

func getVideoAspectRatio(filepath string) (string, error) {
	ffprobe := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filepath)
	var buffer bytes.Buffer
	// set commands output to the initialized buffer
	ffprobe.Stdout = &buffer
	err := ffprobe.Run()
	if err != nil {
		log.Printf("failed to run ffprobe to get video aspect ratio: %v\n", err)
		return "", err
	}
	ffJSON := ffJSON{}
	if err := json.Unmarshal(buffer.Bytes(), &ffJSON); err != nil {
		log.Printf("failed unmarshal: %v\n", err)
		return "", ErrFailedUnmarshallingVideoMetadata
	}

	if len(ffJSON.Streams) == 0 {
		return "", ErrNoStreamsFound
	}

	width := ffJSON.Streams[0].Width
	height := ffJSON.Streams[0].Height
	aspectRatio := float64(width) / float64(height)
	if math.Abs(aspectRatio - sixteenToNineAspectRatio) < threshold {
		return "16:9", nil
	} else if math.Abs(aspectRatio - nineToSixteenAspectRatio) < threshold {
		return "9:16", nil
	} else {
		return "other", nil
	}
}

func processVideoForFastStart(filePath string) (string, error) {
	outputFilePath := filePath + ".processing"
	ffmpeg := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", outputFilePath)
	var buf bytes.Buffer
	ffmpeg.Stdout = &buf
	err := ffmpeg.Run()
	if err != nil {
		log.Printf("failed to run ffmpeg to set the moov atom: %v\n", err)
		return "", err
	}

	return outputFilePath, nil
}