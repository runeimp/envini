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

var dataMap map[string]map[string]string

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

func stringIsTruthy(s string) bool {
	s = strings.ToLower(s)
	switch s {
	case "true", "t", "yes", "y", "1":
		return true
	}
	return false
}

func Unmarshal(data []byte, config interface{}) (err error) {
	// log.Printf("EnvINI.Unmarshal() | data: %q | config: %#v\n", data, config)

	dataMap = bytesToDataMap(data)
	log.Printf("EnvINI.Unmarshal() | dataMap: %q\n", dataMap)

	// err = walkStruct("GLOBAL", config)
	err = walkerTexasRanger("GLOBAL", config)
	log.Printf("EnvINI.Unmarshal() | config: %#v\n", config)

	// rv := reflect.ValueOf(config) // reflect.Value
	// copy := reflect.New(rv.Type()).Elem()
	// err = walkObject("GLOBAL", 0, copy, rv)
	// if err != nil {
	// 	log.Printf("EnvINI.Unmarshal() | err: %s\n", err.Error())
	// } else {
	// 	log.Println("EnvINI.Unmarshal() | err: nil")
	// }
	// log.Printf("EnvINI.Unmarshal() | rv: %#v\n", rv)
	// log.Printf("EnvINI.Unmarshal() | copy: %#v\n", copy)

	return err
}

func walkerTexasRanger(section string, config interface{}) (err error) {
	log.Printf("EnvINI.walkerTexasRanger() | -Call-  | section: %q\n", section)
	reflectType := reflect.TypeOf(config).Elem()
	rv := reflect.ValueOf(config) // reflect.Value
	reflectValue := rv.Elem()

	/*
	 *
	 *
	 */

	// isPointer := rv.Kind() == reflect.Ptr
	// log.Printf("EnvINI.walkerTexasRanger() | rv.Kind(): %#v | isPointer: %t | rv.IsNil(): %t\n", rv.Kind(), isPointer, rv.IsNil())
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(config)}
	}

	log.Printf("EnvINI.walkerTexasRanger() | rv.Kind(): %s           | rv.IsValid(): %t\n", rv.Kind(), rv.IsValid())
	// log.Printf("EnvINI.walkerTexasRanger() | rv.Interface(): %#v\n", rv.Interface())
	log.Printf("EnvINI.walkerTexasRanger() | reflectType.Kind(): %s | reflectValue.Kind(): %s | reflectValue.IsValid(): %t\n", reflectType.Kind(), reflectValue.Kind(), reflectValue.IsValid())
	log.Printf("EnvINI.walkerTexasRanger() | reflectValue.Interface(): %#v\n", reflectValue.Interface())

	for i := 0; i < reflectType.NumField(); i++ {
		fieldName := reflectType.Field(i).Name     // Struct field name
		tag := reflectType.Field(i).Tag.Get("ini") // Go struct tag value for the "ini" key
		fieldType := reflectValue.Field(i).Type()  // Go value type
		// fieldInterface := reflectValue.Field(i).Interface() // Actual field value

		log.Printf("EnvINI.walkerTexasRanger() | fieldName: %-12s | %-7s | tag: %s\n", fieldName, fieldType.Kind(), tag)

		if fieldType.Kind() == reflect.Struct {
			v := reflectValue.FieldByName(fieldName)
			log.Printf("EnvINI.walkerTexasRanger() | Struct  | v.IsValid(): %t\n", v.IsValid())
			if v.IsValid() {
				// vv := v.Interface()
				// err = walkerTexasRanger(fieldName, &vv)
				// err = walkerTexasRanger(fieldName, &fieldInterface)
				// err = walkerTexasRanger(fieldName, &reflectValue)
				// fieldValue := reflectValue.Field(i) // Actual field value?
				// err = walkerTexasRanger(fieldName, &fieldValue)
				err = walkerTexasRanger(fieldName, v.Addr().Interface())
			}
			section = tag
			continue
		}

		if dataValue, ok := dataMap[section][tag]; ok || fieldName == "" && len(tag) > 0 {
			// log.Printf("EnvINI.walkerTexasRanger() | field.Name: %q\n", field.Name)
			v := reflectValue.FieldByName(fieldName)
			// log.Printf("EnvINI.walkerTexasRanger() | fieldName: %-12s | fieldInterface: %#-9v | dataValue: %#-9v | %-7s | v.IsValid(): %t\n", fieldName, fieldInterface, dataValue, fieldType.Kind(), v.IsValid())
			// log.Printf("EnvINI.walkerTexasRanger() | fieldName: %-12s | dataValue: %#-9v | %-7s | v.IsValid(): %t\n", fieldName, dataValue, fieldType.Kind(), v.IsValid())

			switch fieldType.Kind() {
			case reflect.Bool:
				if v.IsValid() {
					v.SetBool(stringIsTruthy(dataValue))
				}
			// case reflect.Float:
			// 	if v.IsValid() {
			// 		floatValue, err := strconv.ParseFloat(dataValue, 64)
			// 		if err == nil {
			// 			v.SetFloat(floatValue)
			// 		}
			// 	}
			case reflect.Float32:
				if v.IsValid() {
					floatValue, err := strconv.ParseFloat(dataValue, 64)
					if err == nil {
						v.SetFloat(floatValue)
					} else {
						log.Printf("EnvINI.walkerTexasRanger() | Float32 | err: %s\n", err.Error())
					}
				}
			case reflect.Float64:
				if v.IsValid() {
					floatValue, err := strconv.ParseFloat(dataValue, 64)
					if err == nil {
						v.SetFloat(floatValue)
					} else {
						log.Printf("EnvINI.walkerTexasRanger() | Float64 | err: %s\n", err.Error())
					}
				}
			case reflect.Int:
				if v.IsValid() {
					intValue, err := strconv.ParseInt(dataValue, 10, 64)
					if err == nil {
						v.SetInt(intValue)
					}
				}
			case reflect.Uint8:
				if v.IsValid() {
					intValue, err := strconv.ParseUint(dataValue, 10, 8)
					if err == nil {
						v.SetUint(intValue)
					}
				}
			case reflect.String:
				if v.IsValid() {
					v.SetString(dataValue)
				}
			case reflect.Ptr:
				log.Printf("EnvINI.walkerTexasRanger() | Pointer | v.IsValid(): %t\n", v.IsValid())
				if v.IsValid() {
					// v.SetString(dataValue)
					// reflectValue
					t := reflect.TypeOf(v).Elem() // Dereference the interface{} to a reflect.Value
					err = walkerTexasRanger(fieldName, &t)
				}
			// case reflect.Struct:
			// 	log.Printf("EnvINI.walkerTexasRanger() | Struct  | v.IsValid(): %t\n", v.IsValid())
			// 	// if v.IsValid() {
			// 	// 	// v.SetString(dataValue)
			// 	// 	t := reflect.TypeOf(v).Elem() // Dereference the interface{} to a reflect.Value
			// 	// 	err = walkerTexasRanger(fieldName, &t)
			// 	// }
			// 	// v := valueValue := reflectValue.Field(i).Addr()
			// 	if v.IsValid() {
			// 		err = walkerTexasRanger(fieldName, v.Addr().Interface())
			// 	}
			default: // struct
				log.Printf("EnvINI.walkerTexasRanger() | DEFAULT | v.IsValid(): %t\n", v.IsValid())
				if v.IsValid() {
					// rv := reflect.ValueOf(v).Elem()
					// vv := v.Interface()
					// vv := v.Elem()
					// err = walkerTexasRanger(fieldName, &vv)
				}
			}
		}

		// switch reflectValue.Field(i).Kind() {
		// case reflect.String:
		// 	log.Printf("EnvINI.walkerTexasRanger() | String  | %s: %q (%s) [%s]\n", fieldName, fieldInterface, fieldType, tag)
		// case reflect.Int32:
		// 	log.Printf("EnvINI.walkerTexasRanger() | Int32   | %s: %i (%s) [%s]\n", fieldName, fieldInterface, fieldType, tag)
		// case reflect.Struct:
		// 	log.Printf("EnvINI.walkerTexasRanger() | Struct  | %q is %s [%s]\n", fieldName, fieldType, tag)
		// 	walkerTexasRanger(fieldName, reflectValue.Field(i).Addr().Interface())
		// default:
		// 	log.Printf("EnvINI.walkerTexasRanger() | Default | %s: %v (%s) [%s]\n", fieldName, fieldInterface, fieldType, tag)
		// }
	}
	return err
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
