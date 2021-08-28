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
	writeTChanChan := make(chan [2]chan messages.Message)

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

func channelWriter(aChanMessage chan messages.Message, aChanConnexion chan net.Conn, toClient bool, instance uint) {
	aConn := <-aChanConnexion

	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("ChannelWriter: closed: %v, to client: %t, for conn: %v", err, toClient, aConn)
		}
	}()

	for {
		select {
		case msg := <-aChanMessage:
			if aConn == nil {
				aConn = <-aChanConnexion
			}
			_, err := aConn.Write(pack.Write(msg, toClient, instance))
			bug := handleErrReadWrite(err)

			if bug {
				break
			}
		case aConn = <-aChanConnexion:
		}
	}
}

func loginListener(wg *sync.WaitGroup,
	instanceChan chan uint,
	connToMyClientChanChan chan chan net.Conn,
	writeTChanChan chan [2]chan messages.Message) {
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

			writeInMyClientChan := make(chan messages.Message)
			writeToOfficialServerChan := make(chan messages.Message)

			myClientContinueChan := make(chan bool)
			officialServerContinueChan := make(chan bool)

			go channelWriter(writeInMyClientChan, myConnToMyClientChan, true, instance)
			myConnToMyClientChan <- myConnToMyClient

			go func() {
				var myWg sync.WaitGroup
				myWg.Add(1)
				callbackOfficialServer := handlingAuth(writeInMyClientChan, writeToOfficialServerChan, officialServerContinueChan, instance)
				go launchLoginClientToOfficialSocket(&myWg, writeToOfficialServerChan, officialServerContinueChan, callbackOfficialServer, instance, myConnToOfficialChan)
				myWg.Wait()
				callbackOfficialServer = handlingGame(writeInMyClientChan, writeToOfficialServerChan, officialServerContinueChan, instance)
				go launchGameClientToOfficialSocket(nil, officialServerContinueChan, callbackOfficialServer, instance, myConnToOfficialChan)
			}()

			go func() {
				callbackInMyClient := handlingMyClient(writeInMyClientChan, writeToOfficialServerChan, officialServerContinueChan, instance)

				var myWg sync.WaitGroup
				myWg.Add(1)
				go launchServerForMyClientSocket(nil, myConnToMyClient, myClientContinueChan, callbackInMyClient, instance)
				myWg.Wait()

				log.Println("myConnToMyClientChan is CLOSE")
				myConnToMyClientChan <- nil
				_ = myConnToMyClient.Close()
				log.Println("myConnToMyClientChan is CLOSE")

				instanceChan <- instance
				connToMyClientChanChan <- myConnToMyClientChan
				writeTChanChan <- [2]chan messages.Message{writeInMyClientChan, writeToOfficialServerChan}
			}()
		}()
	}
}

func gameListener(wg *sync.WaitGroup,
	instanceChan chan uint,
	connToMyClientChanChan chan chan net.Conn,
	writeTChanChan chan [2]chan messages.Message) {
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

			callbackInMyClient := handlingMyClient(writeChan[0], writeChan[1], officialServerContinueChan, instance)

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
	myReadServer func([]byte) bool, myPipeline *pack.Pipe,
	officialServerContinueChan chan bool,
	callBack func(pipe *pack.Pipe), instance uint) {

	defer func(conn_ net.Conn) {
		_ = conn_.Close()
	}(myConnServer)

	myLecture := make([]byte, 1024)

	next := true
	for next {
		select {
		case next = <-officialServerContinueChan:
		default:
		}

		n, err := myConnServer.Read(myLecture)
		bug := handleErrReadWrite(err)

		if bug {
			break
		}

		if n == 0 {
			continue
		}
		fmt.Printf("Server n°%d: %d octets reçu\n", instance, n)

		_ = myReadServer(myLecture[:n])

		if len(myPipeline.Wefts) > 0 {
			callBack(myPipeline)
		}
	}

	log.Printf("Dial: server official lost from instance n°%d !\n", instance)
}

func launchGameClientToOfficialSocket(wg *sync.WaitGroup,
	officialServerContinueChan chan bool,
	callBack func(pipe *pack.Pipe),
	instance uint, myConnToOfficialChan chan net.Conn) {

	defer func() {
		myConnToOfficialChan <- nil
	}()
	if wg != nil {
		defer wg.Done()
	}

	myReadServer, myPipeline := pack.NewServerReader()

	selectedServerDataExtended := messages.GetSelectedServerDataExtendedNOA(instance)
	gameRAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", selectedServerDataExtended.SSD.Address, selectedServerDataExtended.SSD.Ports[0]))
	if err != nil {
		log.Println(err)
		return
	}
	myConnServer, err := net.DialTCP("tcp", nil, gameRAddr)

	if err != nil {
		log.Printf("Failed to dial: %v", err)
		return
	} else {
		log.Printf("Connexion to server instance n°%d OK.\n", instance)
	}

	myConnToOfficialChan <- myConnServer
	factoryServerClientToOfficial(myConnServer, myReadServer, myPipeline, officialServerContinueChan, callBack, instance)
}

func launchLoginClientToOfficialSocket(wg *sync.WaitGroup,
	writeToOfficialServerChan chan messages.Message,
	officialServerContinueChan chan bool,
	callBack func(pipe *pack.Pipe),
	instance uint, myConnToOfficialChan chan net.Conn) {

	defer func() {
		myConnToOfficialChan <- nil
	}()
	if wg != nil {
		defer wg.Done()
	}

	myReadServer, myPipeline := pack.NewServerReader()
	myConnServer, err := net.DialTCP("tcp", nil, rAddr)

	if err != nil {
		log.Printf("Failed to dial with server for instance n°%d OK.\n", instance)
		return
	} else {
		log.Printf("Connexion to server for instance n°%d OK.\n", instance)
	}

	go channelWriter(writeToOfficialServerChan, myConnToOfficialChan, false, instance)
	myConnToOfficialChan <- myConnServer

	factoryServerClientToOfficial(myConnServer, myReadServer, myPipeline, officialServerContinueChan, callBack, instance)
}

func launchServerForMyClientSocket(wg *sync.WaitGroup, myConnToMyClient net.Conn, myClientContinueChan chan bool, callBack func(pipe *pack.Pipe), instance uint) {
	if wg != nil {
		defer wg.Done()
	}

	myReadClient, myPipeline := pack.NewClientReader()

	lecture := make([]byte, 1024)
	next := true
	for next {
		select {
		case next = <-myClientContinueChan:
		default:
		}

		err := myConnToMyClient.SetWriteDeadline(time.Now().Add(100 * time.Millisecond))
		if err != nil {
			break
		}
		_, err = myConnToMyClient.Write(pack.Write(messages.GetBasicPongNOA(instance), true, instance))
		if err != nil {
			break
		}

		bug := handleErrReadWrite(err)
		if bug {
			break
		}

		err = myConnToMyClient.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		if err != nil {
			return
		}
		n, err := myConnToMyClient.Read(lecture)
		bug = handleErrReadWrite(err)

		if bug {
			break
		}

		if n > 0 {
			_ = myReadClient(lecture[:n])
			if len(myPipeline.Wefts) > 0 {
				callBack(myPipeline)
			}
		}
	}
	log.Printf("Listener: Go client lost from instance n°%d !\n", instance)
}

func handleErrReadWrite(err error) bool {
	if _, ok := err.(net.Error); ok {
		return true
	}

	return false
}
