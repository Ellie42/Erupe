package stringsupport

import (
	"bytes"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

func MustConvertShiftJISToUTF8(text string) string {
	result, err := ConvertShiftJISToUTF8(text)

	if err != nil {
		panic(err)
	}

	return result
}

func MustConvertUTF8ToShiftJIS(text string) string {
	result, err := ConvertUTF8ToShiftJIS(text)

	if err != nil {
		panic(err)
	}

	return result
}

func ConvertShiftJISToUTF8(text string) (string, error) {
	r := bytes.NewBuffer([]byte(text))
	decoded, err := ioutil.ReadAll(transform.NewReader(r, japanese.ShiftJIS.NewDecoder()))

	if err != nil {
		return "", err
	}

	return string(decoded), nil
}

func ConvertUTF8ToShiftJIS(text string) (string, error) {
	r := bytes.NewBuffer([]byte(text))
	encoded, err := ioutil.ReadAll(transform.NewReader(r, japanese.ShiftJIS.NewEncoder()))

	if err != nil {
		return "", err
	}

	return string(encoded), nil
}
