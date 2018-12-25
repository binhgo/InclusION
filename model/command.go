package model

import (
	"log"
	"strings"
)

// find user
// action: findUser

// create room 11
// action: createRoom11
// arg1: username1
// arg2: username2

// create room many
// action: createRoomMany
// arg1: username1
// arg2: room name

// notify other users
// action: notify
// arg1: username
// arg2: data or channel id


type Command struct {
	Action string
	Arg1 string
	Arg2 string
	Arg3 string
}


func (c *Command) ParseCommand(cmdJson string) (error, Command) {

	// escape string
	cmdStr := strings.Replace(cmdJson, "\\", "" , -1)
	log.Printf("cmd: %s", cmdStr)

	// json to object
	err, cmd := DecodeRequestIntoCommand(cmdStr)
	if err != nil {
		return err, cmd
	}

	return nil, cmd
}




