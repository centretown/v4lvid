package ha

import (
	"encoding/json"
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
)

func (data *HomeData) ParseResponse(buf []byte) (err error) {
	resp := &Response{}
	err = decode(buf, resp)
	if err != nil {
		data.Err = err
		return
	}

	switch resp.Type {
	case "event":
		data.parseEvent(buf)
	case "result":
		data.parseResult(buf)
	}

	return
}

func (data *HomeData) parseResult(buf []byte) {
	result := &PartialResult{}
	err := decode(buf, result)
	if err != nil {
		data.Err = err
		return
	}

	if result.ID == data.loadStatesID {
		result := &StateResult{}
		err = decode(buf, result)
		if err != nil {
			log.Println(err, "parseResult")
			return
		}

		for _, entity := range result.Entities {
			data.Entities[entity.EntityID] = entity
			data.Consume(entity.EntityID, entity)
		}
		// data.loaded.Set(true)
	} else if !result.Success {
		showYaml(result)
	}
}

func (data *HomeData) parseEvent(buf []byte) {
	idResult := &EventResult[DataState]{}
	data.Err = decode(buf, idResult)
	if data.Err != nil {
		return
	}

	if idResult.Event.EventType == "state_changed" {
		entityID := idResult.Event.Data.EntityID
		result := &EventResult[DataStateChange[json.RawMessage]]{}
		data.Err = decode(buf, result)
		if data.Err != nil {
			return
		}

		newState := result.Event.Data.NewState
		data.Entities[entityID] = newState
		data.Consume(entityID, newState)
	}
}

func decode(buf []byte, lresult any) (err error) {
	err = json.Unmarshal(buf, lresult)
	if err != nil {
		log.Println(string(buf))
		log.Println("parseAny Unmarshal", err)
	}
	return
}

func showYaml(entity any) {
	out, err := yaml.Marshal(entity)
	if err != nil {
		log.Println("Marshal yaml", err)
		return
	}
	fmt.Println(string(out))
}
