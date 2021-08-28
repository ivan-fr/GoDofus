package pack

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Weft struct {
	PackId     uint16
	LengthType uint16
	instanceID uint32
	Length     uint32
	Message    []byte
	waitLength bool
}

type lastSignal struct {
	request        int
	TypeRequest    int
	containForType []byte
	containNoType  []byte
}

type Pipe struct {
	Wefts []*Weft
	index int
}

const (
	headerTwoFirstBytes = iota
	headerLength
	headerInstance
	messageLength
	noType
)

func (p *Pipe) append(w *Weft) {
	p.Wefts = append(p.Wefts, w)
}

func (p *Pipe) Get() *Weft {
	if len(p.Wefts) == 0 {
		return nil
	}

	w := p.Wefts[0]
	p.Wefts = p.Wefts[1:]
	return w
}

func (p *Pipe) Contains(packetIDs map[uint16]bool) bool {
	var found []uint16
	for _, weftInPipe := range p.Wefts {
		_, ok := packetIDs[weftInPipe.PackId]

		if ok {
			found = append(found, weftInPipe.PackId)

			if len(found) == len(packetIDs) {
				return true
			}
		}
	}

	return false
}

func (lS *lastSignal) update(request int, typeRequest int, containForType []byte, containNoType []byte) {
	if typeRequest == lS.TypeRequest && lS.request < 0 && request <= lS.request {
		return
	}

	if typeRequest == lS.TypeRequest && lS.request < 0 && request > lS.request {
		lS.request = request
		lS.containForType = append(lS.containForType, containForType...)
	} else if lS.request == 0 || typeRequest == noType {
		lS.request = request
		lS.TypeRequest = typeRequest
		lS.containForType = containForType
		lS.containNoType = nil
	} else {
		panic("incoherence from lastSignal")
	}

	if lS.request > 0 {
		if containNoType == nil {
			panic("incoherence contain no type can't be nil")
		}

		lS.containNoType = containNoType
	}
}

func commit(aPipeline *Pipe, aLWeft **Weft, aLSignal *lastSignal) {
	if aLSignal.request >= 0 {
		if *aLWeft == nil {
			return
		}

		if (*aLWeft).LengthType == 0 && (*aLWeft).waitLength {
			(*aLWeft).waitLength = false
			aPipeline.append(*aLWeft)
			*aLWeft = nil
			return
		}

		if len((*aLWeft).Message) > 0 {
			return
		}

		if aLSignal.TypeRequest == messageLength {
			(*aLWeft).Message = aLSignal.containForType
			aPipeline.append(*aLWeft)
			*aLWeft = nil
		}
	}
}

func min(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}

func tryRead(aLSignal *lastSignal, reader *bytes.Reader, step int, bytesWanted uint) bool {
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

	aLSignal.update(request, step, containForType, containNoType)
	return ok
}

func readHeaderTwoFirstBytes(aLWeft **Weft, aLSignal *lastSignal, reader *bytes.Reader) bool {
	ok := tryRead(aLSignal, reader, headerTwoFirstBytes, 2)

	if !ok {
		return false
	}

	twoBytes := binary.BigEndian.Uint16(aLSignal.containForType)
	packetId := twoBytes >> 2
	lengthType := twoBytes & 0b11

	*aLWeft = &Weft{PackId: packetId, LengthType: lengthType, waitLength: true}
	return true
}

func readHeaderInstance(aLWeft *Weft, aLSignal *lastSignal, reader *bytes.Reader) bool {
	ok := tryRead(aLSignal, reader, headerInstance, 4)

	if !ok {
		return false
	}

	instanceID := binary.BigEndian.Uint32(aLSignal.containForType)

	aLWeft.instanceID = instanceID
	return true
}

func readHeaderLength(aLWeft *Weft, aLSignal *lastSignal, reader *bytes.Reader) bool {
	if aLWeft == nil {
		panic("incoherence last weft can't be nil")
	}

	ok := tryRead(aLSignal, reader, headerLength, uint(aLWeft.LengthType))

	if !ok {
		return false
	}

	switch aLWeft.LengthType {
	case 3:
		var specialCaseReader = bytes.NewReader(aLSignal.containForType)
		var firstByte uint8
		_ = binary.Read(specialCaseReader, binary.BigEndian, &firstByte)
		var secondByte uint8
		_ = binary.Read(specialCaseReader, binary.BigEndian, &secondByte)
		var thirdByte uint8
		_ = binary.Read(specialCaseReader, binary.BigEndian, &thirdByte)
		aLWeft.Length = (uint32(firstByte) << 16) + (uint32(secondByte) << 8) + (uint32(thirdByte) & 255)
	case 2:
		aLWeft.Length = uint32(binary.BigEndian.Uint16(aLSignal.containForType))
	case 1:
		aLWeft.Length = uint32(aLSignal.containForType[0])
	case 0:
		aLWeft.Length = 0
	default:
		panic("wrong length type")
	}

	aLWeft.waitLength = false
	return true
}

func read(aPipeline *Pipe, aLWeft **Weft, aLSignal *lastSignal, toClient bool, bytesPack []byte) bool {
	switch aLSignal.TypeRequest {
	case messageLength:
		switch {
		case aLSignal.request == 0:
			commit(aPipeline, aLWeft, aLSignal)
			aLSignal.update(0, noType, nil, nil)
			return read(aPipeline, aLWeft, aLSignal, toClient, bytesPack)
		case aLSignal.request > 0:
			commit(aPipeline, aLWeft, aLSignal)
			newBytesPack := append(aLSignal.containNoType, bytesPack...)
			aLSignal.update(0, noType, nil, nil)
			return read(aPipeline, aLWeft, aLSignal, toClient, newBytesPack)
		case aLSignal.request < 0:
			reader := bytes.NewReader(bytesPack)
			ok := tryRead(aLSignal, reader, messageLength, uint(-aLSignal.request))

			if !ok {
				return false
			}

			commit(aPipeline, aLWeft, aLSignal)
			newBytesPack := aLSignal.containNoType
			aLSignal.update(0, noType, nil, nil)
			return read(aPipeline, aLWeft, aLSignal, toClient, newBytesPack)
		}
	case headerTwoFirstBytes:
		switch {
		case aLSignal.request < 0:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderTwoFirstBytes(aLWeft, aLSignal, reader)
			if !ok {
				return false
			}

			newBytesPack := aLSignal.containNoType
			aLSignal.update(0, noType, nil, nil)
			return read(aPipeline, aLWeft, aLSignal, toClient, newBytesPack)
		default:
			panic("incoherence, last signal can't be positive")
		}
	case headerLength:
		switch {
		case aLSignal.request < 0:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderLength(*aLWeft, aLSignal, reader)
			if !ok {
				return false
			}

			newBytesPack := aLSignal.containNoType
			aLSignal.update(0, noType, nil, nil)
			return read(aPipeline, aLWeft, aLSignal, toClient, newBytesPack)
		default:
			panic("incoherence, last signal can't be positive")
		}
	case headerInstance:
		switch {
		case aLSignal.request < 0:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderInstance(*aLWeft, aLSignal, reader)
			if !ok {
				return false
			}

			newBytesPack := aLSignal.containNoType
			aLSignal.update(0, noType, nil, nil)
			return read(aPipeline, aLWeft, aLSignal, toClient, newBytesPack)
		default:
			panic("incoherence, last signal can't be positive")
		}
	case noType:
		switch {
		case bytesPack == nil:
			commit(aPipeline, aLWeft, aLSignal)
			return true
		case (*aLWeft) == nil:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderTwoFirstBytes(aLWeft, aLSignal, reader)
			if !ok {
				return false
			}

			newBytesPack := aLSignal.containNoType
			aLSignal.update(0, noType, nil, nil)
			return read(aPipeline, aLWeft, aLSignal, toClient, newBytesPack)
		case (*aLWeft).instanceID == 0 && toClient:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderInstance(*aLWeft, aLSignal, reader)
			if !ok {
				return false
			}

			newBytesPack := aLSignal.containNoType
			aLSignal.update(0, noType, nil, nil)
			return read(aPipeline, aLWeft, aLSignal, toClient, newBytesPack)
		case (*aLWeft).waitLength:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderLength(*aLWeft, aLSignal, reader)
			if !ok {
				return false
			}

			newBytesPack := aLSignal.containNoType
			aLSignal.update(0, noType, nil, nil)
			return read(aPipeline, aLWeft, aLSignal, toClient, newBytesPack)
		default:
			reader := bytes.NewReader(bytesPack)
			_ = tryRead(aLSignal, reader, messageLength, uint((*aLWeft).Length))
			return read(aPipeline, aLWeft, aLSignal, toClient, nil)
		}
	default:
		panic("program don't know the step.")
	}

	return false
}

func NewServerReader() (func([]byte) bool, *Pipe) {
	var aLSignal = &lastSignal{TypeRequest: noType}
	var aPipeline = new(Pipe)
	var aLWeft *Weft = nil
	return func(bytesPack []byte) bool { return read(aPipeline, &aLWeft, aLSignal, false, bytesPack) }, aPipeline
}

func NewClientReader() (func([]byte) bool, *Pipe) {
	var aLSignal = &lastSignal{TypeRequest: noType}
	var aPipeline = new(Pipe)
	var aLWeft *Weft = nil
	return func(bytesPack []byte) bool { return read(aPipeline, &aLWeft, aLSignal, true, bytesPack) }, aPipeline
}
