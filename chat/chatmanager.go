package chat

import (
	"github.com/InclusION/model"
	"log"
)

func FindUser(username string) model.User {

	u := model.User{Username:username}
	err, user := u.QueryByUsername()

	if err != nil {
		log.Fatal(err)
	}

	return user
}


func CreateChannel(room model.Room) model.Room {

	err, r := room.CreateRoom1To1()
	if err != nil {
		log.Fatal(err)
	}

	// save channel id to database
	r.SaveToDB()

	return r
}

// Join channel is client process so that no need to implement here
func JoinChannel() {
}

// call Firebase to notify other users
func NotifyOtherUsers() {
}

