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
var Callback func([]byte, int)

var Address string
var currentAddress string

var blockServerReadLinear bool
var blockLinearReadGo = make(chan bool)
var blockGoBlock = make(chan bool)

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

	blockServerReadLinear = true
	blockGoBlock <- true

	lecture := make([]byte, 1024)
	fmt.Println("go Lecture listener")

	blockGo := false
	for connListener != nil {
		select {
		case blockGo = <-blockLinearReadGo:
		default:
		}

		if blockGo {
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

	blockServerReadLinear = false
	blockGoBlock <- false
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

	msg := messages.GetBasicPongNOA()
	packToWrite := pack.Write(msg, true)

	block := false
	for connListener != nil {
		select {
		case block = <-blockGoBlock:
		default:
		}

		if block {
			continue
		}

		time.Sleep(time.Second * 2)
		blockLinearReadGo <- true
		_, err = connListener.Write(packToWrite)

		if err != nil {
			connListener = nil
		}

		blockLinearReadGo <- false
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

	blockGo := false
	for connServer != nil {
		if stop {
			break
		}

		waitMyClient()

		select {
		case blockGo = <-blockLinearReadGo:
		default:
		}

		if blockServerReadLinear || blockGo {
			continue
		}

		blockGoBlock <- true
		n, err := connServer.Read(lecture)
		blockGoBlock <- false

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
