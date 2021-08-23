package main

import (
	"GoDofus/pack"
	"GoDofus/signer"
	"GoDofus/socket"
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

	pack.ToggleClient()

	hex1 := "8c97ea1eef6c0433c29c96fe08004500014460f8400040060155c0a8014e3f22d64ece2a15b3081e1f5d84017b04501801f5d89d00002b3e0000003b011401023c030000000c00000266728002665621969ad6a7c4ff3ddc26bbb1a7913bd54896bf8160da73d4005d37cf438c5d607c673dd3b2f77c7cba53964decad6d193045326ccbf9d970a2a4358ad36b3a850a20d5ef883b43144d754dc4faabc5b7ca723d63be0583d39b0cbfe2e848e433b06f1fd77ae70107fa7a720aacb353f151236d3070edb282bf461ba7ca748bb6d7756425e1d5351998ae1ad5c6ad22b77da87950bdca02e86e4c440c546c17631aa2fb1025461bf95051eb3480ec08a980ba0c6cbbaf352c76e1177f4e0d99904c334bef5827c234149f869664197637665c89a0eac6b07bac54decd1c389d1a41feb3d8aceadb3f58491c4a32c60c2ce3f64744129a8a4de0531a0ab7130000000000"
	b1, _ := hex.DecodeString(hex1)
	ok := pack.Read(b1)

	fmt.Println(pack.GetPipeline().Get(), ok)
}
