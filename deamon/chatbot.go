package main

// Connect, subscribe on channel, publish into channel, read presence and history info.

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/centrifugal/centrifuge-go"
)

var url = "ws://localhost:8080/connection/websocket"
var chatbotChannel = "chatbot"
var clientID = ""

type eventHandler struct{}
type subEventHandler struct{}

func (h *eventHandler) OnConnect(c *centrifuge.Client, e centrifuge.ConnectEvent) {
	clientID = e.ClientID
	log.Printf("chatbot connected, clientID: %s", clientID)
	log.Printf("chatbot is listening...")
}

func (h *eventHandler) OnDisconnect(c *centrifuge.Client, e centrifuge.DisconnectEvent) {
	log.Println("client diconnected")
}


func (h *subEventHandler) OnJoin(sub *centrifuge.Subscription, e centrifuge.JoinEvent) {
	log.Println(fmt.Sprintf("User %s (client ID %s) joined channel %s", e.User, e.Client, sub.Channel()))
}

func (h *subEventHandler) OnLeave(sub *centrifuge.Subscription, e centrifuge.LeaveEvent) {
	log.Println(fmt.Sprintf("User %s (client ID %s) left channel %s", e.User, e.Client, sub.Channel()))
}


func (h *subEventHandler) OnPublish(sub *centrifuge.Subscription, e centrifuge.PublishEvent) {
	log.Println(fmt.Sprintf("New publication received from channel %s: %s", sub.Channel(), string(e.Data)))

	if e.GetInfo().Client != clientID {
		out, err := exec.Command("python3", "../chat/chatbot.py", "-q", string(e.Data)).Output()
		if err != nil {
			log.Printf("ERROR Chatbot: %s", err)
		}
		fmt.Printf("Chatbot response: %s", out)

		dataBytes, _ := json.Marshal(string(out))
		err = sub.Publish(dataBytes)
		if err != nil {
			log.Printf("ERROR Publish: %s", err)
		}
	}
}

func waitExitSignal() {
	wait := make(chan int)
	<-wait
}


func main() {
	started := time.Now()

	c := centrifuge.New(url, centrifuge.DefaultConfig())
	defer c.Close()

	eventHandler := &eventHandler{}
	c.OnConnect(eventHandler)
	c.OnDisconnect(eventHandler)

	err := c.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	sub, err := c.NewSubscription(chatbotChannel)
	if err != nil {
		log.Fatalln(err)
	}

	subEventHandler := &subEventHandler{}
	sub.OnPublish(subEventHandler)
	sub.OnJoin(subEventHandler)
	sub.OnLeave(subEventHandler)

	err = sub.Subscribe()
	if err != nil {
		log.Fatalln(err)
	}

	// presence, err := sub.Presence()
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// log.Printf("%d clients in channel %s", len(presence), sub.Channel())

	waitExitSignal()
	log.Printf("END: %s", time.Since(started))

}
