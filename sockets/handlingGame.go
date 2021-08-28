package sockets

import (
	"GoDofus/messages"
	"GoDofus/pack"
	"bytes"
	"fmt"
	"time"
)

func handlingGame(writeInMyClientChan, writeToAnkamaServerChan chan messages.Message, myClientContinueChan, ankamaServerContinueChan chan bool) func(*pack.Pipe) {
	return func(pipe *pack.Pipe) {
		for weft := pipe.Get(); weft != nil; weft = pipe.Get() {
			switch weft.PackId {
			case messages.ProtocolID:
				msg := messages.GetProtocolNOA()
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
			case messages.HelloGameID:
				msg := messages.GetHelloGameNOA()
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				time.Sleep(time.Millisecond * 150)
				msg2 := messages.GetAuthenticationTicketNOA()
				writeToAnkamaServerChan <- msg2
			case messages.RawDataID:
				msg := messages.GetRawDataNOA()
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			case messages.AuthenticationTicketAcceptedID:
				msg := messages.GetAuthenticationTicketAcceptedNOA()
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				writeInMyClientChan <- msg
			default:
				fmt.Printf("Client: there is no traitment for %d ID\n", weft.PackId)
			}
		}
	}
}
