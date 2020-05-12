package helpers

import (
	"encoding/json"
	"github.com/albertsundjaja/frankie/models"
)

// check if a string is already inside the list
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// return marshaled ErrorObject with error message and status
func HttpErrorResponder(status int, message string) []byte {
	var httpError models.ErrorObject
	httpError.Code = status
	httpError.Message = message
	errorResp, _ := json.Marshal(httpError)
	return errorResp
}