package message

type Input struct {
	left     bool
	right    bool
	up       bool
	down     bool
	spacebar bool
}

type InputPayload struct {
	Inputs []*Input
}
