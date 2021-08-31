package sockets

import (
	"GoDofus/messages"
	"GoDofus/pack"
	"bytes"
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
					sendChanMsg(writeInMyClientChan, msg, true, instance)
					continue
				}

				sendChanWeft(writeInMyClientChan, weft, true, instance)
			}
		}
	}
}
