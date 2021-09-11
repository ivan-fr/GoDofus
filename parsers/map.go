package parsers

import (
	"GoDofus/utils"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

func getDecryptionKey() []byte {
	decodeString := []byte("649ae451ca33ec53bbcbcc33becf15f4")
	return decodeString
}

var decryptionKey = getDecryptionKey()

type cell struct {
	cellId        int16
	elementsCount int16
	myLayer       *layer
}

type cellData struct {
	cellId                  int16
	_floor                  byte
	_mov                    bool
	_nonWalkableDuringFight bool
	_nonWalkableDuringRP    bool
	_los                    bool
	_havenbagCell           bool
	_losmov                 byte
	speed                   byte
	mapChangeData           byte
	moveZone                byte
	_linkedZone             byte
	_farmCell               bool
	myMap                   *Map_
}

type layer struct {
	layerId    int32
	cellsCount uint16
	cells      []*cell
	myMap      *Map_
}

type Map_ struct {
	header            byte
	mapVersion        byte
	id                uint32
	encrypted         bool
	encryptionVersion byte
	dataLen           uint32
	encryptedData     []byte
	relativeId        uint32
	mapType           byte
	subareaId         int32
	topNeighbourId    int32
	bottomNeighbourId int32
	leftNeighbourId   int32
	rightNeighbourId  int32
	layersCount       byte
	layers            []*layer
	cells             []*cellData
	topArrowCell      []int16
	bottomArrowCell   []int16
	rightArrowCell    []int16
	leftArrowCell     []int16
}

func decodeMap(reader *bytes.Reader) *Map_ {
	var err error
	newMap := &Map_{}

	_ = binary.Read(reader, binary.BigEndian, &newMap.header)
	_ = binary.Read(reader, binary.BigEndian, &newMap.mapVersion)
	_ = binary.Read(reader, binary.BigEndian, &newMap.id)

	if newMap.mapVersion >= 7 {
		_ = binary.Read(reader, binary.BigEndian, &newMap.encrypted)
		_ = binary.Read(reader, binary.BigEndian, &newMap.encryptionVersion)
		_ = binary.Read(reader, binary.BigEndian, &newMap.dataLen)

		if newMap.encrypted {
			newMap.encryptedData = make([]byte, newMap.dataLen)
			_ = binary.Read(reader, binary.BigEndian, newMap.encryptedData)

			for i := uint32(0); i < newMap.dataLen; i++ {
				newMap.encryptedData[i] ^= decryptionKey[i%uint32(len(decryptionKey))]
			}
			reader = bytes.NewReader(newMap.encryptedData)
		}
	}

	_ = binary.Read(reader, binary.BigEndian, &newMap.relativeId)
	_ = binary.Read(reader, binary.BigEndian, &newMap.mapType)
	_ = binary.Read(reader, binary.BigEndian, &newMap.subareaId)

	_ = binary.Read(reader, binary.BigEndian, &newMap.topNeighbourId)
	_ = binary.Read(reader, binary.BigEndian, &newMap.bottomNeighbourId)
	_ = binary.Read(reader, binary.BigEndian, &newMap.leftNeighbourId)
	_ = binary.Read(reader, binary.BigEndian, &newMap.rightNeighbourId)

	_, err = reader.Seek(4, io.SeekCurrent)
	if err != nil {
		panic("wrong seek")
	}

	if newMap.mapVersion >= 9 {
		_, err = reader.Seek(4+4, io.SeekCurrent)
		if err != nil {
			panic("wrong seek")
		}
	} else if newMap.mapVersion >= 3 {
		_, err = reader.Seek(1+1+1, io.SeekCurrent)
		if err != nil {
			panic("wrong seek")
		}
	}

	if newMap.mapVersion >= 4 {
		_, err = reader.Seek(2+2+2, io.SeekCurrent)
		if err != nil {
			panic("wrong seek")
		}
	}

	if newMap.mapVersion > 10 {
		_, err = reader.Seek(4, io.SeekCurrent)
		if err != nil {
			panic("wrong seek")
		}
	}

	var backgroundsCount byte
	_ = binary.Read(reader, binary.BigEndian, &backgroundsCount)

	for i := 0; i < int(backgroundsCount); i++ {
		_, err = reader.Seek(4+2+2+2+2+2+1+1+1+1, io.SeekCurrent)
		if err != nil {
			panic("wrong seek")
		}
	}

	var foregroundsCount byte
	_ = binary.Read(reader, binary.BigEndian, &foregroundsCount)
	for i := 0; i < int(foregroundsCount); i++ {
		_, err = reader.Seek(4+2+2+2+2+2+1+1+1+1, io.SeekCurrent)
		if err != nil {
			panic("wrong seek")
		}
	}

	_, err = reader.Seek(4+4, io.SeekCurrent)
	if err != nil {
		panic("wrong seek")
	}

	var layersCount byte
	_ = binary.Read(reader, binary.BigEndian, &layersCount)
	newMap.layers = make([]*layer, layersCount)
	for i := 0; i < int(layersCount); i++ {
		newMap.layers[i] = new(layer)

		if newMap.mapVersion >= 9 {
			var smallLayerId byte
			_ = binary.Read(reader, binary.BigEndian, &smallLayerId)
			newMap.layers[i].layerId = int32(smallLayerId)
		} else {
			_ = binary.Read(reader, binary.BigEndian, &newMap.layers[i].layerId)
		}
		_ = binary.Read(reader, binary.BigEndian, &newMap.layers[i].cellsCount)

		if newMap.layers[i].cellsCount > 0 {
			newMap.layers[i].cells = make([]*cell, newMap.layers[i].cellsCount)
			for j := 0; j < int(newMap.layers[i].cellsCount); j++ {
				newMap.layers[i].cells[j] = new(cell)

				_ = binary.Read(reader, binary.BigEndian, &newMap.layers[i].cells[j].cellId)

				_ = binary.Read(reader, binary.BigEndian, &newMap.layers[i].cells[j].elementsCount)
				for z := 0; z < int(newMap.layers[i].cells[j].elementsCount); z++ {
					var elementType byte
					_ = binary.Read(reader, binary.BigEndian, &elementType)

					switch elementType {
					case 2:
						_, err = reader.Seek(4+1+1+1+1+1+1, io.SeekCurrent)
						if err != nil {
							panic("wrong seek")
						}

						if newMap.mapVersion <= 4 {
							_, err = reader.Seek(1+1, io.SeekCurrent)
							if err != nil {
								panic("wrong seek")
							}
						} else {
							_, err = reader.Seek(2+2, io.SeekCurrent)
							if err != nil {
								panic("wrong seek")
							}
						}
						_, err = reader.Seek(1+4, io.SeekCurrent)
						if err != nil {
							panic("wrong seek")
						}
					case 33:
						_, err = reader.Seek(4+2+4+4+2+2, io.SeekCurrent)
						if err != nil {
							panic("wrong seek")
						}
					}
				}

				newMap.layers[i].cells[j].myLayer = newMap.layers[i]
			}
		}
	}

	var mapCellsCount = 560
	newMap.cells = make([]*cellData, mapCellsCount)
	for i := 0; i < mapCellsCount; i++ {
		newMap.cells[i] = new(cellData)
		newMap.cells[i].myMap = newMap

		var topArrow bool
		var bottomArrow bool
		var rightArrow bool
		var leftArrow bool

		newMap.cells[i].cellId = int16(i)
		_ = binary.Read(reader, binary.BigEndian, &newMap.cells[i]._floor)

		if int16(newMap.cells[i]._floor)*10 == -1280 {
			continue
		}

		if newMap.mapVersion >= 9 {
			var tmpbytesv9 int16
			_ = binary.Read(reader, binary.BigEndian, &tmpbytesv9)
			newMap.cells[i]._mov = (tmpbytesv9 & 1) == 0
			newMap.cells[i]._nonWalkableDuringFight = (tmpbytesv9 & 2) != 0
			newMap.cells[i]._nonWalkableDuringRP = (tmpbytesv9 & 4) != 0
			newMap.cells[i]._los = (tmpbytesv9 & 8) == 0
			newMap.cells[i]._farmCell = (tmpbytesv9 & 128) != 0

			if newMap.mapVersion >= 10 {
				newMap.cells[i]._havenbagCell = (tmpbytesv9 & 256) != 0
				topArrow = (tmpbytesv9 & 512) != 0
				bottomArrow = (tmpbytesv9 & 1024) != 0
				rightArrow = (tmpbytesv9 & 2048) != 0
				leftArrow = (tmpbytesv9 & 4096) != 0
			} else {
				topArrow = (tmpbytesv9 & 256) != 0
				bottomArrow = (tmpbytesv9 & 512) != 0
				rightArrow = (tmpbytesv9 & 1024) != 0
				leftArrow = (tmpbytesv9 & 2048) != 0
			}

			if topArrow {
				newMap.topArrowCell = append(newMap.topArrowCell, int16(i))
			}
			if bottomArrow {
				newMap.bottomArrowCell = append(newMap.bottomArrowCell, int16(i))
			}
			if rightArrow {
				newMap.rightArrowCell = append(newMap.rightArrowCell, int16(i))
			}
			if leftArrow {
				newMap.leftArrowCell = append(newMap.leftArrowCell, int16(i))
			}
		} else {
			_ = binary.Read(reader, binary.BigEndian, &newMap.cells[i]._losmov)
			newMap.cells[i]._los = (newMap.cells[i]._losmov&2)>>1 == 1
			newMap.cells[i]._mov = (newMap.cells[i]._losmov & 1) == 1
			newMap.cells[i]._nonWalkableDuringRP = (newMap.cells[i]._losmov&128)>>7 == 1
			newMap.cells[i]._nonWalkableDuringFight = (newMap.cells[i]._losmov&4)>>2 == 1
			newMap.cells[i]._farmCell = (newMap.cells[i]._losmov&32)>>5 == 1
		}

		_ = binary.Read(reader, binary.BigEndian, &newMap.cells[i].speed)
		_ = binary.Read(reader, binary.BigEndian, &newMap.cells[i].mapChangeData)

		if newMap.mapVersion > 5 {
			_ = binary.Read(reader, binary.BigEndian, &newMap.cells[i].moveZone)
		}

		if newMap.mapVersion > 10 &&
			(newMap.cells[i]._mov && !newMap.cells[i]._farmCell ||
				newMap.cells[i]._mov && !newMap.cells[i]._nonWalkableDuringFight && !newMap.cells[i]._farmCell && !newMap.cells[i]._havenbagCell) {
			_ = binary.Read(reader, binary.BigEndian, &newMap.cells[i]._linkedZone)
		}

		if newMap.mapVersion > 7 && newMap.mapVersion < 9 {
			var tmpBits byte
			_ = binary.Read(reader, binary.BigEndian, &tmpBits)
			var arrow = 15 & tmpBits

			if (arrow & 1) != 0 {
				newMap.topArrowCell = append(newMap.topArrowCell, int16(i))
			}
			if (arrow & 2) != 0 {
				newMap.bottomArrowCell = append(newMap.bottomArrowCell, int16(i))
			}
			if (arrow & 4) != 0 {
				newMap.rightArrowCell = append(newMap.rightArrowCell, int16(i))
			}
			if (arrow & 8) != 0 {
				newMap.leftArrowCell = append(newMap.leftArrowCell, int16(i))
			}
		}
	}

	return newMap
}

var properties = make(map[string]string)
var indexes = make(map[string]map[string]interface{})

func DecodeAllD2p(d2pPath string) {
	raw, err := os.ReadFile(d2pPath)
	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(raw)
	var vMax byte
	var vMin byte
	_ = binary.Read(reader, binary.BigEndian, &vMax)
	_ = binary.Read(reader, binary.BigEndian, &vMin)

	if vMax != 2 || vMin != 1 {
		return
	}

	_, err = reader.Seek(-24, io.SeekEnd)
	if err != nil {
		return
	}

	var dataOffset uint32
	var dataCount uint32
	var indexOffset uint32
	var indexCount uint32
	var propertiesOffset int32
	var propertiesCount uint32

	_ = binary.Read(reader, binary.BigEndian, &dataOffset)
	_ = binary.Read(reader, binary.BigEndian, &dataCount)
	_ = binary.Read(reader, binary.BigEndian, &indexOffset)
	_ = binary.Read(reader, binary.BigEndian, &indexCount)
	_ = binary.Read(reader, binary.BigEndian, &propertiesOffset)
	_ = binary.Read(reader, binary.BigEndian, &propertiesCount)

	_, err = reader.Seek(int64(propertiesOffset), io.SeekStart)
	if err != nil {
		return
	}

	var newFile string

	for i := uint32(0); i < propertiesCount; i++ {
		propertyName := string(utils.ReadUTF(reader))
		propertyValue := string(utils.ReadUTF(reader))
		properties[propertyName] = propertyValue

		if propertyName == "link" {
			index := strings.LastIndex(d2pPath, "\\")
			if index != -1 {
				newFile = d2pPath[:index+1] + propertyValue
			} else {
				panic("wrong path")
			}
		}
	}

	_, err = reader.Seek(int64(indexOffset), io.SeekStart)
	if err != nil {
		return
	}

	for i := uint32(0); i < indexCount; i++ {
		filePath := string(utils.ReadUTF(reader))
		var fileOffset int32
		var fileLength int32
		_ = binary.Read(reader, binary.BigEndian, &fileOffset)
		_ = binary.Read(reader, binary.BigEndian, &fileLength)
		indexes[filePath] = map[string]interface{}{"o": int64(fileOffset) + int64(dataOffset), "l": fileLength, "stream": reader}
	}

	if newFile != "" {
		DecodeAllD2p(newFile)
	}
}

func LoadMap(mapId uint64) *Map_ {
	index := indexes[fmt.Sprintf("%d/%d.dlm", mapId%10, mapId)]
	reader := index["stream"].(*bytes.Reader)
	_, err := reader.Seek(index["o"].(int64), io.SeekStart)
	if err != nil {
		panic("wrong seek")
	}

	data := make([]byte, index["l"].(int32))
	_, _ = io.ReadFull(reader, data)
	readerData := bytes.NewReader(data)

	open, err := zlib.NewReader(readerData)
	if err != nil {
		panic(err)
	}
	all, err := io.ReadAll(open)
	if err != nil {
		panic(err)
	}

	return decodeMap(bytes.NewReader(all))
}

func (m *Map_) printToConsole() {

}
