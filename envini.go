package envini

import (
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "json: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "json: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "json: Unmarshal(nil " + e.Type.String() + ")"
}

func GetConfig(configPath string, config interface{}) {
	log.Printf("EnvINI.GetConfig() | configPath: %q | config: %v\n", configPath, config)

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := Unmarshal(data, config); err != nil {
		// return err
	}
}

func Unmarshal(data []byte, config interface{}) error {
	log.Printf("EnvINI.Unmarshal() | data: %q | config: %#v\n", data, config)

	dataMap := bytesToDataMap(data)
	log.Printf("EnvINI.Unmarshal() | dataMap: %q\n", dataMap)

	rv := reflect.ValueOf(config) // reflect.Value
	// isPointer := rv.Kind() == reflect.Ptr
	// log.Printf("EnvINI.Unmarshal() | rv.Kind(): %#v | isPointer: %t | rv.IsNil(): %t\n", rv.Kind(), isPointer, rv.IsNil())
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(config)}
	}

	log.Printf("EnvINI.Unmarshal() | rv.IsValid(): %t\n", rv.IsValid())

	t := reflect.TypeOf(config).Elem() // Dereference the interface{} to a reflect.Value
	// field, _ := t.FieldByName("ProjectName") //alternative
	// log.Printf("EnvINI.Unmarshal() | ProjectName Field Tag: %q\n", field.Tag.Get("ini"))

	// log.Printf("EnvINI.Unmarshal() | ini Field Tag: %q\n", field.Tag.Get("ini"))
	// log.Printf("EnvINI.Unmarshal() | t.NumField(): %d\n", t.NumField())

	section := "GLOBAL"
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("ini")
		// log.Printf("EnvINI.Unmarshal() | Config Field Tag: %q\n", field.Tag.Get("ini"))
		// log.Printf("EnvINI.Unmarshal() | %d | field.Name: %v (%v) | field.Tag: %v | tag: %q\n", i, field.Name, field.Type.Name(), field.Tag, tag)
		log.Printf("EnvINI.Unmarshal() | %d | field.Name: %s (%s) | tag: %q\n", i, field.Name, field.Type.Name(), tag)
		if value, ok := dataMap[section][tag]; ok {
			switch field.Type.Name() {
			case "bool":
				boolValue := false
				value = strings.ToLower(value)
				switch value {
				case "true", "t", "yes", "y", "1":
					boolValue = true
				}
				v := rv.Elem().FieldByName(field.Name)
				if v.IsValid() {
					v.SetBool(boolValue)
				}
			case "int":
				v := rv.Elem().FieldByName(field.Name)
				if v.IsValid() {
					intValue, err := strconv.ParseInt(value, 10, 64)
					if err == nil {
						v.SetInt(intValue)
					}
				}
			case "string":
				v := rv.Elem().FieldByName(field.Name)
				if v.IsValid() {
					v.SetString(value)
				}
			}
		}
	}

	// switch rv.Kind() {
	// case reflect.Struct:
	// 	log.Println("EnvINI.Unmarshal() | Struct")
	// case reflect.Slice:
	// 	log.Println("EnvINI.Unmarshal() | Slice")
	// case reflect.Array:
	// 	log.Println("EnvINI.Unmarshal() | Array")
	// case reflect.Map:
	// 	log.Println("EnvINI.Unmarshal() | Map")
	// case reflect.Ptr:
	// 	log.Println("EnvINI.Unmarshal() | Ptr")
	// case reflect.Interface:
	// 	log.Println("EnvINI.Unmarshal() | Interface")
	// default:
	// 	log.Printf("EnvINI.Unmarshal() | Unhandled Kind: %#v\n", rv.Kind())
	// }

	return nil
}

func bytesToDataMap(data []byte) map[string]map[string]string {
	line := ""
	lastC := '\n'
	i := 0
	dataLength := len(data)

	dataMap := make(map[string]map[string]string)
	section := "GLOBAL"

	for i < dataLength {
		r, w := utf8.DecodeRune(data[i:])

		switch r {
		case '\n':
			if lastC != '\n' && len(line) > 0 {
				k, v, s := lineToSectionKV(line)
				if s {
					section = v
				} else {
					if dataMap[section] == nil {
						dataMap[section] = make(map[string]string)
					}
					dataMap[section][k] = v
				}
			}
			line = ""
		default:
			lastC = r
			line += string(r)
			if i+w == dataLength {
				k, v, s := lineToSectionKV(line)
				if s {
					section = v
				} else {
					if dataMap[section] == nil {
						dataMap[section] = make(map[string]string)
					}
					dataMap[section][k] = v
				}
			}
		}

		i += w
	}

	return dataMap
}

func lineToSectionKV(line string) (string, string, bool) {
	line = strings.TrimSpace(line)
	if line[0] == '[' {
		// Section
		section := strings.Trim(line, `[]`)
		section = strings.TrimSpace(section)
		// log.Printf("EnvINI.lineToSectionKV() | Section: %q\n", section)

		return "Section", section, true
	}

	// Key/Value
	kv := strings.SplitN(line, "=", 2)
	k := strings.TrimSpace(kv[0])
	v := strings.TrimSpace(kv[1])
	v = strings.Trim(v, `"`)
	// log.Printf("EnvINI.lineToSectionKV() | Key: %q | Value: %q\n", k, v)

	return k, v, false
}
