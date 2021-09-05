package main

import (
	"GoDofus/generates"
	"GoDofus/parsers"
	"GoDofus/sockets"
	"flag"
	"log"
)

func main() {
	rsa_ := flag.Bool("get_rsa", false, "Génére le private/public key et la signature DofusPublicKey. (-get_rsa=true/false)")
	hello := flag.Bool("get_hello", false, "Génére le private/public key pour helloConnect. (-get_hello=true/false)")
	hosts := flag.String("hosts", "", "Génére la signature des hosts pour config.xml. Entrez localhost,127.0.0.1 par exemple.")
	XMLSPath := flag.String("xmls", "", "Génére une signature en-tête pour xmls (donner le chemin absolu).")
	launchClient := flag.Bool("client", false, "Lance le socket côte client (-client=true/false).")
	buildRegistry := flag.Bool("buildR", false, "build Registry")
	msgName := flag.String("msgName", "", "msgName")
	msgId := flag.Uint("msgId", 0, "packet id")

	flag.Parse()

	parsers.DecodeAllD2p("C:\\Users\\Ivan\\Desktop\\Dofus\\content\\maps\\maps0.d2p")

	if *buildRegistry {
		generates.BuildRegistry()
	}

	if *rsa_ {
		err := generates.Signature()
		if err != nil {
			panic(err)
		}
		log.Println("RSA généré.")
	}

	if *hello {
		err := generates.HelloConnectPair()
		if err != nil {
			panic(err)
		}
		log.Println("RSA Hello généré.")
	}

	if *hosts != "" {
		generates.GenerateHostsSignature(*hosts)
		log.Println("HOSTS généré.")
	}

	if *XMLSPath != "" {
		err := generates.GenerateXMLSignature(*XMLSPath)
		if err != nil {
			return
		}
		log.Println("XMLS généré.")
	}

	if *launchClient {
		sockets.GoSocket()
	}

	if *msgName != "" && *msgId != 0 {
		generates.GenerateMessage(*msgName, uint32(*msgId))
	}
}
