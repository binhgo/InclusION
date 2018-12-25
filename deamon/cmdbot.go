package main

// Connect, subscribe on channel, publish into channel, read presence and history info.
import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/centrifugal/centrifuge-go"
	"github.com/InclusION/model"
)

var clientID = ""
var channelId = "main6868"
var url = "ws://localhost:8000/connection/websocket"

type eventHandler struct{}

func (h *eventHandler) OnConnect(c *centrifuge.Client, e centrifuge.ConnectEvent) {
	clientID = e.ClientID
	log.Printf("chatbot connected, clientID: %s", clientID)
	log.Printf("chatbot is listening...")
}

func (h *eventHandler) OnDisconnect(c *centrifuge.Client, e centrifuge.DisconnectEvent) {
	log.Println("client diconnected")
}

type subEventHandler struct{}

func (h *subEventHandler) OnPublish(sub *centrifuge.Subscription, e centrifuge.PublishEvent) {
	log.Println(fmt.Sprintf("New publication received from channel %s: %s", sub.Channel(), string(e.Data)))

	if e.GetInfo().Client != clientID {
		if sub.Channel() == channelId {
			str := string(e.Data)
			cmd := str[1 : len(str)-1]

			handleCmd(cmd, sub)
		}
	}
}

func (h *subEventHandler) OnJoin(sub *centrifuge.Subscription, e centrifuge.JoinEvent) {
	log.Println(fmt.Sprintf("User %s (client ID %s) joined channel %s", e.User, e.Client, sub.Channel()))
}

func (h *subEventHandler) OnLeave(sub *centrifuge.Subscription, e centrifuge.LeaveEvent) {
	log.Println(fmt.Sprintf("User %s (client ID %s) left channel %s", e.User, e.Client, sub.Channel()))
}

func main() {

	exit := make(chan bool)
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

	sub, err := c.NewSubscription(channelId)
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

	<-exit
	log.Printf("END: %s", time.Since(started))
}


func handleCmd(cmdJson string, sub *centrifuge.Subscription) {

	var response string

	c := model.Command{}
	err, cmd := c.ParseCommand(cmdJson)
	if err != nil {
		log.Println("Error in ParseCommand")
		response = "sorry, I only understand json"
	}

	if cmd.Action == "cmd allUsers" {
		u := model.User{}
		allUsers := u.QueryAll()
		var allUsernames string

		for _, u := range allUsers {
			allUsernames += u.Username + "-"
		}

		response = allUsernames

	} else if cmd.Action == "cmd createRoom11" {

		u1 := cmd.Arg1
		u2 := cmd.Arg2

		// create random channel id and return
		r := model.Room{Username1:u1, Username2:u2}
		err, room := r.CreateRoom1To1()
		if err != nil {
			log.Println("Error in CreateRoom1To1")
			response = err.Error()
		}

		res, err := json.Marshal(room)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(res)

		response = room.ChannelId

	} else {
		// return cmd not found
		response = "cmd not found"
	}

	dataBytes, _ := json.Marshal(response)
	err = sub.Publish(dataBytes)

	if err != nil {
		log.Printf("ERROR Publish: %s", err)
	}

}