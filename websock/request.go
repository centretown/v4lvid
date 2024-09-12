package websock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const host = "http://melon:8123"
const api = "/api/"

func Get(cmd string) (buf []byte, err error) {

	log.Printf("GET %s", cmd)
	req, err := http.NewRequest("GET", host+api+cmd, nil)
	if err != nil {
		log.Println(err, "GET")
		return
	}

	buf, err = Request(req)
	fmt.Println(string(buf))
	return
}

func Post(cmd string, body string) ([]byte, error) {
	log.Printf("POST %s %s\n", cmd, body)
	buf := bytes.NewBuffer(([]byte)(body))

	req, err := http.NewRequest("POST", host+api+cmd, buf)
	if err != nil {
		log.Println(err, "POST")
		return buf.Bytes(), err
	}

	req.Header.Add("Content-Type", "application/json")
	return Request(req)
}

func Request(req *http.Request) (buf []byte, err error) {
	client := &http.Client{}
	req.Header.Add("Authorization", "Bearer "+Token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err, "PROCESS")
		return
	}

	buf, err = io.ReadAll(resp.Body)

	if err != nil {
		if err.Error() != "EOF" {
			log.Println(err, "READ")
			return
		}
	}

	var v any
	err = json.Unmarshal(buf, &v)
	if err != nil {
		log.Println(string(buf), err, "UNMARSHAL")
		return
	}

	return
}
