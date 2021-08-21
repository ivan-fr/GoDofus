package paquets

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
	typeRest       int
	containForType []byte
	containNoType  []byte
}

type pipe struct {
	wefts []*weft
}

var lSignal = &lastSignal{typeRest: noType}
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
	if typeRest == lSignal.typeRest && lSignal.request < 0 && request > lSignal.request {
		lSignal.request = request
		lSignal.containForType = append(lSignal.containForType, containForType...)
	} else if lSignal.request == 0 || lSignal.typeRest == noType {
		lSignal.request = request
		lSignal.typeRest = typeRest
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
	if lSignal.request >= 0 && lSignal.typeRest == messageLength {
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

func tryRead(reader *bytes.Reader, step int, bytesWanted uint) ([]byte, bool) {
	nbBytesToRead := min(uint(reader.Len()), bytesWanted)
	var lectureBytes = make([]byte, nbBytesToRead)
	_, _ = io.ReadFull(reader, lectureBytes)

	bytesRest := int(nbBytesToRead - bytesWanted) // <= 0
	ok := false
	var containNoType []byte

	if bytesRest >= 0 {
		ok = true
	}

	if bytesRest > 0 {
		containNoType, _ = io.ReadAll(reader)
	}

	lSignal.update(bytesRest, step, lectureBytes, containNoType)
	return lectureBytes, ok
}

func readHeaderTwoFirstBytes(reader *bytes.Reader) bool {
	result, ok := tryRead(reader, headerTwoFirstBytes, 2)

	if !ok {
		return false
	}

	twoBytes := binary.BigEndian.Uint16(result)
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

	result, ok := tryRead(reader, headerLength, uint(lastWeft.lengthType))

	if !ok {
		return false
	}

	lastWeft.length = binary.BigEndian.Uint32(result)
	return true
}

func Read(bytesPack []byte) {
	switch lSignal.typeRest {
	case messageLength:
		switch {
		case lSignal.request == 0:
			commit()
			lSignal.update(0, noType, nil, nil)
			Read(bytesPack)
		case lSignal.request > 0:
			commit()
			newBytesPack := append(lSignal.containNoType, bytesPack...)
			lSignal.update(0, noType, nil, nil)
			Read(newBytesPack)
		case lSignal.request < 0:
			reader := bytes.NewReader(bytesPack)
			_, ok := tryRead(reader, messageLength, uint(-lSignal.request))

			if !ok {
				break
			}

			commit()
			newBytePack := lSignal.containNoType
			lSignal.update(0, noType, nil, nil)
			Read(newBytePack)
		}
	case headerTwoFirstBytes:
		switch {
		case lSignal.request < 0:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderTwoFirstBytes(reader)
			if !ok {
				break
			}

			newBytePack := lSignal.containNoType
			lSignal.update(0, noType, nil, nil)
			Read(newBytePack)
		default:
			panic("incoherence, last signal can't be positive")
		}
	case headerLength:
		switch {
		case lSignal.request < 0:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderLength(reader)
			if !ok {
				break
			}

			newBytePack := lSignal.containNoType
			lSignal.update(0, noType, nil, nil)
			Read(newBytePack)
		default:
			panic("incoherence, last signal can't be positive")
		}
	case noType:
		switch {
		case bytesPack == nil:
			return
		case lastWeft == nil:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderTwoFirstBytes(reader)
			if !ok {
				break
			}

			newBytePack := lSignal.containNoType
			lSignal.update(0, noType, nil, nil)
			Read(newBytePack)
		case lastWeft.length == 0:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderLength(reader)
			if !ok {
				break
			}

			newBytePack := lSignal.containNoType
			lSignal.update(0, noType, nil, nil)
			Read(newBytePack)
		default:
			reader := bytes.NewReader(bytesPack)
			_, ok := tryRead(reader, messageLength, uint(lastWeft.length))

			if ok && lSignal.request == 0 {
				Read(nil)
			}
		}
	}
}
