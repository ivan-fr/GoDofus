package sockets

import (
	"GoDofus/messages"
	"GoDofus/pack"
	"bytes"
	"fmt"
)

func handlingMyClient(writeInMyClientChan, writeToOfficialServerChan chan messages.Message, myClientContinueChan, officialServerContinueChan chan bool, instance uint) func(chan *pack.Weft) {
	return func(weftChan chan *pack.Weft) {
		for {
			weft := <-weftChan

			if weft == nil {
				break
			}

			switch weft.PackId {
			case messages.IdentificationID:
				continue
			case messages.CharactersListRequestID, messages.AuthenticationTicketID, messages.HaapiApiKeyRequestID:
				msg := messages.Types_[int(weft.PackId)].GetNOA(instance)
				go sendChanMsg(writeToOfficialServerChan, msg)
			default:
				msg, ok := messages.Types_[int(weft.PackId)]

				if ok {
					msg = msg.GetNOA(instance)
					msg.Deserialize(bytes.NewReader(weft.Message))
					fmt.Println(msg)
					go sendChanMsg(writeToOfficialServerChan, msg)
					continue
				}
				fmt.Printf("Listener: Instance n°%d there is no traitment for %d ID\n", instance, weft.PackId)
			}
		}
	}
}
