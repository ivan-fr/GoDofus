package generates

import (
	"GoDofus/utils"
	"bytes"
	"crypto"
	"crypto/md5"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	_ "io"
	"log"
	"math/rand"
	"os"
	"time"
)

func HelloConnectPair() error {
	privateKey, err := rsa.GenerateKey(cryptoRand.Reader, 1024)
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.PublicKey

	pki, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}

	publicPemData := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pki,
	})

	privatePemData := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	_ = os.WriteFile("./sign/hello_private_key.pem", privatePemData, 0644)
	_ = os.WriteFile("./sign/hello_public_key.pem", publicPemData, 0644)
	return err
}

func Signature() error {
	privateKey, err := rsa.GenerateKey(cryptoRand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.PublicKey
	var nStr = fmt.Sprintf("%x", publicKey.N)
	var eStr = fmt.Sprintf("%x", publicKey.E)

	buff := new(bytes.Buffer)

	stringArray := [3]string{
		"DofusPublicKey",
		nStr,
		eStr,
	}

	for _, value := range stringArray {
		utils.WriteUTF(buff, []byte(value))
		if err != nil {
			panic(err)
		}
	}

	pki, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}

	privatePemData := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	_ = os.WriteFile("./sign/private_key.pem", privatePemData, 0644)
	_ = os.WriteFile("./sign/public_key.pem",
		[]byte(
			fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----",
				base64.StdEncoding.EncodeToString(pki))),
		0644)
	_ = os.WriteFile("./sign/signature.bin", buff.Bytes(), 0644)

	return err
}

func getSignatureBuffer(byteToEncode []byte) *bytes.Buffer {
	privatePam, _ := os.ReadFile("./sign/private_key.pem")
	block, _ := pem.Decode(privatePam)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		log.Fatal("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomByte := byte(random.Intn(254) + 1)

	hashBuffer := new(bytes.Buffer)

	md5_ := md5.Sum(byteToEncode)
	hexBytes := []byte(hex.EncodeToString(md5_[:]))

	_ = binary.Write(hashBuffer, binary.BigEndian, randomByte)
	_ = binary.Write(hashBuffer, binary.BigEndian, uint32(len(byteToEncode)))
	_ = binary.Write(hashBuffer, binary.BigEndian, hexBytes)

	hashBytes := hashBuffer.Bytes()

	for i := 2; i < len(hashBytes); i++ {
		hashBytes[i] ^= randomByte
	}

	signedData, _ := rsa.SignPKCS1v15(nil, privateKey, crypto.Hash(0), hashBytes)
	startString := "AKSF"

	finalBuff := new(bytes.Buffer)
	utils.WriteUTF(finalBuff, []byte(startString))
	_ = binary.Write(finalBuff, binary.BigEndian, int16(1))
	_ = binary.Write(finalBuff, binary.BigEndian, int32(len(signedData)))
	_ = binary.Write(finalBuff, binary.BigEndian, signedData)

	return finalBuff
}

func GenerateXMLSignature(XMLPath string) error {
	configXML, err := os.ReadFile(XMLPath)
	if err != nil {
		panic(err)
	}

	finalBuff := getSignatureBuffer(configXML)
	_ = binary.Write(finalBuff, binary.BigEndian, configXML)
	_ = os.WriteFile("./sign/signature.xmls", finalBuff.Bytes(), 0644)

	return err
}

func GenerateHostsSignature(hosts string) {
	bytesHosts := []byte(hosts)
	finalBuff := getSignatureBuffer(bytesHosts)
	encoded := base64.StdEncoding.EncodeToString(finalBuff.Bytes())

	log.Println(encoded)
}
