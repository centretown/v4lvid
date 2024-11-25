package web

import (
	"io"
	"log"
	"net/http"
	"path"
	"v4lvid/camera"
)

type IpwebcamHandler struct {
	ipwebcam *camera.Ipcam
}

func NewIpwebcamHandler(ipwebcam *camera.Ipcam) *IpwebcamHandler {
	return &IpwebcamHandler{
		ipwebcam: ipwebcam,
	}
}

func (handler *IpwebcamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	url := path.Join(handler.ipwebcam.Path(), r.RequestURI)
	req, err := http.NewRequest(r.Method, url, nil)
	if err != nil {
		log.Println("IpwebcamHandler NewRequest", url, err)
		return
	}

	buf, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println("IpwebcamHandler ReadAll", url, err)
		return
	}

	log.Println(url, "REQUEST RETURNED", string(buf))
}
