package auth

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	coreFakes "github.com/opctl/opctl/sdks/go/node/core/fakes"
)

var _ = Context("Adder", func() {
	Context("Invoke", func() {
		It("should call apiClient.Invoke w/ expected args", func() {
			/* arrange */
			fakeCore := new(coreFakes.FakeCore)

			providedCtx := context.TODO()

			expectedCtx := providedCtx
			expectedReq := model.AddAuthReq{
				Resources: "Resources",
				Creds: model.Creds{
					Username: "username",
					Password: "password",
				},
			}

			objectUnderTest := _adder{core: fakeCore}

			/* act */
			err := objectUnderTest.Add(
				expectedCtx,
				expectedReq.Resources,
				expectedReq.Username,
				expectedReq.Password,
			)

			/* assert */
			actualReq := fakeCore.AddAuthArgsForCall(0)
			Expect(err).To(BeNil())
			Expect(actualReq).To(BeEquivalentTo(expectedReq))
		})
		Context("apiClient.Invoke errors", func() {
			It("should return expected error", func() {
				/* arrange */
				fakeCore := new(coreFakes.FakeCore)
				expectedError := errors.New("dummyError")
				fakeCore.AddAuthReturns(expectedError)

				objectUnderTest := _adder{core: fakeCore}

				/* act */
				err := objectUnderTest.Add(context.TODO(), "", "", "")

				/* assert */
				Expect(err).To(MatchError(expectedError))
			})
		})
	})
})
