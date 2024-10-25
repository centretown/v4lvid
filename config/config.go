package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"v4lvid/camera"
)

type Action struct {
	Name string
	Icon string
}

type Config struct {
	Output  string
	Cameras []*camera.VideoConfig
	HttpUrl string
	Actions []*Action
	WsUrl   string
	Drivers map[string][]*camera.ControlKey
}

var DefaultConfig = Config{
	Output: "/mnt/molly/output/",
	Cameras: []*camera.VideoConfig{
		{
			CameraType: camera.V4L_CAMERA,
			Path:       "/dev/video0",
			Codec:      "MJPG",
			Width:      1920,
			Height:     1080,
			FPS:        30,
		},
		{
			CameraType: camera.IP_CAMERA,
			Path:       "http://192.168.10.30:8080",
			Codec:      "MJPG",
			Width:      1024,
			Height:     768,
			FPS:        2,
		},
	},
	HttpUrl: "192.168.10.7:9000",
	Actions: []*Action{
		{Name: "camera", Icon: "settings_video_camera"},
		{Name: "sun", Icon: "wb_twilight"},
		{Name: "weather", Icon: "routine"},
		{Name: "wifi", Icon: "network_wifi"},
		{Name: "lights", Icon: "backlight_high"},
	},

	Drivers: map[string][]*camera.ControlKey{
		// DeviceName NexiGo N660 FHD Webcam: NexiGo  DriverName uvcvideo
		"uvcvideo": {
			{Key: "Zoom, Absolute", Controls: []*camera.Control{
				{Url: "/zoomin", Icon: "zoom_in", Multiplier: 1},
				{Url: "/zoomout", Icon: "zoom_out", Multiplier: -1},
			}},

			{Key: "Pan, Absolute", Controls: []*camera.Control{
				{Url: "/panleft", Icon: "arrow_back", Multiplier: -1},
				{Url: "/panright", Icon: "arrow_forward", Multiplier: 1},
			}},

			{Key: "Tilt, Absolute", Controls: []*camera.Control{
				{Url: "/tiltup", Icon: "arrow_upward", Multiplier: 1},
				{Url: "/tiltdown", Icon: "arrow_downward", Multiplier: -1},
			}},

			{Key: "Brightness", Controls: []*camera.Control{
				{Url: "/brightnessup", Icon: "brightness_high", Multiplier: 10},
				{Url: "/brightnessdown", Icon: "brightness_low", Multiplier: -10},
			}},

			{Key: "Contrast", Controls: []*camera.Control{
				{Url: "/contrastup", Icon: "contrast_square", Multiplier: 10},
				{Url: "/contrastdown", Icon: "exposure", Multiplier: -10},
			}},

			{Key: "Saturation", Controls: []*camera.Control{
				{Url: "/saturationup", Icon: "backlight_high", Multiplier: 10},
				{Url: "/saturationdown", Icon: "backlight_low", Multiplier: -10},
			}},
		}},
	WsUrl: "192.168.10.7:9900",
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

func (cfg *Config) NewCameraServers() (cameraServers []*camera.Server) {
	cameraServers = make([]*camera.Server, 0)
	var (
		err    error
		source camera.VideoSource
	)
	for _, vcfg := range cfg.Cameras {
		switch vcfg.CameraType {
		case camera.V4L_CAMERA:
			source = camera.NewWebcam(vcfg.Path)
		case camera.IP_CAMERA:
			source = camera.NewIpcam(vcfg.Path)
		default:
			continue
		}
		err = source.Open(vcfg)
		if err != nil {
			log.Println(err)
			continue
		}
		cameraServers = append(cameraServers,
			camera.NewVideoServer(source, vcfg))
	}
	return
}
