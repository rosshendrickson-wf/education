package main

import (
	"log"
	"strconv"
	"sync"
	"time"
)

func main() {
	a := generateStrings(100)
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

	s := make([]string, n)

	var wg sync.WaitGroup
	wg.Add(n)
	go func() {
		for i := 0; i < n; i++ {

			go func(a int) {
				if a%2 == 0 {
					s[a] = "A"
				} else {
					s[a] = "H"
				}
				time.Sleep(5000 * time.Millisecond)
				wg.Done()
			}(i)
		}
	}()
	wg.Wait()

	result := ""
	for _, val := range s {
		result += val
	}

	return result
}
