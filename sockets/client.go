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

var rAddr = getRAddr()

func getRAddr() *net.TCPAddr {
	rAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", settings.Settings.ServerAnkamaAddress, settings.Settings.ServerAnkamaPort))
	if err != nil {
		panic(err)
	}

	return rAddr
}

func channelWriter(aChan chan messages.Message, aConn net.Conn, toClient bool) {
	for {
		msg := <-aChan
		_, err := aConn.Write(pack.Write(msg, toClient))
		closed := handleErrReadWrite(aConn, err)

		if closed {
			break
		}
	}
}

func syncStop(myClientContinueChan, ankamaServerContinueChan chan bool) {
	verif := make([]bool, 2)
	for {
		select {
		case next := <-myClientContinueChan:
			if !next {
				verif[0] = true
			}
		case next := <-ankamaServerContinueChan:
			if !next {
				verif[1] = true
			}
		}

		var ok bool

		for i := 0; i < len(verif); i++ {
			if !verif[i] {
				ok = false
				break
			}
		}

		if ok {
			break
		}
	}
}

func loginListener(instanceChan chan uint) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", settings.Settings.LocalAddress, settings.Settings.LocalLoginPort))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Login listener ready.")

	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)

	var instance uint
	for {
		myConnToMyClient, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		instance++
		log.Printf("Login listener: connexion instance n°%d\n", instance)

		go func() {
			myClientContinueChan := make(chan bool)
			ankamaServerContinueChan := make(chan bool)

			writeInMyClientChan := make(chan messages.Message)
			writeToAnkamaServerChan := make(chan messages.Message)

			callback := HandlingAuth(writeInMyClientChan, writeToAnkamaServerChan, myClientContinueChan, ankamaServerContinueChan)

			go channelWriter(writeInMyClientChan, myConnToMyClient, true)
			go launchServerForMyClientSocket(myConnToMyClient, myClientContinueChan)
			go launchLoginClientToAnkamaSocket(writeToAnkamaServerChan, ankamaServerContinueChan, callback, instance)
			syncStop(myClientContinueChan, ankamaServerContinueChan)
			instanceChan <- instance
			_ = myConnToMyClient.Close()
		}()
	}
}

func gameListener(instanceChan chan uint) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", settings.Settings.LocalAddress, settings.Settings.LocalGamePort))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Game listener ready.")

	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)

	var instance uint
	for {
		myConnToMyClient, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		instance = <-instanceChan
		log.Printf("Game listener: connexion instance n°%d\n", instance)

		go func() {
			myClientContinueChan := make(chan bool)
			ankamaServerContinueChan := make(chan bool)

			writeInMyClientChan := make(chan messages.Message)
			writeToAnkamaServerChan := make(chan messages.Message)

			callback := HandlingAuth(writeInMyClientChan, writeToAnkamaServerChan, myClientContinueChan, ankamaServerContinueChan)

			go channelWriter(writeInMyClientChan, myConnToMyClient, true)
			go launchServerForMyClientSocket(myConnToMyClient, myClientContinueChan)
			go launchGameClientToAnkamaSocket(writeToAnkamaServerChan, ankamaServerContinueChan, callback, instance)
			syncStop(myClientContinueChan, ankamaServerContinueChan)
			_ = myConnToMyClient.Close()
		}()
	}
}

func GoSocket() {
	instanceChan := make(chan uint)
	go loginListener(instanceChan)
	go gameListener(instanceChan)
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

	go channelWriter(writeToAnkamaServerChan, myConnServer, false)

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
}

func launchGameClientToAnkamaSocket(writeToAnkamaServerChan chan messages.Message, ankamaServerContinueChan chan bool, callBack func(pipe *pack.Pipe), instance uint) {
	myReadServer, myPipeline := pack.NewServerReader()
	myConnServer, err := net.DialTCP("tcp", nil, rAddr)
	factoryServerClientToAnkama(myConnServer, err, myReadServer, myPipeline, writeToAnkamaServerChan, ankamaServerContinueChan, callBack, instance)
}

func launchLoginClientToAnkamaSocket(writeToAnkamaServerChan chan messages.Message, ankamaServerContinueChan chan bool, callBack func(pipe *pack.Pipe), instance uint) {
	myReadServer, myPipeline := pack.NewServerReader()
	myConnServer, err := net.DialTCP("tcp", nil, rAddr)
	factoryServerClientToAnkama(myConnServer, err, myReadServer, myPipeline, writeToAnkamaServerChan, ankamaServerContinueChan, callBack, instance)
}

func launchServerForMyClientSocket(myConnToMyClient net.Conn, myClientContinueChan chan bool) {
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
				handlingMyClient()
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
