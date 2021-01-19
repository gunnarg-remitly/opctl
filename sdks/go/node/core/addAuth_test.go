package core

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	fakes "github.com/opctl/opctl/sdks/go/node/core/internal/fakes"
)

var _ = Context("core", func() {
	Context("AddAuth", func() {
		It("should call opAdder.Add w/ expected args", func() {

			/* arrange */
			fakeStore := new(fakes.FakeStateStore)
			objectUnderTest := core{stateStore: fakeStore}
			providedReq := model.AddAuthReq{
				Creds: model.Creds{
					Username: "username",
					Password: "password",
				},
				Resources: "resources",
			}

			/* act */
			result := objectUnderTest.AddAuth(
				context.Background(),
				providedReq,
			)

			/* assert */
			Expect(result).To(BeNil())
			Expect(fakeStore.AddAuthArgsForCall(0)).To(Equal(model.AuthAdded{
				Auth: model.Auth{
					Resources: providedReq.Resources,
					Creds:     providedReq.Creds,
				},
			}))
		})
	})
})
