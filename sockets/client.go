package sockets

import (
	"GoDofus/managers"
	"GoDofus/messages"
	"GoDofus/pack"
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"net"
	"os"
	"time"
)

var conn *net.TCPConn
var listener net.Listener
var connListener net.Conn = nil
var toogleContinueToTryConnListener bool

var stop bool
var Callback func([]byte, int)
var Address string
var currentAddress string
var blockServerRead bool

func reloadConnListenerWrite(msg messages.Message) {
	if listener == nil {
		return
	}

	connListener = nil

	tryConnListener(time.Second * 2)

	_, err := connListener.Write(pack.Write(msg, true))
	if err != nil {
		panic(err)
	}
}

func reloadConnListenerRead(lecture []byte) int {
	if listener == nil {
		return 0
	}

	connListener = nil

	tryConnListener(time.Second * 2)

	n, err := connListener.Read(lecture)
	if err != nil {
		panic(err)
	}

	return n
}

func writeInListener(msg messages.Message, waitResponse bool) {
	_, err := connListener.Write(pack.Write(msg, true))
	if err != nil {
		reloadConnListenerWrite(msg)
	}

	if waitResponse {
		readInListener()
	}
}

func readInListener() {
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

		ok := pack.ReadClient(lecture[:n])

		if ok {
			handlingListener(n)
			break
		}
	}

	fmt.Println("end lecture listener")

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

			writeInListener(msg, true)

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

type Server struct {
	Address string `yaml:"localServer"`
}

var myServer = getConf()

func getConf() *Server {
	var server = &Server{}

	yamlFile, err := os.ReadFile("./settings.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, server)
	if err != nil {
		panic(err)
	}

	return server
}

func launchServerSocket() {
	if connListener != nil {
		panic("un connexion listener est déjà active")
	}

	var err error

	if listener == nil {
		listener, err = net.Listen("tcp", myServer.Address)
		if err != nil {
			log.Fatal(err)
		}
	}

	for toogleContinueToTryConnListener {
		connListener, err = listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Listener: Go client !")
		break
	}

	if !toogleContinueToTryConnListener {
		toogleContinueToTryConnListener = true
		return
	}

	for connListener != nil {
	}

	fmt.Println("Listener: Go client lost !")
}

func tryConnListener(duration time.Duration) {
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
			toogleContinueToTryConnListener = false
			break
		default:
			continue
		}
	}

	if connListener == nil {
		panic("tryConnListener: timout")
	}
}

func waitMyClient() {
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
	conn, err = net.DialTCP("tcp", nil, rAddr)

	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	} else {
		log.Println("La connexion to server OK.")
	}
	currentAddress = Address
	stop = false

	defer func(conn_ net.Conn) {
		err = conn.Close()
		if err != nil {
			log.Fatal(err)
		}
		conn = nil

		if Address == currentAddress {
			err = connListener.Close()
			if err != nil {
			}
			connListener = nil

			err = listener.Close()
			if err != nil {
				return
			}
			listener = nil
		} else {
			time.Sleep(time.Millisecond * 500)
		}
	}(conn)

	lecture := make([]byte, 1024)

	for {
		if stop {
			break
		}

		waitMyClient()

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
