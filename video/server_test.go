package video

import (
	"fmt"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	webcam := NewWebcam("/dev/video0")
	config := &VideoConfig{
		Codec:  "MJPG",
		Width:  1920,
		Height: 1080,
		FPS:    30,
	}

	server := NewVideoServer(webcam, config)

	err := webcam.Open(config)
	if err != nil {
		t.Fatal(err)
	}

	if !webcam.isOpened {
		t.Fatal(fmt.Errorf("Not isOpen"))
	}

	go server.Serve()

	time.Sleep(1 * time.Second)
	server.Quit <- 1

	time.Sleep(100 * time.Millisecond)
	if server.Source.IsOpened() {
		server.Source.Close()
		t.Fatal("source still open")
	}
}
