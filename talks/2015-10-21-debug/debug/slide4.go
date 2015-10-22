package main

import (
	"log"
	"strconv"
)

func main() {

	a := generateStrings(10)
	sendToLog(a)
}

func mutate(s string) string {
	l := sumLength(s)
	s += strconv.Itoa(l)
	return s
}

func sendToLog(s string) {
	b := mutate(s)
	log.Print(b[len(s)+10])
}

func sumLength(s string) int {
	return len(s)
}

func generateStrings(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			s += "A"
		} else {
			s += "H"
		}
	}
	return s
}
