package message

import (
	"encoding/json"
)

const (
	ball = 0
	line = 1

	maxState = 3
)

type State struct {
	Kind     int
	UUID     int
	Position Vec
	Rotation float32
}

func StatesToMessages(states []*State) []*Message {

	//	result := make([]*Message, 0)
	//	payload := StatesPayload(states)
	//	m := &Message{Type: StateUpdate, Payload: payload}
	//	result = append(result, m)
	//	return result
	//
	if len(states) <= maxState {

		//	log.Printf("Vec: %+v", vectors)
		//	log.Printf("Vec: %+v", vectors[0])

		results := make([]*Message, 1)
		vec := StatesPayload(states)
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
		vec := StatesPayload(chunk)
		payload, _ := json.Marshal(vec)
		m.Payload = payload
		m.Type = StateUpdate
		results[i] = m
	}
	return results

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
