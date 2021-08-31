package commands

type Move struct {
	X       int16
	Y       int16
	worldId uint32
	MapId   float64
}

var WorldIdMax = 2 << 12
var MapCoordsMax = 2 << 8

func abs(a int16) int16 {
	if a < 0 {
		return -a
	}
	return a
}

func (m *Move) SetFromMapId() {
	m.worldId = uint32(uint64(m.MapId) & 1073479680 >> 18)
	m.X = int16(int64(m.MapId) >> 9 & 511)
	m.Y = int16(int64(m.MapId) & 511)
	if (m.X & 256) == 256 {
		m.X = -(m.X & 255)
	}
	if (m.Y & 256) == 256 {
		m.Y = -(m.Y & 255)
	}
}

func (m *Move) SetFromCoords() bool {
	if int(m.X) > MapCoordsMax || int(m.Y) > MapCoordsMax || int(m.worldId) > WorldIdMax {
		return false
	}
	var worldValue = m.worldId & 4095
	var xValue = abs(m.X) & 255
	if m.X < 0 {
		xValue |= 256
	}
	var yValue = abs(m.Y) & 255
	if m.Y < 0 {
		yValue |= 256
	}
	m.MapId = float64(worldValue<<18 | (uint32(xValue)<<9 | uint32(yValue)))

	return true
}
