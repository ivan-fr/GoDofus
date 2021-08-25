package pack

import (
	"bytes"
	"encoding/binary"
	"io"
)

type weft struct {
	PackId     uint16
	LengthType uint16
	instanceID uint32
	Length     uint32
	Message    []byte
	waitLength bool
}

type lastSignal struct {
	request        int
	typeRequest    int
	containForType []byte
	containNoType  []byte
}

type pipe struct {
	Wefts []*weft
	index int
}

var lSignal = &lastSignal{typeRequest: noType}
var pipeline = new(pipe)
var lastWeft *weft = nil
var client = false

const (
	headerTwoFirstBytes = iota
	headerLength
	headerInstance
	messageLength
	noType
)

func ToggleClient() {
	client = !client
}

func GetPipeline() *pipe {
	return pipeline
}

func (p *pipe) append(w *weft) {
	p.Wefts = append(p.Wefts, w)
}

func (p *pipe) Get() *weft {
	if len(p.Wefts) == 0 {
		return nil
	}

	w := p.Wefts[0]
	p.Wefts = p.Wefts[1:]
	return w
}

func (lSignal *lastSignal) update(request int, typeRequest int, containForType []byte, containNoType []byte) {
	if typeRequest == lSignal.typeRequest && lSignal.request < 0 && request <= lSignal.request {
		return
	}

	if typeRequest == lSignal.typeRequest && lSignal.request < 0 && request > lSignal.request {
		lSignal.request = request
		lSignal.containForType = append(lSignal.containForType, containForType...)
	} else if lSignal.request == 0 || typeRequest == noType {
		lSignal.request = request
		lSignal.typeRequest = typeRequest
		lSignal.containForType = containForType
		lSignal.containNoType = nil
	} else {
		panic("incoherence from lastSignal")
	}

	if lSignal.request > 0 {
		if containNoType == nil {
			panic("incoherence contain no type can't be nil")
		}

		lSignal.containNoType = containNoType
	}
}

func commit() {
	if lSignal.request >= 0 && lSignal.typeRequest == messageLength {
		if lastWeft == nil {
			return
		}

		if lastWeft.LengthType == 0 && lastWeft.waitLength {
			lastWeft.waitLength = false
			pipeline.append(lastWeft)
			lastWeft = nil
			return
		}

		if len(lastWeft.Message) > 0 {
			return
		}

		lastWeft.Message = lSignal.containForType
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

	lastWeft = &weft{PackId: packetId, LengthType: lengthType, waitLength: true}
	return true
}

func readHeaderInstance(reader *bytes.Reader) bool {
	ok := tryRead(reader, headerInstance, 4)

	if !ok {
		return false
	}

	instanceID := binary.BigEndian.Uint32(lSignal.containForType)

	lastWeft.instanceID = instanceID
	return true
}

func readHeaderLength(reader *bytes.Reader) bool {
	if lastWeft == nil {
		panic("incoherence last weft can't be nil")
	}

	ok := tryRead(reader, headerLength, uint(lastWeft.LengthType))

	if !ok {
		return false
	}

	switch lastWeft.LengthType {
	case 3:
		var specialCaseReader = bytes.NewReader(lSignal.containForType)
		var firstByte uint8
		_ = binary.Read(specialCaseReader, binary.BigEndian, &firstByte)
		var twoBytes uint16
		_ = binary.Read(specialCaseReader, binary.BigEndian, &twoBytes)
		lastWeft.Length = uint32(firstByte)<<16 | uint32(twoBytes)
	case 2:
		lastWeft.Length = uint32(binary.BigEndian.Uint16(lSignal.containForType))
	case 1:
		lastWeft.Length = uint32(lSignal.containForType[0])
	case 0:
		lastWeft.Length = 0
	default:
		panic("wrong length type")
	}

	lastWeft.waitLength = false
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
	case headerInstance:
		switch {
		case lSignal.request < 0:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderInstance(reader)
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
		case lastWeft.instanceID == 0 && client:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderInstance(reader)
			if !ok {
				return false
			}

			newBytePack := lSignal.containNoType
			lSignal.update(0, noType, nil, nil)
			return Read(newBytePack)
		case lastWeft.waitLength:
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
			_ = tryRead(reader, messageLength, uint(lastWeft.Length))
			return Read(nil)
		}
	default:
		panic("program don't know the step.")
	}

	return false
}
