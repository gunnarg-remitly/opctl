package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

type fakeStateStore struct {
	request model.AuthAdded
}

func (ss fakeStateStore) AddAuth(req model.AuthAdded) error {
	ss.request = req
	return nil
}

func (ss fakeStateStore) TryGetAuth(resource string) *model.Auth {
	return nil
}

var _ = Context("core", func() {
	Context("AddAuth", func() {
		It("should call opAdder.Add w/ expected args", func() {

			/* arrange */
			store := fakeStateStore{}
			objectUnderTest := _core{stateStore: store}
			providedReq := model.AddAuthReq{
				Creds: model.Creds{
					Username: "username",
					Password: "password",
				},
				Resources: "resources",
			}

			/* act */
			result := objectUnderTest.AddAuth(providedReq)

			/* assert */
			Expect(result).To(BeNil())
			Expect(store.request.Auth.Resources).To(Equal("resources"))
		})
	})
})
