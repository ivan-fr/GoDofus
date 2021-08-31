package sockets

import (
	"GoDofus/messages"
	"GoDofus/pack"
	"bytes"
)

func handlingAuth(writeInMyClientChan, writeToOfficialServerChan chan []byte, myClientContinueChan, officialServerContinueChan chan bool, instance uint) func(chan *pack.Weft) {
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
				sendChanMsg(writeInMyClientChan, msg, true, instance)
				idMessage := messages.Types_[messages.IdentificationID].GetNOA(instance).(*messages.Identification)
				idMessage.InitIdentificationMessage()
				sendChanMsg(writeToOfficialServerChan, idMessage, false, instance)
			case messages.SelectedServerDataExtendedID:
				msg := messages.Types_[int(weft.PackId)].GetNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				writeInMyClientChan <- pack.Write(msg, true, instance)
				myClientContinueChan <- false
				officialServerContinueChan <- false
			default:
				msg, ok := messages.Types_[int(weft.PackId)]

				if ok {
					msg = msg.GetNOA(instance)
					msg.Deserialize(bytes.NewReader(weft.Message))
					sendChanMsg(writeInMyClientChan, msg, true, instance)
					continue
				}

				sendChanWeft(writeInMyClientChan, weft, true, instance)
			}
		}
	}
}
