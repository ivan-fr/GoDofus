package commands

type Move struct {
	X       uint32
	Y       uint32
	worldId uint32
	MapId   float64
}

var WorldIdMax = uint32(2 << 12)
var MapCoordsMax = uint32(2 << 8)

func abs(a uint32) uint32 {
	if a < 0 {
		return -a
	}
	return a
}

func (m *Move) SetFromMapId() {
	m.worldId = uint32(int64(m.MapId) & 1073479680 >> 18)
	m.X = uint32(int64(m.MapId) >> 9 & 511)
	m.Y = uint32(int64(m.MapId) & 511)
	if (m.X & 256) == 256 {
		m.X = -(m.X & 255)
	}
	if (m.Y & 256) == 256 {
		m.Y = -(m.Y & 255)
	}
}

func (m *Move) SetFromCoords() bool {
	if m.X > MapCoordsMax || m.Y > MapCoordsMax || m.worldId > WorldIdMax {
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
	m.MapId = float64(worldValue<<18 | (xValue<<9 | yValue))
	return true
}
