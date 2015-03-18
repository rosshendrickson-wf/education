package message

import (
	//	"bytes"
	"encoding/json"
	//	"os"
)

const (
	ball = 0
	line = 1

	maxState = 3
)

type State struct {
	Kind     int
	UUID     int
	Position *Vec
	Rotation float32
}

type StatePayload struct {
	States []*State
}

func StatesToMessages(states []*State) []*Message {

	if len(states) <= maxState {

		//	log.Printf("Vec: %+v", vectors)
		//	log.Printf("Vec: %+v", vectors[0])

		results := make([]*Message, 1)
		vec := StatePayload{states}
		payload, _ := json.Marshal(vec)
		results[0] = &Message{Type: StateUpdate, Payload: payload}
		//	println("HERE")
		return results
	}

	numMessages := len(states) / maxState
	results := make([]*Message, numMessages)

	j := 0
	for i := range results {
		m := &Message{}
		chunk := make([]*State, 0)
		for len(chunk) < maxState && j < len(states)-1 {
			chunk = append(chunk, states[j])
			j++
		}
		vec := StatePayload{chunk}
		payload, _ := json.Marshal(vec)
		m.Payload = payload
		m.Type = StateUpdate
		results[i] = m
	}
	return results

}

func PayloadToStates(payload []byte) StatePayload {

	var sPayload StatePayload
	err := json.Unmarshal(payload, &sPayload)

	if err != nil {
		panic(err)
	}

	//var out bytes.Buffer
	//	json.Indent(&out, payload, "=", "\t")
	//	out.WriteTo(os.Stdout)

	return sPayload
}
