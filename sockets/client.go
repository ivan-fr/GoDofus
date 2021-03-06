package sockets

import (
	"GoDofus/messages"
	"GoDofus/pack"
	"GoDofus/settings"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func GoSocket() {
	instanceChan := make(chan uint)
	connToMyClientChanChan := make(chan chan net.Conn)
	writeTChanChan := make(chan [2]chan []byte)

	var wg sync.WaitGroup
	wg.Add(2)
	go loginListener(&wg, instanceChan, connToMyClientChanChan, writeTChanChan)
	go gameListener(&wg, instanceChan, connToMyClientChanChan, writeTChanChan)
	wg.Wait()
}

var rAddr = getRAddr()

func getRAddr() *net.TCPAddr {
	rAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", settings.Settings.ServerAnkamaAddress, settings.Settings.ServerAnkamaPort))
	if err != nil {
		panic(err)
	}

	return rAddr
}

func channelWriter(aChanMessage chan []byte, aChanConnexion chan net.Conn, instance uint) {
	aConn := <-aChanConnexion

	for {
		select {
		case msg := <-aChanMessage:
			if aConn == nil {
				aConn = <-aChanConnexion
				log.Printf("Writer instance n°%d updated\n", instance)
			}

			_, err := aConn.Write(msg)
			if netErr, ok := err.(net.Error); ok {
				if !netErr.Timeout() {
					panic(netErr)
				}
				continue
			}
		case aConn = <-aChanConnexion:
			log.Printf("Writer instance n°%d updated\n", instance)
		}
	}
}

func loginListener(wg *sync.WaitGroup,
	instanceChan chan uint,
	connToMyClientChanChan chan chan net.Conn,
	writeTChanChan chan [2]chan []byte) {
	defer wg.Done()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", settings.Settings.LocalAddress, settings.Settings.LocalLoginPort))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Login listener ready", fmt.Sprintf("%s:%d", settings.Settings.LocalAddress, settings.Settings.LocalLoginPort))

	defer func(listener net.Listener) {
		_ = listener.Close()
		log.Println("Login listener close.")
	}(listener)

	var instance uint
	for {
		myConnToMyClient, err := listener.Accept()
		if err != nil {
			log.Println(err)
			break
		}

		instance++
		log.Printf("Login listener: connexion instance n°%d\n", instance)

		go func() {
			myConnToMyClientChan := make(chan net.Conn)
			myConnToOfficialChan := make(chan net.Conn)

			writeInMyClientChan := make(chan []byte)
			writeToOfficialServerChan := make(chan []byte)

			myClientContinueChan := make(chan bool)
			officialServerContinueChan := make(chan bool)

			go channelWriter(writeInMyClientChan, myConnToMyClientChan, instance)
			myConnToMyClientChan <- myConnToMyClient

			go func() {
				var myWg sync.WaitGroup
				myWg.Add(1)
				callbackOfficialServer := handlingAuth(writeInMyClientChan, writeToOfficialServerChan, myClientContinueChan, officialServerContinueChan, instance)
				go launchLoginClientToOfficialSocket(&myWg, writeToOfficialServerChan, officialServerContinueChan, callbackOfficialServer, instance, myConnToOfficialChan)
				myWg.Wait()
				callbackOfficialServer = handlingGame(writeInMyClientChan, writeToOfficialServerChan, myClientContinueChan, officialServerContinueChan, instance)
				go launchGameClientToOfficialSocket(nil, officialServerContinueChan, callbackOfficialServer, instance, myConnToOfficialChan)
			}()

			go func() {
				callbackInMyClient := handlingMyClient(writeInMyClientChan, writeToOfficialServerChan, myClientContinueChan, officialServerContinueChan, instance)

				var myWg sync.WaitGroup
				myWg.Add(1)
				go launchServerForMyClientSocket(&myWg, myConnToMyClient, myClientContinueChan, callbackInMyClient, instance)
				myWg.Wait()

				myConnToMyClientChan <- nil
				_ = myConnToMyClient.Close()

				instanceChan <- instance
				connToMyClientChanChan <- myConnToMyClientChan
				writeTChanChan <- [2]chan []byte{writeInMyClientChan, writeToOfficialServerChan}
			}()
		}()
	}
}

func gameListener(wg *sync.WaitGroup,
	instanceChan chan uint,
	connToMyClientChanChan chan chan net.Conn,
	writeTChanChan chan [2]chan []byte) {
	defer wg.Done()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", settings.Settings.LocalAddress, settings.Settings.LocalGamePort))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Game listener ready", fmt.Sprintf("%s:%d", settings.Settings.LocalAddress, settings.Settings.LocalGamePort))

	defer func(listener net.Listener) {
		_ = listener.Close()
		log.Println("Game listener close.")
	}(listener)

	for {
		myConnToMyClient, err := listener.Accept()
		if err != nil {
			log.Println(err)
			break
		}

		log.Printf("Game listener: a connexion was accepted !\n")

		instance := <-instanceChan
		myConnToMyClientChan := <-connToMyClientChanChan
		myConnToMyClientChan <- myConnToMyClient
		writeChan := <-writeTChanChan

		log.Printf("Game listener: connexion instance n°%d\n", instance)

		go func() {
			myClientContinueChan := make(chan bool)
			officialServerContinueChan := make(chan bool)

			callbackInMyClient := handlingMyClient(writeChan[0], writeChan[1], myClientContinueChan, officialServerContinueChan, instance)

			var myWg sync.WaitGroup
			myWg.Add(1)
			go launchServerForMyClientSocket(&myWg, myConnToMyClient, myClientContinueChan, callbackInMyClient, instance)
			myWg.Wait()

			myConnToMyClientChan <- nil

			_ = myConnToMyClient.Close()
		}()
	}
}

func factoryServerClientToOfficial(myConnServer net.Conn,
	myReadServer *pack.Reader,
	officialServerContinueChan chan bool,
	callBack func(chan *pack.Weft), instance uint) {

	defer func(conn_ net.Conn) {
		_ = conn_.Close()
	}(myConnServer)

	myWeftChan := make(chan *pack.Weft)
	myPipelineChan := make(chan bool)

	go callBack(myWeftChan)
	go func(pipelineChan chan bool) {
		for weft := myReadServer.APipeline.Get(); ; weft = myReadServer.APipeline.Get() {
			select {
			case <-pipelineChan:
				return
			default:
				if weft != nil {
					myWeftChan <- weft
				}
			}
		}

	}(myPipelineChan)

	myLecture := make([]byte, 256)

	next := true
	for next {
		select {
		case next = <-officialServerContinueChan:
		default:
		}

		_ = myConnServer.SetReadDeadline(time.Now().Add(time.Millisecond))
		n, err := myConnServer.Read(myLecture)
		if netErr, ok := err.(net.Error); ok {
			if !netErr.Timeout() {
				break
			}
		}

		if n == 0 {
			continue
		}

		myReadServer.Read(false, myLecture[:n])
	}

	myWeftChan <- nil
	myPipelineChan <- false
	log.Printf("Dial: server official lost from instance n°%d !\n", instance)
}

func launchGameClientToOfficialSocket(wg *sync.WaitGroup,
	officialServerContinueChan chan bool,
	callBack func(chan *pack.Weft),
	instance uint, myConnToOfficialChan chan net.Conn) {

	defer func() {
		myConnToOfficialChan <- nil
	}()
	if wg != nil {
		defer wg.Done()
	}

	myReadServer := pack.NewReader()

	selectedServerDataExtended := messages.Types_[messages.SelectedServerDataExtendedID].GetNOA(instance).(*messages.SelectedServerDataExtended)

	var gameRAddr *net.TCPAddr
	var err error

	gameRAddr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", selectedServerDataExtended.SSD.Address, selectedServerDataExtended.SSD.Ports[0]))

	if err != nil {
		log.Println(err)
		return
	}
	myConnServer, err := net.DialTCP("tcp", nil, gameRAddr)

	if err != nil {
		log.Printf("Failed to dial: %v", err)
		return
	} else {
		log.Printf("Connexion to official server GAME instance n°%d OK.\n", instance)
	}

	myConnToOfficialChan <- myConnServer
	factoryServerClientToOfficial(myConnServer, myReadServer, officialServerContinueChan, callBack, instance)
}

func launchLoginClientToOfficialSocket(wg *sync.WaitGroup,
	writeToOfficialServerChan chan []byte,
	officialServerContinueChan chan bool,
	callBack func(chan *pack.Weft),
	instance uint, myConnToOfficialChan chan net.Conn) {

	defer func() {
		myConnToOfficialChan <- nil
	}()
	if wg != nil {
		defer wg.Done()
	}

	myReadServer := pack.NewReader()
	myConnServer, err := net.DialTCP("tcp", nil, rAddr)

	if err != nil {
		log.Printf("Failed to dial with server for instance n°%d OK.\n", instance)
		return
	} else {
		log.Printf("Connexion to official server LOGIN for instance n°%d OK.\n", instance)
	}

	go channelWriter(writeToOfficialServerChan, myConnToOfficialChan, instance)
	myConnToOfficialChan <- myConnServer

	factoryServerClientToOfficial(myConnServer, myReadServer, officialServerContinueChan, callBack, instance)
}

func launchServerForMyClientSocket(wg *sync.WaitGroup, myConnToMyClient net.Conn, myClientContinueChan chan bool, callBack func(chan *pack.Weft), instance uint) {
	if wg != nil {
		defer wg.Done()
	}

	myReadClient := pack.NewReader()

	myWeftChan := make(chan *pack.Weft)
	pipelineChan := make(chan bool)

	go callBack(myWeftChan)
	go func(pipelineChan chan bool) {
		for weft := myReadClient.APipeline.Get(); ; weft = myReadClient.APipeline.Get() {
			select {
			case <-pipelineChan:
				return
			default:
				if weft != nil {
					myWeftChan <- weft
				}
			}
		}
	}(pipelineChan)

	myLecture := make([]byte, 256)

	next := true
	for next {
		select {
		case next = <-myClientContinueChan:
		default:
		}

		_ = myConnToMyClient.SetReadDeadline(time.Now().Add(time.Millisecond))
		n, err := myConnToMyClient.Read(myLecture)
		if netErr, ok := err.(net.Error); ok {
			if !netErr.Timeout() {
				break
			}
		}

		if n == 0 {
			continue
		}

		myReadClient.Read(true, myLecture[:n])
	}

	myWeftChan <- nil
	pipelineChan <- false
	log.Printf("Listener: Go client lost from instance n°%d !\n", instance)
}
