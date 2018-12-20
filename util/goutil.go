package util

import (
	"crypto/sha256"
	"fmt"
	"encoding/json"
	"log"
	"strings"
	"net/http"
	"github.com/InclusION/model"
)


func Hash(input string) string {
	h := sha256.New()

	h.Write([]byte(input))
	//fmt.Printf("%x", h.Sum(nil))
	hashString := fmt.Sprintf("%x", h.Sum(nil))
	return hashString
}



func CheckBodyNil(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
}

func DecodeRequestIntoUser(w http.ResponseWriter, r *http.Request) (error, model.User) {

	decoder := json.NewDecoder(r.Body)

	var user model.User
	err := decoder.Decode(&user)
	if err != nil {
		return err, user
	}

	return nil, user
}

func DecodeRequestIntoRoom(w http.ResponseWriter, r *http.Request) (error, model.Room) {

	decoder := json.NewDecoder(r.Body)

	var room model.Room
	err := decoder.Decode(&room)
	if err != nil {
		return err, room
	}

	return nil, room
}

func DecodeRequestIntoHealth(w http.ResponseWriter, r *http.Request) (error, model.Health) {

	decoder := json.NewDecoder(r.Body)

	var health model.Health
	err := decoder.Decode(&health)
	if err != nil {
		return err, health
	}

	return nil, health
}

func CheckAuth(username string, token string) bool {

	// query user information
	u := model.NewUser(username)
	err, user := u.QueryByUsername()
	if err != nil {
		return false
	}

	log.Println(user.Username)
	log.Println(user.Password)
	log.Println(user.Email)
	log.Println(user.LoginNonce)


	// hash then compare with current hash
	t := Hash(fmt.Sprintf("%s%s%s%d", user.Username, user.Password, user.Email, user.LoginNonce))
	log.Printf("complied token: ", t)
	log.Printf("client token: ", token)
	isEqual := strings.Compare(t, token)

	if isEqual != 0 {
		return false
	}

	return true
}


func HideSensitiveUser(user *model.User) {
	user.Password = ""
	user.LoginNonce = 0
}

func HideSensitiveHealth(health *model.Health) {
	health.Token = ""
}

func HideSensitiveData(input interface{}) interface{} {

	// if model = User
	if user, ok := input.(model.User) ; ok {
		log.Println("hideSensitiveData: User")
		user.Token = ""
		user.Password = ""
		user.LoginNonce = 00

		return user
	}

	// if model = Health
	if health, ok := input.(model.Health) ; ok {
		log.Println("hideSensitiveData: Health")
		health.Token = ""

		return health
	}

	return nil
}
