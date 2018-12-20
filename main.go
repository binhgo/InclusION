package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"time"
	"github.com/InclusION/static"
	"strconv"
	"github.com/InclusION/model"
	"encoding/json"
	"github.com/globalsign/mgo/bson"
	"fmt"
	"github.com/InclusION/util"
	"github.com/InclusION/mdb"
	"github.com/gorilla/websocket"
	"github.com/centrifugal/centrifuge"
	"os"
	"os/signal"
	"syscall"
	"context"
	"github.com/InclusION/chat"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)

// Configure the upgrader
var upgrader = websocket.Upgrader{}

type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	// run  goroutine first
	go handleMessages()

	log.Printf("Server started. Listening on port %s", static.PORT)
	log.Printf("UTC Time: %s", time.Now().UTC())

	router := mux.NewRouter()

	//router.HandleFunc("/", oklah).Methods(static.HTTP_GET)
	router.HandleFunc("/TestConnection", testConnection).Methods(static.HTTP_GET)
	router.HandleFunc("/Register", register).Methods(static.HTTP_POST)
	router.HandleFunc("/Login", login).Methods(static.HTTP_POST)
	router.HandleFunc("/SyncHealth", syncHealth).Methods(static.HTTP_POST)
	router.HandleFunc("/GetLastHealth", getLastHealth).Methods(static.HTTP_POST)
	router.HandleFunc("/UpdateProfile", updateProfile).Methods(static.HTTP_POST)
	router.HandleFunc("/Blog/page/{no}", getAllBlogWithPaging).Methods(static.HTTP_GET)
	router.HandleFunc("/Blog/{id}", getBlogById).Methods(static.HTTP_GET)

	//chat
	router.HandleFunc("/Chat/FindUser", findUser).Methods(static.HTTP_POST)
	router.HandleFunc("/Chat/CreateChannel", createChannel11).Methods(static.HTTP_POST)

	//chat gorilla
	router.HandleFunc("/ws", handleConnections)

	// centrifuge chat
	node := initCentrifuge()
	router.Handle("/connection/websocket", centrifuge.NewWebsocketHandler(node, centrifuge.WebsocketConfig{}))

	// files
	fs := http.FileServer(http.Dir("./chat"))
	router.PathPrefix("/chat").Handler(http.StripPrefix("/chat", fs))

	// Start HTTP server.
	go func() {
		//if err := http.ListenAndServe(":8000", nil); err != nil {
		//	panic(err)
		//}

		// start listening
		err := http.ListenAndServe(static.PORT, router)
		if err != nil {
			log.Fatal(err)
		}

	}()

	// Run program until interrupted.
	waitExitSignal(node)
}

func handleLog(e centrifuge.LogEntry) {
	log.Printf("%s: %v", e.Message, e.Fields)
}

// Wait until program interrupted. When interrupted gracefully shutdown Node.
func waitExitSignal(n *centrifuge.Node) {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		n.Shutdown(ctx)
		done <- true
	}()
	<-done
}

func initCentrifuge() *centrifuge.Node {

	cfg := centrifuge.DefaultConfig
	cfg.ClientInsecure = true
	cfg.Publish = true
	node, _ := centrifuge.New(cfg)

	node.On().Connect(func(ctx context.Context, client *centrifuge.Client, e centrifuge.ConnectEvent) centrifuge.ConnectReply {

		client.On().Subscribe(func(e centrifuge.SubscribeEvent) centrifuge.SubscribeReply {
			log.Printf("client subscribes on channel %s", e.Channel)
			return centrifuge.SubscribeReply{}
		})

		client.On().Publish(func(e centrifuge.PublishEvent) centrifuge.PublishReply {
			log.Printf("client publishes into channel %s: %s", e.Channel, string(e.Data))
			return centrifuge.PublishReply{}
		})

		// Set Disconnect Handler to react on client disconnect events.
		client.On().Disconnect(func(e centrifuge.DisconnectEvent) centrifuge.DisconnectReply {
			log.Printf("client disconnected")
			return centrifuge.DisconnectReply{}
		})

		// In our example transport will always be Websocket but it can also be SockJS.
		transportName := client.Transport().Name()
		// In our example clients connect with JSON protocol but it can also be Protobuf.
		transportEncoding := client.Transport().Encoding()

		log.Printf("client connected via %s (%s)", transportName, transportEncoding)
		return centrifuge.ConnectReply{}
	})

	node.SetLogHandler(centrifuge.LogLevelDebug, handleLog)

	if err := node.Run(); err != nil {
		panic(err)
	}

	return node
}


func handleConnections(w http.ResponseWriter, r *http.Request) {

	log.Println("here")

	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {

		// Grab the next message from the broadcast channel
		msg := <-broadcast

		// Send it out to every client that is currently connected
		for client := range clients {

			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
			}

		}
	}
}

func testConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, your connection is fine. %s!", r.URL.Path[1:])
}


func register(w http.ResponseWriter, r *http.Request) {

	util.CheckBodyNil(w, r)

	err, user := util.DecodeRequestIntoUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	user.MongoID = bson.NewObjectId()
	err = user.Insert()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	uu := mdb.QueryById(static.TBL_USERS, user.MongoID)

	json.NewEncoder(w).Encode(&uu)

}


func login(w http.ResponseWriter, r *http.Request) {

	util.CheckBodyNil(w, r)

	err, user := util.DecodeRequestIntoUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	log.Println(user)

	err, u := user.QueryByUsernameAndPassword()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	log.Println(u)

	// if Token expired, create new token
	// hash the username and password and timestamp
	// then insert the hash into user.Token, and user.TokenExpiryTime
	if u.TokenExpiryTime.Before(time.Now().UTC()) || u.TokenExpiryTime.IsZero() {

		u.LoginNonce++
		u.TokenExpiryTime = time.Now().UTC().Add(time.Hour * 24 * 10)

		err := u.UpdateById()
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}


	log.Println(u.Username)
	log.Println(u.Password)
	log.Println(u.Email)
	log.Println(u.LoginNonce)


	// create token and return to client
	// because login, so that we have to return Token, cannot hide it
	u.Token = util.Hash(fmt.Sprintf("%s%s%s%d", u.Username, u.Password, u.Email, u.LoginNonce))

	util.HideSensitiveUser(&u)

	json.NewEncoder(w).Encode(&u)
}


func syncHealth(w http.ResponseWriter, r *http.Request) {

	util.CheckBodyNil(w, r)

	err, health := util.DecodeRequestIntoHealth(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	// test
	log.Println(health)

	isAuth := util.CheckAuth(health.Username, health.Token)
	if isAuth == false {
		http.Error(w, "Authentication fail.", 400)
		return
	}

	health.MongoID = bson.NewObjectId()
	err = health.Insert()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	util.HideSensitiveHealth(&health)

	json.NewEncoder(w).Encode(&health)
}


func getLastHealth(w http.ResponseWriter, r *http.Request) {

	util.CheckBodyNil(w, r)

	err, u := util.DecodeRequestIntoUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	// test
	log.Println(u)

	isAuth := util.CheckAuth(u.Username, u.Token)
	if isAuth == false {
		http.Error(w, "Authentication fail.", 400)
		return
	}

	health := model.Health{}

	err, h := health.QueryLastHealthByUser(u)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	log.Println(u)
	log.Println(h)

	json.NewEncoder(w).Encode(&h)
}


func updateProfile(w http.ResponseWriter, r *http.Request) {

	util.CheckBodyNil(w, r)

	err, user := util.DecodeRequestIntoUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	isAuth := util.CheckAuth(user.Username, user.Token)
	if isAuth == false {
		http.Error(w, "Authentication fail.", 400)
		return
	}

	if len(user.MongoID) > 0 {
		err = user.UpdateById()
		log.Println("UpdateById")
	} else {
		err = user.UpdateByUsername()
		log.Println("UpdateByUsername")
	}

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	util.HideSensitiveUser(&user)

	json.NewEncoder(w).Encode(&user)
}

func getAllBlogWithPaging(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	pageNum := params["no"]
	log.Printf("page %s", pageNum)

	i, e := strconv.Atoi(pageNum)
	if e != nil {
		http.Error(w, e.Error(), 400)
		return
	}

	b := model.Blog{}
	err, blogs := b.QueryAllPaging(i)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(&blogs)
}

func getBlogById(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	blogId := params["id"]
	//log.Printf("blog id %s", blogId)

	b := model.NewBlog(bson.ObjectIdHex(blogId))

	err, blog := b.QueryById()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(&blog)
}

func findUser(w http.ResponseWriter, r *http.Request) {

	util.CheckBodyNil(w, r)

	err, user := util.DecodeRequestIntoUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	//err, u := user.QueryByUsername()
	//if err != nil {
	//	http.Error(w, err.Error(), 400)
	//	return
	//}

	u := chat.FindUser(user.Username)

	json.NewEncoder(w).Encode(&u)
}


func createChannel11(w http.ResponseWriter, r *http.Request) {
	util.CheckBodyNil(w, r)

	err, room := util.DecodeRequestIntoRoom(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	room11 := chat.CreateChannel(room)

	//err, room11 := room.CreateRoom1To1()
	//if err != nil {
	//	http.Error(w, err.Error(), 400)
	//	return
	//}

	json.NewEncoder(w).Encode(&room11)
}