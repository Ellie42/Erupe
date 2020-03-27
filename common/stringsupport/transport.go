package stringsupport

import (
	"fmt"
	"reflect"
)

// Iterate a struct and convert all strings from UTF8 to null-terminated Shift-JIS
func PrepareStringsForTransport(obj interface{}) error {
	reflected := reflect.ValueOf(obj)

	element := reflected.Elem()

	for i := 0; i < element.NumField(); i++ {
		field := element.Field(i)

		if field.Kind() != reflect.String || !field.CanSet() {
			continue
		}

		text, err := ConvertUTF8ToShiftJIS(field.String())

		if err != nil {
			return err
		}

		field.SetString(fmt.Sprintf("%s\x00", text))
	}

	return nil
}

func ConvertUTF8ToNullTerminatedShiftJIS(str string) (string, error) {
	text, err := ConvertUTF8ToShiftJIS(str)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\x00", text), nil
}
