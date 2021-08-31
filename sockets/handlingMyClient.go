package sockets

import (
	"GoDofus/commands"
	"GoDofus/messages"
	"GoDofus/pack"
	"bytes"
	"fmt"
	"log"
	"time"
)

func handlingMyClient(writeInMyClientChan, writeToOfficialServerChan chan []byte, myClientContinueChan, officialServerContinueChan chan bool, instance uint) func(chan *pack.Weft) {
	chanCrypt := make(chan []byte)
	go func(chanCrypt chan []byte) {
		for {
			currentMap := messages.Types_[messages.CurrentMapID].GetNOA(instance).(*messages.CurrentMap)
			log.Println(currentMap.MapId)
			move := commands.Move{MapId: currentMap.MapId}
			move.SetFromMapId()

			var aCommand string
			fmt.Println("Command!")
			_, err := fmt.Scanln(&aCommand)
			if err != nil {
				continue
			}
			switch aCommand {
			case "d":
				move.X += 1
			case "q":
				move.X -= 1
			case "s":
				move.Y += 1
			case "z":
				move.Y -= 1
			default:
				continue
			}

			if move.SetFromCoords() {
				msg := messages.Types_[messages.ChangeMapID].GetNOA(instance).(*messages.ChangeMap)
				msg.MapId = currentMap.MapId
				log.Println(currentMap.MapId)
				sendChanMsg(writeToOfficialServerChan, msg, false, instance)
				time.Sleep(time.Second * 2)
				msg2 := messages.Types_[messages.MapInformationsRequestID].GetNOA(instance).(*messages.MapInformationsRequest)
				msg2.MapId = move.MapId
				sendChanMsg(writeToOfficialServerChan, msg2, false, instance)
				log.Println(move.MapId)
				fmt.Println("Command! OK!")
			} else {
				fmt.Println("Command! No ok!")
			}
		}
	}(chanCrypt)

	return func(weftChan chan *pack.Weft) {
		for {
			weft := <-weftChan

			if weft == nil {
				break
			}

			switch weft.PackId {
			case messages.RawDataMessageID:
				chanCrypt <- weft.Message
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
					sendChanMsg(writeToOfficialServerChan, msg, false, instance)
					continue
				}

				sendChanWeft(writeToOfficialServerChan, weft, false, instance)
			}
		}
	}
}
