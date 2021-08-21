package signer

import (
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

func Signature() error {
	privateKey, err := rsa.GenerateKey(cryptoRand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.PublicKey
	var nStr = fmt.Sprintf("%x", publicKey.N)
	var eStr = fmt.Sprintf("%x", publicKey.E)

	buf := new(bytes.Buffer)

	bytesArray := [3][]byte{
		[]byte("DofusPublicKey"),
		[]byte(nStr),
		[]byte(eStr),
	}

	for _, bytesValues := range bytesArray {
		err = binary.Write(buf, binary.BigEndian, uint16(len(bytesValues)))
		err = binary.Write(buf, binary.LittleEndian, bytesValues)
		if err != nil {
			panic(err)
		}
	}

	pki, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}

	fo, err := os.Create("./sign/public_key.pem")
	if err != nil {
		panic(err)
	}

	_, err = fo.WriteString(
		fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----",
			base64.StdEncoding.EncodeToString(pki),
		),
	)
	if err != nil {
		return err
	}

	if err != nil {
		panic(err)
	}

	err = fo.Close()
	if err != nil {
		panic(err)
	}

	privatePemData := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	_ = os.WriteFile("./sign/private_key.pem", privatePemData, 0644)
	_ = os.WriteFile("./sign/signature.bin", buf.Bytes(), 0644)

	return err
}

func getSignatureBuffer(byteToEncode []byte) *bytes.Buffer {
	privatePam, _ := os.ReadFile("./sign/private_key.pem")
	block, _ := pem.Decode(privatePam)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		log.Fatal("failed to decode PEM block containing public key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomInt := byte(random.Intn(254) + 1)

	hashBuffer := new(bytes.Buffer)

	md5_ := md5.Sum(byteToEncode)
	hexBytes := []byte(hex.EncodeToString(md5_[:]))

	_ = binary.Write(hashBuffer, binary.LittleEndian, randomInt)
	_ = binary.Write(hashBuffer, binary.BigEndian, uint32(len(byteToEncode)))
	_ = binary.Write(hashBuffer, binary.LittleEndian, hexBytes)

	hashBytes := hashBuffer.Bytes()

	for i := 2; i < len(hashBytes); i++ {
		hashBytes[i] ^= randomInt
	}

	signedData, _ := rsa.SignPKCS1v15(nil, privateKey, crypto.Hash(0), hashBytes)
	startString := []byte("AKSF")

	finalBuff := new(bytes.Buffer)
	_ = binary.Write(finalBuff, binary.BigEndian, uint16(len(startString)))
	_ = binary.Write(finalBuff, binary.LittleEndian, startString)
	_ = binary.Write(finalBuff, binary.BigEndian, int16(1))
	_ = binary.Write(finalBuff, binary.BigEndian, int32(len(signedData)))
	_ = binary.Write(finalBuff, binary.LittleEndian, signedData)

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
