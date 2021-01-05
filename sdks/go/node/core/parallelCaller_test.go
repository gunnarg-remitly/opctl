package core

import (
	"context"
	"errors"
	"fmt"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uniquestringFakes "github.com/opctl/opctl/sdks/go/internal/uniquestring/fakes"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/core/internal/fakes"
)

var _ = Context("parallelCaller", func() {
	Context("newParallelCaller", func() {
		It("should return parallelCaller", func() {
			/* arrange/act/assert */
			Expect(newParallelCaller(
				new(FakeCaller),
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call caller for every parallelCall w/ expected args", func() {
			/* arrange */
			providedCallID := "dummyCallID"
			providedInboundScope := map[string]*model.Value{}
			providedRootCallID := "dummyRootCallID"
			providedOpPath := "providedOpPath"
			providedCallSpecParallelCalls := []*model.CallSpec{
				{
					Container: &model.ContainerCallSpec{},
				},
				{
					Op: &model.OpCallSpec{},
				},
				{
					Parallel: &[]*model.CallSpec{},
				},
				{
					Serial: &[]*model.CallSpec{},
				},
			}

			mtx := sync.Mutex{}

			fakeCaller := new(FakeCaller)
			eventChannel := make(chan model.Event, 100)
			callerCallIndex := 0
			fakeCaller.CallStub = func(
				context.Context,
				string,
				map[string]*model.Value,
				*model.CallSpec,
				string,
				*string,
				string,
			) (
				map[string]*model.Value,
				error,
			) {
				mtx.Lock()
				eventChannel <- model.Event{
					CallEnded: &model.CallEnded{
						Call: model.Call{
							ID: fmt.Sprintf("%v", callerCallIndex),
						},
					},
				}

				callerCallIndex++

				mtx.Unlock()

				return nil, nil
			}

			fakeUniqueStringFactory := new(uniquestringFakes.FakeUniqueStringFactory)
			uniqueStringCallIndex := 0
			expectedChildCallIDs := []string{}
			fakeUniqueStringFactory.ConstructStub = func() (string, error) {
				defer func() {
					uniqueStringCallIndex++
				}()
				childCallID := fmt.Sprintf("%v", uniqueStringCallIndex)
				expectedChildCallIDs = append(expectedChildCallIDs, fmt.Sprintf("%v", uniqueStringCallIndex))
				return childCallID, nil
			}

			objectUnderTest := _parallelCaller{
				caller:              fakeCaller,
				uniqueStringFactory: fakeUniqueStringFactory,
			}

			/* act */
			objectUnderTest.Call(
				context.Background(),
				providedCallID,
				providedInboundScope,
				providedRootCallID,
				providedOpPath,
				providedCallSpecParallelCalls,
			)

			/* assert */
			for callIndex := range providedCallSpecParallelCalls {
				_,
					actualNodeID,
					actualChildOutboundScope,
					actualCallSpec,
					actualOpPath,
					actualParentCallID,
					actualRootCallID := fakeCaller.CallArgsForCall(callIndex)

				Expect(actualChildOutboundScope).To(Equal(providedInboundScope))
				Expect(actualOpPath).To(Equal(providedOpPath))
				Expect(actualParentCallID).To(Equal(&providedCallID))
				Expect(actualRootCallID).To(Equal(providedRootCallID))

				// handle unordered asserts because call order can't be relied on within go statement
				Expect(expectedChildCallIDs).To(ContainElement(actualNodeID))
				Expect(providedCallSpecParallelCalls).To(ContainElement(actualCallSpec))
			}
		})
		It("can error", func() {
			/* arrange */
			providedCallID := "dummyCallID"
			providedInboundScope := map[string]*model.Value{}
			providedRootCallID := "dummyRootCallID"
			providedOpPath := "providedOpPath"
			providedCallSpecParallelCalls := []*model.CallSpec{
				{
					Container: &model.ContainerCallSpec{},
				},
				{
					Op: &model.OpCallSpec{},
				},
				{
					Parallel: &[]*model.CallSpec{},
				},
				{
					Serial: &[]*model.CallSpec{},
				},
			}

			expectedErr := errors.New("errorMessage")

			fakeCaller := new(FakeCaller)
			fakeCaller.CallStub = func(
				context.Context,
				string,
				map[string]*model.Value,
				*model.CallSpec,
				string,
				*string,
				string,
			) (
				map[string]*model.Value,
				error,
			) {
				return nil, expectedErr
			}

			fakeUniqueStringFactory := new(uniquestringFakes.FakeUniqueStringFactory)
			uniqueStringCallIndex := 0
			expectedChildCallIDs := []string{}
			fakeUniqueStringFactory.ConstructStub = func() (string, error) {
				defer func() {
					uniqueStringCallIndex++
				}()
				childCallID := fmt.Sprintf("%v", uniqueStringCallIndex)
				expectedChildCallIDs = append(expectedChildCallIDs, childCallID)
				return childCallID, nil
			}

			objectUnderTest := _parallelCaller{
				caller:              fakeCaller,
				uniqueStringFactory: fakeUniqueStringFactory,
			}

			/* act */
			actualOutputs, actualErr := objectUnderTest.Call(
				context.Background(),
				providedCallID,
				providedInboundScope,
				providedRootCallID,
				providedOpPath,
				providedCallSpecParallelCalls,
			)

			/* assert */
			Expect(actualOutputs).To(BeNil())
			Expect(actualErr).To(MatchError(expectedErr))
		})
		Context("caller doesn't error", func() {
			It("shouldn't exit until all childCalls complete & not error", func() {
				/* arrange */
				providedCallID := "dummyCallID"
				providedInboundScope := map[string]*model.Value{}
				providedRootCallID := "dummyRootCallID"
				providedOpPath := "providedOpPath"
				providedCallSpecParallelCalls := []*model.CallSpec{
					{
						Container: &model.ContainerCallSpec{},
					},
					{
						Op: &model.OpCallSpec{},
					},
					{
						Parallel: &[]*model.CallSpec{},
					},
					{
						Serial: &[]*model.CallSpec{},
					},
				}

				fakeCaller := new(FakeCaller)
				fakeCaller.CallStub = func(
					context.Context,
					string,
					map[string]*model.Value,
					*model.CallSpec,
					string,
					*string,
					string,
				) (
					map[string]*model.Value,
					error,
				) {
					return nil, nil
				}

				fakeUniqueStringFactory := new(uniquestringFakes.FakeUniqueStringFactory)
				uniqueStringCallIndex := 0
				expectedChildCallIDs := []string{}
				fakeUniqueStringFactory.ConstructStub = func() (string, error) {
					defer func() {
						uniqueStringCallIndex++
					}()
					childCallID := fmt.Sprintf("%v", uniqueStringCallIndex)
					expectedChildCallIDs = append(expectedChildCallIDs, fmt.Sprintf("%v", uniqueStringCallIndex))
					return childCallID, nil
				}

				objectUnderTest := _parallelCaller{
					caller:              fakeCaller,
					uniqueStringFactory: fakeUniqueStringFactory,
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					providedCallID,
					providedInboundScope,
					providedRootCallID,
					providedOpPath,
					providedCallSpecParallelCalls,
				)

				/* assert */
				for callIndex := range providedCallSpecParallelCalls {
					_,
						actualNodeID,
						actualChildOutboundScope,
						actualCallSpec,
						actualOpPath,
						actualParentCallID,
						actualRootCallID := fakeCaller.CallArgsForCall(callIndex)

					Expect(actualChildOutboundScope).To(Equal(providedInboundScope))
					Expect(actualOpPath).To(Equal(providedOpPath))
					Expect(actualParentCallID).To(Equal(&providedCallID))
					Expect(actualRootCallID).To(Equal(providedRootCallID))

					// handle unordered asserts because call order can't be relied on within go statement
					Expect(expectedChildCallIDs).To(ContainElement(actualNodeID))
					Expect(providedCallSpecParallelCalls).To(ContainElement(actualCallSpec))
				}
			})
		})
	})
})
