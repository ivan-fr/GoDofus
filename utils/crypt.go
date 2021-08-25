package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"io"
	"math/big"
	"strconv"
)

var AnkamaSignedFileHeader = "AKSF"
var SignatureHeader = "AKSD"
var PublicKeyHeader string = "DofusPublicKey"

var PrivateKeyHeader string = "DofusPrivateKey"

var publicV2Pem = ReadRSA("./binaryData/verify_key.bin")
var blockV2 = DecodePem(publicV2Pem)
var publicKeyV2 = PublicKeyOf(blockV2)
var publicKeyV1 = getPublicKeyFromV1()

func v2(reader *bytes.Reader, headerPosition int64) []byte {
	_, _ = reader.Seek(headerPosition-4, io.SeekStart)

	var signedDataLength int16
	_ = binary.Read(reader, binary.BigEndian, &signedDataLength)
	_, _ = reader.Seek(headerPosition-4-int64(signedDataLength), io.SeekStart)

	var cryptData = make([]byte, signedDataLength)
	_ = binary.Read(reader, binary.BigEndian, cryptData)

	sigData := RsaPublicDecrypt(publicKeyV2, cryptData)
	readerSig := bytes.NewReader(sigData)

	sigHeader := ReadUTF(readerSig)

	if string(sigHeader) != SignatureHeader {
		panic("wrong header from sign")
	}

	var sigVersion byte
	_ = binary.Read(readerSig, binary.BigEndian, &sigVersion)

	var sigFileLength int32
	_ = binary.Read(readerSig, binary.BigEndian, &sigFileLength)
	_ = binary.Read(readerSig, binary.BigEndian, &sigFileLength)
	_ = binary.Read(readerSig, binary.BigEndian, &sigFileLength)

	if int64(sigFileLength) != headerPosition-4-int64(signedDataLength) {
		panic("wrong structure")
	}

	var hashType byte
	_ = binary.Read(readerSig, binary.BigEndian, &hashType)
	signHash := ReadUTF(readerSig)

	_, _ = reader.Seek(0, io.SeekStart)
	var output = make([]byte, headerPosition-4-int64(signedDataLength))
	_ = binary.Read(reader, binary.BigEndian, output)

	var contentHash []byte
	switch hashType {
	case 0:
		sum := md5.Sum(output)
		contentHash = sum[:]
	case 1:
		sum := sha256.Sum256(output)
		contentHash = sum[:]
	default:
		panic("wrong hashType")
	}

	if bytes.Compare(signHash, contentHash) != 0 {
		panic("wrong hash")
	}

	return output
}

func getPublicKeyFromV1() *rsa.PublicKey {
	v1_ := ReadRSA("./binaryData/V1.bin")
	readerV1 := bytes.NewReader(v1_)

	var header = string(ReadUTF(readerV1))
	if header != PublicKeyHeader && header != PrivateKeyHeader {
		panic("wrong hedear v1 decode")
	}

	if header == PublicKeyHeader {
		pub := new(rsa.PublicKey)
		N := ReadUTF(readerV1)
		E := ReadUTF(readerV1)

		n := new(big.Int)
		n, ok := n.SetString(string(N), 16)

		if !ok {
			panic("can't get big int")
		}

		e, err := strconv.ParseInt(string(E), 16, 32)
		if err != nil {
			panic("can't get E")
		}

		pub.N = n
		pub.E = int(e)

		return pub
	}

	panic("the program don't have private key")
}

func v1(reader *bytes.Reader) []byte {
	var formatVersion int16
	_ = binary.Read(reader, binary.BigEndian, &formatVersion)

	var len_ int32
	_ = binary.Read(reader, binary.BigEndian, &len_)

	var sigData = make([]byte, len_)
	_ = binary.Read(reader, binary.BigEndian, sigData)

	decryptedHash := RsaPublicDecrypt(publicKeyV1, sigData)

	readerHash := bytes.NewReader(decryptedHash)
	var randomPart byte
	_ = binary.Read(readerHash, binary.BigEndian, &randomPart)

	for i := 2; i < len(decryptedHash); i++ {
		decryptedHash[i] ^= randomPart
	}

	var contentLen int32
	_ = binary.Read(readerHash, binary.BigEndian, &contentLen)
	testedContentLen := reader.Len()

	signHash := make([]byte, readerHash.Len())
	signHash = signHash[1:]

	var output = make([]byte, reader.Len())
	_ = binary.Read(reader, binary.BigEndian, output)

	contentHash := md5.Sum(output)
	contentHash_ := contentHash[1:]

	if bytes.Compare(signHash, contentHash_) == 0 && contentLen == int32(testedContentLen) {
		return output
	}

	panic("wrong result")
}

func DecryptV(reader *bytes.Reader) []byte {
	var headerSize uint16
	_ = binary.Read(reader, binary.BigEndian, &headerSize)

	if headerSize != uint16(len(AnkamaSignedFileHeader)) {
		var headerPosition = reader.Size() - int64(len(AnkamaSignedFileHeader))
		_, _ = reader.Seek(headerPosition, io.SeekStart)
		var header = make([]byte, len(AnkamaSignedFileHeader))
		_ = binary.Read(reader, binary.BigEndian, &header)

		if string(header) == AnkamaSignedFileHeader {
			return v2(reader, headerPosition)
		}
	} else {
		var header = make([]byte, len(AnkamaSignedFileHeader))
		_ = binary.Read(reader, binary.BigEndian, &header)

		if string(header) == AnkamaSignedFileHeader {
			return v1(reader)
		}
	}

	panic("wrong header")
}
