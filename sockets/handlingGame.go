package sockets

import (
	"GoDofus/messages"
	"GoDofus/pack"
	"bytes"
	"fmt"
	"time"
)

func handlingGame(writeInMyClientChan, writeToOfficialServerServerChan chan messages.Message, myClientContinueChan, officialServerContinueChan chan bool, instance uint) func(chan *pack.Weft) {
	return func(weftChan chan *pack.Weft) {
		for {
			weft := <-weftChan

			if weft == nil {
				break
			}
			switch weft.PackId {
			case messages.ProtocolID:
				msg := messages.GetProtocolNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
			case messages.HelloGameID:
				msg := messages.GetHelloGameNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
				time.Sleep(time.Millisecond * 150)
				msg2 := messages.GetAuthenticationTicketNOA(instance)
				writeToOfficialServerServerChan <- msg2
			case messages.RawDataID:
				msg := messages.GetRawDataNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			case messages.BasicTimeID:
				msg := messages.GetBasicTimeNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			case messages.ServerSettingsID:
				msg := messages.GetServerSettingsNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			case messages.ServerOptionalFeaturesID:
				msg := messages.GetServerOptionalFeaturesNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			case messages.ServerSessionConstantsID:
				msg := messages.GetServerSessionConstantsNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			case messages.AccountCapabilitiesID:
				msg := messages.GetAccountCapabilitiesNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			case messages.TrustStatusID:
				msg := messages.GetTrustStatusNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			case messages.HaapiSessionID:
				msg := messages.GetHaapiSessionNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			case messages.CharactersListRequestID:
				msg := messages.GetCharactersListRequestNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			default:
				fmt.Printf("Client: there is no traitment for %d ID\n", weft.PackId)
			}
		}
	}
}
