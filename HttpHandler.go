package main

import (
	"log"
	"time"
	"net/http"
	"fmt"
	"github.com/InclusION/util"
	"github.com/InclusION/model"
	"github.com/globalsign/mgo/bson"
	"github.com/InclusION/mdb"
	"github.com/InclusION/static"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
	"github.com/InclusION/chat"
	"github.com/InclusION/fcm"
)


//**********************************************************************************//
// http requests
func testConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, your connection is fine. %s!", r.URL.Path[1:])
}


func register(w http.ResponseWriter, r *http.Request) {

	util.CheckBodyNil(w, r)

	err, user := model.DecodeRequestIntoUser(w, r)
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

	err, user := model.DecodeRequestIntoUser(w, r)
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
	// create token and return to client
	// because login, so that we have to return Token, cannot hide it
	u.Token = util.Hash(fmt.Sprintf("%s%s%s%d", u.Username, u.Password, u.Email, u.LoginNonce))

	model.HideSensitiveUser(&u)

	json.NewEncoder(w).Encode(&u)
}


func syncHealth(w http.ResponseWriter, r *http.Request) {

	util.CheckBodyNil(w, r)

	err, health := model.DecodeRequestIntoHealth(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	// test
	log.Println(health)

	isAuth := model.CheckAuth(health.Username, health.Token)
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

	model.HideSensitiveHealth(&health)

	json.NewEncoder(w).Encode(&health)
}


func getLastHealth(w http.ResponseWriter, r *http.Request) {

	util.CheckBodyNil(w, r)

	err, u := model.DecodeRequestIntoUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	// test
	log.Println(u)

	isAuth := model.CheckAuth(u.Username, u.Token)
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

	err, rUser := model.DecodeRequestIntoUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	isAuth := model.CheckAuth(rUser.Username, rUser.Token)
	if isAuth == false {
		http.Error(w, "Authentication fail.", 400)
		return
	}

	if len(rUser.MongoID) > 0 {
		err = rUser.UpdateById()
		log.Println("UpdateById")
	} else {
		err = rUser.UpdateByUsername()
		log.Println("UpdateByUsername")
	}

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	model.HideSensitiveUser(&rUser)

	json.NewEncoder(w).Encode(&rUser)
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

	rBlog := model.Blog{}
	err, blogs := rBlog.QueryAllPaging(i)
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

	rBlog := model.NewBlog(bson.ObjectIdHex(blogId))

	err, blog := rBlog.QueryById()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(&blog)
}

func findUser(w http.ResponseWriter, r *http.Request) {

	util.CheckBodyNil(w, r)

	err, rUser := model.DecodeRequestIntoUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err, user := chat.FindUser(rUser.Username)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(&user)
}


func createChannel11(w http.ResponseWriter, r *http.Request) {

	util.CheckBodyNil(w, r)

	err, rRoom := model.DecodeRequestIntoRoom(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err, room11 := chat.CreateChannel11(rRoom)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// go routine push notification
	go fcm.NotifyAllDevicesOfUser(rRoom.Username2, room11.ChannelId)

	json.NewEncoder(w).Encode(&room11)
}

func addToken(w http.ResponseWriter, r *http.Request) {

	util.CheckBodyNil(w, r)

	err, p := model.DecodeRequestIntoPhone(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	p.MongoID = bson.NewObjectId()

	err = p.Insert()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	phone := mdb.QueryById(static.TBL_DEVICES, p.MongoID)

	json.NewEncoder(w).Encode(&phone)
}


func removeToken(w http.ResponseWriter, r *http.Request) {

	util.CheckBodyNil(w, r)

	err, p := model.DecodeRequestIntoPhone(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = p.HardDelete()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(&p)
}


func pushToDevice(w http.ResponseWriter, r *http.Request) {
	// device token
	// data to push

	util.CheckBodyNil(w, r)

	err, mess := model.DecodeRequestIntoFcmMessage(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if len(mess.DeviceToken) <= 0 {
		http.Error(w, "Device token cannot be nil", 400)
		return
	}

	err = fcm.Notify1Device(mess.DeviceToken, mess.Content)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(&mess)
}


func pushToUser(w http.ResponseWriter, r *http.Request) {
	// device token
	// data to push

	util.CheckBodyNil(w, r)

	err, mess := model.DecodeRequestIntoFcmMessage(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if len(mess.Username) <= 0 {
		http.Error(w, "Username cannot be nil", 400)
		return
	}

	err = fcm.NotifyAllDevicesOfUser(mess.Username, mess.Content)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(&mess)
}

// http requests
//**********************************************************************************//
