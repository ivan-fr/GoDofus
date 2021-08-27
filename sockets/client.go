package sockets

import (
	"GoDofus/messages"
	"GoDofus/pack"
	"GoDofus/settings"
	"fmt"
	"log"
	"net"
	"time"
)

var connServer *net.TCPConn
var listener net.Listener
var connListener net.Conn = nil

var stop bool
var Callback func()

var Address string
var currentAddress string

var blockClientToAnkamaLinear bool
var blockServerToMyClientLinear bool
var blockThreadServerToToMyClient = make(chan bool)

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
	blockClientToAnkamaLinear = true
	blockThreadServerToToMyClient <- true
	_, err := connListener.Write(pack.Write(msg, true))
	if err != nil {
		reloadConnListenerWrite(msg)
	}
	blockClientToAnkamaLinear = false
	blockThreadServerToToMyClient <- false

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

	blockClientToAnkamaLinear = true
	blockThreadServerToToMyClient <- true

	lecture := make([]byte, 1024)
	fmt.Println("go Lecture listener")

	for connListener != nil {
		if blockServerToMyClientLinear {
			continue
		}

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

	blockClientToAnkamaLinear = false
	blockThreadServerToToMyClient <- false
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

	block := false

	go func() {
		for {
			select {
			case block = <-blockThreadServerToToMyClient:
			default:
			}
			if connListener == nil {
				break
			}
		}
	}()

	lecture := make([]byte, 1)
	for connListener != nil {
		time.Sleep(time.Second)

		if block {
			continue
		}

		_ = connListener.SetReadDeadline(time.Now().Add(time.Second * 3))

		blockServerToMyClientLinear = true
		blockClientToAnkamaLinear = true
		n, err := connListener.Read(lecture)
		if err != nil {
			err = connListener.Close()
			if err != nil {
			}
			connListener = nil
		}

		if n > 0 {
			_ = pack.ReadClient(lecture[:n])
		}
		blockServerToMyClientLinear = false
		blockClientToAnkamaLinear = false
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
				_ = listener.Close()
				listener = nil
			}
			if connServer != nil {
				_ = connServer.Close()
				connServer = nil
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
	} else {
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

		if connListener == nil {
			tryReloadConnListener(time.Second * 8)
		}

		if blockClientToAnkamaLinear {
			continue
		}

		blockThreadServerToToMyClient <- true
		blockServerToMyClientLinear = true
		n, err := connServer.Read(lecture)
		blockThreadServerToToMyClient <- false
		blockServerToMyClientLinear = false

		if err != nil {
			panic(err)
		}

		if n == 0 {
			continue
		}
		fmt.Printf("Server: %d octets reçu\n", n)

		ok := pack.ReadServer(lecture[:n])

		if ok {
			Callback()
		}
	}
}
