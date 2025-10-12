package responses

import (
	"encoding/json"
	"net/http"
	"reflect"
)

type WriteableResponse interface {
	Write(w http.ResponseWriter)
}

type BaseResponse struct{}

func (BaseResponse) Write(w http.ResponseWriter, self any) {
	// NOTE: maybe include other? in case Accept header is different
	w.Header().Set("Content-Type", "application/json")

	v := reflect.ValueOf(self)
	if v.Kind() == reflect.Struct {
		statusField := v.FieldByName("Status")
		if statusField.IsValid() && statusField.Kind() == reflect.Int {
			if status := int(statusField.Int()); status != 0 {
				w.WriteHeader(status)
			}
		}
	}

	json.NewEncoder(w).Encode(self)
}
