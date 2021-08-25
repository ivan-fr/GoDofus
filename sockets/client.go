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
var stop bool
var Callback func([]byte, int)
var Address string
var currentAddress string

func handlingGame(lecture []byte, n int) {
	fmt.Printf("%d octets reçu\n", n)

	ok := pack.Read(lecture[:n])

	if ok {
		pipe := pack.GetPipeline()
		for weft := pipe.Get(); weft != nil; weft = pipe.Get() {
			switch weft.PackId {
			case messages.ProtocolID:
				protocol := messages.GetProtocolNOA()
				protocol.Deserialize(bytes.NewReader(weft.Message))
				fmt.Println(protocol)
			default:
				fmt.Printf("there is no traitment for %d ID\n", weft.PackId)
			}
		}
	} else {
		fmt.Println("paquet incomplet")
	}
}

func HandlingAuth(lecture []byte, n int) {
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
				stop = true
				Address = fmt.Sprintf("%s:%d", msg.SSD.Address, msg.SSD.Ports[0])
				Callback = handlingGame
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
	conn, err = d.DialContext(ctx, "tcp", Address)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	} else {
		log.Println("La connexion au serveur d'authentification est réussie.")
	}
	currentAddress = Address
	stop = false

	defer func(conn_ net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
		conn = nil
		if Address != currentAddress {
			LaunchClientSocket()
		}
	}(conn)

	for {
		lecture := make([]byte, 1024)
		n, err := conn.Read(lecture)

		if err != nil {
			log.Fatal(err)
		}

		if stop {
			break
		}

		if n == 0 {
			continue
		}

		Callback(lecture, n)
	}
}
