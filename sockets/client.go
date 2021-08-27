package sockets

import (
	"GoDofus/managers"
	"GoDofus/messages"
	"GoDofus/pack"
	"GoDofus/settings"
	"bytes"
	"fmt"
	"log"
	"net"
	"time"
)

var connServer *net.TCPConn
var listener net.Listener
var connListener net.Conn = nil

var stop bool
var Callback func([]byte, int)
var Address string
var currentAddress string
var blockServerRead bool

func reloadConnListenerWrite(msg messages.Message) {
	if listener == nil {
		return
	}

	tryReloadConnListener(time.Second * 10)

	_, err := connListener.Write(pack.Write(msg, true))
	if err != nil {
		panic(err)
	}
}

func reloadConnListenerRead(lecture []byte) int {
	if listener == nil {
		return 0
	}

	tryReloadConnListener(time.Second * 10)

	n, err := connListener.Read(lecture)
	if err != nil {
		panic(err)
	}

	return n
}

func writeInMyClient(msg messages.Message, waitResponses []int) {
	_, err := connListener.Write(pack.Write(msg, true))
	if err != nil {
		reloadConnListenerWrite(msg)
	}

	if waitResponses != nil {
		readInListener(waitResponses)
	}
}

func readInListener(responses []int) {
	var packetIds = make(map[uint16]bool)
	pipe := pack.GetClientPipeline()

	for _, packetId := range responses {
		packetIds[uint16(packetId)] = true
	}

	blockServerRead = true

	lecture := make([]byte, 1024)
	fmt.Println("go Lecture listener")

	for connListener != nil {
		n, err := connListener.Read(lecture)

		if err != nil {
			n = reloadConnListenerRead(lecture)
		}

		if n == 0 {
			continue
		}
		fmt.Printf("Listener: %d octets reçu\n", n)

		ok := pack.ReadClient(lecture[:n])

		if len(responses) == 0 && ok {
			fmt.Println("end lecture listener")
			handlingListener()
			break
		} else if pipe.Contains(packetIds) {
			fmt.Println("end lecture listener")
			handlingListener()
			break
		}
	}

	blockServerRead = false
}

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
		default:
			fmt.Printf("Client: there is no traitment for %d ID\n", weft.PackId)
		}
	}
}

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
			writeInMyClient(msg, nil)
			fmt.Println(msg)
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

func launchServerSocket() {
	if connListener != nil {
		panic("a listener connexion is already active")
	}

	var err error

	if listener == nil {
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", settings.Settings.LocalAddress, settings.Settings.LocalPort))
		if err != nil {
			log.Fatal(err)
		}
	}

	connListener, err = listener.Accept()
	if err != nil {
		if connListener != nil {
			err = connListener.Close()
			if err != nil {
			}
			connListener = nil
		}

		return
	}
	fmt.Println("Listener: Go client !")

	for connListener != nil {
	}

	fmt.Println("Listener: Go client lost !")
}

func tryReloadConnListener(duration time.Duration) {
	connListener = nil

	fmt.Println("Try connect a listener...")

	if connListener == nil {
		go launchServerSocket()
	}

	done := make(chan bool)
	go func() {
		time.Sleep(duration)
		done <- true
	}()

	for connListener == nil {
		select {
		case <-done:
			if listener != nil {
				err := listener.Close()
				listener = nil
				if err != nil {
					panic(err)
				}
			}
			if connServer != nil {
				err := connServer.Close()
				connServer = nil
				if err != nil {
					panic(err)
				}
			}
			log.Fatal("Listener: Kill instant ! From Timeout.")
		default:
			continue
		}
	}
}

func waitMyClient() {
	if connListener != nil {
		return
	}

	if connListener == nil {
		go launchServerSocket()
	}

	for connListener == nil {
		continue
	}
}

func LaunchClientSocket() {
	waitMyClient()

	rAddr, err := net.ResolveTCPAddr("tcp", Address)
	connServer, err = net.DialTCP("tcp", nil, rAddr)

	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	} else {
		log.Println("Connexion to server OK.")
	}
	currentAddress = Address
	stop = false

	defer func(conn_ net.Conn) {
		if connServer != nil {
			err = connServer.Close()
			if err != nil {
				log.Fatal(err)
			}
			connServer = nil
		}

		if Address == currentAddress {
			if listener == nil {
				return
			}

			if connListener != nil {
				err = connListener.Close()
				if err != nil {
				}
				connListener = nil
			}

			err = listener.Close()
			if err != nil {
			}
			listener = nil
		} else {
			LaunchClientSocket()
		}
	}(connServer)

	lecture := make([]byte, 1024)

	for connServer != nil {
		if stop {
			break
		}

		waitMyClient()

		if blockServerRead {
			continue
		}

		n, err := connServer.Read(lecture)

		if err != nil {
			panic(err)
		}

		if n == 0 {
			continue
		}
		fmt.Printf("Server: %d octets reçu\n", n)

		Callback(lecture, n)
	}
}
