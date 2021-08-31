package sockets

import (
	"GoDofus/messages"
	"GoDofus/pack"
)

func sendChanMsg(channel chan []byte, msg messages.Message, toClient bool, instance uint) {
	packer := pack.Write(msg, toClient, instance)
	go func() {
		channel <- packer
	}()
}
