package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"time"
	"InclusION/static"
)


func main() {

	log.Printf("Server started. Listening on port %s", static.PORT)
	log.Printf("UTC Time: %s", time.Now().UTC())

	router := mux.NewRouter()

	router.HandleFunc("/TestConnection", testConnection).Methods(static.HTTP_GET)
	router.HandleFunc("/Register", register).Methods(static.HTTP_POST)
	router.HandleFunc("/Login", login).Methods(static.HTTP_POST)
	router.HandleFunc("/SyncHealth", syncHealth).Methods(static.HTTP_POST)
	router.HandleFunc("/UpdateProfile", updateProfile).Methods(static.HTTP_POST)
	router.HandleFunc("/Blog/page/{no}", getAllBlogWithPaging).Methods(static.HTTP_GET)
	router.HandleFunc("/Blog/{id}", getBlogById).Methods(static.HTTP_GET)

	//chat

	err := http.ListenAndServe(static.PORT, router)
	if err != nil {
		log.Fatal(err)
	}




	// test
	//coll := "pserson"
	//
	//id := bson.NewObjectId()
	//user := model.User{id, 1, "hnb2018", "123456Aa@", "hb@gmail.com", time.Now()}
	//mdb.Insert(coll, user)
	//
	//log.Println("before")
	//users := mdb.QueryAll(coll)
	//for _, v := range users {
	//	log.Println(v)
	//}
	//log.Println("------")
	//
	//
	//log.Println("update")
	//mdb.Update(coll)
	//log.Println("------")
	//
	//
	//log.Println("after")
	//users = mdb.QueryAll(coll)
	//for _, v := range users {
	//	log.Println(v)
	//}
	//log.Println("------")
	//
	//
	//log.Println("end...")
	// test
}