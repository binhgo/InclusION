package chat

import (
	"github.com/InclusION/model"
	"github.com/InclusION/fcm"
	"log"
	"errors"
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

	// check if room exist or not, if exist then return without creating new room
	ok, room := r.CheckRoomExist()
	if ok {
		return nil, room
	}

	// no room found, create new room id (channel id)
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
func NotifyUser(username string, content string) {
	// go routine push notification
	err := fcm.NotifyAllDevicesOfUser(username, content)
	if err != nil {
		log.Fatal(err)
	}
}


func CreateChannelWithHashTag(r model.Room) (error, model.Room) {

	// check if room exist or not, if exist then return without creating new room
	ok, room := r.CheckRoomExist()
	if ok {
		return nil, room
	}

	// no room found, create new room id (channel id)
	err, room := r.CreateRoomWithHashTag()
	if err != nil {
		return err, room
	}

	room.SaveToDB()

	return nil, room
}


func InviteUserToJoinGroup(user model.User, room model.Room) {
	// push notification to user
	NotifyUser(user.Username, room.ChannelId)
}


func FindRoomWithHashTag(r model.Room) (error, model.Room) {

	// check if room exist or not, if exist then return without creating new room
	ok, room := r.CheckRoomExist()
	if ok {
		return nil, room
	}

	return errors.New("error no room"), room
}


func ValidateUserJoinRoom(username string) bool {

	// from centriguge client id -> query username, then check room info to see if user valid


	return false
}






