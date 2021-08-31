package sockets

import (
	"GoDofus/messages"
	"GoDofus/pack"
	"bytes"
	"fmt"
)

func handlingGame(writeInMyClientChan, writeToOfficialServerChan chan []byte, myClientContinueChan, officialServerContinueChan chan bool, instance uint) func(chan *pack.Weft) {
	return func(weftChan chan *pack.Weft) {
		for {
			weft := <-weftChan

			if weft == nil {
				break
			}
			switch weft.PackId {
			default:
				msg, ok := messages.Types_[int(weft.PackId)]

				if ok {
					msg = msg.GetNOA(instance)
					msg.Deserialize(bytes.NewReader(weft.Message))
					fmt.Println(msg)
					go sendChanMsg(writeInMyClientChan, msg, true, instance)
					continue
				}
				fmt.Printf("Client: Instance n°%d there is no traitment for %d ID\n", instance, weft.PackId)
			}
		}
	}
}
