package config

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"text/template"
	"v4lvid/camera"
	// "honnef.co/go/tools/config"
)

func TestConfig(t *testing.T) {
	cfg, err := Load("../config.json")
	if err != nil {
		t.Fatal("Config Load", err)
	}

	for k, v := range cfg.IPWCCommands {
		t.Log(k, v)
	}
}

func TestStatus(t *testing.T) {
	var (
		url    = "http://192.168.10.92:8080/status.json?show_avail=1"
		client = &http.Client{}
		req    *http.Request
		resp   *http.Response
		buf    []byte
		err    error
	)

	req, err = http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		t.Fatal("NewRequest", url, err)
	}

	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Do Request", url, err)
	}

	buf, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("ReadAll", url, err)
	}

	t.Log(string(buf))

	var status camera.IPWebcamVariant
	err = json.Unmarshal(buf, &status)
	if err != nil {
		t.Fatal("Unmarshal", url, err)
	}
	t.Log(status)

	var statusM camera.IPWebcamStatus
	err = json.Unmarshal(buf, &statusM)
	if err != nil {
		t.Fatal("Unmarshal", url, err)
	}
	// t.Log(statusM)
	t.Log("statusM")

	for k, v := range statusM.Options {
		r := statusM.OptionMap[k]
		t.Log(k, v, r)
	}
}

const tmpltext = `
    {{range .}}
		<div class="form-entry">
			<span class="symbols-small">flare</span>
			<label for="{{.Key}}">{{.Key}}</label>
			<select name="{{.Key}}" id="{{.Key}}" class="form-input"
				hx-put="/ipwc/{{.Key}}" hx-swap=none>
				{{$eff:=.Value}}
				{{range .Options}}
					<option class="form-option" value="{{.}}"
						{{if eq $eff .}}selected{{end}}>{{.}}
					</option>
				{{end}}
			</select>
		</div>
    {{end}}
`

func TestLoadStatus(t *testing.T) {
	url := "http://192.168.10.92:8080"
	cfg, err := Load("../config.json")
	if err != nil {
		t.Fatal("Config Load", err)
	}

	ipcw := camera.NewIpWebCam()
	err = ipcw.Load(url, cfg.IPWCCommands)
	if err != nil {
		t.Fatal("Unmarshal", url, err)
	}
	tmpl := template.Must(template.New("name").Parse(tmpltext))
	err = tmpl.Execute(os.Stdout, ipcw.Properties)
	if err != nil {
		t.Fatal("Execute", url, err)
	}
}

func TestLoadTemplate(t *testing.T) {
	cfg, err := Load("../config.json")
	if err != nil {
		t.Fatal("Config Load", err)
	}

	url := "http://192.168.10.92:8080"
	ipcw := camera.NewIpWebCam()
	err = ipcw.Load(url, cfg.IPWCCommands)
	if err != nil {
		t.Fatal("Load", url, err)
	}
	const pattern = "../www/*.html"
	glob := template.Must(template.ParseGlob(pattern))
	tmpl := glob.Lookup("layout.ipwebcam")

	err = tmpl.Execute(os.Stdout,
		&IPWCCameraData{
			Action:   &Action{Name: "camera", Title: "Camera Settings", Icon: "settings_video_camera"},
			IPWebcam: ipcw,
		})
	if err != nil {
		t.Fatal("Execute", url, err)
	}
	// second time
	err = ipcw.Load(url, cfg.IPWCCommands)
	if err != nil {
		t.Fatal("Load second time", url, err)
	}

	for k, v := range ipcw.Properties {
		t.Log(k, v.Key, v.Value, v.Options, v.Command)
	}
}
