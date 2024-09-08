package utilities

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func SetString(value string) *string {
	return &value
}

func SetBool(value bool) *bool {
	return &value
}

func StringIntArray(str string, arr []string) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == str {
			return true
		}
	}

	return false
}

func StringToFloat64(value string) float64 {
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Println("err", err)
		return float64(0)
	}

	return result
}

func WriteJSON(w http.ResponseWriter, code int, obj interface{}) error {
	w.WriteHeader(code)
	WriteContentType(w, []string{"application/json; charset=utf-8"})

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	w.Write(jsonBytes)
	return nil
}

func WriteContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
