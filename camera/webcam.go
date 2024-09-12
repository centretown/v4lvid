package camera

import (
	"log"
	"strings"

	"github.com/korandiz/v4l"
)

var _ VideoSource = (*Webcam)(nil)

type Webcam struct {
	path     string
	config   *VideoConfig
	Device   *v4l.Device
	Info     v4l.DeviceInfo
	Buffer   []byte
	Controls map[string]v4l.ControlInfo
	Configs  []v4l.DeviceConfig
	// Config     v4l.DeviceConfig
	isOpened bool
}

func NewWebcam(path string) *Webcam {
	cam := &Webcam{
		path:     path,
		Buffer:   make([]byte, 0),
		Controls: make(map[string]v4l.ControlInfo, 0),
		Configs:  make([]v4l.DeviceConfig, 0),
	}
	return cam
}

func (cam *Webcam) Path() string {
	return cam.path
}

func (cam *Webcam) Config() *VideoConfig {
	return cam.config
}

func (cam *Webcam) Open(config *VideoConfig) error {

	device, err := v4l.Open(cam.path)
	if err != nil {
		log.Println("Open", cam.path, err)
		return err
	}
	cam.config = config
	cam.isOpened = true
	cam.Device = device
	cam.Info, err = device.DeviceInfo()
	if err != nil {
		log.Println("DeviceInfo", cam.path, err)
		return err
	}

	err = cam.LoadConfigs()
	if err != nil {
		log.Println("LoadConfigs", cam.path, err)
		return err
	}

	controls, err := device.ListControls()
	if err != nil {
		log.Println("ListControls", cam.path, err)
		return err
	}

	log.Println("Controls:")
	for _, c := range controls {
		cam.Controls[strings.ToLower(c.Name)] = c
		log.Printf("CID='%v', Name='%s', Type=%v, Default=%v, Max=%v, Min=%v, Step=%v\n",
			c.CID, c.Name, c.Type, c.Default, c.Max, c.Min, c.Step)
	}

	// cam.SetControl("brightness", 128)
	// cam.SetControl("Pan, Absolute", 0)
	// cam.SetControl("Tilt, Absolute", 0)
	// cam.SetControl("Zoom, Absolute", 10)

	err = cam.SetConfig(config)
	if err != nil {
		log.Println("SetConfig", cam.path, err)
		return err
	}

	// device.Close()
	// time.Sleep(time.Millisecond * 100)

	// device, err = v4l.Open(cam.path)
	// if err != nil {
	// 	log.Println("Open", cam.path, err)
	// 	return err
	// }
	// cam.Device = device
	// cam.Device.TurnOn()
	// cam.Close()

	// device, err = v4l.Open(cam.path)
	// if err != nil {
	// 	log.Println("Open", cam.path, err)
	// 	return err
	// }
	// cam.Device = device
	return nil
}

func (cam *Webcam) SetConfig(videoParams *VideoConfig) (err error) {
	cam.Device.TurnOff()

	preferred := &v4l.DeviceConfig{
		Format: ToFourCC(videoParams.Codec),
		Width:  videoParams.Width, Height: videoParams.Height,
		FPS: v4l.Frac{N: videoParams.FPS, D: 1},
	}

	err = cam.Device.SetConfig(*cam.findConfig(preferred))
	if err != nil {
		log.Println("SetConfig", cam.path, err)
		return
	}

	bufferInfo, err := cam.Device.BufferInfo()
	if err != nil {
		log.Println("BufferInfo", cam.path, err)
		return
	}

	cam.Buffer = make([]byte, bufferInfo.BufferSize)

	err = cam.Device.TurnOn()
	if err != nil {
		log.Println("TurnOn", cam.path, err)
		return
	}

	return
}

func (cam *Webcam) GetConfig() (config v4l.DeviceConfig) {
	config, _ = cam.Device.GetConfig()
	return
}

// CID=9963776, Name=Brightness, Type=int, Default=128, Max=255, Min=0, Step=1
// CID=9963777, Name=Contrast, Type=int, Default=128, Max=255, Min=0, Step=1
// CID=9963778, Name=Saturation, Type=int, Default=128, Max=255, Min=0, Step=1
// CID=9963788, Name=White Balance, Automatic, Type=bool, Default=1, Max=1, Min=0, Step=1
// CID=9963795, Name=Gain, Type=int, Default=5, Max=100, Min=0, Step=1
// CID=9963800, Name=Power Line Frequency, Type=enum, Default=1, Max=2, Min=0, Step=1
// CID=9963802, Name=White Balance Temperature, Type=int, Default=4650, Max=6500, Min=2600, Step=1
// CID=9963803, Name=Sharpness, Type=int, Default=128, Max=255, Min=0, Step=1
// CID=10094849, Name=Auto Exposure, Type=enum, Default=0, Max=3, Min=0, Step=1
// CID=10094850, Name=Exposure Time, Absolute, Type=int, Default=100, Max=6500, Min=0, Step=1
// CID=10094856, Name=Pan, Absolute, Type=int, Default=0, Max=36000, Min=-36000, Step=3600
// CID=10094857, Name=Tilt, Absolute, Type=int, Default=0, Max=36000, Min=-36000, Step=3600
// CID=10094860, Name=Focus, Automatic Continuous, Type=bool, Default=0, Max=1, Min=0, Step=1
// CID=10094861, Name=Zoom, Absolute, Type=int, Default=10, Max=20, Min=10, Step=1

func (cam *Webcam) GetControlInfo(key string) (info v4l.ControlInfo, err error) {
	control, ok := cam.Controls[strings.ToLower(key)]
	if !ok {
		log.Println("unknown control", key)
		return
	}
	return cam.Device.ControlInfo(control.CID)
}

func (cam *Webcam) GetControlValue(key string) (value int32) {
	control, ok := cam.Controls[strings.ToLower(key)]
	if !ok {
		log.Println("unknown control", key, value)
		return
	}

	value, err := cam.Device.GetControl(control.CID)
	if err != nil {
		log.Println("GetControl", key, value, err)
		return
	}

	return
}

func (cam *Webcam) SetValue(key string, value int32) {
	control, ok := cam.Controls[strings.ToLower(key)]
	if !ok {
		log.Println("unknown control", key, value)
		return
	}

	err := cam.Device.SetControl(control.CID, value)
	if err != nil {
		log.Println("SetControl", key, value, err)
		return
	}

	log.Println("SetControl", key, value)
}

func (cam *Webcam) LoadConfigs() (err error) {
	cam.Configs, err = cam.Device.ListConfigs()
	if err != nil {
		log.Println("ListConfigs", err)
		return err
	}

	for _, c := range cam.Configs {
		log.Println(FourCC(c.Format), c.Width, c.Height, c.FPS.N)
	}
	return
}

func (cam *Webcam) IsOpened() bool {
	return cam.isOpened
}

func (cam *Webcam) Close() {
	cam.Device.TurnOff()
	cam.Device.Close()
	cam.isOpened = false
}

func (cam *Webcam) Read() (buf []byte, err error) {
	buf = cam.Buffer
	var (
		vbuf  *v4l.Buffer
		count int
	)
	vbuf, err = cam.Device.Capture()
	if err != nil {
		log.Println("Webcam Capture", err)
		return
	}

	count, err = vbuf.Read(buf)
	if err != nil {
		log.Println("Webcam Read", err)
		return
	}
	// log.Println(count, "bytes read")
	buf = buf[:count]
	return
}

func (cam *Webcam) findConfig(b *v4l.DeviceConfig) *v4l.DeviceConfig {
	var (
		selected int
		lowest   int = 1_000_000
		score    int
	)

	for i := range cam.Configs {
		score = scoreConfig(b, &cam.Configs[i])
		if score < lowest {
			selected = i
			lowest = score
		}
	}
	// fmt.Println("lowest", lowest, "selected", selected)
	return &cam.Configs[selected]
}

func scoreConfig(a, b *v4l.DeviceConfig) (score int) {
	abs := func(a int) int {
		if a < 0 {
			return -a
		}
		return a
	}
	if a.Format != b.Format {
		score += 100
	}
	if a.Width != b.Width {
		score += abs(a.Width - b.Width)
	}
	if a.Height != b.Height {
		score += abs(a.Height - b.Height)
	}
	if a.FPS != b.FPS {
		score += abs(int(a.FPS.N) - int(b.FPS.N))
	}
	return
}
