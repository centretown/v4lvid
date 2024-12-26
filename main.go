package main

import (
	"flag"
	"log"
	"path/filepath"
	"v4lvid/config"
	"v4lvid/namer"
	"v4lvid/web"
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

		namer.OutputBase, err = filepath.Abs(*base)
		if err != nil {
			log.Println("Abs", err)
		}

		log.Println("namer.VideoBase", namer.OutputBase)
	}

	web.Run(cfg)
}
