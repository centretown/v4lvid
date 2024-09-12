package camera

import (
	"testing"
	"time"
)

func TestFileName(t *testing.T) {
	u := time.Now()
	foldername := u.Format("2006-01-02")
	filename := u.Format("15:04:05-0700")
	t.Logf("foldername '%s' filename '%s'\n", foldername, filename)
}

func TestNextFileName(t *testing.T) {
	name, err := NextFileName(VideoBase, "mp4")
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	t.Log(name)
}
