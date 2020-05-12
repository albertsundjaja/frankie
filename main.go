package main

import (
	"encoding/json"
	"fmt"
	"github.com/albertsundjaja/frankie/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"strconv"
)

var uniqueSessionKeys []string

func init() {
	uniqueSessionKeys = make([]string, 0)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/isgood", isGoodHandler)
	http.ListenAndServe(":8888", r)
}

func isGoodHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var body models.DeviceCheckDetailsObjectCollection
	var httpError models.ErrorObject

	err:= decoder.Decode(&body)
	if err != nil {
		httpError.Code = 500
		httpError.Message = "JSON Payload structure error"
		errorResp, _ := json.Marshal(httpError)
		w.WriteHeader(500)
		w.Write(errorResp)
		return
	}

	// validate all the DeviceCheckDetailsObject
	for _, item := range body {
		checkTypeOk := validateCheckType(item.CheckType)
		if !checkTypeOk {
			httpError.Code = 500
			httpError.Message = "CheckType error"
			errorResp, _ := json.Marshal(httpError)
			w.WriteHeader(500)
			w.Write(errorResp)
			return
		}

		activityTypeOk := validateActivityType(item.ActivityType)
		if !activityTypeOk {
			httpError.Code = 500
			httpError.Message = "ActivityType error"
			errorResp, _ := json.Marshal(httpError)
			w.WriteHeader(500)
			w.Write(errorResp)
			return
		}

		sessionKeyOk := validateSessionKey(item.CheckSessionKey)
		if !sessionKeyOk {
			httpError.Code = 500
			httpError.Message = "SessionKey error"
			errorResp, _ := json.Marshal(httpError)
			w.WriteHeader(500)
			w.Write(errorResp)
			return
		}
		for _, activityData := range(item.ActivityData) {
			uniqueKvpKeys := make([]string, 0)
			if stringInSlice(activityData.KvpKey, uniqueKvpKeys) {
				httpError.Code = 500
				httpError.Message = "KvpKey in each object must be unique"
				errorResp, _ := json.Marshal(httpError)
				w.WriteHeader(500)
				w.Write(errorResp)
				return
			}
			kvpTypeOk := validateKvpType(activityData.KvpType)
			if !kvpTypeOk{
				httpError.Code = 500
				httpError.Message = "KvpType error"
				errorResp, _ := json.Marshal(httpError)
				w.WriteHeader(500)
				w.Write(errorResp)
				return
			}
			kvpValueOk := validateKvpValue(activityData.KvpValue, activityData.KvpType)
			if !kvpValueOk{
				httpError.Code = 500
				httpError.Message = "KvpValue is of different type than claimed"
				errorResp, _ := json.Marshal(httpError)
				w.WriteHeader(500)
				w.Write(errorResp)
				return
			}
			uniqueKvpKeys = append(uniqueKvpKeys, activityData.KvpKey)
		}
		// ensure that each session keys only appear once
		uniqueSessionKeys = append(uniqueSessionKeys, item.CheckSessionKey)
	}

	resp := models.PuppyObject{Puppy: true}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(jsonResp)
}

// check if checkType is of valid enum values
func validateCheckType(checkType models.EnumCheckType) bool {
	switch checkType {
		case models.CheckTypeDevice, models.CheckTypeBiometric, models.CheckTypeCombo:
			return true
	}
	return false
}

// check if activityType is of valid enum values
func validateActivityType(activityType models.EnumActivityType) bool {
	switch activityType {
		case models.ActivityTypeLogin, models.ActivityTypeSignup, models.ActivityTypeConfirmation, models.ActivityTypePayment:
			return true
	}
	return false
}

// check if sessionKey already used
func validateSessionKey(sessionKey string) bool {
	return !stringInSlice(sessionKey, uniqueSessionKeys)
}

// check if kvpType is of valid enum values
func validateKvpType(kvpType models.EnumKVPType) bool {
	switch kvpType {
		case models.EnumKVPTypeGeneralBool, models.EnumKVPTypeGeneralFloat, models.EnumKVPTypeGeneralInteger, models.EnumKVPTypeGeneralString:
			return true
	}
	return false
}

// check if the kvpValue can be parsed into claimedType data type
func validateKvpValue(kvpValue string, claimedType models.EnumKVPType) bool {
	switch claimedType {
		case models.EnumKVPTypeGeneralString:
			return true
		case models.EnumKVPTypeGeneralInteger:
			_, err := strconv.ParseInt(kvpValue, 10, 64)
			if err != nil {
				return false
			}
		case models.EnumKVPTypeGeneralFloat:
			_, err := strconv.ParseFloat(kvpValue, 64)
			if err != nil {
				return false
			}
		case models.EnumKVPTypeGeneralBool:
			_, err := strconv.ParseBool(kvpValue)
			if err != nil {
				return false
			}
	}
	return true
}

// check if a string is already inside the list
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}