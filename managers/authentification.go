package managers

import (
	"GoDofus/settings"
	"bytes"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"encoding/binary"
	"math/rand"
	"time"
)

type loginAction struct {
	username         string
	password         string
	AutoSelectServer bool
	serverId         uint16
	host             string
}

type Authentification struct {
	AESKey    []byte
	LA        *loginAction
	Lang      []byte
	PublicKey *rsa.PublicKey
	Salt      []byte
	instance  uint
}

var AESLength = uint(32)

var authentificationMap = make(map[uint]*Authentification)

func GetAuthentificationManager(instance uint) *Authentification {
	authentificationMap_, ok := authentificationMap[instance]

	if ok {
		return authentificationMap_
	}

	authentificationMap[instance] = &Authentification{AESKey: generateAESKey(), Lang: []byte("fr"), instance: instance}
	return authentificationMap[instance]
}

func (a *Authentification) InitLoginAction() {
	la := &loginAction{AutoSelectServer: true}
	la.username = settings.Settings.Ndc
	la.password = settings.Settings.Pass
	a.LA = la
}

func (a *Authentification) GetCipher(pubKey *rsa.PublicKey, salt []byte) []byte {
	buff := new(bytes.Buffer)
	_ = binary.Write(buff, binary.BigEndian, salt)
	_ = binary.Write(buff, binary.BigEndian, a.AESKey)
	_ = binary.Write(buff, binary.BigEndian, byte(len(a.LA.username)))
	_ = binary.Write(buff, binary.BigEndian, []byte(a.LA.username))
	_ = binary.Write(buff, binary.BigEndian, []byte(a.LA.password))

	credentials, err := rsa.EncryptPKCS1v15(cryptoRand.Reader, pubKey, buff.Bytes())
	if err != nil {
		panic(err)
	}
	return credentials
}

func generateAESKey() []byte {
	aes := make([]byte, AESLength)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := uint(0); i < AESLength; i++ {
		aes[i] = byte(random.Intn(255))
	}

	return aes
}
