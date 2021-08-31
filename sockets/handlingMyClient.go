package sockets

import (
	"GoDofus/messages"
	"GoDofus/pack"
	"bytes"
	"fmt"
)

func handlingMyClient(writeInMyClientChan, writeToOfficialServerChan chan []byte, myClientContinueChan, officialServerContinueChan chan bool, instance uint) func(chan *pack.Weft) {
	return func(weftChan chan *pack.Weft) {
		for {
			weft := <-weftChan

			if weft == nil {
				break
			}

			switch weft.PackId {
			case messages.IdentificationID:
				continue
			case messages.AuthenticationTicketID:
				msg := messages.Types_[int(weft.PackId)].GetNOA(instance)
				sendChanMsg(writeToOfficialServerChan, msg, false, instance)
			default:
				msg, ok := messages.Types_[int(weft.PackId)]

				if ok {
					msg = msg.GetNOA(instance)
					msg.Deserialize(bytes.NewReader(weft.Message))
					fmt.Println(msg)
					sendChanMsg(writeToOfficialServerChan, msg, false, instance)
					continue
				}

				sendChanWeft(writeToOfficialServerChan, weft, false, instance)
				fmt.Printf("Listener: Instance n°%d there is no traitment for %d ID\nNatural Weft sended.\n", instance, weft.PackId)
			}
		}
	}
}
