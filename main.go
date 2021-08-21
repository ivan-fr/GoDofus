package main

import (
	"GoDofus/signer"
	"bytes"
	"encoding/binary"
	"flag"
	"io"
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

	var lol = []byte{98}
	p := bytes.NewReader(lol)
	var bb [4]byte

	t := p.Len()
	err := binary.Read(p, binary.LittleEndian, &bb)
	log.Println(bb)

	if err == io.ErrUnexpectedEOF {
		log.Println(len(bb)-t, bb)
	}
}
