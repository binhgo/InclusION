package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"log"
	"time"
	"github.com/gorilla/mux"
	"strconv"
	"github.com/InclusION/util"
	"github.com/InclusION/mdb"
	"github.com/InclusION/model"
	"github.com/InclusION/static"
)


func oklah(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, your connection is fine. %s!", r.URL.Path[1:])
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
