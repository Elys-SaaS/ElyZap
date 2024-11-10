package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type VideoPackage struct {
	VideoBytes  []byte
	VideoName   string
	RandomId    string
	ChunkNumber int
	ChunkCount  int
}

func main() {
	// map of randomId -> bytes[]
	videoMap := make(map[string][][]byte)
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!!!")
	})
	e.Use(middleware.CORS())
	e.POST("/upload-video", func(c echo.Context) error {
		// get bytes, id, chunk number
		var videoPackage VideoPackage
		err := c.Bind(&videoPackage)
		if err != nil {
			return c.String(400, "Error binding video package")
		}
		// if chunk number is -1 that means video is done uploading
		if videoPackage.ChunkNumber == -1 {
			// combine all chunks into one
			videoBytes := []byte{}
			for _, chunk := range videoMap[videoPackage.RandomId] {
				videoBytes = append(videoBytes, chunk...)
			}
			// delete from map
			delete(videoMap, videoPackage.RandomId)
			// save to disk
			err := saveVideo(videoBytes, videoPackage.RandomId, videoPackage.VideoName)
			if err != nil {
				return c.String(500, "Error saving video")
			}
			return c.String(200, "Video uploaded")
		}

		// insert into map if not already there
		if _, ok := videoMap[videoPackage.RandomId]; !ok {
			videoMap[videoPackage.RandomId] = make([][]byte, videoPackage.ChunkCount)
		}

		// insert into map
		videoMap[videoPackage.RandomId][videoPackage.ChunkNumber] = videoPackage.VideoBytes

		return c.JSON(200,
			map[string]interface{}{
				"message":     "Chunk uploaded",
				"chunkNumber": videoPackage.ChunkNumber,
				"chunkCount":  videoPackage.ChunkCount,
			},
		)
	})
	e.Logger.Fatal(e.Start(":8080"))
}

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
