package helpers_test

import (
	"encoding/json"
	"github.com/albertsundjaja/frankie/models"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"

	. "github.com/albertsundjaja/frankie/helpers"
)

var _ = Describe("Helpers", func() {
	Describe("testing StringInSlice", func() {
		It("should return true if string in slice", func() {
			exists := StringInSlice("a", []string{"a", "b"})
			Expect(exists).To(BeTrue())
		})
		It("should return true if string not in slice", func() {
			exists := StringInSlice("c", []string{"a", "b"})
			Expect(exists).To(BeFalse())
		})
	})

	Describe("testing HttpError", func() {
		It("should return the correct []byte containing the error data", func() {
			expected := models.ErrorObject{Code: http.StatusInternalServerError, Message: "test message"}

			errReturnedByte := HttpErrorResponder(expected.Code, expected.Message)
			var errReturned models.ErrorObject
			err := json.Unmarshal(errReturnedByte, &errReturned)
			Expect(err).To(BeNil())
			Expect(cmp.Equal(errReturned, expected)).To(BeTrue())
		})
	})
})
