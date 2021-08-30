package generates

import (
	"GoDofus/messages"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var packageTemplate = template.Must(template.New("").Parse(`// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// {{ .Timestamp }}

package messages

import (
	"bytes"
	"fmt"
)

type {{ .Name }} struct {
	PacketId uint32
{{range $index, $element := .StructFields}}	{{$element}}
{{end}}}

var {{ .Name }}Map = make(map[uint]*{{ .Name }})

func Get{{ .NameCapFirst }}NOA(instance uint) *{{ .Name }} {
	{{ .Name }}_, ok := {{ .Name }}Map[instance]

	if ok {
		return {{ .Name }}_
	}

	{{ .Name }}Map[instance] = &{{ .Name }}{PacketId: {{ .NameCapFirst }}ID}
	return {{ .Name }}Map[instance]
}

func ({{.FistLetter}} *{{.Name}}) Serialize(buff *bytes.Buffer) {
{{range $index, $element := .SerializerString}}	{{$element}}
{{end}}}

func ({{.FistLetter}} *{{.Name}}) Deserialize(reader *bytes.Reader) {
{{range $index, $element := .DeserializerString}}	{{$element}}
{{end}}}

func ({{.FistLetter}} *{{.Name}}) GetPacketId() uint32 {
	return {{ .FistLetter }}.PacketId
}

func ({{.FistLetter}} *{{.Name}}) String() string {
	return fmt.Sprintf("packetId: %d\n", {{ .FistLetter }}.PacketId)
}
`))

var packageIDSTemplate = template.Must(template.New("").Parse(`// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// {{ .Timestamp }}

package messages

const (
	{{ .ConstBefore }}	{{ .NameCapFirst }}ID = {{ .Id }}
)
`))

var structFields []string
var serializerString []string
var deserializerString []string

var instance uint

type varInt16 bool
type varInt32 bool
type varInt64 bool
type varUInt64 bool
type flag bool
type utf bool

func putStringSimpleType(firstLetter, variableName, variableType string) {
	structFields = append(structFields, fmt.Sprintf("%s %s", variableName, variableType))
	serializerString = append(serializerString, fmt.Sprintf(
		`_ = binary.Write(buff, binary.BigEndian, %s.%s)`, firstLetter, variableName))
	deserializerString = append(deserializerString, fmt.Sprintf(
		`_ = binary.Read(reader, binary.BigEndian, &%s.%s)`, firstLetter, variableName))
}

func putStringFlag(firstLetter, variableName string) {
	for _, name := range strings.Split(variableName, ",") {
		structFields = append(structFields, fmt.Sprintf("%s bool", name))
	}

	var stringerS = fmt.Sprintf("var box%d uint32", instance)
	for i, name := range strings.Split(variableName, ",") {
		stringerS = fmt.Sprintf("%s\n	box%d = utils.SetFlag(box%d, %d, %s.%s)", stringerS, instance, instance, i, firstLetter, name)
	}
	stringerS = fmt.Sprintf("%s\n	_ = binary.Write(buff, binary.BigEndian, byte(box%d))", stringerS, instance)

	serializerString = append(serializerString, stringerS)

	var stringerD = fmt.Sprintf(`var box%d byte
	_ = binary.Read(reader, binary.BigEndian, &box%d)`, instance, instance)

	for i, name := range strings.Split(variableName, ",") {
		stringerD = fmt.Sprintf("%s\n	%s.%s = utils.GetFlag(uint32(box%d), %d)", stringerD, firstLetter, name, instance, i)
	}

	deserializerString = append(deserializerString, stringerD)
}

func putStringSimpleVarSlice(firstLetter, variableName, variableVarType, variableType string) {
	var caster string
	if variableType == "float64" {
		caster = fmt.Sprintf("float64(utils.Read%s(reader))", variableVarType)
	} else {
		caster = fmt.Sprintf("utils.Read%s(reader)", variableVarType)
	}

	structFields = append(structFields, fmt.Sprintf("%s []%s", variableName, variableType))
	serializerString = append(serializerString, fmt.Sprintf(`_ = binary.Write(buff, binary.BigEndian, uint16(len(%s.%s)))
	for i := 0; i < len(%s.%s); i++ {
		utils.Write%s(buff, %s.%s[i])
	}`, firstLetter, variableName, firstLetter, variableName, variableVarType, firstLetter, variableName))

	deserializerString = append(deserializerString, fmt.Sprintf(`var len%d_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len%d_)
	%s.%s = nil
	for i := 0; i < int(len%d_); i++ {
		%s.%s = append(%s.%s, %s)
	}`, instance, instance, firstLetter, variableName, instance, firstLetter, variableName, firstLetter, variableName, caster))
}

func putStringMessageSlice(firstLetter, messageName string) {
	messageName = string(instance) + messageName
	structFields = append(structFields, fmt.Sprintf("%s []*%s", messageName, messageName))
	serializerString = append(serializerString, fmt.Sprintf(`_ = binary.Write(buff, binary.BigEndian, uint16(len(%s.%s)))
	for i := 0; i < len(%s.%s); i++ {
		%s.%s[i].Serialize(buff)
	}`, firstLetter, messageName, firstLetter, messageName, firstLetter, messageName))

	deserializerString = append(deserializerString, fmt.Sprintf(`var len%d_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len%d_)
	%s.%s = nil
	for i := 0; i < int(len%d_); i++ {
		aMessage%d := new(%s)
		aMessage%d.Deserialize(reader)
		%s.%s = append(%s.%s, aMessage%d)
	}`, instance, instance, firstLetter, messageName, instance, instance, messageName, instance, firstLetter, messageName, firstLetter, messageName, instance))
}

func putStringSimpleVarType(firstLetter, variableName, variableVarType, variableType string) {
	structFields = append(structFields, fmt.Sprintf("%s %s", variableName, variableType))
	serializerString = append(serializerString, fmt.Sprintf(
		`utils.Write%s(buff, %s.%s)`, variableVarType, firstLetter, variableName))

	var caster string
	if variableType == "float64" {
		caster = fmt.Sprintf("float64(utils.Read%s(reader))", variableVarType)
	} else {
		caster = fmt.Sprintf("utils.Read%s(reader)", variableVarType)
	}

	deserializerString = append(deserializerString, fmt.Sprintf(
		`%s.%s = %s`, firstLetter, variableName, caster))
}

func putStringSimpleSlice(firstLetter, variableName, variableType string) {
	structFields = append(structFields, fmt.Sprintf("%v []%s", variableName, variableType))
	serializerString = append(serializerString, fmt.Sprintf(`_ = binary.Write(buff, binary.BigEndian, uint16(len(%s.%s)))
	_ = binary.Write(buff, binary.BigEndian, %s.%s)`, firstLetter, variableName, firstLetter, variableName))
	deserializerString = append(deserializerString, fmt.Sprintf(`var len%d_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len%d_)
	%s.%s = make([]%s, len%d_)
	_ = binary.Read(reader, binary.BigEndian, %s.%s)`, instance, instance, firstLetter, variableName, variableType, instance, firstLetter, variableName))
}

func scan(sentence string, random bool) string {
	var type_ string
	fmt.Println(sentence)
	_, err := fmt.Scanln(&type_)
	if err != nil {
		if random {
			return fmt.Sprintf("var%d", instance)
		}
		return ""
	}
	return type_
}

func dispatchSerializer() interface{} {
	type_ := scan("Enter variable type: ", false)
	switch type_ {
	case "-varUInt64":
		return []varUInt64{}
	case "-varInt64":
		return []varInt64{}
	case "-varInt32":
		return []varInt32{}
	case "-varInt16":
		return []varInt16{}
	case "varUInt64":
		var t varUInt64
		return t
	case "varInt64":
		var t varInt64
		return t
	case "varInt32":
		var t varInt32
		return t
	case "varInt16":
		var t varInt16
		return t
	case "flag":
		var t flag
		return t
	case "bool":
		var t bool
		return t
	case "byte":
		var t byte
		return t
	case "int16":
		var t int16
		return t
	case "uint16":
		var t uint16
		return t
	case "int32":
		var t int32
		return t
	case "uint32":
		var t uint32
		return t
	case "int64":
		var t int64
		return t
	case "uint64":
		var t uint64
		return t
	case "float64":
		var t float64
		return t
	case "utf":
		var t utf
		return t
	case "-byte":
		return []byte{}
	case "-int16":
		return []int16{}
	case "-uint16":
		return []uint16{}
	case "-int32":
		return []int32{}
	case "-uint32":
		return []uint32{}
	case "-int64":
		return []int64{}
	case "-uint64":
		return []uint64{}
	case "-float64":
		return []float64{}
	case "-utf":
		return []utf{}
	case "":
		return nil
	default:
		var isSlice bool
		if type_[0] == '-' {
			type_ = type_[1:]
			isSlice = true
		}

		parseInt, err := strconv.ParseInt(type_, 10, 32)
		if err != nil {
			fmt.Println("not found")
			return dispatchSerializer()
		}

		message, ok := messages.Types_[int(parseInt)].(messages.Message)
		if !ok {
			fmt.Println("not found")
			return dispatchSerializer()
		}

		if isSlice {
			return []messages.Message{message}
		}
		return message
	}
}

func serializer(i interface{}, firstLetter string, variableName string) {
	switch i.(type) {
	case []varUInt64:
		putStringSimpleVarSlice(firstLetter, variableName, "VarUInt16", "float64")
	case []varInt64:
		putStringSimpleVarSlice(firstLetter, variableName, "VarInt64", "float64")
	case []varInt32:
		putStringSimpleVarSlice(firstLetter, variableName, "VarInt32", "int32")
	case []varInt16:
		putStringSimpleVarSlice(firstLetter, variableName, "VarInt16", "int32")
	case []utf:
		putStringSimpleVarSlice(firstLetter, variableName, "UTF", "[]byte")
	case varUInt64:
		putStringSimpleVarType(firstLetter, variableName, "VarUInt64", "float64")
	case varInt64:
		putStringSimpleVarType(firstLetter, variableName, "VarInt64", "float64")
	case varInt32:
		putStringSimpleVarType(firstLetter, variableName, "VarInt32", "int32")
	case varInt16:
		putStringSimpleVarType(firstLetter, variableName, "VarInt16", "int32")
	case utf:
		putStringSimpleVarType(firstLetter, variableName, "UTF", "[]byte")
	case flag:
		putStringFlag(firstLetter, variableName)
	case bool:
		putStringSimpleType(firstLetter, variableName, "bool")
	case byte:
		putStringSimpleType(firstLetter, variableName, "byte")
	case int16:
		putStringSimpleType(firstLetter, variableName, "int16")
	case uint16:
		putStringSimpleType(firstLetter, variableName, "uint16")
	case int32:
		putStringSimpleType(firstLetter, variableName, "int32")
	case uint32:
		putStringSimpleType(firstLetter, variableName, "uint32")
	case int64:
		putStringSimpleType(firstLetter, variableName, "int64")
	case uint64:
		putStringSimpleType(firstLetter, variableName, "uint64")
	case float64:
		putStringSimpleType(firstLetter, variableName, "float64")
	case []byte:
		putStringSimpleSlice(firstLetter, variableName, "byte")
	case []int16:
		putStringSimpleSlice(firstLetter, variableName, "int16")
	case []uint16:
		putStringSimpleSlice(firstLetter, variableName, "uint16")
	case []int32:
		putStringSimpleSlice(firstLetter, variableName, "int32")
	case []uint32:
		putStringSimpleSlice(firstLetter, variableName, "uint32")
	case []int64:
		putStringSimpleSlice(firstLetter, variableName, "int64")
	case []uint64:
		putStringSimpleSlice(firstLetter, variableName, "uint64")
	case []float64:
		putStringSimpleSlice(firstLetter, variableName, "float64")
	case []messages.Message:
		t := reflect.TypeOf(i.([]messages.Message)[0])
		messageName := t.Elem().Name()
		fmt.Printf("%s\n", messageName)
		putStringMessageSlice(firstLetter, messageName)
	case messages.Message:
		t := reflect.TypeOf(i.(messages.Message))
		messageName := t.Elem().Name()
		fmt.Printf("%s\n", messageName)
		messageName = string(instance) + messageName
		structFields = append(structFields, fmt.Sprintf("%s *%s", messageName, messageName))
		serializerString = append(serializerString, fmt.Sprintf(`%s.%s.Serialize(buff)`, firstLetter, messageName))
		deserializerString = append(deserializerString, fmt.Sprintf(`%s.%s = new(%s)
	%s.%s.Deserialize(reader)`, firstLetter, messageName, messageName, firstLetter, messageName))
	}

	instance++
}

func GenerateMessage(name string, packetId uint32) {
	name = strings.ToLower(name[:1]) + name[1:]

	for interfaceType := dispatchSerializer(); interfaceType != nil; interfaceType = dispatchSerializer() {
		name_ := scan("Enter variable name:", true)
		serializer(interfaceType, (name)[:2], name_)
	}

	constIDS, _ := os.ReadFile("./messages/IDS.go")

	constIDString := string(constIDS)
	var rgx = regexp.MustCompile(`^[^(]+\([^a-zA-z]*([^)]+).*`)
	rs := rgx.FindStringSubmatch(constIDString)

	f, err := os.Create(fmt.Sprintf("./messages/%s.go", name))
	if err != nil {
		panic(err)
	}

	f2, err := os.Create("./messages/IDS.go")
	if err != nil {
		panic(err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
		err = f2.Close()
		if err != nil {
			panic(err)
		}
	}()

	err = packageIDSTemplate.Execute(f2, struct {
		Timestamp    time.Time
		NameCapFirst string
		Id           uint32
		ConstBefore  string
	}{
		Timestamp:    time.Now(),
		NameCapFirst: strings.Title(name),
		Id:           packetId,
		ConstBefore:  rs[1],
	})
	if err != nil {
		panic(err)
	}

	err = packageTemplate.Execute(f, struct {
		Timestamp          time.Time
		NameCapFirst       string
		Name               string
		FistLetter         string
		StructFields       []string
		SerializerString   []string
		DeserializerString []string
	}{
		Timestamp:          time.Now(),
		Name:               name,
		NameCapFirst:       strings.Title(name),
		FistLetter:         (name)[:2],
		StructFields:       structFields,
		SerializerString:   serializerString,
		DeserializerString: deserializerString,
	})
	if err != nil {
		panic(err)
	}

	BuildRegistry()
}
