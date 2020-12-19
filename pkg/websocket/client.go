package websocket

import (
	//"flag"
	"log"
	//"net/url"
	"os"
	"os/signal"
	// "time"

	"github.com/gorilla/websocket"	
)

type Client struct {
	URL string
}

func NewClient(url string) *Client {
	return &Client{URL: url}
}
func (client *Client) Send(data map[string]interface{}, cb func(word []byte, done chan string)) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial(client.URL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan string)

	// send
	// err = c.WriteMessage(websocket.TextMessage, []byte(word))
	err = c.WriteJSON(data)

	if err != nil {
		log.Println("write:", err)
			return
	}

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			cb(message, done)
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			}
			return
		}
	}
}