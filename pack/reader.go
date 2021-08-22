package pack

import (
	"bytes"
	"encoding/binary"
	"io"
)

type weft struct {
	packId     uint16
	lengthType uint16
	length     uint32
	message    []byte
}

type lastSignal struct {
	request        int
	typeRequest    int
	containForType []byte
	containNoType  []byte
}

type pipe struct {
	wefts []*weft
}

var lSignal = &lastSignal{typeRequest: noType}
var pipeline = new(pipe)
var lastWeft *weft = nil

const (
	headerTwoFirstBytes = iota
	headerLength
	messageLength
	noType
)

func GetPipeline() *pipe {
	return pipeline
}

func (p *pipe) append(w *weft) {
	p.wefts = append(p.wefts, w)
	/*	if len(p.wefts) > 3 {
			p.wefts = p.wefts[2:]
		}
	*/
}

func (lSignal *lastSignal) update(request int, typeRest int, containForType []byte, containNoType []byte) {
	if typeRest == lSignal.typeRequest && lSignal.request < 0 && request <= lSignal.request {
		return
	}

	if typeRest == lSignal.typeRequest && lSignal.request < 0 && request > lSignal.request {
		lSignal.request = request
		lSignal.containForType = append(lSignal.containForType, containForType...)
	} else if lSignal.request == 0 || lSignal.typeRequest == noType {
		lSignal.request = request
		lSignal.typeRequest = typeRest
		lSignal.containForType = containForType
		lSignal.containForType = nil
	} else {
		panic("incoherence from lastSignal")
	}

	if lSignal.request > 0 {
		if containNoType == nil {
			panic("incoherence contain no type can't be nil")
		}

		lSignal.containNoType = containForType
	}
}

func commit() {
	if lSignal.request >= 0 && lSignal.typeRequest == messageLength {
		if lastWeft == nil {
			panic("lastWeft must not be nil.")
		}

		if len(lastWeft.message) > 0 {
			panic("lastWeft wasn't purge.")
		}

		lastWeft.message = lSignal.containForType
		pipeline.append(lastWeft)
		lastWeft = nil
	}
}

func min(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}

func tryRead(reader *bytes.Reader, step int, bytesWanted uint) bool {
	readerLen := uint(reader.Len())
	nbBytesToRead := min(uint(reader.Len()), bytesWanted) // <= bytesWanted && <= reader.Len
	var containForType = make([]byte, nbBytesToRead)
	_, _ = io.ReadFull(reader, containForType)

	request := int(readerLen - bytesWanted)
	ok := false
	var containNoType []byte

	if request >= 0 {
		ok = true
	}

	if request > 0 {
		containNoType, _ = io.ReadAll(reader)
	}

	lSignal.update(request, step, containForType, containNoType)
	return ok
}

func readHeaderTwoFirstBytes(reader *bytes.Reader) bool {
	ok := tryRead(reader, headerTwoFirstBytes, 2)

	if !ok {
		return false
	}

	twoBytes := binary.BigEndian.Uint16(lSignal.containForType)
	packetId := twoBytes >> 2
	lengthType := twoBytes & 0b11

	*lastWeft = weft{packId: packetId, lengthType: lengthType}
	return true
}

func readHeaderLength(reader *bytes.Reader) bool {
	if lastWeft == nil {
		panic("incoherence last weft can't be nil")
	}

	if lastWeft.lengthType == 0 {
		panic("incoherence last weft can't have length type equals zero")
	}

	ok := tryRead(reader, headerLength, uint(lastWeft.lengthType))

	if !ok {
		return false
	}

	switch lastWeft.lengthType {
	case 3:
		var specialCaseReader = bytes.NewReader(lSignal.containForType)
		var firstByte uint8
		_ = binary.Read(specialCaseReader, binary.BigEndian, firstByte)
		var twoBytes uint16
		_ = binary.Read(specialCaseReader, binary.BigEndian, twoBytes)
		lastWeft.length = uint32(firstByte)<<16 | uint32(twoBytes)
	case 2:
		lastWeft.length = uint32(binary.BigEndian.Uint16(lSignal.containForType))
	case 1:
		lastWeft.length = uint32(lSignal.containForType[0])
	}

	return true
}

func Read(bytesPack []byte) bool {
	switch lSignal.typeRequest {
	case messageLength:
		switch {
		case lSignal.request == 0:
			commit()
			lSignal.update(0, noType, nil, nil)
			return Read(bytesPack)
		case lSignal.request > 0:
			commit()
			newBytesPack := append(lSignal.containNoType, bytesPack...)
			lSignal.update(0, noType, nil, nil)
			return Read(newBytesPack)
		case lSignal.request < 0:
			reader := bytes.NewReader(bytesPack)
			ok := tryRead(reader, messageLength, uint(-lSignal.request))

			if !ok {
				return false
			}

			commit()
			newBytePack := lSignal.containNoType
			lSignal.update(0, noType, nil, nil)
			return Read(newBytePack)
		}
	case headerTwoFirstBytes:
		switch {
		case lSignal.request < 0:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderTwoFirstBytes(reader)
			if !ok {
				return false
			}

			newBytePack := lSignal.containNoType
			lSignal.update(0, noType, nil, nil)
			return Read(newBytePack)
		default:
			panic("incoherence, last signal can't be positive")
		}
	case headerLength:
		switch {
		case lSignal.request < 0:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderLength(reader)
			if !ok {
				return false
			}

			newBytePack := lSignal.containNoType
			lSignal.update(0, noType, nil, nil)
			return Read(newBytePack)
		default:
			panic("incoherence, last signal can't be positive")
		}
	case noType:
		switch {
		case bytesPack == nil:
			return true
		case lastWeft == nil:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderTwoFirstBytes(reader)
			if !ok {
				return false
			}

			newBytePack := lSignal.containNoType
			lSignal.update(0, noType, nil, nil)
			return Read(newBytePack)
		case lastWeft.length == 0:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderLength(reader)
			if !ok {
				return false
			}

			newBytePack := lSignal.containNoType
			lSignal.update(0, noType, nil, nil)
			return Read(newBytePack)
		default:
			reader := bytes.NewReader(bytesPack)
			_ = tryRead(reader, messageLength, uint(lastWeft.length))
			return Read(nil)
		}
	default:
		panic("program don't know the step.")
	}

	return true
}
