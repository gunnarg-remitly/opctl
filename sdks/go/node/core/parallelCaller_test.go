package core

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dgraph-io/badger/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	containerRuntimeFakes "github.com/opctl/opctl/sdks/go/node/core/containerruntime/fakes"
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

		Context("caller errors", func() {

			It("should return expected results", func() {
				/* arrange */
				dbDir, err := ioutil.TempDir("", "")
				if nil != err {
					panic(err)
				}

				db, err := badger.Open(
					badger.DefaultOptions(dbDir).WithLogger(nil),
				)
				if nil != err {
					panic(err)
				}
				pubSub := pubsub.New(db)

				objectUnderTest := _parallelCaller{
					caller: newCaller(
						newContainerCaller(
							new(containerRuntimeFakes.FakeContainerRuntime),
							pubSub,
							newStateStore(
								context.Background(),
								db,
								pubSub,
							),
						),
						dbDir,
						pubSub,
					),
					pubSub: pubSub,
				}

				/* act */
				_, actualErr := objectUnderTest.Call(
					context.Background(),
					"callID",
					map[string]*model.Value{},
					"rootCallID",
					"opPath",
					[]*model.CallSpec{
						{
							// intentionally invalid
							Container: &model.ContainerCallSpec{},
						},
					},
				)

				/* assert */
				Expect(actualErr.Error()).To(Equal("child call failed"))
			})
		})

		It("should start each child as expected", func() {

			/* arrange */
			dbDir, err := ioutil.TempDir("", "")
			if nil != err {
				panic(err)
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
			providedOpRef := "providedOpRef"
			providedParentID := "providedParentID"
			providedRootID := "providedRootID"
			childOpRef := filepath.Join(wd, "testdata/parallelCaller")
			input1Key := "input1"
			childOp1Path := filepath.Join(childOpRef, "op1")
			childOp2Path := filepath.Join(childOpRef, "op2")

			ctx := context.Background()

			fakeContainerRuntime := new(containerRuntimeFakes.FakeContainerRuntime)
			fakeContainerRuntime.RunContainerStub = func(
				ctx context.Context,
				req *model.ContainerCall,
				rootCallID string,
				eventPublisher pubsub.EventPublisher,
				stdOut io.WriteCloser,
				stdErr io.WriteCloser,
			) (*int64, error) {

				stdErr.Close()
				stdOut.Close()

			objectUnderTest := _parallelCaller{
				caller:              fakeCaller,
				uniqueStringFactory: fakeUniqueStringFactory,
			}

			eventChannel, err := pubSub.Subscribe(
				ctx,
				model.EventFilter{},
			)
			if nil != err {
				panic(err)
			}

			input1Value := "input1Value"
			providedInboundScope := map[string]*model.Value{
				input1Key: {String: &input1Value},
			}

			objectUnderTest := _parallelCaller{
				caller: newCaller(
					newContainerCaller(
						fakeContainerRuntime,
						pubSub,
						newStateStore(
							ctx,
							db,
							pubSub,
						),
					),
					dbDir,
					pubSub,
				),
				pubSub: pubSub,
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
			providedOpRef := "providedOpRef"
			providedParentID := "providedParentID"
			providedRootID := "providedRootID"
			childOpRef := filepath.Join(wd, "testdata/parallelCaller")
			input1Key := "input1"
			childOp1Path := filepath.Join(childOpRef, "op1")
			childOp2Path := filepath.Join(childOpRef, "op2")

			ctx := context.Background()

			fakeContainerRuntime := new(containerRuntimeFakes.FakeContainerRuntime)
			fakeContainerRuntime.RunContainerStub = func(
				ctx context.Context,
				req *model.ContainerCall,
				rootCallID string,
				eventPublisher pubsub.EventPublisher,
				stdOut io.WriteCloser,
				stdErr io.WriteCloser,
			) (*int64, error) {

				stdErr.Close()
				stdOut.Close()

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
			if nil != err {
				panic(err)
			}

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
