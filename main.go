package main

import (
	"log"
	"v4lvid/video"
	"v4lvid/web"
)

func main() {
	webcam := video.NewWebcam("/dev/video0")
	config := &video.VideoConfig{
		Codec:  "MJPG",
		Width:  1920,
		Height: 1080,
		FPS:    30,
	}
	iconfig := &video.VideoConfig{
		Codec:  "MJPG",
		Width:  1024,
		Height: 768,
		FPS:    2,
	}
	cserve := video.NewVideoServer(webcam, config)
	err := webcam.Open(config)
	if err != nil {
		log.Fatal(err)
	}

	ipcam := video.NewIpcam("http://192.168.0.28:8080")
	iserve := video.NewVideoServer(ipcam, iconfig)
	err = ipcam.Open(iconfig)
	if err != nil {
		log.Fatal(err)
	}

	web.Serve([]*video.Server{cserve, iserve})

}

// path := "/dev/video0"

// func old_main1() {
// 	info := v4l.FindDevices()
// 	for _, d := range info {
// 		cam := video.NewWebcam(d.Path)
// 		for _, k := range cam.Configs {
// 			buf, _ := json.Marshal(&k)
// 			log.Println(string(buf))
// 		}
// 	}

// }

// func old_main() {
// 	// var devices []*v4l.Device = make([]*v4l.Device, 0)
// 	paths := make([]string, 6)
// 	for i := range paths {
// 		paths[i] = fmt.Sprintf("/dev/video%d", i)
// 	}

// 	for _, path := range paths {
// 		// path := fmt.Sprintf(path)
// 		device, err := v4l.Open(path)
// 		if err != nil {
// 			log.Println("Open", path, err)
// 			continue
// 		}

// 		showInfo(device, path)

// 	}
// }

// func showInfo(device *v4l.Device, path string) {
// 	defer device.Close()
// 	var (
// 		formatBytes []byte = make([]byte, 4)
// 	)

// 	showConfig := func(c v4l.DeviceConfig) {
// 		formatBytes[0] = byte(c.Format)
// 		formatBytes[1] = byte(c.Format >> 8)
// 		formatBytes[2] = byte(c.Format >> 16)
// 		formatBytes[3] = byte(c.Format >> 24)
// 		log.Printf("FPS=%d,Format=%x %s Width=%d Height=%d\n",
// 			c.FPS,
// 			c.Format, string(formatBytes),
// 			c.Width,
// 			c.Height)
// 	}

// 	info, err := device.DeviceInfo()
// 	if err != nil {
// 		log.Println("DeviceInfo", path, err)
// 		return
// 	}
// 	log.Printf("%s %s %s [%d.%d.%d] %s\n", path,
// 		info.DeviceName,
// 		info.DriverName,
// 		info.DriverVersion[0],
// 		info.DriverVersion[1],
// 		info.DriverVersion[2],
// 		info.BusInfo,
// 	)

// 	configs, err := device.ListConfigs()
// 	if err != nil {
// 		log.Println("ListConfigs", path, err)
// 		return
// 	}

// 	var newConfig v4l.DeviceConfig
// 	for _, c := range configs {
// 		showConfig(c)
// 		if c.Width == 1920 && c.Height == 1080 && c.FPS.N == 30 {
// 			newConfig = c
// 		}
// 	}

// 	config, err := device.GetConfig()
// 	if err != nil {
// 		log.Println("GetConfig", path, err)
// 		return
// 	}

// 	log.Println("Current Config:")
// 	showConfig(config)

// 	err = device.SetConfig(newConfig)
// 	if err != nil {
// 		log.Println("SetConfig", path, err)
// 		return
// 	}

// 	log.Println("New Config:")
// 	showConfig(newConfig)

// 	controls, err := device.ListControls()
// 	if err != nil {
// 		log.Println("ListControls", path, err)
// 		return
// 	}

// 	log.Println("Controls:")
// 	for _, c := range controls {
// 		log.Printf("CID=%v, Name=%v, Type=%v, Default=%v, Max=%v, Min=%v, Step=%v\n",
// 			c.CID, c.Name, c.Type, c.Default, c.Max, c.Min, c.Step)
// 	}

// 	err = device.TurnOn()
// 	if err != nil {
// 		log.Println("TurnOn", path, err)
// 		return
// 	}

// 	buf, err := device.Capture()
// 	if err != nil {
// 		log.Println("Capture", path, err)
// 		return
// 	}

// 	log.Println("Buffer Length", buf.Len())

// }
