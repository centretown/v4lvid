package video

// {"Format":1448695129,"Width":1280,"Height":720,"FPS":{"N":10,"D":1}}
type VideoConfig struct {
	Codec  string
	Width  int
	Height int
	FPS    uint32
}

type VideoSource interface {
	Path() string
	Open(*VideoConfig) error
	IsOpened() bool
	Close()
	Read() ([]byte, error)
	Config() *VideoConfig
}
