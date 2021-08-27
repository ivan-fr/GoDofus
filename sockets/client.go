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
var connListener net.Conn = nil
var stop bool
var Callback func([]byte, int)
var Address string
var currentAddress string
var blockServerRead bool

func writeInListener(msg messages.Message, waitResponse bool) {
	_, err := connListener.Write(pack.Write(msg, true))
	if err != nil {
		panic(err)
	}

	if waitResponse {
		readInListener()
	}
}

func readInListener() {
	blockServerRead = true

	lecture := make([]byte, 1024)
	fmt.Println("Lecture listener")

	for connListener != nil {
		n, err := connListener.Read(lecture)

		if err != nil {
			panic(err)
		}

		if n == 0 {
			continue
		}

		ok := pack.ReadClient(lecture[:n])

		if ok {
			handlingListener(n)
			break
		}
	}

	fmt.Println("fin lecture listener")

	blockServerRead = false
}

func handlingListener(n int) {
	fmt.Printf("Listener: %d octets reçu\n", n)

	pipe := pack.GetClientPipeline()
	for weft := pipe.Get(); weft != nil; weft = pipe.Get() {
		switch weft.PackId {
		case messages.CheckIntegrityID:
			msg := messages.GetCheckIntegrityNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			_, err := conn.Write(pack.Write(msg, false))
			if err != nil {
				panic(err)
			}
		case messages.ClientKeyID:
			msg := messages.GetClientKeyNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			_, err := conn.Write(pack.Write(msg, false))
			if err != nil {
				panic(err)
			}
		default:
			fmt.Printf("Listener: there is no traitment for %d ID\n", weft.PackId)
		}
	}
}

func handlingGame(lecture []byte, n int) {
	fmt.Printf("%d octets reçu\n", n)

	ok := pack.ReadServer(lecture[:n])

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
			_, err := conn.Write(pack.Write(msg2, false))
			if err != nil {
				panic(err)
			}
		case messages.RawDataID:
			msg := messages.GetRawDataNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			writeInListener(msg, true)
		default:
			fmt.Printf("Client: there is no traitment for %d ID\n", weft.PackId)
		}
	}
}

func HandlingAuth(lecture []byte, n int) {
	fmt.Printf("%d octets reçu\n", n)

	ok := pack.ReadServer(lecture[:n])

	if !ok {
		return
	}

	pipe := pack.GetPipeline()
	for weft := pipe.Get(); weft != nil; weft = pipe.Get() {
		switch weft.PackId {
		case messages.HelloConnectID:
			msg := messages.GetHelloConnectNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)

			writeInListener(msg, false)
			readInListener()

			fmt.Println("======= GO Identification =======")
			mAuth := managers.GetAuthentification()
			mAuth.InitIdentificationMessage()

			authMessage := messages.GetIdentificationNOA()
			_, err := conn.Write(pack.Write(authMessage, false))
			if err != nil {
				panic(err)
			}
		case messages.ProtocolID:
			msg := messages.GetProtocolNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			writeInListener(msg, false)
			fmt.Println(msg)
		case messages.IdentificationFailedForBadVersionID:
			msg := messages.GetIdentificationFailedForBadVersionNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			writeInListener(msg, false)
		case messages.IdentificationFailedID:
			msg := messages.GetIdentificationFailedNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			writeInListener(msg, false)
		case messages.LoginQueueID:
			msg := messages.GetLoginQueueStatusNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			writeInListener(msg, false)
		case messages.IdentificationSuccessID:
			msg := messages.GetIdentificationSuccessNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			writeInListener(msg, false)
		case messages.SelectedServerDataExtendedID:
			msg := messages.GetSelectedServerDataExtendedNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			writeInListener(msg, false)
			stop = true
			Address = fmt.Sprintf("%s:%d", msg.SSD.Address, msg.SSD.Ports[0])
			Callback = handlingGame
		case messages.CredentialsAcknowledgementID:
			msg := messages.GetCredentialsAcknowledgementNOA()
			msg.Deserialize(bytes.NewReader(weft.Message))
			fmt.Println(msg)
			writeInListener(msg, false)
		default:
			fmt.Printf("Client: there is no traitment for %d ID\n", weft.PackId)
		}
	}
}

func LaunchServerSocket() {
	if connListener != nil {
		panic("un connexion listener est déjà active")
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
	}(l)
	for {
		connListener, err = l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Listener: connexion avec l'esclave angagé")
		break
	}

	for connListener != nil {
	}

	fmt.Println("Listener: connexion avec l'esclave perdu.")
}

func LaunchClientSocket() {
	if connListener == nil {
		go LaunchServerSocket()
	}

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
			connListener = nil
		}
	}(conn)

	lecture := make([]byte, 1024)

	for {
		if stop {
			break
		}

		if blockServerRead {
			continue
		}

		n, err := conn.Read(lecture)

		if err != nil {
			panic(err)
		}

		if n == 0 {
			continue
		}

		Callback(lecture, n)
	}
}
