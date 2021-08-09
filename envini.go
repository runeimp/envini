package envini

import (
	"errors"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/runeimp/envini/inidata"
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

var (
	dataMap              *inidata.DataMap
	dataPath             string
	ErrorValueUnsettable = errors.New("reflect.Value unsettable")
	ErrorValueInvalid    = errors.New("reflect.Value invalid")
)

func configWalker(section string, config interface{}) (err error) {
	// log.Printf("EnvINI.configWalker()  | -Call-  | section: %q\n", section)
	rt := reflect.TypeOf(config).Elem()
	rv := reflect.ValueOf(config) // reflect.Value
	rvE := rv.Elem()

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(config)}
	}

	// log.Printf("EnvINI.configWalker()  | rv.Kind(): %-16s | rv.IsValid(): %t\n", rv.Kind(), rv.IsValid())
	// log.Printf("EnvINI.configWalker()  | rv.Interface(): %#v\n", rv.Interface())
	// log.Printf("EnvINI.configWalker()  | rt.Kind(): %-16s | rvE.Kind(): %-11s | rvE.IsValid(): %t\n", rt.Kind(), rvE.Kind(), rvE.IsValid())
	// log.Printf("EnvINI.configWalker()  | rvE.Interface(): %#v\n", rvE.Interface())

	for i := 0; i < rt.NumField(); i++ {
		fieldName := rt.Field(i).Name                // Struct field name
		tagEnv := rt.Field(i).Tag.Get("env")         // Go struct tagINI value for the "ini" key
		tagINI := rt.Field(i).Tag.Get("ini")         // Go struct tagINI value for the "ini" key
		tagDefault := rt.Field(i).Tag.Get("default") // Go struct tagINI value for the "default" key
		fieldType := rvE.Field(i).Type()             // Go value type
		// fieldInterface := rvE.Field(i).Interface() // Actual field value

		// log.Printf("EnvINI.configWalker()  | fieldName: %-16q | Kind: %-17s | tagEnv: %-14q | tagINI: %q | tagDefault: %q\n", fieldName, fieldType.Kind(), tagEnv, tagINI, tagDefault)

		if fieldType.Kind() == reflect.Struct {
			v := rvE.FieldByName(fieldName)
			// log.Printf("EnvINI.configWalker()  | Struct  | v.IsValid(): %t\n", v.IsValid())
			if v.IsValid() {
				err = configWalker(tagINI, v.Addr().Interface())
			}
			section = tagINI
			continue
		}

		if len(fieldName) == 0 {
			panic(errors.New("fieldName is zero length"))
		}

		dataValue, ok := dataMap.GetKey(tagINI, section)
		// log.Printf("EnvINI.configWalker()  | key: %-22q | section: %-14q | ok: %-19t | dataValue: %q\n", tagINI, section, ok, dataValue)
		if ok == false && len(tagDefault) > 0 {
			v := rvE.FieldByName(fieldName)
			err = setFieldValue(v, fieldName, tagEnv, tagDefault)
		} else if ok || fieldName == "" && len(tagINI) > 0 {
			// log.Printf("EnvINI.configWalker()  | fieldName: %-16q | tagINI: %q\n", fieldName, tagINI)
			v := rvE.FieldByName(fieldName)
			err = setFieldValue(v, fieldName, tagEnv, dataValue)
		}

		// Check error and end look if there is a problem
		if err != nil {
			return err
		}
	}

	if err != nil {
		log.Printf("EnvINI.configWalker()  | err: %v\n", err)
	}
	return err
}

// GetConfig takes a string reference to an config file (INI file) and a pointer to a struct and populates the struct with data from the config file
func GetConfig(configPath string, config interface{}) error {
	// log.Printf("EnvINI.GetConfig() | configPath: %q | config: %v\n", configPath, config)

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	dataPath = configPath

	if err := Unmarshal(data, config); err != nil {
		return err
	}

	return nil
}

func GetConfigJSON(configPath string) (jsonStr string, err error) {

	if len(dataPath) == 0 {
		err = errors.New("you must call envini.GetConfig() first")
	}
	if err == nil {
		jsonStr = dataMap.String()
	}

	return jsonStr, err
}

func stringIsTruthy(s string) bool {
	s = strings.ToLower(s)
	switch s {
	case "true", "t", "yes", "y", "1":
		return true
	}
	return false
}

func setFieldValue(v reflect.Value, fieldName, envName, dataValue string) (err error) {
	// log.Printf("EnvINI.setFieldValue() | fieldName: %-16q | envName: %-14q | dataValue: %q\n", fieldName, envName, dataValue)
	if len(envName) > 0 {
		env := os.Getenv(envName)
		if len(env) > 0 {
			dataValue = env
		}
		// log.Printf("EnvINI.setFieldValue() | fieldName: %-16q | env: %-15q    | dataValue: %q\n", fieldName, env, dataValue)
	}

	if v.IsValid() == false {
		return ErrorValueInvalid
	}
	if v.CanSet() == false {
		return ErrorValueUnsettable
	}

	switch v.Type().Kind() {
	case reflect.Bool:
		v.SetBool(stringIsTruthy(dataValue))
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(dataValue, 64)
		if err != nil {
			return err
		}
		v.SetFloat(floatValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(dataValue, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(intValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(dataValue, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(uintValue)
	case reflect.String:
		v.SetString(dataValue)
	case reflect.Ptr:
		// log.Printf("EnvINI.setFieldValue() | Pointer | v.IsValid(): %t\n", v.IsValid())
		// v.SetString(dataValue)
		// reflectValue
		t := reflect.TypeOf(v).Elem() // Dereference the interface{} to a reflect.Value
		err = configWalker(fieldName, &t)
	default: // array, slice?
		// rv := reflect.ValueOf(v).Elem()
		// vv := v.Interface()
		// vv := v.Elem()
		// err = configWalker(fieldName, &vv)
	}

	return err
}

// Unmarshal takes byte slice (ostensibly raw data from an INI file) and a pointer to a struct and populates the struct with the data from the byte slice
func Unmarshal(data []byte, config interface{}) (err error) {
	// log.Printf("EnvINI.Unmarshal() | data: %q | config: %#v\n", data, config)

	dataMap = inidata.NewDataMap()
	err = dataMap.ParseBytes(data)
	if err != nil {
		return err
	}

	// log.Printf("EnvINI.Unmarshal() | dataMap: %s\n", dataMap)

	err = configWalker("GLOBAL", config)

	return err
}
