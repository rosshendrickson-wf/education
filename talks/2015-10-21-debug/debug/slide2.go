package main

import "strings"

func main() {

	a := strings.Split("HAHAHAHAHAHAHAHAHAHAH", "A")

	for i, char := range "AHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHH BATMAN" {

		go func() {
			print(i)
			print(char)

			if i%2 == 0 {
				print(a[10001010])
			}
		}()
	}
	print(a[100000])
}
