package main

import "github.com/labstack/echo/v4"

func (app *App) UploadVideo(c echo.Context) error {
	// get bytes, id, chunk number
	var videoPackage VideoPackage
	err := c.Bind(&videoPackage)
	if err != nil {
		return c.String(400, "Error binding video package")
	}
	if app.assignedIds[videoPackage.RandomId] != "" && app.assignedIds[videoPackage.RandomId] != c.RealIP() {
		return c.String(400, "Random id already assigned to another")
	}
	// ip of client
	ip := c.RealIP()
	app.assignedIds[videoPackage.RandomId] = ip

	// if chunk number is -1 that means video is done uploading
	if videoPackage.ChunkNumber == -1 {
		// combine all chunks into one
		videoBytes := []byte{}
		for _, chunk := range app.videoMap[videoPackage.RandomId] {
			videoBytes = append(videoBytes, chunk...)
		}
		// delete from map
		delete(app.videoMap, videoPackage.RandomId)
		delete(app.assignedIds, videoPackage.RandomId)
		// save to disk
		err := saveVideo(videoBytes, videoPackage.RandomId, videoPackage.VideoName)
		if err != nil {
			return c.String(500, "Error saving video")
		}
		return c.String(200, "Video uploaded")
	}

	// insert into map if not already there
	if _, ok := app.videoMap[videoPackage.RandomId]; !ok {
		app.videoMap[videoPackage.RandomId] = make([][]byte, videoPackage.ChunkCount)
	}

	// insert into map
	app.videoMap[videoPackage.RandomId][videoPackage.ChunkNumber] = videoPackage.VideoBytes

	return c.JSON(200,
		map[string]interface{}{
			"message":     "Chunk uploaded",
			"chunkNumber": videoPackage.ChunkNumber,
			"chunkCount":  videoPackage.ChunkCount,
		},
	)
}
