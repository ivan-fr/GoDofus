package sockets

import "GoDofus/messages"

func sendChanMsg(channel chan messages.Message, msg messages.Message) {
	channel <- msg
}

func sendChanBool(channel chan bool, b bool) {
	channel <- b
}
