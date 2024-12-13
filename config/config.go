package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"v4lvid/camera"
)

type ActionGroup int

const (
	Camera ActionGroup = iota
	Home
	Record
	Chat
)

type Action struct {
	Name  string
	Title string
	Icon  string
	Group ActionGroup
}

type Config struct {
	Output        string
	HttpUrl       string
	Cameras       []*camera.VideoConfig
	ActionsCamera []*Action
	ActionsHome   []*Action
	ActionsChat   []*Action
	Drivers       map[string][]*camera.ControlKey
	IPWCCommands  map[string]*camera.IPCWConfig
}

type IPWCCameraData struct {
	Action   *Action
	IPWebcam *camera.IPWebcam
}

func (cfg *Config) NewActionMap() (m map[string]*Action) {
	m = make(map[string]*Action)
	for _, action := range cfg.ActionsCamera {
		m[action.Name] = action
	}
	for _, action := range cfg.ActionsHome {
		m[action.Name] = action
	}
	for _, action := range cfg.ActionsChat {
		m[action.Name] = action
	}
	return
}

func (cfg *Config) AddCamera(vc *camera.VideoConfig) int {
	cfg.Cameras = append(cfg.Cameras, vc)
	return len(cfg.Cameras) - 1
}

func Load(filename string) (cfg *Config, err error) {
	cfg = &Config{}
	var f *os.File
	f, err = os.Open(filename)
	if err != nil {
		log.Println("config.Load Open", err)
		return
	}
	defer f.Close()
	var buf []byte
	buf, err = io.ReadAll(f)
	if err != nil {
		log.Println("config.Load ReadAll", err)
		return
	}
	err = json.Unmarshal(buf, cfg)
	if err != nil {
		log.Println("config.Load Unmarshal", err)
		return
	}
	return
}

func Save(cfg *Config, filename string) (err error) {
	var buf []byte
	buf, err = json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Println("config.Save", err)
		return
	}
	var f *os.File
	f, err = os.Create("config.json")
	if err != nil {
		log.Println("config.Save", err)
		return
	}
	defer f.Close()
	f.Write(buf)
	return
}
