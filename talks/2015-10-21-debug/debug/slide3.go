package main

import (
	"log"
	"math/rand"
	"strings"
)

func main() {

	a := strings.Split("HAHAHAHAHAHAHAHAHAHAH", "A")

	for i, char := range "AHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHH BATMAN" {
		for j, b := range "wakakakakakakakaka" {
			go func() {
				log.Print(i)
				print(char)
				print(j)
				print(b)
				chance := rand.Float64()
				if chance > 0.9 {
					print(a[10001010])
				}
			}()
		}
	}
}
