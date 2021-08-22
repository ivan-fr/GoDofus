package socket

import (
	"GoDofus/messages"
	"GoDofus/pack"
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

func LaunchClientSocket() {
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", "34.252.21.81:5555")
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	} else {
		log.Println("La connexion au serveur d'authentification est réussie.")
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	for {
		lecture := make([]byte, 1024)
		n, err := conn.Read(lecture)

		if err != nil {
			log.Fatal(err)
		}

		if n == 0 {
			continue
		}

		fmt.Printf("%d bytes reçu\n", n)

		ok := pack.Read(lecture[:n])

		if ok {
			pipe := pack.GetPipeline()
			lastWeft := pipe.Wefts[len(pipe.Wefts)-1]

			if lastWeft.PackId == 1030 {
				hConnect := messages.GetHelloConnectNOA()
				hConnect.Deserialize(bytes.NewReader(lastWeft.Message))
				fmt.Println(hConnect)
			} else if lastWeft.PackId == 9546 {
				protocol := messages.GetProtocolNOA()
				protocol.Deserialize(bytes.NewReader(lastWeft.Message))
				fmt.Println(protocol.Version)
			}
		} else {
			fmt.Println("paquet incomplet")
		}
	}
}
