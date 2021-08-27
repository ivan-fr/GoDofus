package sockets

import (
	"GoDofus/messages"
	"GoDofus/pack"
	"GoDofus/settings"
	"fmt"
	"io"
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
var blockThreadServerToToMyClient = make(chan bool)

func writeInMyClient(msg messages.Message) {
	if connListener == nil {
		tryReloadConnListener(time.Second * 8)
	}

	blockClientToAnkamaLinear = true
	blockThreadServerToToMyClient <- true

	_ = connListener.SetWriteDeadline(time.Now().Add(time.Second))
	_, _ = connListener.Write(pack.Write(msg, true))

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

	lecture := make([]byte, 1024)
	for connListener != nil {

		if block {
			continue
		}

		_ = connListener.SetReadDeadline(time.Now().Add(time.Second * 1))
		blockClientToAnkamaLinear = true
		n, err := connListener.Read(lecture)
		if errNet, ok := err.(net.Error); ok {
			if !errNet.Timeout() {
				err := connListener.Close()
				if err != nil {
				}
				connListener = nil
			}
		} else if err == io.EOF {
			err := connListener.Close()
			if err != nil {
			}
			connListener = nil
		}

		if n > 0 {
			_ = pack.ReadClient(lecture[:n])
			if len(pack.GetClientPipeline().Wefts) > 0 {
				handlingMyClient()
			}
		}
		blockClientToAnkamaLinear = false
	}

	fmt.Println("Listener: Go client lost !")
}

func tryReloadConnListener(duration time.Duration) {
	if connListener != nil {
		return

	}

	fmt.Println("Try connect a listener...")
	go launchServerSocket()

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
				log.Println("ICI")
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
		n, err := connServer.Read(lecture)
		if err != nil {
			panic(err)
		}

		if n == 0 {
			continue
		}
		fmt.Printf("Server: %d octets reçu\n", n)

		_ = pack.ReadServer(lecture[:n])
		blockThreadServerToToMyClient <- false

		if len(pack.GetServerPipeline().Wefts) > 0 {
			Callback()
		}
	}
}
