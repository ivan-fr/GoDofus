package main

import (
	"GoDofus/messages"
	"GoDofus/signer"
	"GoDofus/socket"
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
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

	b2, _ := hex.DecodeString("01023c040000000d000002667280020c6f52e597455920b273738ddfd9a5aa63adf9eb0febcf4e274133a862d4a2af929f71e1992e58d25d446a324f379d3931d91b638ddeca1017bc385cd265f8a8434445a51530d8c4c06d79e4d43c91b955a2b16ac7f186a33ee941e3915c76f311bf3aac6255d23aecdfc45a29e70ddcaa8084d9fe0dcdb4dcf7e2fd69912c9de36eec434ccc2f1e2773da57b3d0c6f40b2d5599a0d94ef184b08bb4d59ddceaf6b760f02ced77e7a6cf9772ea68cf292b4487f9de13be51194530a81238a43726f0a2afba6a24efd887840dc0475f27619615e8dc8e6de4d45b446b912770432116ded4e699eefd1e5b35509128ff22205e8048b827b8d44d74e24d08fbce900000000000")
	r2 := bytes.NewReader(b2)
	messages.GetIdentificationNOA().Deserialize(r2)
	fmt.Println(messages.GetIdentificationNOA(), len(b2))

	b2, _ = hex.DecodeString("02023c040000000d0000026672800240bf48254b73228bf0fc5b826f1d5488c75b96470e2b49a05df18a8a0282542f1bb85f479942ac346a0302a8ebb01e29df8857b764b931aefc6445790258525c0304d984d76208e2ff19a3942dd887047b11c9c1de3dbb9016c4fc82221826a66a15fcaf01701c07a7cd4032448c2602b4d2897d2529fa9510fa7d0437ae43438bb813ca062d5b3d411fe8ce29b329c248a13640970cf2e529bd79eca76c2847f309932b402cba3edca42d735dd845f24ce728065ad31e486652229310e1afccced1a2da9b1e9f53566b29735e22e7fa906d70160cfca06a612710be9c452bc18d29b12ff5b1695c8887ba1738682e7f1d89e184266159b7d955e6cb5773a60c0000000000")
	r2 = bytes.NewReader(b2)
	messages.GetIdentificationNOA().Deserialize(r2)
	fmt.Println(messages.GetIdentificationNOA(), len(b2))
}
