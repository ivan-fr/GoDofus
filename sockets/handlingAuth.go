package sockets

import (
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

				go sendChanMsg(writeInMyClientChan, msg)

				fmt.Println("======= GO Identification =======")
				idMessage := messages.Types_[messages.IdentificationID].GetNOA(instance).(*messages.Identification)
				idMessage.InitIdentificationMessage()
				go sendChanMsg(writeToOfficialServerChan, idMessage)
			case messages.SelectedServerDataExtendedID:
				msg := messages.Types_[int(weft.PackId)].GetNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				go sendChanMsg(writeInMyClientChan, msg)
				go sendChanBool(myClientContinueChan, false)
				go sendChanBool(officialServerContinueChan, false)
			default:
				msg, ok := messages.Types_[int(weft.PackId)]

				if ok {
					msg = msg.GetNOA(instance)
					msg.Deserialize(bytes.NewReader(weft.Message))
					fmt.Println(msg)
					go sendChanMsg(writeInMyClientChan, msg)
					continue
				}

				fmt.Printf("Client: Instance n°%d there is no traitment for %d ID\n", instance, weft.PackId)
			}
		}
	}
}
