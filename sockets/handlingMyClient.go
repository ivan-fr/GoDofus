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
				msg2 := messages.GetAuthenticationTicketNOA(instance)
				writeToOfficialServerChan <- msg2
			case messages.HaapiApiKeyRequestID:
				msg2 := messages.GetHaapiApiKeyRequestNOA(instance)
				writeToOfficialServerChan <- msg2
			case messages.CharactersListRequestID:
				msg := messages.GetCharactersListRequestNOA(instance)
				writeToOfficialServerChan <- msg
			case messages.CheckIntegrityID:
				msg := messages.GetCheckIntegrityNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeToOfficialServerChan <- msg
			case messages.ClientKeyID:
				msg := messages.GetClientKeyNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeToOfficialServerChan <- msg
			case messages.CharacterSelectionID:
				msg := messages.GetCharacterSelectionNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeToOfficialServerChan <- msg
			default:
				fmt.Printf("Listener: there is no traitment for %d ID\n", weft.PackId)
			}
		}
	}
}
