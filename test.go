package main

import "fmt"

func main() {
	message := make(chan string)

	go func() {
		message <- "aligm 79"
	}()

	go func() {
		message <- "gm ali"
	}()
	
	fmt.Print(message)
	msg := <- message
	fmt.Print(msg)
}