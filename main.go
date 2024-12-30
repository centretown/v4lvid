package main

import (
	"flag"
	"log"
	"path/filepath"
	"v4lvid/config"
	"v4lvid/web"

	"github.com/centretown/avcam"
)

func main() {
	flag.Parse()
	// cfg := &config.DefaultConfig
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatal("config.Load", err)
	}
	base := flag.String("output", cfg.Output, "Output folder")

	if len(*base) > 0 {
		var err error

		avcam.OutputBase, err = filepath.Abs(*base)
		if err != nil {
			log.Println("Abs", err)
		}

		log.Println("camera.VideoBase", avcam.OutputBase)
	}

	web.Run(cfg)
}
