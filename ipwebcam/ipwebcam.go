package ipwebcam

import (
	"encoding/json"
	"io"
)

// curl http://192.168.10.177:8080/status.json
const (
	PTZ              = "ptz"
	SETTINGS         = "settings"
	ZOOM             = "zoom"
	QUALITY          = "quality"
	SET              = "set"
	ON               = "on"
	OFF              = "off"
	AUTO             = "auto"
	FFC              = "ffc"
	NIGHT_VISION     = "night_vision"
	MOTION_DETECT    = "motion_detect"
	ISO              = "iso"
	WHITEBALANCE     = "whitebalance"
	INCANDESCENT     = "incandescent"
	FLUORESCENT      = "fluorescent"
	WARM_FLUORESCENT = "warm-fluorescent"
	FLASHMODE        = "flashmode"
	RED_EYE          = "red-eye"
	ANTIBANDING      = "antibanding"
	FREQ_50HZ        = "50hz"
	FREQ_60HZ        = "60hz"
	FOCUS_DISTANCE   = "focus_distance"
	SCENEMODE        = "scenemode"
	LANDSCAPE        = "landscape"
	PORTRAIT         = "portrait"
	FACE_PRIORITY    = "face_priority"
	ACTION           = "action"
	NIGHT            = "night"
	NIGHT_PORTRAIT   = "night-portrait"
	THEATRE          = "theatre"
	BEACH            = "beach"
	SNOW             = "snow"
	SUNSET           = "sunset"
	STEADYPHOTO      = "steadyphoto"
)

type RecStatus struct {
	Manual   any `json:"manual"`
	Circular any `json:"circular"`
	Modet    any `json:"modet"`
}

type VideoStatus struct {
	Result    string    `json:"result"`
	Enabled   bool      `json:"enabled"`
	VRec      RecStatus `json:"vrec"`
	ARec      RecStatus `json:"arec"`
	VoRec     RecStatus `json:"vorec"`
	AudioOnly bool      `json:"audio_only"`
	Mode      string    `json:"mode"`
}

type CurrentValues struct {
	Orientation        string  `json:"orientation"`
	Idle               string  `json:"idle"`
	AudioOnly          string  `json:"audio_only"`
	Overlay            string  `json:"overlay"`
	Quality            string  `json:"quality"`
	FocusHoming        string  `json:"focus_homing"`
	IpAddress          string  `json:"ip_address"`
	MainPort           string  `json:"main_port"`
	Ipv6Address        string  `json:"ipv6_address"`
	MotionLimit        int64   `json:"motion_limit"`
	AdetLimit          int64   `json:"adet_limit"`
	NightVision        string  `json:"night_vision"`
	NightVisionAverage float32 `json:"night_vision_average"`
	NightVisionGain    float32 `json:"night_vision_gain"`
	VideoAcquisition   string  `json:"video_acquisition"`
	MotionDetect       string  `json:"motion_detect"`
	MotionDisplay      string  `json:"motion_display"`
	VideoChunkLen      int64   `json:"video_chunk_len"`
	GpsActive          string  `json:"gps_active"`
	VideoSize          string  `json:"video_size"`
	MirrorFlip         string  `json:"mirror_flip"`
	Ffc                string  `json:"ffc"`
	RtspVideoFormats   string  `json:"rtsp_video_formats"`
	RtspAudioFormats   string  `json:"rtsp_audio_formats"`
	VideoConnections   int64   `json:"video_connections"`
	AudioConnections   int64   `json:"audio_connections"`
	IvideonStreaming   string  `json:"ivideon_streaming"`
	Zoom               float32 `json:"zoom"`
	CropX              float32 `json:"crop_x"`
	CropY              float32 `json:"crop_y"`
	Coloreffect        string  `json:"coloreffect"`
	Scenemode          string  `json:"scenemode"`
	Focusmode          string  `json:"focusmode"`
	Whitebalance       string  `json:"whitebalance"`
	Flashmode          string  `json:"flashmode"`
	Antibanding        string  `json:"antibanding"`
	Torch              string  `json:"torch"`
	FocusDistance      float32 `json:"focus_distance"`
	FocalLength        float32 `json:"focal_length"`
	Aperture           float32 `json:"aperture"`
	FilterDensity      float32 `json:"filter_density"`
	ExposureNs         int64   `json:"exposure_ns"`
	FrameDuration      int64   `json:"frame_duration"`
	Iso                int64   `json:"iso"`
	ManualSensor       string  `json:"manual_sensor"`
	PhotoSize          string  `json:"photo_size"`
	PhotoRotation      int64   `json:"photo_rotation"`
}

type DeviceInfo struct {
	FreeSpaceBytes      string  `json:"freeSpaceBytes"`
	UserSuppliedId      string  `json:"userSuppliedId"`
	BatteryPercent      float32 `json:"batteryPercent"`
	BatteryVoltage      float32 `json:"batteryVoltage"`
	BatteryTemperatureC float32 `json:"batteryTemperatureC"`
	BatteryCharging     string  `json:"batteryCharging"`
}

type Status struct {
	VideoConnections int               `json:"video_connections"`
	AudioConnections int               `json:"audio_connections"`
	VideoStatus      VideoStatus       `json:"video_status"`
	CurrentValues    map[string]string `json:"curvals"`
	DeviceInfo       DeviceInfo        `json:"deviceInfo"`
}

type Setting struct {
	Code    string
	Value   any
	Choices []string
}

func (st *Status) Load(r io.Reader) (err error) {
	var buf []byte
	buf, err = io.ReadAll(r)
	if err != nil {
		return
	}
	err = json.Unmarshal(buf, st)
	return
}

/*
Request URL:
http://192.168.10.177:8080/settings/motion_detect?set=off
Request Method:
POST
Status Code:
200 OK
Remote Address:
192.168.10.177:8080
Referrer Policy:
strict-origin-when-cross-origin
*/

/* http://192.168.10.177:8080/ptz?zoom=15 */
// http://192.168.10.177:8080/settings/quality?set=49
// http://192.168.10.177:8080/settings/ffc?set=on
// http://192.168.10.177:8080/settings/ffc?set=off
// http://192.168.10.177:8080/settings/night_vision?set=on
// http://192.168.10.177:8080/settings/motion_detect?set=off
// http://192.168.10.177:8080/settings/iso?set=100 to 1600
// http://192.168.10.177:8080/settings/whitebalance?set=auto
// http://192.168.10.177:8080/settings/whitebalance?set=off
// http://192.168.10.177:8080/settings/whitebalance?set=incandescent
// http://192.168.10.177:8080/settings/whitebalance?set=fluorescent
// http://192.168.10.177:8080/settings/whitebalance?set=warm-fluorescent
// http://192.168.10.177:8080/settings/flashmode?set=auto
// http://192.168.10.177:8080/settings/flashmode?set=off
// http://192.168.10.177:8080/settings/flashmode?set=on
// http://192.168.10.177:8080/settings/flashmode?set=red-eye
// http://192.168.10.177:8080/settings/antibanding?set=off
// http://192.168.10.177:8080/settings/antibanding?set=auto
// http://192.168.10.177:8080/settings/antibanding?set=50hz
// http://192.168.10.177:8080/settings/antibanding?set=60hz
// http://192.168.10.177:8080/settings/focus_distance?set=0.00 to 10.00
// http://192.168.10.177:8080/settings/scenemode?set=landscape
// http://192.168.10.177:8080/settings/scenemode?set=auto
// http://192.168.10.177:8080/settings/scenemode?set=portrait
// http://192.168.10.177:8080/settings/scenemode?set=face_priority
// http://192.168.10.177:8080/settings/scenemode?set=action
// http://192.168.10.177:8080/settings/scenemode?set=night
// http://192.168.10.177:8080/settings/scenemode?set=night-portrait
// http://192.168.10.177:8080/settings/scenemode?set=theatre
// http://192.168.10.177:8080/settings/scenemode?set=beach
// http://192.168.10.177:8080/settings/scenemode?set=snow
// http://192.168.10.177:8080/settings/scenemode?set=sunset
// http://192.168.10.177:8080/settings/scenemode?set=steadyphoto
// http://192.168.10.177:8080/settings/scenemode?set=sunset
