package sockets

import (
	"GoDofus/messages"
	"GoDofus/pack"
)

func sendChanMsg(channel chan []byte, msg messages.Message, toClient bool, instance uint) {
	channel <- pack.Write(msg, toClient, instance)
}

func sendChanWeft(channel chan []byte, msg *pack.Weft, toClient bool, instance uint) {
	channel <- pack.WriteWeft(msg, toClient, instance)
}
