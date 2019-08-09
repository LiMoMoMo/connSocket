package models

import (
	"fmt"
	"reflect"
)

func init() {
	// init global maps
	reportMap = make(map[ReportType]getObj)
	commandMap = make(map[CommandType]getObj)
	// register Type_Login
	RegisterReport(Type_Register, func() BaseContent { return &Register{} })
	RegisterCommand(Command_Start, func() BaseContent { return &Start{} })
}

func RegisterReport(t ReportType, f getObj) {
	reportMap[t] = f
}

func RegisterCommand(t CommandType, f getObj) {
	commandMap[t] = f
}

// BaseContent as interfae of Content
type BaseContent interface {
	Fill(m map[string]interface{})
}

type getObj func() BaseContent

var reportMap map[ReportType]getObj
var commandMap map[CommandType]getObj

// SetField set value to field
func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	fieldVal := structValue.FieldByName(name)

	if !fieldVal.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !fieldVal.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	val := reflect.ValueOf(value)

	if fieldVal.Type() != val.Type() {

		if m, ok := value.(map[string]interface{}); ok {

			// if field value is struct
			if fieldVal.Kind() == reflect.Struct {
				return FillStruct(m, fieldVal.Addr().Interface())
			}

			// if field value is a pointer to struct
			if fieldVal.Kind() == reflect.Ptr && fieldVal.Type().Elem().Kind() == reflect.Struct {
				if fieldVal.IsNil() {
					fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
				}
				// fmt.Printf("recursive: %v %v\n", m,fieldVal.Interface())
				return FillStruct(m, fieldVal.Interface())
			}

		}

		return fmt.Errorf("Provided value type didn't match obj field type")
	}

	fieldVal.Set(val)
	return nil

}

// FillStruct fill struct as right values
func FillStruct(m map[string]interface{}, s interface{}) error {
	for k, v := range m {
		err := SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
