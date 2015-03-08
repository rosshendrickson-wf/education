package message

import (
	"encoding/json"
)

const (
	ball = 0
	line = 1

	maxState = 15
)

type State struct {
	Kind     int
	UUID     int
	Position Vec
	Rotation float32
}

func StatesToMessages(states []*State) []*Message {

	result := make([]*Message, 0)
	payload := StatesPayload(states)
	m := &Message{Type: StateUpdate, Payload: payload}
	result = append(result, m)
	return result
}

func StatesPayload(states []*State) []byte {
	payload, _ := json.Marshal(states)
	return payload
}

func PayloadToStates(payload []byte) []*State {

	var sPayload []*State
	json.Unmarshal(payload, &sPayload)

	return sPayload
}
