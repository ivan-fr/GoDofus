package main

import (
	"GoDofus/signer"
	"GoDofus/socket"
	"flag"
	"log"
)

func main() {
	rsa_ := flag.Bool("get_rsa", false, "Génére le private/public key et la signature DofusPublicKey. (-get_rsa=true/false)")
	hosts := flag.String("hosts", "", "Génére la signature des hosts pour config.xml. Entrez localhost,127.0.0.1 par exemple.")
	XMLSPath := flag.String("xmls", "", "Génére une signature en-tête pour xmls (donner le chemin absolu).")
	launchClient := flag.Bool("client", false, "Lance le socket côte client (-client=true/false).")

	flag.Parse()

	if *rsa_ {
		err := signer.Signature()
		if err != nil {
			panic(err)
		}
		log.Println("RSA généré.")
	}

	if *hosts != "" {
		signer.GenerateHostsSignature(*hosts)
		log.Println("HOSTS généré.")
	}

	if *XMLSPath != "" {
		err := signer.GenerateXMLSignature(*XMLSPath)
		if err != nil {
			return
		}
		log.Println("XMLS généré.")
	}

	if *launchClient {
		socket.LaunchClientSocket()
	}

}
