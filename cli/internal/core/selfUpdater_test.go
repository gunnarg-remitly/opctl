package core

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/updater"
)

var _ = Context("_selfUpdateInvoker", func() {

	Context("SelfUpdate", func() {
		Context("invalid channel", func() {
			It("should return expected error", func() {
				/* arrange */
				providedReleaseChannel := "invalidChannel"

				objectUnderTest := newSelfUpdater()

				/* act */
				_, err := objectUnderTest.SelfUpdate(providedReleaseChannel)

				/* assert */
				Expect(err).To(MatchError(fmt.Sprintf(
					"%v is not an available release channel. "+
						"Available release channels are 'alpha', 'beta', and 'stable'.", providedReleaseChannel)))
			})
		})
		Context("valid channel", func() {
			It("should call updater.GetUpdateIfExists w/ expected args", func() {
				/* arrange */
				fakeUpdater := new(updater.Fake)

				objectUnderTest := _selfUpdateInvoker{
					updater: fakeUpdater,
				}

				providedChannel := "beta"

				/* act */
				_, err := objectUnderTest.SelfUpdate(providedChannel)

				/* assert */
				Expect(err).To(BeNil())
				Expect(fakeUpdater.GetUpdateIfExistsArgsForCall(0)).To(Equal(providedChannel))
			})
			Context("updater.GetUpdateIfExists errors", func() {
				It("should return expected error", func() {
					/* arrange */
					returnedError := errors.New("dummyError")

					fakeUpdater := new(updater.Fake)
					fakeUpdater.GetUpdateIfExistsReturns(&updater.Update{}, returnedError)

					objectUnderTest := _selfUpdateInvoker{
						updater: fakeUpdater,
					}

					/* act */
					_, err := objectUnderTest.SelfUpdate("beta")

					/* assert */
					Expect(err).To(MatchError(returnedError))
				})
			})
			Context("updater.GetUpdateIfExists doesn't error", func() {
				Context("update doesn't exist", func() {
					It("should return expected error", func() {
						/* arrange */
						fakeUpdater := new(updater.Fake)
						fakeUpdater.GetUpdateIfExistsReturns(nil, nil)

						objectUnderTest := _selfUpdateInvoker{
							updater: fakeUpdater,
						}

						/* act */
						message, err := objectUnderTest.SelfUpdate("beta")

						/* assert */
						Expect(err).To(BeNil())
						Expect(message).To(Equal("No update available, already at the latest version!"))
					})
				})
				Context("update exists", func() {
					It("should call updater.ApplyUpdate w/ expected args", func() {
						/* arrange */
						fakeUpdater := new(updater.Fake)
						returnedUpdate := &updater.Update{Version: "dummyVersion"}

						fakeUpdater.GetUpdateIfExistsReturns(returnedUpdate, nil)

						objectUnderTest := _selfUpdateInvoker{
							updater: fakeUpdater,
						}

						objectUnderTest.SelfUpdate("beta")
					})
					Context("updater.ApplyUpdate errors", func() {
						It("should return expected error", func() {
							/* arrange */
							returnedError := errors.New("dummyError")

							fakeUpdater := new(updater.Fake)

							fakeUpdater.GetUpdateIfExistsReturns(&updater.Update{Version: "dummyVersion"}, nil)

							fakeUpdater.ApplyUpdateReturns(returnedError)

							objectUnderTest := _selfUpdateInvoker{
								updater: fakeUpdater,
							}

							/* act */
							_, err := objectUnderTest.SelfUpdate("beta")

							/* assert */
							Expect(err).To(MatchError(returnedError))
						})
					})
				})
			})
		})
	})
})
