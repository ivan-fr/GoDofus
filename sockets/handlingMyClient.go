package sockets

import (
	"GoDofus/messages"
	"GoDofus/pack"
	"bytes"
	"fmt"
)

func handlingMyClient(writeInMyClientChan, writeToAnkamaServerChan chan messages.Message, myClientContinueChan, ankamaServerContinueChan chan bool) func(*pack.Pipe) {
	return func(pipe *pack.Pipe) {
		for weft := pipe.Get(); weft != nil; weft = pipe.Get() {
			switch weft.PackId {
			case messages.CheckIntegrityID:
				msg := messages.GetCheckIntegrityNOA()
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeToAnkamaServerChan <- msg
			case messages.ClientKeyID:
				msg := messages.GetClientKeyNOA()
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeToAnkamaServerChan <- msg
			default:
				fmt.Printf("Listener: there is no traitment for %d ID\n", weft.PackId)
			}
		}
	}
}
