package sockets

import (
	"GoDofus/managers"
	"GoDofus/messages"
	"GoDofus/pack"
	"bytes"
	"fmt"
)

func handlingAuth(writeInMyClientChan, writeToAnkamaServerChan chan messages.Message, myClientContinueChan, ankamaServerContinueChan chan bool) func(*pack.Pipe) {
	return func(pipe *pack.Pipe) {
		for weft := pipe.Get(); weft != nil; weft = pipe.Get() {
			switch weft.PackId {
			case messages.HelloConnectID:
				msg := messages.GetHelloConnectNOA()
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)

				writeInMyClientChan <- msg

				fmt.Println("======= GO Identification =======")
				mAuth := managers.GetAuthentification()
				mAuth.InitIdentificationMessage()

				authMessage := messages.GetIdentificationNOA()
				writeToAnkamaServerChan <- authMessage
			case messages.ProtocolID:
				msg := messages.GetProtocolNOA()
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			case messages.IdentificationFailedForBadVersionID:
				msg := messages.GetIdentificationFailedForBadVersionNOA()
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			case messages.IdentificationFailedID:
				msg := messages.GetIdentificationFailedNOA()
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			case messages.LoginQueueID:
				msg := messages.GetLoginQueueStatusNOA()
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			case messages.IdentificationSuccessID:
				msg := messages.GetIdentificationSuccessNOA()
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			case messages.SelectedServerDataExtendedID:
				msg := messages.GetSelectedServerDataExtendedNOA()
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
				ankamaServerContinueChan <- false
			case messages.CredentialsAcknowledgementID:
				msg := messages.GetCredentialsAcknowledgementNOA()
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			default:
				fmt.Printf("Client: there is no traitment for %d ID\n", weft.PackId)
			}
		}
	}
}
