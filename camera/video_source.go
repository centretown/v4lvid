package camera

type VideoSource interface {
	Open(*VideoConfig) error
	IsOpened() bool
	Close()
	Path() string
	Read() ([]byte, error)
}
