package sockets

import (
	"GoDofus/messages"
	"GoDofus/pack"
	"GoDofus/settings"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

func GoSocket() {
	instanceChan := make(chan uint)

	var wg sync.WaitGroup
	wg.Add(2)
	go loginListener(&wg, instanceChan)
	go gameListener(&wg, instanceChan)
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

func channelWriter(wg *sync.WaitGroup, aChan chan messages.Message, aConn net.Conn, toClient bool) {
	defer wg.Done()

	for {
		msg := <-aChan
		_, err := aConn.Write(pack.Write(msg, toClient))
		closed := handleErrReadWrite(aConn, err)

		if closed {
			break
		}
	}
}

func loginListener(wg *sync.WaitGroup, instanceChan chan uint) {
	defer wg.Done()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", settings.Settings.LocalAddress, settings.Settings.LocalLoginPort))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Login listener ready.")

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
			myClientContinueChan := make(chan bool)
			ankamaServerContinueChan := make(chan bool)

			writeInMyClientChan := make(chan messages.Message)
			writeToAnkamaServerChan := make(chan messages.Message)

			callbackAnkamaServer := handlingAuth(writeInMyClientChan, writeToAnkamaServerChan, myClientContinueChan, ankamaServerContinueChan, instance)
			callbackInMyClient := handlingMyClient(writeInMyClientChan, writeToAnkamaServerChan, myClientContinueChan, ankamaServerContinueChan, instance)

			var myWg sync.WaitGroup
			myWg.Add(3)
			go channelWriter(&myWg, writeInMyClientChan, myConnToMyClient, true)
			go launchServerForMyClientSocket(&myWg, myConnToMyClient, myClientContinueChan, callbackInMyClient)
			go launchLoginClientToAnkamaSocket(&myWg, writeToAnkamaServerChan, ankamaServerContinueChan, callbackAnkamaServer, instance)
			myWg.Wait()

			instanceChan <- instance
			_ = myConnToMyClient.Close()
		}()
	}
}

func gameListener(wg *sync.WaitGroup, instanceChan chan uint) {
	defer wg.Done()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", settings.Settings.LocalAddress, settings.Settings.LocalGamePort))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Game listener ready.")

	defer func(listener net.Listener) {
		_ = listener.Close()
		log.Println("Game listener close.")
	}(listener)

	var instance uint
	for {
		myConnToMyClient, err := listener.Accept()
		if err != nil {
			log.Println(err)
			break
		}
		instance = <-instanceChan
		log.Printf("Game listener: connexion instance n°%d\n", instance)

		go func() {
			myClientContinueChan := make(chan bool)
			ankamaServerContinueChan := make(chan bool)

			writeInMyClientChan := make(chan messages.Message)
			writeToAnkamaServerChan := make(chan messages.Message)

			callbackAnkamaServer := handlingAuth(writeInMyClientChan, writeToAnkamaServerChan, myClientContinueChan, ankamaServerContinueChan, instance)
			callbackInMyClient := handlingMyClient(writeInMyClientChan, writeToAnkamaServerChan, myClientContinueChan, ankamaServerContinueChan, instance)

			var myWg sync.WaitGroup
			myWg.Add(3)
			go channelWriter(&myWg, writeInMyClientChan, myConnToMyClient, true)
			go launchServerForMyClientSocket(&myWg, myConnToMyClient, myClientContinueChan, callbackInMyClient)
			go launchGameClientToAnkamaSocket(&myWg, writeToAnkamaServerChan, ankamaServerContinueChan, callbackAnkamaServer, instance)
			myWg.Wait()

			_ = myConnToMyClient.Close()
		}()
	}
}

func factoryServerClientToAnkama(myConnServer net.Conn, err error, myReadServer func([]byte) bool, myPipeline *pack.Pipe,
	writeToAnkamaServerChan chan messages.Message,
	ankamaServerContinueChan chan bool,
	callBack func(pipe *pack.Pipe), instance uint) {

	if err != nil {
		log.Printf("Failed to dial: %v", err)
		return
	} else {
		log.Printf("Connexion to server instance n°%d OK.\n", instance)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go channelWriter(&wg, writeToAnkamaServerChan, myConnServer, false)

	defer func(conn_ net.Conn) {
		_ = myConnServer.Close()
	}(myConnServer)

	myLecture := make([]byte, 1024)

	next := true
	for next {
		select {
		case next = <-ankamaServerContinueChan:
			ankamaServerContinueChan <- false
		default:
		}

		err := myConnServer.SetReadDeadline(time.Now().Add(time.Second * 2))
		if err != nil {
			break
		}
		n, err := myConnServer.Read(myLecture)
		closed := handleErrReadWrite(myConnServer, err)

		if closed {
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

	wg.Wait()
}

func launchGameClientToAnkamaSocket(wg *sync.WaitGroup, writeToAnkamaServerChan chan messages.Message, ankamaServerContinueChan chan bool, callBack func(pipe *pack.Pipe), instance uint) {
	defer wg.Done()
	myReadServer, myPipeline := pack.NewServerReader()

	selectedServerDataExtended := messages.GetSelectedServerDataExtendedNOA(instance)
	gameRAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", selectedServerDataExtended.SSD.Address, selectedServerDataExtended.SSD.Ports[0]))
	if err != nil {
		log.Println(err)
		return
	}
	myConnServer, err := net.DialTCP("tcp", nil, gameRAddr)
	factoryServerClientToAnkama(myConnServer, err, myReadServer, myPipeline, writeToAnkamaServerChan, ankamaServerContinueChan, callBack, instance)
}

func launchLoginClientToAnkamaSocket(wg *sync.WaitGroup, writeToAnkamaServerChan chan messages.Message, ankamaServerContinueChan chan bool, callBack func(pipe *pack.Pipe), instance uint) {
	defer wg.Done()
	myReadServer, myPipeline := pack.NewServerReader()
	myConnServer, err := net.DialTCP("tcp", nil, rAddr)
	factoryServerClientToAnkama(myConnServer, err, myReadServer, myPipeline, writeToAnkamaServerChan, ankamaServerContinueChan, callBack, instance)
}

func launchServerForMyClientSocket(wg *sync.WaitGroup, myConnToMyClient net.Conn, myClientContinueChan chan bool, callBack func(pipe *pack.Pipe)) {
	defer wg.Done()

	myReadClient, myPipeline := pack.NewClientReader()

	lecture := make([]byte, 1024)
	next := true
	for next {
		select {
		case next = <-myClientContinueChan:
			myClientContinueChan <- false
		default:
		}
		_ = myConnToMyClient.SetReadDeadline(time.Now().Add(time.Millisecond * 500))
		n, err := myConnToMyClient.Read(lecture)
		closed := handleErrReadWrite(myConnToMyClient, err)

		if n > 0 {
			_ = myReadClient(lecture[:n])
			if len(myPipeline.Wefts) > 0 {
				callBack(myPipeline)
			}
		}

		if closed {
			break
		}
	}
	log.Println("Listener: Go client lost !")
}

func handleErrReadWrite(conn net.Conn, err error) bool {
	if errNet, ok := err.(net.Error); ok {
		if !errNet.Timeout() {
			_ = conn.Close()
			return true
		}
	} else if err == io.EOF {
		_ = conn.Close()
		return true
	}

	return false
}
