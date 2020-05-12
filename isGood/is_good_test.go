package isGood_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/albertsundjaja/frankie/models"
	"github.com/go-chi/chi"
	. "github.com/onsi/ginkgo"
	"github.com/albertsundjaja/frankie/isGood"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"github.com/google/go-cmp/cmp"
)

var _ = Describe("IsGood", func() {

	Describe("checkType validation", func() {
		It("should validate if checkType is in enums available values", func() {
			checkOk := isGood.ValidateCheckType(models.CheckTypeDevice)
			Expect(checkOk).To(BeTrue())
		})

		It("should validate if checkType is in enums available values", func() {
			checkOk := isGood.ValidateCheckType(models.CheckTypeBiometric)
			Expect(checkOk).To(BeTrue())
		})

		It("should validate if checkType is in enums available values", func() {
			checkOk := isGood.ValidateCheckType(models.CheckTypeCombo)
			Expect(checkOk).To(BeTrue())
		})

		It("should INVALIDATE if checkType is not in enums available values", func() {
			checkOk := isGood.ValidateCheckType("valueThatIsNotThere")
			Expect(checkOk).To(BeFalse())
		})
	})

	Describe("activityType validation", func() {
		It("should validate if activityType is in enums available values", func() {
			checkOk := isGood.ValidateActivityType(models.ActivityTypeConfirmation)
			Expect(checkOk).To(BeTrue())
		})

		It("should INVALIDATE if activityType is not in enums available values", func() {
			checkOk := isGood.ValidateCheckType("valueThatIsNotThere")
			Expect(checkOk).To(BeFalse())
		})
	})

	Describe("kvpType validation", func() {
		It("should validate if kvpType is in enums available values", func() {
			checkOk := isGood.ValidateKvpType(models.EnumKVPTypeGeneralBool)
			Expect(checkOk).To(BeTrue())
		})

		It("should INVALIDATE if kvpType is not in enums available values", func() {
			checkOk := isGood.ValidateKvpType("valueThatIsNotThere")
			Expect(checkOk).To(BeFalse())
		})
	})

	Describe("kvpValue validation", func() {
		It("should validate if kvpValue can be parsed to float", func() {
			checkOk := isGood.ValidateKvpValue("10.2312", models.EnumKVPTypeGeneralFloat)
			Expect(checkOk).To(BeTrue())
		})

		It("should validate if kvpValue can be parsed to integer", func() {
			checkOk := isGood.ValidateKvpValue("10", models.EnumKVPTypeGeneralInteger)
			Expect(checkOk).To(BeTrue())
		})

		It("should validate if kvpValue can be parsed to bool", func() {
			checkOk := isGood.ValidateKvpValue("true", models.EnumKVPTypeGeneralBool)
			Expect(checkOk).To(BeTrue())
		})

		It("should validate if kvpValue can be parsed to string", func() {
			checkOk := isGood.ValidateKvpValue("some string", models.EnumKVPTypeGeneralString)
			Expect(checkOk).To(BeTrue())
		})

		It("should INVALIDATE if kvpValue cannot be parsed to float", func() {
			checkOk := isGood.ValidateKvpValue("some string", models.EnumKVPTypeGeneralFloat)
			Expect(checkOk).To(BeFalse())
		})

		It("should INVALIDATE if kvpValue cannot be parsed to integer", func() {
			checkOk := isGood.ValidateKvpValue("some string", models.EnumKVPTypeGeneralInteger)
			Expect(checkOk).To(BeFalse())
		})

		It("should INVALIDATE if kvpValue cannot be parsed to bool", func() {
			checkOk := isGood.ValidateKvpValue("some string", models.EnumKVPTypeGeneralBool)
			Expect(checkOk).To(BeFalse())
		})

	})

	Describe("sessionKey validation", func() {
		It("should validate if sessionKey not yet used", func() {
			checkOk := isGood.ValidateSessionKey("sessionkey")
			Expect(checkOk).To(BeTrue())
		})

		It("should INVALIDATE if sessionKey is already stored in slice", func() {
			checkOk := isGood.ValidateKvpType("sessionkey")
			Expect(checkOk).To(BeFalse())
		})
	})

	Describe("testing IsGoodHandler", func() {
		It("should response a lovely puppy if everything is ok", func() {
			router := chi.NewRouter()
			router.Post("/isgood", isGood.IsGoodHandler)

			activityData_1 := models.KeyValuePairObject{
				KvpKey:   "ip.address",
				KvpType:  models.EnumKVPTypeGeneralBool,
				KvpValue: "true",
			}
			deviceCheckDetail_1 := models.DeviceCheckDetailsObject{
				ActivityData:    []*models.KeyValuePairObject{&activityData_1},
				ActivityType:    models.ActivityTypeConfirmation,
				CheckSessionKey: "sessionkey2",
				CheckType:       models.CheckTypeDevice,
			}

			jsonRequest, err := json.Marshal(models.DeviceCheckDetailsObjectCollection{&deviceCheckDetail_1})
			Expect(err).To(BeNil())

			req, err := http.NewRequest("POST", "/isgood", bytes.NewReader(jsonRequest))
			Expect(err).To(BeNil())

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusOK))

			expectedResponse := models.PuppyObject{Puppy: true}
			decoder := json.NewDecoder(rr.Body)
			var body models.PuppyObject
			err = decoder.Decode(&body)
			if err != nil {
				fmt.Println(err)
			}
			Expect(err).To(BeNil())
			Expect(cmp.Equal(body, expectedResponse)).To(BeTrue())
		})

		It("should response with 500 if KvpValue can't be parsed", func() {
			router := chi.NewRouter()
			router.Post("/isgood", isGood.IsGoodHandler)

			activityData_1 := models.KeyValuePairObject{
				KvpKey:   "ip.address",
				KvpType:  models.EnumKVPTypeGeneralBool,
				KvpValue: "thisisnotabool",
			}
			deviceCheckDetail_1 := models.DeviceCheckDetailsObject{
				ActivityData:    []*models.KeyValuePairObject{&activityData_1},
				ActivityType:    models.ActivityTypeConfirmation,
				CheckSessionKey: "sessionkey3",
				CheckType:       models.CheckTypeDevice,
			}

			jsonRequest, err := json.Marshal(models.DeviceCheckDetailsObjectCollection{&deviceCheckDetail_1})
			Expect(err).To(BeNil())

			req, err := http.NewRequest("POST", "/isgood", bytes.NewReader(jsonRequest))
			Expect(err).To(BeNil())

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should response with 500 if KvpType is not in enums list", func() {
			router := chi.NewRouter()
			router.Post("/isgood", isGood.IsGoodHandler)

			activityData_1 := models.KeyValuePairObject{
				KvpKey:   "ip.address",
				KvpType:  "somethingthatisnotinthelist",
				KvpValue: "thisisnotabool",
			}
			deviceCheckDetail_1 := models.DeviceCheckDetailsObject{
				ActivityData:    []*models.KeyValuePairObject{&activityData_1},
				ActivityType:    models.ActivityTypeConfirmation,
				CheckSessionKey: "sessionkey4",
				CheckType:       models.CheckTypeDevice,
			}

			jsonRequest, err := json.Marshal(models.DeviceCheckDetailsObjectCollection{&deviceCheckDetail_1})
			Expect(err).To(BeNil())

			req, err := http.NewRequest("POST", "/isgood", bytes.NewReader(jsonRequest))
			Expect(err).To(BeNil())

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should response with 500 if KvpKey is duplicated", func() {
			router := chi.NewRouter()
			router.Post("/isgood", isGood.IsGoodHandler)

			activityData_1 := models.KeyValuePairObject{
				KvpKey:   "ip.address",
				KvpType:  models.EnumKVPTypeGeneralBool,
				KvpValue: "true",
			}
			activityData_2 := models.KeyValuePairObject{
				KvpKey:   "ip.address",
				KvpType:  models.EnumKVPTypeGeneralBool,
				KvpValue: "true",
			}
			deviceCheckDetail_1 := models.DeviceCheckDetailsObject{
				ActivityData:    []*models.KeyValuePairObject{&activityData_1, &activityData_2},
				ActivityType:    models.ActivityTypeConfirmation,
				CheckSessionKey: "sessionkey5",
				CheckType:       models.CheckTypeDevice,
			}

			jsonRequest, err := json.Marshal(models.DeviceCheckDetailsObjectCollection{&deviceCheckDetail_1})
			Expect(err).To(BeNil())

			req, err := http.NewRequest("POST", "/isgood", bytes.NewReader(jsonRequest))
			Expect(err).To(BeNil())

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should response with 500 if checkType is not valid", func() {
			router := chi.NewRouter()
			router.Post("/isgood", isGood.IsGoodHandler)

			activityData_1 := models.KeyValuePairObject{
				KvpKey:   "ip.address",
				KvpType:  models.EnumKVPTypeGeneralBool,
				KvpValue: "true",
			}
			activityData_2 := models.KeyValuePairObject{
				KvpKey:   "ip.address2",
				KvpType:  models.EnumKVPTypeGeneralBool,
				KvpValue: "true",
			}
			deviceCheckDetail_1 := models.DeviceCheckDetailsObject{
				ActivityData:    []*models.KeyValuePairObject{&activityData_1, &activityData_2},
				ActivityType:    models.ActivityTypeConfirmation,
				CheckSessionKey: "sessionkey6",
				CheckType:       "somethinginvalid",
			}

			jsonRequest, err := json.Marshal(models.DeviceCheckDetailsObjectCollection{&deviceCheckDetail_1})
			Expect(err).To(BeNil())

			req, err := http.NewRequest("POST", "/isgood", bytes.NewReader(jsonRequest))
			Expect(err).To(BeNil())

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should response with 500 if activityType is not valid", func() {
			router := chi.NewRouter()
			router.Post("/isgood", isGood.IsGoodHandler)

			activityData_1 := models.KeyValuePairObject{
				KvpKey:   "ip.address",
				KvpType:  models.EnumKVPTypeGeneralBool,
				KvpValue: "true",
			}
			activityData_2 := models.KeyValuePairObject{
				KvpKey:   "ip.address2",
				KvpType:  models.EnumKVPTypeGeneralBool,
				KvpValue: "true",
			}
			deviceCheckDetail_1 := models.DeviceCheckDetailsObject{
				ActivityData:    []*models.KeyValuePairObject{&activityData_1, &activityData_2},
				ActivityType:    "somethinginvalid",
				CheckSessionKey: "sessionkey7",
				CheckType:       models.CheckTypeDevice,
			}

			jsonRequest, err := json.Marshal(models.DeviceCheckDetailsObjectCollection{&deviceCheckDetail_1})
			Expect(err).To(BeNil())

			req, err := http.NewRequest("POST", "/isgood", bytes.NewReader(jsonRequest))
			Expect(err).To(BeNil())

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should response with 500 if sessionKey is used", func() {
			router := chi.NewRouter()
			router.Post("/isgood", isGood.IsGoodHandler)

			activityData_1 := models.KeyValuePairObject{
				KvpKey:   "ip.address",
				KvpType:  models.EnumKVPTypeGeneralBool,
				KvpValue: "true",
			}
			deviceCheckDetail_1 := models.DeviceCheckDetailsObject{
				ActivityData:    []*models.KeyValuePairObject{&activityData_1},
				ActivityType:    models.ActivityTypeConfirmation,
				CheckSessionKey: "sessionkey2",
				CheckType:       models.CheckTypeDevice,
			}

			jsonRequest, err := json.Marshal(models.DeviceCheckDetailsObjectCollection{&deviceCheckDetail_1})
			Expect(err).To(BeNil())

			req, err := http.NewRequest("POST", "/isgood", bytes.NewReader(jsonRequest))
			Expect(err).To(BeNil())

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should response with 500 if request is malformed", func() {
			router := chi.NewRouter()
			router.Post("/isgood", isGood.IsGoodHandler)

			randomData := struct {
				Something string `json:"something"`
			}{"something nasty"}

			jsonRequest, err := json.Marshal(randomData)
			Expect(err).To(BeNil())

			req, err := http.NewRequest("POST", "/isgood", bytes.NewReader(jsonRequest))
			Expect(err).To(BeNil())

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
