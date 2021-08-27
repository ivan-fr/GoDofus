package sockets

import (
	"GoDofus/managers"
	"GoDofus/messages"
	"GoDofus/pack"
	"bytes"
	"fmt"
)

func HandlingAuth(lecture []byte, n int) {
	ok := pack.ReadServer(lecture[:n])

	if !ok {
		return
	}

	pipe := pack.GetServerPipeline()
	for weft := pipe.Get(); weft != nil; weft = pipe.Get() {
		switch weft.PackId {
		case messages.HelloConnectID:
			msg := messages.GetHelloConnectNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)

			writeInMyClient(msg, []int{messages.ClientKeyID})

			fmt.Println("======= GO Identification =======")
			mAuth := managers.GetAuthentification()
			mAuth.InitIdentificationMessage()

			authMessage := messages.GetIdentificationNOA()
			_, err := connServer.Write(pack.Write(authMessage, false))
			if err != nil {
				panic(err)
			}
		case messages.ProtocolID:
			msg := messages.GetProtocolNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			writeInMyClient(msg, nil)
		case messages.IdentificationFailedForBadVersionID:
			msg := messages.GetIdentificationFailedForBadVersionNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			writeInMyClient(msg, nil)
		case messages.IdentificationFailedID:
			msg := messages.GetIdentificationFailedNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			writeInMyClient(msg, nil)
		case messages.LoginQueueID:
			msg := messages.GetLoginQueueStatusNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			writeInMyClient(msg, nil)
		case messages.IdentificationSuccessID:
			msg := messages.GetIdentificationSuccessNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			writeInMyClient(msg, nil)
		case messages.SelectedServerDataExtendedID:
			msg := messages.GetSelectedServerDataExtendedNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			writeInMyClient(msg, nil)
			stop = true
			Address = fmt.Sprintf("%s:%d", msg.SSD.Address, msg.SSD.Ports[0])
			Callback = handlingGame
		case messages.CredentialsAcknowledgementID:
			msg := messages.GetCredentialsAcknowledgementNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			writeInMyClient(msg, nil)
		default:
			fmt.Printf("Client: there is no traitment for %d ID\n", weft.PackId)
		}
	}
}
