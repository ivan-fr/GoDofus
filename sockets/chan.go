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

func sendChanWeft(channel chan []byte, msg *pack.Weft, toClient bool, instance uint) {
	packer := pack.WriteWeft(msg, toClient, instance)
	go func() {
		channel <- packer
	}()
}
