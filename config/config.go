package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"v4lvid/camera"
)

type Action struct {
	Name     string
	Title    string
	Icon     string
	HomeData bool
}

type Config struct {
	Output       string
	HttpUrl      string
	Cameras      []*camera.VideoConfig
	Actions      []*Action
	ActionsHome  []*Action
	ActionsRight []*Action
	Drivers      map[string][]*camera.ControlKey
}

func (cfg *Config) NewActionMap() (m map[string]*Action) {
	m = make(map[string]*Action)
	for _, action := range cfg.Actions {
		m[action.Name] = action
	}
	return
}

var DefaultConfig = Config{
	Output:  "/mnt/molly/output/",
	HttpUrl: "192.168.10.7:9000",
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
	Actions: []*Action{
		{Name: "camera", Title: "Camera Settings", Icon: "settings_video_camera"},
		{Name: "cameraadd", Title: "Add Camera", Icon: "linked_camera"},
		{Name: "camera_list", Title: "List Cameras", Icon: "view_list"},
		{Name: "sun", Title: "Next Sun", Icon: "wb_twilight", HomeData: true},
		{Name: "weather", Title: "Forecast Home", Icon: "routine", HomeData: true},
		{Name: "wifi", Title: "WIFI Signals", Icon: "network_wifi", HomeData: true},
		{Name: "lights", Title: "LED Lights", Icon: "backlight_high", HomeData: true},
	},
	ActionsHome: []*Action{},
	ActionsRight: []*Action{
		{Name: "chat", Title: "Chat", Icon: "chat"},
		{Name: "resetcontrols", Title: "Reset Camera", Icon: "reset_settings"},
		{Name: "record", Title: "Record", Icon: "radio_button_checked"},
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

func (cfg *Config) NewCameraServers(indicator camera.StreamIndicator) (cameraServers []*camera.Server) {
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
			camera.NewVideoServer(source, vcfg, indicator))
	}
	return
}

func (cfg *Config) NewCameraMap() map[string]*camera.VideoConfig {
	cm := make(map[string]*camera.VideoConfig)
	for _, cam := range cfg.Cameras {
		cm[cam.Path] = cam
	}
	return cm
}

func (cfg *Config) AddCamera(vc *camera.VideoConfig) {
	cfg.Cameras = append(cfg.Cameras, vc)
}
