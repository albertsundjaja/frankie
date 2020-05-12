package isGood

import (
	"encoding/json"
	"github.com/albertsundjaja/frankie/models"
	"net/http"
	"strconv"
	"github.com/albertsundjaja/frankie/helpers"
)

var uniqueSessionKeys []string

func init() {
	uniqueSessionKeys = make([]string, 0)
}

func IsGoodHandler(w http.ResponseWriter, r *http.Request) {

	// read the request
	decoder := json.NewDecoder(r.Body)
	var body models.DeviceCheckDetailsObjectCollection
	err:= decoder.Decode(&body)
	if err != nil {
		errorResponse := helpers.HttpErrorResponder(http.StatusInternalServerError, "JSON Payload structure error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorResponse)
		return
	}

	// validate all the DeviceCheckDetailsObject
	for _, item := range body {
		checkTypeOk := validateCheckType(item.CheckType)
		if !checkTypeOk {
			errorResponse := helpers.HttpErrorResponder(http.StatusInternalServerError, "checkType error")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(errorResponse)
			return
		}

		activityTypeOk := validateActivityType(item.ActivityType)
		if !activityTypeOk {
			errorResponse := helpers.HttpErrorResponder(http.StatusInternalServerError, "activityType error")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(errorResponse)
			return
		}

		sessionKeyOk := validateSessionKey(item.CheckSessionKey)
		if !sessionKeyOk {
			errorResponse := helpers.HttpErrorResponder(http.StatusInternalServerError, "sessionKey error")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(errorResponse)
			return
		}

		// must not contain duplicate activityData KvpKey
		uniqueKvpKeys := make([]string, 0)
		for _, activityData := range(item.ActivityData) {
			if helpers.StringInSlice(activityData.KvpKey, uniqueKvpKeys) {
				errorResponse := helpers.HttpErrorResponder(http.StatusInternalServerError, "kvpKey in each object must be unique")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(errorResponse)
				return
			}
			kvpTypeOk := validateKvpType(activityData.KvpType)
			if !kvpTypeOk{
				errorResponse := helpers.HttpErrorResponder(http.StatusInternalServerError, "kvpType error")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(errorResponse)
				return
			}
			kvpValueOk := validateKvpValue(activityData.KvpValue, activityData.KvpType)
			if !kvpValueOk{
				errorResponse := helpers.HttpErrorResponder(http.StatusInternalServerError, "kvpValue is of different data type than claimed")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(errorResponse)
				return
			}
			uniqueKvpKeys = append(uniqueKvpKeys, activityData.KvpKey)
		}
		// ensure that each session keys only appear once
		uniqueSessionKeys = append(uniqueSessionKeys, item.CheckSessionKey)
	}

	// everything is ok response with lovely puppy
	resp := models.PuppyObject{Puppy: true}
	jsonResp, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
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
	return !helpers.StringInSlice(sessionKey, uniqueSessionKeys)
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