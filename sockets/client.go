package sockets

import (
	"GoDofus/managers"
	"GoDofus/messages"
	"GoDofus/pack"
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

var conn net.Conn

func handling(lecture []byte, n int) {
	fmt.Printf("%d octets reçu\n", n)

	ok := pack.Read(lecture[:n])

	if ok {
		pipe := pack.GetPipeline()
		for weft := pipe.Get(); weft != nil; weft = pipe.Get() {
			switch weft.PackId {
			case messages.HelloConnectID:
				hConnect := messages.GetHelloConnectNOA()
				hConnect.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(hConnect)

				fmt.Println("======= GO Identification =======")
				mAuth := managers.GetAuthentification()
				mAuth.InitIdentificationMessage()

				authMessage := messages.GetIdentificationNOA()
				_, err := conn.Write(pack.Write(authMessage))
				if err != nil {
					panic(err)
				}
				clientKey := messages.GetClientKeyNOA()
				_, err = conn.Write(pack.Write(clientKey))
				if err != nil {
					panic(err)
				}
			case messages.ProtocolID:
				protocol := messages.GetProtocolNOA()
				protocol.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(protocol)
			case messages.IdentificationFailedForBadVersionID:
				msg := messages.GetIdentificationFailedForBadVersionNOA()
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
			case messages.IdentificationFailedID:
				msg := messages.GetIdentificationFailedNOA()
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
			case messages.LoginQueueID:
				msg := messages.GetLoginQueueStatusNOA()
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
			case messages.IdentificationSuccessID:
				msg := messages.GetIdentificationSuccessNOA()
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
			case messages.SelectedServerDataExtendedID:
				msg := messages.GetSelectedServerDataExtendedNOA()
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
			case messages.CredentialsAcknowledgementID:
				msg := messages.GetCredentialsAcknowledgementNOA()
				msg.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(msg)
			default:
				fmt.Printf("there is no traitment for %d ID\n", weft.PackId)
			}
		}
	} else {
		fmt.Println("paquet incomplet")
	}
}

func LaunchClientSocket() {
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var err error
	conn, err = d.DialContext(ctx, "tcp", "52.17.231.202:5555")
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	} else {
		log.Println("La connexion au serveur d'authentification est réussie.")
	}
	defer func(conn_ net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
		conn = nil
	}(conn)

	for {
		lecture := make([]byte, 1024)
		n, err := conn.Read(lecture)

		if err != nil {
			log.Fatal(err)
		}

		if n == 0 {
			continue
		}

		handling(lecture, n)
	}
}
