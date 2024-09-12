package camera

import (
	"log"

	"github.com/mattn/go-mjpeg"
)

var _ VideoSource = (*Ipcam)(nil)

type Ipcam struct {
	path     string
	config   *VideoConfig
	decoder  *mjpeg.Decoder
	Buffer   []byte
	isOpened bool
}

func NewIpcam(path string) *Ipcam {
	ipc := &Ipcam{
		path: path,
	}
	return ipc
}

func (ipc *Ipcam) Path() string {
	return ipc.path
}

func (ipc *Ipcam) Config() *VideoConfig {
	return ipc.config
}

func (ipc *Ipcam) Close() {
	ipc.isOpened = false
}

func (ipc *Ipcam) IsOpened() bool {
	return ipc.isOpened
}

func (ipc *Ipcam) Open(config *VideoConfig) (err error) {
	ipc.decoder, err = mjpeg.NewDecoderFromURL(ipc.path)
	ipc.config = config
	if err != nil {
		log.Println("NewDecoderFromURL", err)
	} else {
		ipc.isOpened = true
	}
	return
}

// NewDecoderFromURL Get "http://192.168.10.30:8080": dial tcp 192.168.10.30:8080: connect: no route to host
// 2024/09/09 10:37:18 Open Error 'http://192.168.10.30:8080', Get "http://192.168.10.30:8080": dial tcp 192.168.10.30:8080: connect: no route to host
// 2024/09/09 10:37:19 Closed 'http://192.168.10.30:8080'
// panic: runtime error: invalid memory address or nil pointer dereference
// [signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x9d2bf2]

// goroutine 35 [running]:
// github.com/mattn/go-mjpeg.(*Decoder).DecodeRaw(0x0?)
//
//	/home/dave/go/pkg/mod/github.com/mattn/go-mjpeg@v0.0.3/mjpeg.go:65 +0x12
//
// v4lvid/video.(*Ipcam).Read(0xc000c0bef0?)
//
//	/home/dave/src/v4lvid/video/ipcam.go:54 +0x25
//
// v4lvid/video.(*Server).Serve(0xc000000000)
//
//	/home/dave/src/v4lvid/video/server.go:224 +0x1eb
//
// created by v4lvid/web.Serve in goroutine 1
//
//	/home/dave/src/v4lvid/web/serve.go:50 +0x208
//
// exit status 2
// dave@yeller:~/src/v4lvid$

func (ipc *Ipcam) Read() (buf []byte, err error) {
	buf, err = ipc.decoder.DecodeRaw()
	if err != nil {
		log.Println("DecodeRaw", err)
	}

	return
}

func (ipc *Ipcam) SetControl(key string, value int32) {
}
