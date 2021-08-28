package sockets

import (
	"GoDofus/messages"
	"GoDofus/pack"
	"bytes"
	"fmt"
	"time"
)

func handlingGame(writeInMyClientChan, writeToAnkamaServerChan chan messages.Message, ankamaServerContinueChan chan bool, instance uint) func(*pack.Pipe) {
	return func(pipe *pack.Pipe) {
		for weft := pipe.Get(); weft != nil; weft = pipe.Get() {
			switch weft.PackId {
			case messages.ProtocolID:
				msg := messages.GetProtocolNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
			case messages.HelloGameID:
				msg := messages.GetHelloGameNOA(instance)
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
				time.Sleep(time.Millisecond * 150)
				msg2 := messages.GetAuthenticationTicketNOA(instance)
				writeToAnkamaServerChan <- msg2
			case messages.RawDataID:
				msg := messages.GetRawDataNOA(instance)
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
