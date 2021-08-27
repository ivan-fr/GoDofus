package sockets

import (
	"GoDofus/messages"
	"GoDofus/pack"
	"bytes"
	"fmt"
)

func handlingListener() {
	pipe := pack.GetClientPipeline()
	for weft := pipe.Get(); weft != nil; weft = pipe.Get() {
		switch weft.PackId {
		case messages.CheckIntegrityID:
			msg := messages.GetCheckIntegrityNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			_, err := connServer.Write(pack.Write(msg, false))
			if err != nil {
				panic(err)
			}
		case messages.ClientKeyID:
			msg := messages.GetClientKeyNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			_, err := connServer.Write(pack.Write(msg, false))
			if err != nil {
				panic(err)
			}
		default:
			fmt.Printf("Listener: there is no traitment for %d ID\n", weft.PackId)
		}
	}
}
