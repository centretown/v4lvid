package main

import (
	"flag"
	"log"
	"path/filepath"
	"v4lvid/camera"
	"v4lvid/config"
	"v4lvid/web"
)

func main() {
	// cfg := &config.DefaultConfig
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatal("config.Load", err)
	}
	base := flag.String("output", cfg.Output, "Output folder")
	flag.Parse()

	if len(*base) > 0 {
		var err error

		camera.VideoBase, err = filepath.Abs(*base)
		if err != nil {
			log.Println("Abs", err)
		}

		log.Println("video.VideoBase", camera.VideoBase)
	}

	web.Serve(cfg)
}
