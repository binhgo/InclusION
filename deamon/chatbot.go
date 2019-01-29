package main

// Connect, subscribe on channel, publish into channel, read presence and history info.

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/centrifugal/centrifuge-go"
	"github.com/InclusION/model"
	"github.com/InclusION/chat"
	"strings"
)

type ChatRoom struct{
	Room model.Room
	Sub *centrifuge.Subscription
}

type Message struct {
	ChannelId string
	ClientId string
}

type ChatRequest struct {
	Name string
	Email string
}

type eventHandler struct{}
type subEventHandler struct{}

var url = "ws://localhost:8080/connection/websocket"
var chatbotChannel = "chatbot"
var botClientID = ""

var rooms []ChatRoom

var client *centrifuge.Client

var subHandler = &subEventHandler{}


func (h *eventHandler) OnConnect(c *centrifuge.Client, e centrifuge.ConnectEvent) {
	botClientID = e.ClientID
	log.Printf("chatbot connected, clientID: %s, listening...", botClientID)
}

func (h *eventHandler) OnDisconnect(c *centrifuge.Client, e centrifuge.DisconnectEvent) {
	log.Println("client diconnected")
}

func (h *subEventHandler) OnJoin(sub *centrifuge.Subscription, e centrifuge.JoinEvent) {
	log.Println(fmt.Sprintf("User %s (client ID %s) joined channel %s", e.User, e.Client, sub.Channel()))

	//if e.Client != botClientID {
	//	go spawnAndSubscribeNewChannel(sub, e.Client)
	//}
}


func (h *subEventHandler) OnLeave(sub *centrifuge.Subscription, e centrifuge.LeaveEvent) {
	log.Println(fmt.Sprintf("User %s (client ID %s) left channel %s", e.User, e.Client, sub.Channel()))
}


func (h *subEventHandler) OnPublish(sub *centrifuge.Subscription, e centrifuge.PublishEvent) {
	log.Println(fmt.Sprintf("New publication received from channel %s: %s", sub.Channel(), string(e.Data)))

	if e.GetInfo().Client != botClientID {

		r := strings.NewReader(string(e.Data))
		decoder := json.NewDecoder(r)

		var chatReq ChatRequest
		err := decoder.Decode(&chatReq)

		if err != nil {
			log.Fatal(err)
		} else {
			if len(chatReq.Name) > 0 {
				go spawnAndSubscribeNewChannel(sub, e.GetInfo().Client)
			}
		}

		go findRoomAndReply(e.GetInfo().Client, string(e.Data))
	}
}


func waitExitSignal() {
	wait := make(chan int)
	<-wait
}


func main() {
	started := time.Now()

	client = centrifuge.New(url, centrifuge.DefaultConfig())
	defer client.Close()

	eventHandler := &eventHandler{}
	client.OnConnect(eventHandler)
	client.OnDisconnect(eventHandler)

	err := client.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	sub, err := client.NewSubscription(chatbotChannel)
	if err != nil {
		log.Fatalln(err)
	}

	//subEventHandler := &subEventHandler{}
	sub.OnPublish(subHandler)
	sub.OnJoin(subHandler)
	sub.OnLeave(subHandler)

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


func findRoomAndReply(username1 string, question string) {
	// find correct session and send the message to that only channel
	for _, room := range rooms {

		if room.Room.Username1 ==  username1 {

			var data string
			answer1, ok := filterQuestion(question)
			if ok {
				// if ok, then return the answer1 to end user
				data = answer1
			} else {
				// if not ok, then call chatbot python to answer the question
				out, err := exec.Command("python3", "../chat/chatbot.py", "-q", question).Output()
				if err != nil {
					log.Printf("ERROR Chatbot: %s", err)
				}

				data = filterAnswer(string(out))
			}

			dataBytes, _ := json.Marshal(data)
			err := room.Sub.Publish(dataBytes)
			if err != nil {
				log.Printf("ERROR Publish: %s", err)
			}

			break
		}
	}
}


func spawnAndSubscribeNewChannel(sub *centrifuge.Subscription, clientId1 string) {

	r := model.Room{Username1: clientId1, Username2:botClientID}
	err, room := chat.CreateChannel11(r)
	if err != nil {
		log.Printf("ERROR CreateChannel11: %s", err)
	}
	//return channel id to client
	//subscribe a new sub into this channleid
	subscription, err := client.NewSubscription(room.ChannelId)
	if err != nil {
		log.Fatalln(err)
	}

	chatRoom := ChatRoom{Room:room, Sub:subscription}
	chatRoom.Sub.OnPublish(subHandler)
	chatRoom.Sub.OnJoin(subHandler)
	chatRoom.Sub.OnLeave(subHandler)

	err = subscription.Subscribe()
	if err != nil {
		log.Fatalln(err)
	}
	rooms = append(rooms, chatRoom)

	mess := &Message{ChannelId:chatRoom.Room.ChannelId, ClientId:chatRoom.Room.Username1}
	jsonMess, _ := json.Marshal(mess)
	strMess := string(jsonMess)

	dataBytes, _ := json.Marshal(strMess)
	err = sub.Publish(dataBytes)
	if err != nil {
		log.Printf("ERROR Publish: %s", err)
	}
}


func filterQuestion(question string) (string, bool) {
	if strings.Contains(question, "event") {
		return returnEvents(), true
	} else if strings.Contains(question, "product") {
		return returnProducts(), true
	} else if strings.Contains(question, "suggest") {
		return returnSuggestions(), true
	} else {
		return "", false
	}
}


func filterAnswer(data string) string {
	if data == "I am sorry, but I do not understand." {
		return returnSuggestions()
	} else {
		return data
	}
}


type Product struct {
	Name string
	Description string
	Image string
}


func returnProducts() string {
	var products []Product

	p1 := Product{"Rocco Trunki", "Vali Trẻ Em Siêu Xe Rocco Trunki 0321-GB01 với thiết kế thông minh tạo sự tiện lợi cho cả mẹ và bé, nhưng không kém phần ngộ nghĩnh, đáng yêu và thân thiện với bé.", "url"}
	p2 := Product{"Rocco Trunki", "Vali Trẻ Em Siêu Xe Rocco Trunki 0321-GB01 với thiết kế thông minh tạo sự tiện lợi cho cả mẹ và bé, nhưng không kém phần ngộ nghĩnh, đáng yêu và thân thiện với bé.", "url"}
	p3 := Product{"Rocco Trunki", "Vali Trẻ Em Siêu Xe Rocco Trunki 0321-GB01 với thiết kế thông minh tạo sự tiện lợi cho cả mẹ và bé, nhưng không kém phần ngộ nghĩnh, đáng yêu và thân thiện với bé.", "url"}
	p4 := Product{"Rocco Trunki", "Vali Trẻ Em Siêu Xe Rocco Trunki 0321-GB01 với thiết kế thông minh tạo sự tiện lợi cho cả mẹ và bé, nhưng không kém phần ngộ nghĩnh, đáng yêu và thân thiện với bé.", "url"}
	p5 := Product{"Rocco Trunki", "Vali Trẻ Em Siêu Xe Rocco Trunki 0321-GB01 với thiết kế thông minh tạo sự tiện lợi cho cả mẹ và bé, nhưng không kém phần ngộ nghĩnh, đáng yêu và thân thiện với bé.", "url"}
	products = append(products, p1, p2, p3, p4, p5)

	result, err := json.Marshal(products)
	if err != nil {
		log.Fatal(err)
	}

	return string(result)
}

type Suggestion struct {
	Type string
	Name string
	Description string
	Image string
}

func returnSuggestions() string {
	var suggestions []Suggestion

	s1 := Suggestion{"Event", "Find some events near me", "", "url"}
	s2 := Suggestion{"Group", "Want to join some groups", "", "url"}
	s3 := Suggestion{"Product", "Look for some products", "", "url"}
	s4 := Suggestion{"Therapy", "Want to find some specialists", "", "url"}
	s5 := Suggestion{"Chat", "Look for a person to chat", "", "url"}
	suggestions = append(suggestions, s1, s2, s3, s4, s5)

	result, err := json.Marshal(suggestions)
	if err != nil {
		log.Fatal(err)
	}

	return string(result)
}


type Event struct {
	Name string
	Description string
	Image string
	Time time.Time
}

func returnEvents() string {
	var events []Event

	e1 := Event{"India Inclusion", "India inclusion summit 2018", "", time.Now()}
	e2 := Event{"Singapore Inclusion", "Singapore Inclusion 2019", "", time.Now()}
	e3 := Event{"Vietnam Inclusion", "Vietnam Inclusion 2019", "", time.Now()}
	e4 := Event{"Thailand Inclusion", "Thailand Inclusion 2019", "", time.Now()}
	e5 := Event{"Lambada Inclusion", "Lambada Inclusion 2019", "", time.Now()}
	events = append(events, e1, e2, e3, e4, e5)

	result, err := json.Marshal(events)
	if err != nil {
		log.Fatal(err)
	}

	return string(result)
}




