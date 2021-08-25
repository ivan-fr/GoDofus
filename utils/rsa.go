package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"os"
)

func RsaPublicDecrypt(pubKey *rsa.PublicKey, data []byte) []byte {
	c := new(big.Int)
	m := new(big.Int)
	m.SetBytes(data)
	e := big.NewInt(int64(pubKey.E))
	c.Exp(m, e, pubKey.N)
	out := c.Bytes()
	skip := 0
	for i := 2; i < len(out); i++ {
		if i+1 >= len(out) {
			break
		}
		if out[i] == 0xff && out[i+1] == 0 {
			skip = i + 2
			break
		}
	}
	return out[skip:]
}

func ReadRSA(path string) []byte {
	publicVerifyPem, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return publicVerifyPem
}

func DecodePem(fileContent []byte) *pem.Block {
	var blockVerify, _ = pem.Decode(fileContent)
	if blockVerify == nil {
		panic("block empty")
	}
	return blockVerify
}

func PublicKeyOf(pem *pem.Block) *rsa.PublicKey {
	publicKeyVerify, err := x509.ParsePKIXPublicKey(pem.Bytes)
	if err != nil {
		panic(err)
	}
	p := publicKeyVerify.(*rsa.PublicKey)
	return p
}
