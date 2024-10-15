package ha

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"slices"
	"strings"
	"time"
	"v4lvid/websock"
)

const (
	AuthCommand      = `{ "type":"auth", "access_token":"%s" }`
	ConfigCommand    = `{ "type":"get_config", "id":%d }`
	StatesCommand    = `{ "type":"get_states", "id":%d }`
	SubscribeCommand = `{ "type":"subscribe_events", "event_type":"state_changed", "id":%d }`
)

type HomeData struct {
	Entities   EntityMap
	EntityKeys []string

	Err           error
	monitoring    bool
	stop          chan int
	loadStatesID  int
	subscriptions map[string][]*Subscription
	sock          *websock.WebSockClient
}

func NewHomeData() *HomeData {
	var data = &HomeData{
		Entities:      make(map[string]*Entity[json.RawMessage]),
		stop:          make(chan int),
		subscriptions: make(map[string][]*Subscription),
	}

	var err error
	data.sock, err = websock.NewWebSockClient()
	if err != nil {
		log.Println("NewHomeData", err)
	}
	return data
}

func (data *HomeData) Subscribe(entityID string, subscription *Subscription) {
	list, ok := data.subscriptions[entityID]
	if !ok {
		list = make([]*Subscription, 1)
		list[0] = subscription
	} else {
		list = append(list, subscription)
	}
	data.subscriptions[entityID] = list
}

func (data *HomeData) Consume(entityID string, newState *Entity[json.RawMessage]) {
	subs, ok := data.subscriptions[entityID]
	if ok {
		for _, sub := range subs {
			if sub.Enabled {
				sub.Consume(newState)
			}
		}
	}
}

func (data *HomeData) EntityList(filters ...string) (list []string) {
	list = make([]string, 0)
	all := len(filters) == 0
	// list = make([]string, 0, len(data.Entities))
	for k := range data.Entities {
		if all {
			list = append(list, k)
			continue
		}
		for _, s := range filters {
			if strings.HasPrefix(k, s) {
				list = append(list, k)
				break
			}
		}
	}
	slices.Sort(list)
	return
}

func (data *HomeData) CallService(cmd string) {
	data.sock.WriteCommandID(cmd)
}

func (data *HomeData) StopMonitor() {
	if data.monitoring {
		log.Println("StopMonitor")
		data.stop <- 1
	}
}

func (data *HomeData) Authorize() (ok bool, err error) {
	var (
		result AuthResult
		buf    []byte
		max    = 4
	)

	cmd := fmt.Sprintf(AuthCommand, websock.Token)

	for i := 0; i < max && !ok; i += 1 {
		err = data.sock.WriteCommand(cmd)
		if err != nil {
			return
		}
		buf, err = data.sock.Read()
		if err != nil {
			return
		}

		err = json.Unmarshal(buf, &result)
		if err != nil {
			return
		}

		ok = result.Type == "auth_ok"
	}

	return
}

func (data *HomeData) BuildEntities() (err error) {

	data.loadStatesID, err = data.sock.WriteCommandID(StatesCommand)
	if err != nil {
		log.Println("BuildEntities", err)
		return
	}
	var (
		buf []byte
	)

	for {
		buf, err = data.sock.Read()
		if err != nil && err != io.EOF {
			log.Println("ReadEntities", err)
			return
		}

		if len(buf) > 0 {
			data.ParseResponse(buf)
			if len(data.Entities) > 0 {
				log.Println("COUNT", len(buf), len(data.Entities))
				break
			}
		}
		time.Sleep(time.Millisecond)
	}

	data.EntityKeys = BuildEntityKeys(data.Entities)
	return
}

func (data *HomeData) Monitor() error {

	defer func() {
		if data.Err != nil {
			log.Println("Monitor", data.Err)
		}
	}()

	data.sock, data.Err = websock.NewWebSockClient()
	if data.Err != nil {
		return data.Err
	}

	// authorize

	go data.monitor()

	data.loadStatesID, _ = data.sock.WriteCommandID(StatesCommand)
	// data.loaded.AddListener(binding.NewDataListener(func() {
	// 	isLoaded, _ := data.loaded.Get()
	// 	if isLoaded {
	// 		data.eventsID, data.Err = data.sock.WriteID(subscribe)
	// 		if data.Err != nil {
	// 			return
	// 		}
	// 	}
	// }))
	return data.Err
}

func (data *HomeData) monitor() {
	log.Println("monitor")
	var (
		errCount int
		delay    time.Duration = time.Millisecond * 5
	)
	data.monitoring = true
	for {
		time.Sleep(delay)

		select {
		case <-data.stop:
			log.Println("STOP RECEIVED")
			data.monitoring = false
			return
		default:

			buf, err := data.sock.Read()
			if err != nil {
				errCount++
				if errCount > 10 {
					log.Fatal(err)
				}
				log.Println(err)
				continue
			}
			errCount = 0
			data.ParseResponse(buf)
		}
	}
}
