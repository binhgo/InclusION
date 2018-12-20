package model

import (
	"net/http"
	"encoding/json"
	"log"
	"fmt"
	"strings"
	"github.com/InclusION/util"
)

func DecodeRequestIntoUser(w http.ResponseWriter, r *http.Request) (error, User) {

	decoder := json.NewDecoder(r.Body)

	var user User
	err := decoder.Decode(&user)
	if err != nil {
		return err, user
	}

	return nil, user
}

func DecodeRequestIntoRoom(w http.ResponseWriter, r *http.Request) (error, Room) {

	decoder := json.NewDecoder(r.Body)

	var room Room
	err := decoder.Decode(&room)
	if err != nil {
		return err, room
	}

	return nil, room
}

func DecodeRequestIntoHealth(w http.ResponseWriter, r *http.Request) (error, Health) {

	decoder := json.NewDecoder(r.Body)

	var health Health
	err := decoder.Decode(&health)
	if err != nil {
		return err, health
	}

	return nil, health
}

func CheckAuth(username string, token string) bool {

	// query user information
	u := NewUser(username)
	err, user := u.QueryByUsername()
	if err != nil {
		return false
	}

	log.Println(user.Username)
	log.Println(user.Password)
	log.Println(user.Email)
	log.Println(user.LoginNonce)


	// hash then compare with current hash
	t := util.Hash(fmt.Sprintf("%s%s%s%d", user.Username, user.Password, user.Email, user.LoginNonce))
	log.Printf("complied token: ", t)
	log.Printf("client token: ", token)
	isEqual := strings.Compare(t, token)

	if isEqual != 0 {
		return false
	}

	return true
}


func HideSensitiveUser(user *User) {
	user.Password = ""
	user.LoginNonce = 0
}

func HideSensitiveHealth(health *Health) {
	health.Token = ""
}


func HideSensitiveData(input interface{}) interface{} {

	// if model = User
	if user, ok := input.(User) ; ok {
		log.Println("hideSensitiveData: User")
		user.Token = ""
		user.Password = ""
		user.LoginNonce = 00

		return user
	}

	// if model = Health
	if health, ok := input.(Health) ; ok {
		log.Println("hideSensitiveData: Health")
		health.Token = ""

		return health
	}

	return nil
}
