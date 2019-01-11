package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"time"
	"github.com/InclusION/static"
	"github.com/centrifugal/centrifuge"
)

func main() {

	log.Printf("Server started. Listening on port %s", static.PORT)
	log.Printf("UTC Time: %s", time.Now().UTC())

	router := mux.NewRouter()
	router.HandleFunc("/TestConnection", testConnection).Methods(static.HTTP_GET)
	router.HandleFunc("/Register", register).Methods(static.HTTP_POST)
	router.HandleFunc("/Login", login).Methods(static.HTTP_POST)
	router.HandleFunc("/SyncHealth", syncHealth).Methods(static.HTTP_POST)
	router.HandleFunc("/GetLastHealth", getLastHealth).Methods(static.HTTP_POST)
	router.HandleFunc("/UpdateProfile", updateProfile).Methods(static.HTTP_POST)
	router.HandleFunc("/Blog/page/{no}", getAllBlogWithPaging).Methods(static.HTTP_GET)
	router.HandleFunc("/Blog/{id}", getBlogById).Methods(static.HTTP_GET)

	// fmc push notification
	router.HandleFunc("/Push/AddToken", addToken).Methods(static.HTTP_POST)
	router.HandleFunc("/Push/RemoveToken", removeToken).Methods(static.HTTP_POST)
	router.HandleFunc("/Push/PushToDevice", pushToDevice).Methods(static.HTTP_POST)
	router.HandleFunc("/Push/PushToUser", pushToUser).Methods(static.HTTP_POST)


	// chat centrifuge http
	router.HandleFunc("/Chat/FindUser", findUser).Methods(static.HTTP_POST)
	router.HandleFunc("/Chat/CreateChannel", createChannel11).Methods(static.HTTP_POST)

	// chat gorilla ws
	go handleMessages()
	router.HandleFunc("/ws", handleConnections)

	// chat centrifuge ws
	node := initCentrifuge()
	router.Handle("/connection/websocket", centrifuge.NewWebsocketHandler(node, centrifuge.WebsocketConfig{}))

	// handle files html, js
	fs := http.FileServer(http.Dir("./chat"))
	router.PathPrefix("/chat").Handler(http.StripPrefix("/chat", fs))

	// Start HTTP server async
	go startHTTPServer(router)

	// Run program until interrupted.
	waitExitSignal(node)
}

// Start HTTP server.
func startHTTPServer(handler http.Handler) {
	err := http.ListenAndServe(static.PORT, handler)
	if err != nil {
		log.Fatal(err)
	}
}







