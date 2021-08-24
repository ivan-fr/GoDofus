package structs

import (
	"strconv"
	"strings"
)

type version struct {
	Major     uint8
	Minor     uint8
	Code      uint8
	Build     uint32
	BuildType uint8
}

const (
	RELEASE = iota
)

var version_ = getVersion("2.60.4-13", RELEASE)

func getVersion(args ...interface{}) *version {
	v := &version{}

	if len(args) > 0 {
		string_0, ok := args[0].(string)
		if !(len(args) == 2 && ok) {
			panic("invalid parameters")
		}

		split := strings.Split(string_0, ".")
		major, _ := strconv.ParseUint(split[0], 10, 8)
		v.Major = uint8(major)

		minor, _ := strconv.ParseUint(split[1], 10, 8)
		v.Minor = uint8(minor)

		codeStr := strings.Split(split[2], "-")
		code, _ := strconv.ParseUint(codeStr[0], 10, 8)
		v.Code = uint8(code)

		build, _ := strconv.ParseUint(codeStr[1], 10, 8)
		v.Build = uint32(build)

		v.BuildType = uint8(args[1].(int))
	}

	return v
}

func GetVersionNOA() *version {
	return version_
}
