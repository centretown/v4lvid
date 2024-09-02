package main

import (
	"flag"
	"log"
	"v4lvid/video"
	"v4lvid/web"
)

func main() {
	base := flag.String("base", "./output/", "Video file base folder")
	// base = flag.String("b", "./output/", "Video file base folder")
	flag.Parse()

	if len(*base) > 0 {
		video.VideoBase = *base + "/"
	}

	webcam := video.NewWebcam("/dev/video0")
	config := &video.VideoConfig{
		Codec:  "MJPG",
		Width:  1920,
		Height: 1080,
		FPS:    30,
	}
	iconfig := &video.VideoConfig{
		Codec:  "MJPG",
		Width:  1024,
		Height: 768,
		FPS:    2,
	}
	cserve := video.NewVideoServer(webcam, config)
	err := webcam.Open(config)
	if err != nil {
		log.Fatal(err)
	}

	ipcam := video.NewIpcam("http://192.168.0.28:8080")
	iserve := video.NewVideoServer(ipcam, iconfig)
	err = ipcam.Open(iconfig)
	if err != nil {
		log.Fatal(err)
	}

	web.Serve([]*video.Server{cserve, iserve})
}
