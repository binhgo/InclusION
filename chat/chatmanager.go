package chat

import (
	"github.com/InclusION/model"
)

func GetAllUsers() []model.User {

	u := model.User{}
	allUsers := u.QueryAll()

	return allUsers
}

func FindUser(username string) (error, model.User) {

	u := model.User{Username:username}
	err, user := u.QueryByUsername()

	if err != nil {
		return err, user
	}

	return nil, user
}


func CreateChannel11(r model.Room) (error, model.Room) {

	err, room := r.CreateRoom1To1()
	if err != nil {
		return err, room
	}
	room.SaveToDB()

	return nil, room
}

func CreateChannelGroup(r model.Room) (error, model.Room) {

	err, room := r.CreateRoomGroup()
	if err != nil {
		return err, room
	}
	room.SaveToDB()

	return nil, room
}

// Join channel is client process so that no need to implement here
func JoinChannel() {
}

// call Firebase to notify other users
func NotifyOtherUsers() {
}

