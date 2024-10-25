package camera

type CameraType int

const (
	V4L_CAMERA CameraType = iota
	IP_CAMERA
)

// {"Format":1448695129,"Width":1280,"Height":720,"FPS":{"N":10,"D":1}}
type VideoConfig struct {
	CameraType CameraType
	Path       string
	Codec      string
	Width      int
	Height     int
	FPS        uint32
}

type VideoSource interface {
	Path() string
	Open(*VideoConfig) error
	IsOpened() bool
	Close()
	Read() ([]byte, error)
}
