// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-25 10:06:45.9347347 +0200 CEST m=+0.002585101

package structs

import (
	"GoDofus/utils"
	"bytes"
)

type AccountTagInformation struct {
	nickname  []byte
	tagNumber []byte
}

func (A *AccountTagInformation) Serialize(buff *bytes.Buffer) {

}

func (A *AccountTagInformation) Deserialize(reader *bytes.Reader) {
	A.nickname = utils.ReadUTF(reader)
	A.tagNumber = utils.ReadUTF(reader)
}
