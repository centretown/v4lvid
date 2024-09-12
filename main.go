package main

import (
	"flag"
	"log"
	"path/filepath"
	"v4lvid/camera"
	"v4lvid/web"
)

func main() {
	base := flag.String("base", "./output/", "Video file base folder")
	// base = flag.String("b", "./output/", "Video file base folder")
	flag.Parse()

	if len(*base) > 0 {
		var err error

		camera.VideoBase, err = filepath.Abs(*base)
		if err != nil {
			log.Println("Abs", err)
		}

		log.Println("video.VideoBase", camera.VideoBase)
	}

	webcam := camera.NewWebcam("/dev/video0")
	config := &camera.VideoConfig{
		Codec:  "MJPG",
		Width:  1920,
		Height: 1080,
		FPS:    30,
	}
	iconfig := &camera.VideoConfig{
		Codec:  "MJPG",
		Width:  1024,
		Height: 768,
		FPS:    2,
	}
	cserve := camera.NewVideoServer(webcam, config)
	err := webcam.Open(config)
	if err != nil {
		log.Fatal(err)
	}

	ipcam := camera.NewIpcam("http://192.168.10.30:8080")
	iserve := camera.NewVideoServer(ipcam, iconfig)
	err = ipcam.Open(iconfig)
	if err != nil {
		log.Println(err)
	}

	web.Serve([]*camera.Server{cserve, iserve})
}
