package sockets

import (
	"GoDofus/messages"
	"GoDofus/pack"
	"bytes"
	"fmt"
)

func handlingMyClient(writeInMyClientChan, writeToAnkamaServerChan chan messages.Message, myClientContinueChan, officialServerContinueChan chan bool, instance uint) func(chan *pack.Weft) {
	return func(weftChan chan *pack.Weft) {
		for {
			weft := <-weftChan

			if weft == nil {
				break
			}

			switch weft.PackId {
			case messages.CheckIntegrityID:
				msg := messages.GetCheckIntegrityNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeToAnkamaServerChan <- msg
			case messages.ClientKeyID:
				msg := messages.GetClientKeyNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeToAnkamaServerChan <- msg
			default:
				fmt.Printf("Listener: there is no traitment for %d ID\n", weft.PackId)
			}
		}
	}
}
