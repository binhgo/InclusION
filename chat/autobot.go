package chat
//
//// Connect, subscribe on channel, publish into channel, read presence and history info.
//import (
//"fmt"
//"log"
//"time"
//
//"github.com/centrifugal/centrifuge-go"
//"encoding/json"
//)
//
//var clientID = ""
//
//
//
//type eventHandler struct{}
//
//func (h *eventHandler) OnConnect(c *centrifuge.Client, e centrifuge.ConnectEvent) {
//	clientID = e.ClientID
//	log.Printf("chatbot connected, clientID: %s", clientID)
//	log.Printf("chatbot is listening...")
//}
//
//func (h *eventHandler) OnDisconnect(c *centrifuge.Client, e centrifuge.DisconnectEvent) {
//	log.Println("client diconnected")
//}
//
//type subEventHandler struct{}
//
//func (h *subEventHandler) OnPublish(sub *centrifuge.Subscription, e centrifuge.PublishEvent) {
//	log.Println(fmt.Sprintf("New publication received from channel %s: %s", sub.Channel(), string(e.Data)))
//
//	if sub.Channel() == "main6868" {
//
//		cmd := string(e.Data)
//
//		var response string
//		if cmd == "cmd findUser" {
//			//return all username
//			response = "this is all user"
//
//		} else if cmd == "cmd createRoom" {
//			// create random channel id and return
//			response = "this is sample random channel id"
//
//		} else {
//			// return cmd not found
//			response = "cmd not found"
//		}
//
//		dataBytes, _ := json.Marshal(response)
//		err := sub.Publish(dataBytes)
//		if err != nil {
//			log.Println("ERROR Publish")
//			log.Println(err)
//		}
//	}
//}
//
//func (h *subEventHandler) OnJoin(sub *centrifuge.Subscription, e centrifuge.JoinEvent) {
//	log.Println(fmt.Sprintf("User %s (client ID %s) joined channel %s", e.User, e.Client, sub.Channel()))
//}
//
//func (h *subEventHandler) OnLeave(sub *centrifuge.Subscription, e centrifuge.LeaveEvent) {
//	log.Println(fmt.Sprintf("User %s (client ID %s) left channel %s", e.User, e.Client, sub.Channel()))
//}
//
//func main() {
//
//	exit := make(chan bool)
//
//	started := time.Now()
//
//	url := "ws://localhost:8080/connection/websocket"
//
//	c := centrifuge.New(url, centrifuge.DefaultConfig())
//	defer c.Close()
//
//	eventHandler := &eventHandler{}
//	c.OnConnect(eventHandler)
//	c.OnDisconnect(eventHandler)
//
//	err := c.Connect()
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	sub, err := c.NewSubscription("main6868")
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	subEventHandler := &subEventHandler{}
//	sub.OnPublish(subEventHandler)
//	sub.OnJoin(subEventHandler)
//	sub.OnLeave(subEventHandler)
//
//	err = sub.Subscribe()
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	<-exit
//	log.Printf("%s", time.Since(started))
//}
//
