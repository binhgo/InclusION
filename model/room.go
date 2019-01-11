package model

import (
	"sort"
	"github.com/InclusION/util"
	"log"
	"time"
	rand2 "math/rand"
	"errors"
	"github.com/InclusION/mdb"
	"github.com/InclusION/static"
	"github.com/globalsign/mgo/bson"
)

type Room struct {
	ChannelId string
	RoomName string

	Username1 string
	Username2 string
	HashTag string
}

func (r *Room) CreateRoom1To1() (error, Room) {

	// check username1
	if len(r.Username1) == 0 {
		return errors.New("username1 nil"), Room{}
	}

	// check username2
	if len(r.Username2) == 0 {
		return errors.New("username2 nil"), Room{}
	}

	usernames := []string{r.Username1, r.Username2}

	r1 := sortThenConcat(usernames)
	channelID := util.Hash(r1)

	return nil, Room{ChannelId:channelID, RoomName:r1, Username1:r.Username1, Username2:r.Username2}
}

func (r *Room) CreateRoomGroup() (error, Room) {

	// check username1
	if len(r.Username1) == 0 {
		return errors.New("username1 nil"), Room{}
	}

	// check room name
	if len(r.RoomName) == 0 {
		return errors.New("roomname nil"), Room{}
	}

	str := r.Username1 + string(rand2.Int63n(100000000)) + time.Now().String()

	channelID := util.Hash(str)
	return nil, Room{ChannelId:channelID, RoomName:r.RoomName}
}


func (r *Room) CreateRandomRoom() (error, Room) {
	str := string(rand2.Int63n(10000000)) + time.Now().String()
	channelID := util.Hash(str)
	return nil, Room{ChannelId:channelID, RoomName:str}
}


func (r *Room) CreateRoomWithHashTag() (error, Room) {
	// check username1
	if len(r.Username1) == 0 {
		return errors.New("username1 nil"), Room{}
	}

	// check room name
	if len(r.HashTag) == 0 {
		return errors.New("HashTag nil"), Room{}
	}

	channelId := util.Hash(r.HashTag)

	return nil, Room{ChannelId:channelId, RoomName:r.RoomName}
}


func sortThenConcat(usernames []string) string {
	sort.Strings(usernames)

	var result string
	for _, s := range usernames {
		result += s
	}

	return result
}

func (r *Room) SaveToDB() error {
	log.Println("Saved Room to DB")

	err := mdb.Insert(static.TBL_ROOMS, r)
	if err  != nil {
		return err
	}

	log.Println("Inserted")
	return nil
}


func (r *Room) FindRoomByChannelId() (error, Room) {

	db := mdb.InitDB()
	c := db.C(static.TBL_ROOMS)

	var result Room
	err := c.Find(bson.M{"Channelid": r.ChannelId}).One(&result)
	if err != nil {
		return err, result
	}

	return nil, result
}


func (r *Room) FindRoomByHashTag() (error, Room) {

	db := mdb.InitDB()
	c := db.C(static.TBL_ROOMS)

	var result Room
	err := c.Find(bson.M{"hashtag": r.HashTag}).One(&result)
	if err != nil {
		return err, result
	}

	return nil, result
}



func (r *Room) CheckRoomExist() (bool, Room) {

	// check if room create by hash tag = group room
	if len(r.HashTag) > 0 {
		// find room by hash tag
		err, room := r.FindRoomByHashTag()
		if err != nil {
			return false, room
		}

		if len(room.ChannelId) == 0 {
			return false, room
		} else {
			return true, room
		}
	} else {

		err, room := r.FindRoomByChannelId()
		if err != nil {
			return false, room
		}

		if len(room.ChannelId) == 0 {
			return false, room
		} else {
			return true, room
		}
	}
}


