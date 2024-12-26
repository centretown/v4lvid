package audio

type AudioSource interface {
	Record(stop chan int)
	IsEnabled() bool
}
