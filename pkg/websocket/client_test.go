package websocket

import (
	"fmt"
	"testing"
)

func TestSend(t *testing.T) {
	client := NewClient("ws://localhost:8080/echo")

	client.Send("当当", func(msg []byte, done chan string) {
		fmt.Println("test received...")
		fmt.Println(string(msg))

		done <- "done"
	})
}