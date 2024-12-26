package audio

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	"v4lvid/namer"

	"github.com/gordonklaus/portaudio"
)

type AudioMgr struct {
	Current *portaudio.DeviceInfo
	Enabled bool
}

func NewAudio() (au *AudioMgr) {
	au = &AudioMgr{}
	err := portaudio.Initialize()
	if err != nil {
		log.Println("Initialize AudioMgr", err)
	} else {
		au.Current, _ = portaudio.DefaultInputDevice()
	}
	return
}

func (au *AudioMgr) CurrentDevice() (device *portaudio.DeviceInfo, err error) {
	device, err = portaudio.DefaultInputDevice()
	return
}

type FindFlag int

const (
	FindAll FindFlag = iota
	FindPrefix
	FindCase
)

func (au *AudioMgr) findDevices(search string) (result []*portaudio.DeviceInfo) {
	result = make([]*portaudio.DeviceInfo, 0)
	hostApis, err := portaudio.HostApis()
	if err != nil {
		return
	}
	for _, api := range hostApis {
		for _, dvc := range api.Devices {
			s := strings.ToLower(dvc.Name)
			if strings.Contains(s, search) {
				result = append(result, dvc)
			}
		}
	}
	return
}
func (au *AudioMgr) ListAllDevices() (list []*portaudio.DeviceInfo) {
	list = make([]*portaudio.DeviceInfo, 0)
	hostApis, err := portaudio.HostApis()
	if err != nil {
		return
	}
	for _, api := range hostApis {
		list = append(list, api.Devices...)
	}
	return
}

func (au *AudioMgr) FindDevices(searches ...string) (list []*portaudio.DeviceInfo) {
	if len(searches) < 1 {
		list = au.ListAllDevices()
		return
	}
	list = make([]*portaudio.DeviceInfo, 0)
	for _, search := range searches {
		search = strings.ToLower(search)
		lst := au.findDevices(search)
		if len(lst) == 0 {
			continue
		}
		list = append(list, lst...)
	}
	return
}

func (au *AudioMgr) FindDevice(search string) (device *portaudio.DeviceInfo, err error) {
	search = strings.ToLower(search)
	var (
		hostApis []*portaudio.HostApiInfo
	)
	hostApis, err = portaudio.HostApis()
	if err != nil {
		return
	}
	for _, api := range hostApis {
		for _, dvc := range api.Devices {
			s := strings.ToLower(dvc.Name)
			if strings.HasPrefix(s, search) {
				device = dvc
				return
			}
		}
	}

	err = fmt.Errorf("device not found")
	return
}

var _ AudioSource = (*AudioMgr)(nil)

func (au *AudioMgr) IsEnabled() bool { return au.Enabled }

func (au *AudioMgr) Record(stop chan int) {
	fname, _ := namer.NextFileName(namer.OutputBase, "aiff")
	au.RecordX(au.Current, fname, stop)
}

func (au *AudioMgr) RecordX(device *portaudio.DeviceInfo, fileName string, stop chan int) {
	log.Println("RecordX")
	var (
		file   *os.File
		err    error
		stream *portaudio.Stream
	)

	defer func() {
		if err != nil {
			log.Println("Recording Error:", err)
		} else {
			log.Println("Recording Stopped.")
		}
	}()

	if !strings.HasSuffix(fileName, ".aiff") {
		fileName += ".aiff"
	}

	file, err = os.Create(fileName)
	if err != nil {
		return
	}

	// log.Println("InitAIFF", device.DefaultSampleRate, int16(device.MaxInputChannels))
	err = InitAIFF(file, device.DefaultSampleRate, int16(1))
	if err != nil {
		return
	}

	sampleCount := 0
	defer func() {
		finalizeAIFF(file, sampleCount)
	}()

	inbuf := make([]int32, 64)
	var param = portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   device,
			Channels: 1,
			Latency:  device.DefaultHighInputLatency,
		},
		Output: portaudio.StreamDeviceParameters{
			Device: nil,
		},
		SampleRate:      device.DefaultSampleRate,
		FramesPerBuffer: len(inbuf),
		// Flags           StreamFlags
	}
	stream, err = portaudio.OpenStream(param, inbuf)
	// stream, err = portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
	if err != nil {
		return
	}
	log.Println("SampleRate", stream.Info().SampleRate)
	defer stream.Close()

	err = stream.Start()
	if err != nil {
		return
	}

	for {
		select {
		case <-stop:
			err = stream.Stop()
			return
		default:
			err = stream.Read()
			if err != nil {
				return
			}
			err = binary.Write(file, binary.BigEndian, inbuf)
			if err != nil {
				return
			}
			sampleCount += len(inbuf)
		}
	}

}
