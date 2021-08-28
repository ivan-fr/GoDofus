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
				msg := messages.GetHelloConnectNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)

				writeInMyClientChan <- msg

				fmt.Println("======= GO Identification =======")
				mAuth := managers.GetAuthentificationManager(instance)
				mAuth.InitIdentificationMessage()

				authMessage := messages.GetIdentificationNOA(instance)
				writeToOfficialServerChan <- authMessage
			case messages.ProtocolID:
				msg := messages.GetProtocolNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			case messages.IdentificationFailedForBadVersionID:
				msg := messages.GetIdentificationFailedForBadVersionNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			case messages.IdentificationFailedID:
				msg := messages.GetIdentificationFailedNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			case messages.LoginQueueStatueID:
				msg := messages.GetLoginQueueStatusNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			case messages.IdentificationSuccessID:
				msg := messages.GetIdentificationSuccessNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			case messages.SelectedServerDataExtendedID:
				msg := messages.GetSelectedServerDataExtendedNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
				myClientContinueChan <- false
				officialServerContinueChan <- false
			case messages.CredentialsAcknowledgementID:
				msg := messages.GetCredentialsAcknowledgementNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			case messages.AuthenticationTicketAcceptedID:
				msg := messages.GetAuthenticationTicketAcceptedNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			default:
				fmt.Printf("Client: there is no traitment for %d ID\n", weft.PackId)
			}
		}
	}
}
