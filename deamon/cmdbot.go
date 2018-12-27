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

var clientID string
var cmdChannel = "main6868"
var url = "ws://localhost:8080/connection/websocket"

type eventHandler struct{}
type subEventHandler struct{}

func (h *eventHandler) OnConnect(c *centrifuge.Client, e centrifuge.ConnectEvent) {
	clientID = e.ClientID
	log.Printf("chatbot connected! chatbot id: %s", clientID)
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
		if sub.Channel() == cmdChannel {
			str := string(e.Data)
			cmd := str[1 : len(str)-1]

			HandleCommand(cmd, sub)
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

	sub, err := c.NewSubscription(cmdChannel)
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

	waitExitSignal()
	log.Printf("END: %s", time.Since(started))
}


var actionGetAllUsers = "cmd allUsers"
var actionCreateRoom11 = "cmd createRoom11"

func HandleCommand(cmdJson string, sub *centrifuge.Subscription) {
	var response string

	c := model.Command{}
	err, cmd := c.ParseCommand(cmdJson)
	if err != nil {
		response = "sorry, I only understand json"
	}

	if cmd.Action == actionGetAllUsers {
		u := model.User{}
		allUsers := u.QueryAll()
		var allUsernames string

		for _, u := range allUsers {
			allUsernames += u.Username + "--"
		}

		response = allUsernames

	} else if cmd.Action == actionCreateRoom11 {

		u1 := cmd.Arg1
		u2 := cmd.Arg2

		// create random channel id and return
		r := model.Room{Username1:u1, Username2:u2}
		err, room := r.CreateRoom1To1()
		if err != nil {
			response = err.Error()
		}

		_, err = json.Marshal(room)
		if err != nil {
			response = err.Error()
		}

		response = room.ChannelId

	} else {
		response = "cmd not found"
	}

	dataBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error in json.Marshal: %s", err)
	}

	err = sub.Publish(dataBytes)
	if err != nil {
		log.Printf("ERROR Publish: %s", err)
	}

}