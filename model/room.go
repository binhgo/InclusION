package model

import (
	"sort"
	"github.com/InclusION/util"
	"log"
	"time"
	rand2 "math/rand"
)

type Room struct {
	ChannelId string
	RoomName string

	Username1 string
	Username2 string
}

func (r *Room) CreateRoom1To1() (error, Room) {
	usernames := []string{r.Username1, r.Username2}

	r1 := sortThenConcat(usernames)
	channelID := util.Hash(r1)

	return nil, Room{ChannelId:channelID, RoomName:r1, Username1:r.Username1, Username2:r.Username2}
}

func (r *Room) CreateRoomMany(users []User) (error, Room) {
	var usernames []string
	for i, u := range users {
		usernames[i] = u.Username
	}

	r1 := sortThenConcat(usernames)
	channelID := util.Hash(r1)

	return nil, Room{ChannelId:channelID, RoomName:r1}
}


func (r *Room) CreateRandomRoom() (error, Room) {

	//max := big.Int{true, 100000000000}

	str := string(rand2.Int63n(10000000)) + time.Now().String()

	channelID := util.Hash(str)

	return nil, Room{ChannelId:channelID, RoomName:str}
}


func sortThenConcat(usernames []string) string {
	sort.Strings(usernames)
	//fmt.Println(s)

	var result string

	for _, s := range usernames {
		result += s
	}

	return result
}

func (r *Room) SaveToDB() {
	log.Println("Saved to db")
}


