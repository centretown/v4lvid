package ipwebcam

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	files := []string{
		"ao3s.json",
		"s5.json",
	}

	for _, fn := range files {
		t.Log(fn)
		f, err := os.Open(fn)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		st := &Status{}
		err = st.Load(f)
		if err != nil {
			t.Fatal(err)
		}

		for k, v := range st.CurrentValues {
			t.Log(k, v)
		}
	}
}
