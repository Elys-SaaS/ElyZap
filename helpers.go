package main

import (
	"fmt"
	"os"
)

func saveVideo(videoBytes []byte, randomId string, videoName string) error {
	// Define the file path with the randomId and a .mp4 extension
	filePath := fmt.Sprintf("./videos/%s.mp4", videoName+randomId)

	// Create the videos directory if it doesn't exist
	err := os.MkdirAll("./videos", os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create a file at the specified path
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write the video bytes to the file
	_, err = file.Write(videoBytes)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
