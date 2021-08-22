package main

import (
	"GoDofus/messages"
	"GoDofus/signer"
	"flag"
	"log"
)

func main() {
	rsa_ := flag.Bool("is_rsa", false, "Génére le private/public key et la signature DofusPublicKey.")
	hosts := flag.String("hosts", "", "Génére la signature des hosts pour config.xml.")
	XMLSPath := flag.String("xmls", "", "Génére une signature en-tête pour xmls.")

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

	log.Printf("%s", messages.GetSalt())
}
