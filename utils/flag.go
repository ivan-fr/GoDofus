package utils

func SetFlag(a uint32, pos uint, b bool) uint32 {
	switch pos {
	case 0:
		if b {
			a |= 1
		} else {
			a &= 255 - 1
		}
	case 1:
		if b {
			a |= 2
		} else {
			a &= 255 - 2
		}
		break
	case 2:
		if b {
			a |= 4
		} else {
			a &= 255 - 4
		}
	case 3:
		if b {
			a |= 8
		} else {
			a &= 255 - 8
		}
	case 4:
		if b {
			a |= 16
		} else {
			a &= 255 - 16
		}
	case 5:
		if b {
			a |= 32
		} else {
			a &= 255 - 32
		}
	case 6:
		if b {
			a |= 64
		} else {
			a &= 255 - 64
		}
	case 7:
		if b {
			a |= 128
		} else {
			a &= 255 - 128
		}
	default:
		panic("Bytebox overflow.")
	}
	return a
}

func getFlag(a uint32, pos uint32) bool {
	switch pos {
	case 0:
		return (a & 1) != 0
	case 1:
		return (a & 2) != 0
	case 2:
		return (a & 4) != 0
	case 3:
		return (a & 8) != 0
	case 4:
		return (a & 16) != 0
	case 5:
		return (a & 32) != 0
	case 6:
		return (a & 64) != 0
	case 7:
		return (a & 128) != 0
	default:
		panic("Bytebox overflow.")
	}
}
