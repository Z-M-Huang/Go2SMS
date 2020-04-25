package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	client := &SMSClient{
		APIKey: os.Getenv("APIKey"),
		Client: &http.Client{},
	}

	id, err := client.SendMessage("+15612920930", "Hello from Github")
	if err != nil {
		log.Panicln(err)
	} else if id == 0 {
		log.Panicln("ID is 0")
	}

	time.Sleep(5 * time.Second)
	sent, err := client.Status(id)
	if err != nil {
		log.Panicln(err)
	}
	if sent {
		log.Println("Sent")
	} else {
		log.Println("Not yet")
	}
}
