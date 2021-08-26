package sockets

import (
	"GoDofus/managers"
	"GoDofus/messages"
	"GoDofus/pack"
	"bytes"
	"fmt"
	"log"
	"net"
	"time"
)

var conn *net.TCPConn
var connListener net.Conn
var stop bool
var Callback func([]byte, int)
var Address string
var currentAddress string

func handlingListener(lecture []byte, n int) {
	fmt.Printf("Listener: %d octets reçu\n", n)

	ok := pack.Read(lecture[:n])

	if !ok {
		return
	}

	pipe := pack.GetPipeline()
	for weft := pipe.Get(); weft != nil; weft = pipe.Get() {
		switch weft.PackId {
		case messages.CheckIntegrityID:
			msg := messages.GetCheckIntegrityNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			_, err := conn.Write(pack.Write(msg))
			if err != nil {
				panic(err)
			}
		default:
			fmt.Printf("Listener: there is no traitment for %d ID\n", weft.PackId)
		}
	}
}

func HandlingGame(lecture []byte, n int) {
	fmt.Printf("%d octets reçu\n", n)

	ok := pack.Read(lecture[:n])

	if !ok {
		return
	}

	pipe := pack.GetPipeline()
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
			_, err := conn.Write(pack.Write(msg2))
			if err != nil {
				panic(err)
			}
		case messages.RawDataID:
			msg := messages.GetRawDataNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			_, err := connListener.Write(pack.Write(msg))
			if err != nil {
				panic(err)
			}
		default:
			fmt.Printf("Client: there is no traitment for %d ID\n", weft.PackId)
		}
	}
}

func HandlingAuth(lecture []byte, n int) {
	fmt.Printf("%d octets reçu\n", n)

	ok := pack.Read(lecture[:n])

	if !ok {
		return
	}

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
			Callback = HandlingGame
		case messages.CredentialsAcknowledgementID:
			msg := messages.GetCredentialsAcknowledgementNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
		default:
			fmt.Printf("Client: there is no traitment for %d ID\n", weft.PackId)
		}
	}
}

func LaunchServerSocket() {
	if connListener != nil {
		fmt.Println("un connexion listener est déjà active")
	}

	l, err := net.Listen("tcp", "127.0.0.1:443")
	if err != nil {
		log.Fatal(err)
	}
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			log.Fatal(err)
		}
		LaunchServerSocket()
	}(l)
	for {
		connListener, err = l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		break
	}

	lecture := make([]byte, 1024)

	for connListener != nil {
		n, err := conn.Read(lecture)

		if err != nil {
			log.Fatal(err)
		}

		if n == 0 {
			continue
		}

		handlingListener(lecture, n)
	}

	err = l.Close()
	if err != nil {
		panic(err)
	}
}

func LaunchClientSocket() {
	LaunchServerSocket()

	fmt.Println("Client esclave connecté...")

	for connListener == nil {
		continue
	}

	rAddr, err := net.ResolveTCPAddr("tcp", Address)
	conn, err = net.DialTCP("tcp", nil, rAddr)

	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	} else {
		log.Println("La connexion au serveur est réussie.")
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
		} else {
			err := connListener.Close()
			if err != nil {
			}
		}
	}(conn)

	lecture := make([]byte, 1024)

	for {
		if stop {
			break
		}

		n, err := conn.Read(lecture)

		if err != nil {
			log.Fatal(err)
		}

		if n == 0 {
			continue
		}

		Callback(lecture, n)
	}
}
