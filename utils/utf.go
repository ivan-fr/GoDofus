package utils

import (
	"bytes"
	"encoding/binary"
	"io"
)

func WriteUTF(buff *bytes.Buffer, value []byte) {
	_ = binary.Write(buff, binary.BigEndian, uint16(len(value)))
	_ = binary.Write(buff, binary.BigEndian, value)
}

func ReadUTF(reader *bytes.Reader) []byte {
	var length uint16
	err := binary.Read(reader, binary.BigEndian, &length)
	if err != nil {
		panic(err)
	}

	bytesString := make([]byte, length)
	_, err = io.ReadFull(reader, bytesString)
	if err != nil {
		panic(err)
	}

	return bytesString
}
