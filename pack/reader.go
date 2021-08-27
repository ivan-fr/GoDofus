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

var lSignalSlave = &lastSignal{typeRequest: noType}
var pipelineSlave = new(pipe)
var lastWeftSlave *weft = nil

const (
	headerTwoFirstBytes = iota
	headerLength
	headerInstance
	messageLength
	noType
)

func GetClientPipeline() *pipe {
	return pipelineSlave
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

func (p *pipe) appendSlave(w *weft) {
	p.Wefts = append(p.Wefts, w)
}

func (p *pipe) GetSalve() *weft {
	if len(p.Wefts) == 0 {
		return nil
	}

	w := p.Wefts[0]
	p.Wefts = p.Wefts[1:]
	return w
}

func (lS *lastSignal) update(request int, typeRequest int, containForType []byte, containNoType []byte) {
	if typeRequest == lS.typeRequest && lS.request < 0 && request <= lS.request {
		return
	}

	if typeRequest == lS.typeRequest && lS.request < 0 && request > lS.request {
		lS.request = request
		lS.containForType = append(lS.containForType, containForType...)
	} else if lS.request == 0 || typeRequest == noType {
		lS.request = request
		lS.typeRequest = typeRequest
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

func commit(aLWeft **weft, aLSignal *lastSignal) {
	if aLSignal.request >= 0 {
		if *aLWeft == nil {
			return
		}

		if (*aLWeft).LengthType == 0 && (*aLWeft).waitLength {
			(*aLWeft).waitLength = false
			pipeline.append(*aLWeft)
			*aLWeft = nil
			return
		}

		if len((*aLWeft).Message) > 0 {
			return
		}

		if aLSignal.typeRequest == messageLength {
			(*aLWeft).Message = aLSignal.containForType
			pipeline.append(*aLWeft)
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

func readHeaderTwoFirstBytes(aLWeft **weft, aLSignal *lastSignal, reader *bytes.Reader) bool {
	ok := tryRead(aLSignal, reader, headerTwoFirstBytes, 2)

	if !ok {
		return false
	}

	twoBytes := binary.BigEndian.Uint16(lSignal.containForType)
	packetId := twoBytes >> 2
	lengthType := twoBytes & 0b11

	*aLWeft = &weft{PackId: packetId, LengthType: lengthType, waitLength: true}
	return true
}

func readHeaderInstance(aLWeft *weft, aLSignal *lastSignal, reader *bytes.Reader) bool {
	ok := tryRead(aLSignal, reader, headerInstance, 4)

	if !ok {
		return false
	}

	instanceID := binary.BigEndian.Uint32(lSignal.containForType)

	aLWeft.instanceID = instanceID
	return true
}

func readHeaderLength(aLWeft *weft, aLSignal *lastSignal, reader *bytes.Reader) bool {
	if aLWeft == nil {
		panic("incoherence last weft can't be nil")
	}

	ok := tryRead(aLSignal, reader, headerLength, uint(aLWeft.LengthType))

	if !ok {
		return false
	}

	switch aLWeft.LengthType {
	case 3:
		var specialCaseReader = bytes.NewReader(lSignal.containForType)
		var firstByte uint8
		_ = binary.Read(specialCaseReader, binary.BigEndian, &firstByte)
		var secondByte uint8
		_ = binary.Read(specialCaseReader, binary.BigEndian, &secondByte)
		var thirdByte uint8
		_ = binary.Read(specialCaseReader, binary.BigEndian, &thirdByte)
		aLWeft.Length = (uint32(firstByte) << 16) + (uint32(secondByte) << 8) + (uint32(thirdByte) & 255)
	case 2:
		aLWeft.Length = uint32(binary.BigEndian.Uint16(lSignal.containForType))
	case 1:
		aLWeft.Length = uint32(lSignal.containForType[0])
	case 0:
		aLWeft.Length = 0
	default:
		panic("wrong length type")
	}

	aLWeft.waitLength = false
	return true
}

func read(aLWeft **weft, aLSignal *lastSignal, isClient bool, bytesPack []byte) bool {
	switch lSignal.typeRequest {
	case messageLength:
		switch {
		case lSignal.request == 0:
			commit(aLWeft, aLSignal)
			lSignal.update(0, noType, nil, nil)
			return read(aLWeft, aLSignal, isClient, bytesPack)
		case lSignal.request > 0:
			commit(aLWeft, aLSignal)
			newBytesPack := append(lSignal.containNoType, bytesPack...)
			lSignal.update(0, noType, nil, nil)
			return read(aLWeft, aLSignal, isClient, newBytesPack)
		case lSignal.request < 0:
			reader := bytes.NewReader(bytesPack)
			ok := tryRead(aLSignal, reader, messageLength, uint(-lSignal.request))

			if !ok {
				return false
			}

			commit(aLWeft, aLSignal)
			newBytesPack := lSignal.containNoType
			lSignal.update(0, noType, nil, nil)
			return read(aLWeft, aLSignal, isClient, newBytesPack)
		}
	case headerTwoFirstBytes:
		switch {
		case lSignal.request < 0:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderTwoFirstBytes(aLWeft, aLSignal, reader)
			if !ok {
				return false
			}

			newBytesPack := lSignal.containNoType
			lSignal.update(0, noType, nil, nil)
			return read(aLWeft, aLSignal, isClient, newBytesPack)
		default:
			panic("incoherence, last signal can't be positive")
		}
	case headerLength:
		switch {
		case lSignal.request < 0:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderLength(*aLWeft, aLSignal, reader)
			if !ok {
				return false
			}

			newBytesPack := lSignal.containNoType
			lSignal.update(0, noType, nil, nil)
			return read(aLWeft, aLSignal, isClient, newBytesPack)
		default:
			panic("incoherence, last signal can't be positive")
		}
	case headerInstance:
		switch {
		case lSignal.request < 0:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderInstance(*aLWeft, aLSignal, reader)
			if !ok {
				return false
			}

			newBytesPack := lSignal.containNoType
			lSignal.update(0, noType, nil, nil)
			return read(aLWeft, aLSignal, isClient, newBytesPack)
		default:
			panic("incoherence, last signal can't be positive")
		}
	case noType:
		switch {
		case bytesPack == nil:
			commit(aLWeft, aLSignal)
			return true
		case lastWeft == nil:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderTwoFirstBytes(aLWeft, aLSignal, reader)
			if !ok {
				return false
			}

			newBytesPack := lSignal.containNoType
			lSignal.update(0, noType, nil, nil)
			return read(aLWeft, aLSignal, isClient, newBytesPack)
		case lastWeft.instanceID == 0 && isClient:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderInstance(*aLWeft, aLSignal, reader)
			if !ok {
				return false
			}

			newBytesPack := lSignal.containNoType
			lSignal.update(0, noType, nil, nil)
			return read(aLWeft, aLSignal, isClient, newBytesPack)
		case lastWeft.waitLength:
			reader := bytes.NewReader(bytesPack)
			ok := readHeaderLength(*aLWeft, aLSignal, reader)
			if !ok {
				return false
			}

			newBytesPack := lSignal.containNoType
			lSignal.update(0, noType, nil, nil)
			return read(aLWeft, aLSignal, isClient, newBytesPack)
		default:
			reader := bytes.NewReader(bytesPack)
			_ = tryRead(lSignal, reader, messageLength, uint(lastWeft.Length))
			return read(aLWeft, aLSignal, isClient, nil)
		}
	default:
		panic("program don't know the step.")
	}

	return false
}

func ReadServer(bytesPack []byte) bool {
	return read(&lastWeft, lSignal, false, bytesPack)
}

func ReadClient(bytesPack []byte) bool {
	return read(&lastWeftSlave, lSignalSlave, true, bytesPack)
}
