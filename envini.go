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

// var dataMap map[string]map[string]string
var (
	dataMap              *inidata.DataMap
	ErrorValueUnsettable = errors.New("reflect.Value unsettable")
	ErrorValueInvalid    = errors.New("reflect.Value invalid")
)

func configWalker(section string, config interface{}) (err error) {
	log.Printf("EnvINI.configWalker()  | -Call-  | section: %q\n", section)
	reflectType := reflect.TypeOf(config).Elem()
	rv := reflect.ValueOf(config) // reflect.Value
	reflectValue := rv.Elem()

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(config)}
	}

	log.Printf("EnvINI.configWalker()  | rv.Kind(): %s           | rv.IsValid(): %t\n", rv.Kind(), rv.IsValid())
	// log.Printf("EnvINI.configWalker()  | rv.Interface(): %#v\n", rv.Interface())
	log.Printf("EnvINI.configWalker()  | reflectType.Kind(): %s | reflectValue.Kind(): %s | reflectValue.IsValid(): %t\n", reflectType.Kind(), reflectValue.Kind(), reflectValue.IsValid())
	log.Printf("EnvINI.configWalker()  | reflectValue.Interface(): %#v\n", reflectValue.Interface())

	for i := 0; i < reflectType.NumField(); i++ {
		fieldName := reflectType.Field(i).Name                // Struct field name
		tagEnv := reflectType.Field(i).Tag.Get("env")         // Go struct tagINI value for the "ini" key
		tagINI := reflectType.Field(i).Tag.Get("ini")         // Go struct tagINI value for the "ini" key
		tagDefault := reflectType.Field(i).Tag.Get("default") // Go struct tagINI value for the "default" key
		fieldType := reflectValue.Field(i).Type()             // Go value type
		// fieldInterface := reflectValue.Field(i).Interface() // Actual field value

		log.Printf("EnvINI.configWalker()  | fieldName: %-13q | Kind: %-7s | tagEnv: %-14q | tagINI: %q | tagDefault: %q\n", fieldName, fieldType.Kind(), tagEnv, tagINI, tagDefault)

		if fieldType.Kind() == reflect.Struct {
			v := reflectValue.FieldByName(fieldName)
			log.Printf("EnvINI.configWalker()  | Struct  | v.IsValid(): %t\n", v.IsValid())
			if v.IsValid() {
				err = configWalker(fieldName, v.Addr().Interface())
			}
			section = tagINI
			continue
		}

		if len(fieldName) == 0 {
			panic(errors.New("fieldName is zero length"))
		}

		dataValue, ok := dataMap.GetKey(tagINI)
		if ok == false && len(tagDefault) > 0 {
			v := reflectValue.FieldByName(fieldName)
			err = setFieldValue(v, fieldName, tagEnv, tagDefault)
		} else if ok || fieldName == "" && len(tagINI) > 0 {
			log.Printf("EnvINI.configWalker()  | fieldName: %-13q | tagINI: %q\n", fieldName, tagINI)
			v := reflectValue.FieldByName(fieldName)
			// log.Printf("EnvINI.configWalker() | fieldName: %-13q | fieldInterface: %#-9v | dataValue: %#-9v | %-7s | v.IsValid(): %t\n", fieldName, fieldInterface, dataValue, fieldType.Kind(), v.IsValid())
			// log.Printf("EnvINI.configWalker() | fieldName: %-12s | dataValue: %#-9v | %-7s | v.IsValid(): %t\n", fieldName, dataValue, fieldType.Kind(), v.IsValid())
			err = setFieldValue(v, fieldName, tagEnv, dataValue)
		}

		// Check error and end look if there is a problem
		if err != nil {
			return err
		}
	}
	log.Printf("EnvINI.configWalker()  | err: %v\n", err)
	return err
}

func GetConfig(configPath string, config interface{}) error {
	log.Printf("EnvINI.GetConfig() | configPath: %q | config: %v\n", configPath, config)

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := Unmarshal(data, config); err != nil {
		return err
	}

	return nil
}

func stringIsTruthy(s string) bool {
	s = strings.ToLower(s)
	switch s {
	case "true", "t", "yes", "y", "1":
		return true
	}
	return false
}

func reflectValueCheck(v reflect.Value) error {
	if v.IsValid() == false {
		return ErrorValueInvalid
	}
	if v.CanSet() == false {
		return ErrorValueUnsettable
	}
	return nil
}

func setFieldValue(v reflect.Value, fieldName, tagEnv, dataValue string) (err error) {
	log.Printf("EnvINI.setFieldValue() | fieldName: %-13q | tagEnv: %-14q | dataValue: %q\n", fieldName, tagEnv, dataValue)
	if len(tagEnv) > 0 {
		env := os.Getenv(tagEnv)
		if len(env) > 0 {
			dataValue = env
		}
		log.Printf("EnvINI.setFieldValue() | fieldName: %-13q | env: %q    | dataValue: %q\n", fieldName, env, dataValue)
	}

	err = reflectValueCheck(v)
	if err != nil {
		return err
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
		intValue, err := strconv.ParseUint(dataValue, 10, 64)
		if err == nil {
			return err
		}
		v.SetUint(intValue)
	case reflect.String:
		v.SetString(dataValue)
	case reflect.Ptr:
		log.Printf("EnvINI.setFieldValue() | Pointer | v.IsValid(): %t\n", v.IsValid())
		// v.SetString(dataValue)
		// reflectValue
		t := reflect.TypeOf(v).Elem() // Dereference the interface{} to a reflect.Value
		err = configWalker(fieldName, &t)
	default: // array, slice, struct?
		// rv := reflect.ValueOf(v).Elem()
		// vv := v.Interface()
		// vv := v.Elem()
		// err = configWalker(fieldName, &vv)
	}

	return err
}

func Unmarshal(data []byte, config interface{}) (err error) {
	// log.Printf("EnvINI.Unmarshal() | data: %q | config: %#v\n", data, config)

	// dataMap = bytesToDataMap(data)
	dataMap = inidata.NewDataMap()
	err = dataMap.ParseBytes(data)
	if err != nil {
		log.Printf("EnvINI.Unmarshal() | err: %s\n", err.Error())
		return err
	}

	log.Printf("EnvINI.Unmarshal() | dataMap: %s\n", dataMap)

	err = configWalker("GLOBAL", config)
	log.Printf("EnvINI.Unmarshal() | configWalker() | err: %q\n", err)
	if err == nil {
		log.Printf("EnvINI.Unmarshal() | config: %#v\n", config)
	}

	return err
}
