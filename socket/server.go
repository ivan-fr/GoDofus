package socket

import (
	"GoDofus/pack"
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
		var lecture []byte
		n, err := conn.Read(lecture)

		if err != nil {
			log.Fatal(err)
		}

		ok := pack.Read(lecture)

		if ok {
			fmt.Printf("%d octet lues\nétat du pipeline:\n%v", n, pack.GetPipeline())
		}
	}
}
