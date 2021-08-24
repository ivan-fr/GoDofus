package managers

import (
	"GoDofus/messages"
	"GoDofus/structs"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

type loginAction struct {
	username         string
	password         string
	autoSelectServer bool
	serverId         uint16
	host             string
}

type authentification struct {
	AESKey    []byte
	lA        *loginAction
	lang      string
	publicKey []byte
	salt      string
}

var authenticate_ = &authentification{AESKey: generateAESKey(), lang: "fr"}
var AESLength = uint(32)

func GetAuthentification() *authentification {
	return authenticate_
}

func (a *authentification) initLoginAction() {
	la := &loginAction{autoSelectServer: true}
	fmt.Println("Entre nom de compte :")
	_, _ = fmt.Scanln(&la.username)
	fmt.Println("Entre le mot de passe :")
	_, _ = fmt.Scanln(&la.password)
	a.lA = la
}

func (a *authentification) getCipher() []byte {
	buff := new(bytes.Buffer)
	mySalt := a.getSalt()
	_ = binary.Write(buff, binary.BigEndian, []byte(mySalt))
	_ = binary.Write(buff, binary.BigEndian, a.AESKey)
	_ = binary.Write(buff, binary.BigEndian, byte(len(a.lA.username)))
	_ = binary.Write(buff, binary.BigEndian, []byte(a.lA.username))
	_ = binary.Write(buff, binary.BigEndian, []byte(a.lA.password))

	_ = os.WriteFile("./sign/cipher.bin", buff.Bytes(), 0644)

	a.getPublicKey()

	args := []string{"rsautl", "-encrypt", "-inkey", "/home/ivan/GolandProjects/GoDofus/sign/publicKeyFromHello.pem",
		"-pubin", "-in", "/home/ivan/GolandProjects/GoDofus/sign/cipher.bin"}
	credentials, err := exec.Command("openssl", args...).Output()
	if err != nil {
		panic(err)
	}

	return credentials
}

func (a *authentification) InitIdentificationMessage() {
	a.initLoginAction()
	identification := messages.GetIdentificationNOA()
	identification.Lang = a.lang
	identification.AutoSelectServer = a.lA.autoSelectServer

	currentVersion := structs.GetVersionNOA()
	identification.Version.Major = currentVersion.Major
	identification.Version.Minor = currentVersion.Minor
	identification.Version.Code = currentVersion.Code
	identification.Version.Build = currentVersion.Build
	identification.Version.BuildType = currentVersion.BuildType

	identification.Credentials = a.getCipher()
}

func (a *authentification) getPublicKey() {
	hc := messages.GetHelloConnectNOA()

	if a.publicKey != nil && hc.Salt == a.salt {
		return
	}

	if hc.Key == nil {
		panic("helloMessage wasn't call")
	}

	_ = os.WriteFile("./sign/keyFromHello.pem", hc.Key, 0644)
	args := []string{"rsautl", "-inkey", "/home/ivan/GolandProjects/GoDofus/binaryData/verify_key.bin",
		"-pubin", "-in", "/home/ivan/GolandProjects/GoDofus/sign/keyFromHello.pem"}
	out, err := exec.Command("openssl", args...).Output()
	if err != nil {
		log.Fatal(err)
	}

	a.publicKey = []byte(fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----",
		base64.StdEncoding.EncodeToString(out)))

	_ = os.WriteFile("./sign/publicKeyFromHello.pem", a.publicKey, 0644)
}

func (a *authentification) getSalt() string {
	hc := messages.GetHelloConnectNOA()
	if hc.Salt == "" {
		panic("helloMessage wasn't call")
	}

	mySalt := hc.Salt
	a.salt = hc.Salt

	if len(mySalt) < 32 {
		for len(mySalt) < 32 {
			mySalt += " "
		}
	}

	return mySalt
}

func generateAESKey() []byte {
	aes := make([]byte, AESLength)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := uint(0); i < AESLength; i++ {
		aes[i] = byte(random.Intn(255))
	}

	return aes
}
