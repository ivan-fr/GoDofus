package sockets

import (
	"GoDofus/managers"
	"GoDofus/messages"
	"GoDofus/pack"
	"bytes"
	"fmt"
)

func handlingAuth(writeInMyClientChan, writeToOfficialServerChan chan messages.Message, myClientContinueChan, officialServerContinueChan chan bool, instance uint) func(chan *pack.Weft) {
	return func(weftChan chan *pack.Weft) {
		for {
			weft := <-weftChan

			if weft == nil {
				break
			}

			switch weft.PackId {
			case messages.HelloConnectID:
				msg := messages.Types_[int(weft.PackId)].GetNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)

				writeInMyClientChan <- msg

				fmt.Println("======= GO Identification =======")
				mAuth := managers.GetAuthentificationManager(instance)
				mAuth.InitIdentificationMessage()

				authMessage := messages.Types_[messages.IdentificationID].GetNOA(instance)
				writeToOfficialServerChan <- authMessage
			case messages.SelectedServerDataExtendedID:
				msg := messages.Types_[int(weft.PackId)].GetNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
				myClientContinueChan <- false
				officialServerContinueChan <- false
			default:
				msg, ok := messages.Types_[int(weft.PackId)]

				if ok {
					msg = msg.GetNOA(instance)
					msg.Deserialize(bytes.NewReader(weft.Message))
					fmt.Println(msg)
					writeInMyClientChan <- msg
					continue
				}

				fmt.Printf("Client: Instance n°%d there is no traitment for %d ID\n", instance, weft.PackId)
			}
		}
	}
}
