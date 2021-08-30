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
			case messages.AuthenticationTicketID:
				msg2 := messages.Types_[int(weft.PackId)].GetNOA(instance)
				writeToOfficialServerChan <- msg2
			case messages.HaapiApiKeyRequestID:
				msg2 := messages.Types_[int(weft.PackId)].GetNOA(instance)
				writeToOfficialServerChan <- msg2
			case messages.CharactersListRequestID:
				msg := messages.Types_[int(weft.PackId)].GetNOA(instance)
				writeToOfficialServerChan <- msg
			case messages.CheckIntegrityID:
				msg := messages.Types_[int(weft.PackId)].GetNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeToOfficialServerChan <- msg
			case messages.ClientKeyID:
				msg := messages.Types_[int(weft.PackId)].GetNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeToOfficialServerChan <- msg
			case messages.CharacterSelectionID:
				msg := messages.Types_[int(weft.PackId)].GetNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeToOfficialServerChan <- msg
			default:
				msg, ok := messages.Types_[int(weft.PackId)]

				if ok {
					msg = msg.GetNOA(instance)
					msg.Deserialize(bytes.NewReader(weft.Message))
					fmt.Println(msg)
					writeToOfficialServerChan <- msg
					return
				}
				fmt.Printf("Listener: Instance n°%d there is no traitment for %d ID\n", instance, weft.PackId)
			}
		}
	}
}
