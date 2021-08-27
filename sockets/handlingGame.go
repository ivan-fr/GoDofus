package sockets

import (
	"GoDofus/messages"
	"GoDofus/pack"
	"bytes"
	"fmt"
	"time"
)

func handlingGame(lecture []byte, n int) {
	ok := pack.ReadServer(lecture[:n])

	if !ok {
		return
	}

	pipe := pack.GetServerPipeline()
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
			_, err := connServer.Write(pack.Write(msg2, false))
			if err != nil {
				panic(err)
			}
		case messages.RawDataID:
			msg := messages.GetRawDataNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			writeInMyClient(msg, []int{messages.CheckIntegrityID})
		case messages.AuthenticationTicketAcceptedID:
			msg := messages.GetAuthenticationTicketAcceptedNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			writeInMyClient(msg, nil)
		default:
			fmt.Printf("Client: there is no traitment for %d ID\n", weft.PackId)
		}
	}
}
