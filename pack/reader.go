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
	Length     int
	Message    []byte
	waitLength bool
}

const (
	headerTwoFirstBytes = iota
	headerLength
	headerInstance
	messageLength
	noType
)

type Signal struct {
	request            int
	Type               int
	contentForType     []byte
	contentForNextType []byte
}

var noSignal = &Signal{Type: noType}

type Pipe struct {
	wefts []*Weft
	index int
}

func (p *Pipe) append(w *Weft) {
	p.wefts = append(p.wefts, w)
}

func (p *Pipe) Get() *Weft {
	if len(p.wefts) == 0 {
		return nil
	}

	w := p.wefts[0]
	p.wefts = p.wefts[1:]
	return w
}

func (p *Pipe) Contains(packetIDs map[uint16]bool) bool {
	var found []uint16
	for _, weftInPipe := range p.wefts {
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

func (lS *Signal) updateWith(newSignal *Signal) {
	if newSignal.Type == lS.Type && lS.request < 0 && newSignal.request <= lS.request {
		panic("incoherence from arguments")
	}

	if newSignal.Type == noType {
		lS.Type = newSignal.Type
		lS.contentForType = nil
	} else if newSignal.Type == lS.Type && lS.request < 0 && newSignal.request > lS.request {
		lS.request = newSignal.request
		lS.contentForType = append(lS.contentForType, newSignal.contentForType...)
		lS.contentForNextType = newSignal.contentForNextType
	} else if lS.request >= 0 && newSignal.Type != lS.Type {
		lS.request = newSignal.request
		lS.Type = newSignal.Type
		lS.contentForType = newSignal.contentForType
		lS.contentForNextType = newSignal.contentForNextType
	} else {
		panic("incoherence from arguments")
	}
}

type Reader struct {
	aLSignal  *Signal
	APipeline *Pipe
	aLWeft    *Weft
}

func (r *Reader) commit() {
	if r.aLSignal.request >= 0 {
		if r.aLWeft == nil {
			panic("cannot commit a nil weft")
		}

		if r.aLWeft.LengthType == 0 && r.aLWeft.waitLength {
			r.aLWeft.waitLength = false
			r.APipeline.append(r.aLWeft)
			r.aLWeft = nil
			return
		}

		if len(r.aLWeft.Message) > 0 {
			panic("incoherence in class")
		}

		if r.aLSignal.Type == messageLength {
			r.aLWeft.Message = r.aLSignal.contentForType
			r.APipeline.append(r.aLWeft)
			r.aLWeft = nil
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (r *Reader) tryRead(hotBytes []byte, step int, bytesWanted int) bool {
	reader := bytes.NewReader(append(r.aLSignal.contentForNextType, hotBytes...))

	nbBytesToRead := min(reader.Len(), bytesWanted) // <= bytesWanted && <= reader.Len
	nSignal := &Signal{contentForType: make([]byte, nbBytesToRead), Type: step}

	n, _ := io.ReadFull(reader, nSignal.contentForType)

	nSignal.request = n - bytesWanted
	ok := nSignal.request >= 0

	if nSignal.request > 0 {
		nSignal.contentForNextType, _ = io.ReadAll(reader)
	}

	r.aLSignal.updateWith(nSignal)
	return ok
}

func (r *Reader) readHeaderTwoFirstBytes(hotBytes []byte) bool {
	return r.readHeaderTwoFirstBytesRest(hotBytes, 2)
}

func (r *Reader) readHeaderTwoFirstBytesRest(hotBytes []byte, rest int) bool {
	ok := r.tryRead(hotBytes, headerTwoFirstBytes, rest)

	if !ok {
		return false
	}

	twoBytes := binary.BigEndian.Uint16(r.aLSignal.contentForType)
	packetId := twoBytes >> 2
	lengthType := twoBytes & 0b11

	r.aLWeft = &Weft{PackId: packetId, LengthType: lengthType, waitLength: true}
	return true
}

func (r *Reader) readHeaderInstance(hotBytes []byte) bool {
	return r.readHeaderInstanceRest(hotBytes, 4)
}

func (r *Reader) readHeaderInstanceRest(hotBytes []byte, rest int) bool {
	ok := r.tryRead(hotBytes, headerInstance, rest)

	if !ok {
		return false
	}

	instanceID := binary.BigEndian.Uint32(r.aLSignal.contentForType)

	r.aLWeft.instanceID = instanceID
	return true
}

func (r *Reader) readHeaderLength(hotBytes []byte) bool {
	return r.readHeaderLengthRest(hotBytes, int(r.aLWeft.LengthType))
}

func (r *Reader) readHeaderLengthRest(hotBytes []byte, rest int) bool {
	if r.aLWeft == nil {
		panic("incoherence last weft can't be nil")
	}

	ok := r.tryRead(hotBytes, headerLength, rest)

	if !ok {
		return false
	}

	switch r.aLWeft.LengthType {
	case 3:
		var firstByte = r.aLSignal.contentForType[0]
		var secondByte = r.aLSignal.contentForType[1]
		var thirdByte = r.aLSignal.contentForType[2]
		r.aLWeft.Length = (int(firstByte) << 16) + (int(secondByte) << 8) + (int(thirdByte) & 255)
	case 2:
		r.aLWeft.Length = int(binary.BigEndian.Uint16(r.aLSignal.contentForType))
	case 1:
		r.aLWeft.Length = int(r.aLSignal.contentForType[0])
	case 0:
		r.aLWeft.Length = 0
	default:
		panic("wrong length type")
	}

	r.aLWeft.waitLength = false
	return true
}

func (r *Reader) Read(toClient bool, hotBytes []byte) {
	switch r.aLSignal.Type {
	case messageLength:
		switch {
		case r.aLSignal.request < 0:
			ok := r.tryRead(hotBytes, messageLength, -r.aLSignal.request)

			if !ok {
				return
			}

			r.commit()
			r.aLSignal.updateWith(noSignal)
			return
		}
	case headerTwoFirstBytes:
		switch {
		case r.aLSignal.request < 0:
			ok := r.readHeaderTwoFirstBytesRest(hotBytes, -r.aLSignal.request)
			if !ok {
				return
			}

			r.aLSignal.updateWith(noSignal)
			r.Read(toClient, nil)
		default:
			panic("incoherence, last signal can't be positive")
		}
	case headerLength:
		switch {
		case r.aLSignal.request < 0:
			ok := r.readHeaderLengthRest(hotBytes, -r.aLSignal.request)
			if !ok {
				return
			}

			r.aLSignal.updateWith(noSignal)
			r.Read(toClient, nil)
		default:
			panic("incoherence, last signal can't be positive")
		}
	case headerInstance:
		switch {
		case r.aLSignal.request < 0:
			ok := r.readHeaderInstanceRest(hotBytes, -r.aLSignal.request)
			if !ok {
				return
			}

			r.aLSignal.updateWith(noSignal)
			r.Read(toClient, nil)
		default:
			panic("incoherence, last signal can't be positive")
		}
	case noType:
		switch {
		case r.aLWeft == nil:
			ok := r.readHeaderTwoFirstBytes(hotBytes)
			if !ok {
				return
			}

			r.aLSignal.updateWith(noSignal)
			r.Read(toClient, nil)
		case r.aLWeft.LengthType == 0 && r.aLWeft.waitLength:
			r.commit()
		case r.aLWeft.instanceID == 0 && toClient:
			ok := r.readHeaderInstance(nil)
			if !ok {
				return
			}

			r.aLSignal.updateWith(noSignal)
			r.Read(toClient, nil)
		case r.aLWeft.waitLength:
			ok := r.readHeaderLength(nil)
			if !ok {
				return
			}

			r.aLSignal.updateWith(noSignal)
			r.Read(toClient, nil)
		default:
			ok := r.tryRead(nil, messageLength, r.aLWeft.Length)

			if !ok {
				return
			}

			r.commit()
			r.aLSignal.updateWith(noSignal)
			return
		}
	default:
		panic("program don't know the step.")
	}
}

func NewReader() *Reader {
	return &Reader{
		aLWeft:    nil,
		aLSignal:  &Signal{Type: noType},
		APipeline: new(Pipe),
	}
}
