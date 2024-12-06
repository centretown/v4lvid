package config

import "v4lvid/camera"

var DefaultConfig = Config{
	Output:  "/mnt/molly/output/",
	HttpUrl: "192.168.10.7:9000",
	Cameras: []*camera.VideoConfig{
		{
			CameraType: camera.LOCAL_CAMERA,
			Path:       "/dev/video0",
			Codec:      "MJPG",
			Width:      1920,
			Height:     1080,
			FPS:        30,
		},
		{
			CameraType: camera.REMOTE_CAMERA,
			Path:       "http://192.168.10.30:8080",
			Codec:      "MJPG",
			Width:      1024,
			Height:     768,
			FPS:        2,
		},
	},
	ActionsCamera: []*Action{
		{Name: "camera", Title: "Camera Settings", Icon: "settings_video_camera"},
		{Name: "cameraadd", Title: "Add Camera", Icon: "linked_camera"},
		{Name: "camera_list", Title: "List Cameras", Icon: "view_list"},
	},
	ActionsHome: []*Action{
		{Name: "sun", Title: "Next Sun", Icon: "wb_twilight", Group: Home},
		{Name: "weather", Title: "Forecast Home", Icon: "routine", Group: Home},
		{Name: "wifi", Title: "WIFI Signals", Icon: "network_wifi", Group: Home},
		{Name: "lights", Title: "LED Lights", Icon: "backlight_high", Group: Home},
	},
	ActionsChat: []*Action{
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
