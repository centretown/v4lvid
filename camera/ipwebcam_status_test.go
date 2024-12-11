package camera

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

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

	var status IPWebcamVariant
	err = json.Unmarshal(buf, &status)
	if err != nil {
		t.Fatal("Unmarshal", url, err)
	}
	t.Log(status)

	var statusM IPWebcamStatus
	err = json.Unmarshal(buf, &statusM)
	if err != nil {
		t.Fatal("Unmarshal", url, err)
	}
	// t.Log(statusM)
	t.Log("statusM")

	for k, v := range statusM.ValueMap {
		r := statusM.RangeMap[k]
		t.Log(k, v, r)
	}
}

func TestLoadStatus(t *testing.T) {
	url := "http://192.168.10.92:8080"
	ipcw, err := LoadIpWebCamStatus(url)
	if err != nil {
		t.Fatal("Unmarshal", url, err)
	}
	vrs := ipcw.ValuesAndRanges()
	for i, vr := range vrs {
		t.Log(vr.Key, vr.Value, vr.Range, "IIII", i)
	}
}
