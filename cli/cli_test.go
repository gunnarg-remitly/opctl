package main

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/clioutput"
	corePkg "github.com/opctl/opctl/cli/internal/core"
	authFakes "github.com/opctl/opctl/cli/internal/core/auth/fakes"
	coreFakes "github.com/opctl/opctl/cli/internal/core/fakes"
	opFakes "github.com/opctl/opctl/cli/internal/core/op/fakes"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
)

var _ = Context("cli", func() {
	Context("Run", func() {

		Context("--no-color", func() {
			It("should set color.NoColor", func() {
				Skip("")
				/* arrange */
				// var cliOutput clioutput.CliOutput

				objectUnderTest, _ := newCli(
					context.Background(),
					func(
						context.Context,
						clioutput.CliOutput,
						clioutput.OpFormatter,
						local.NodeCreateOpts,
					) (corePkg.Core, error) {
						// cliOutput = _cliOutput
						return new(coreFakes.FakeCore), nil
					},
				)

				objectUnderTest.Run([]string{"opctl", "--no-color", "ls"})

				// No colors are applied
				// Expect(cliOutput.DisableColorCallCount()).To(Equal(1))
			})
		})

		Context("auth", func() {

			Context("add", func() {

				It("should call authAddCmd.Invoke w/ expected args", func() {
					/* arrange */
					providedResources := "resources"
					providedUsername := "username"
					providedPassword := "password"

					fakeCore := new(coreFakes.FakeCore)

					fakeAuth := new(authFakes.FakeAuth)
					fakeCore.AuthReturns(fakeAuth)

					objectUnderTest, _ := newCli(
						context.Background(),
						func(
							context.Context,
							clioutput.CliOutput,
							clioutput.OpFormatter,
							local.NodeCreateOpts,
						) (corePkg.Core, error) {
							return fakeCore, nil
						},
					)

					/* act */
					objectUnderTest.Run([]string{"opctl", "auth", "add", providedResources, "-u", providedUsername, "-p", providedPassword})

					/* assert */
					_,
						actualResources,
						actualUsername,
						actualPassword := fakeAuth.AddArgsForCall(0)

					Expect(actualResources).To(Equal(providedResources))
					Expect(actualUsername).To(Equal(providedUsername))
					Expect(actualPassword).To(Equal(providedPassword))
				})

			})

		})

		Context("ls", func() {
			Context("w/ dirRef", func() {

				It("should call coreFakes.Ls w/ expected args", func() {
					/* arrange */
					fakeCore := new(coreFakes.FakeCore)

					expectedDirRef := "dummyPath"
					objectUnderTest, _ := newCli(
						context.Background(),
						func(
							context.Context,
							clioutput.CliOutput,
							clioutput.OpFormatter,
							local.NodeCreateOpts,
						) (corePkg.Core, error) {
							return fakeCore, nil
						},
					)

					/* act */
					objectUnderTest.Run([]string{"opctl", "ls", expectedDirRef})

					/* assert */
					actualCtx,
						actualDirRef := fakeCore.LsArgsForCall(0)

					Expect(actualCtx).To(Equal(context.TODO()))
					Expect(actualDirRef).To(Equal(expectedDirRef))
				})
			})
			Context("w/out dirRef", func() {

				It("should call coreFakes.Ls w/ expected args", func() {
					/* arrange */
					fakeCore := new(coreFakes.FakeCore)

					expectedDirRef := ".opspec"
					objectUnderTest, _ := newCli(
						context.Background(),
						func(
							context.Context,
							clioutput.CliOutput,
							clioutput.OpFormatter,
							local.NodeCreateOpts,
						) (corePkg.Core, error) {
							return fakeCore, nil
						},
					)

					/* act */
					objectUnderTest.Run([]string{"opctl", "ls"})

					/* assert */
					actualCtx,
						actualDirRef := fakeCore.LsArgsForCall(0)

					Expect(actualCtx).To(Equal(context.TODO()))
					Expect(actualDirRef).To(Equal(expectedDirRef))
				})
			})
		})

		Context("op", func() {

			Context("create", func() {
				Context("w/ path", func() {
					It("should call coreFakes.Create w/ expected args", func() {
						/* arrange */
						fakeCore := new(coreFakes.FakeCore)

						fakeOp := new(opFakes.FakeOp)
						fakeCore.OpReturns(fakeOp)

						expectedOpName := "dummyOpName"
						expectedPath := "dummyPath"

						objectUnderTest, _ := newCli(
							context.Background(),
							func(
								context.Context,
								clioutput.CliOutput,
								clioutput.OpFormatter,
								local.NodeCreateOpts,
							) (corePkg.Core, error) {
								return fakeCore, nil
							},
						)

						/* act */
						objectUnderTest.Run([]string{"opctl", "op", "create", "--path", expectedPath, expectedOpName})

						/* assert */
						actualPath, actualOpDescription, actualOpName := fakeOp.CreateArgsForCall(0)
						Expect(actualOpName).To(Equal(expectedOpName))
						Expect(actualOpDescription).To(BeEmpty())
						Expect(actualPath).To(Equal(expectedPath))
					})
				})

				Context("w/out path", func() {
					It("should call coreFakes.Create w/ expected args", func() {
						/* arrange */
						fakeCore := new(coreFakes.FakeCore)

						fakeOp := new(opFakes.FakeOp)
						fakeCore.OpReturns(fakeOp)

						expectedOpName := "dummyOpName"
						expectedPath := ".opspec"

						objectUnderTest, _ := newCli(
							context.Background(),
							func(
								context.Context,
								clioutput.CliOutput,
								clioutput.OpFormatter,
								local.NodeCreateOpts,
							) (corePkg.Core, error) {
								return fakeCore, nil
							},
						)

						/* act */
						objectUnderTest.Run([]string{"opctl", "op", "create", expectedOpName})

						/* assert */
						actualPath, actualOpDescription, actualOpName := fakeOp.CreateArgsForCall(0)
						Expect(actualOpName).To(Equal(expectedOpName))
						Expect(actualOpDescription).To(BeEmpty())
						Expect(actualPath).To(Equal(expectedPath))
					})
				})
				Context("w/ description", func() {
					It("should call coreFakes.Create w/ expected args", func() {
						/* arrange */
						fakeCore := new(coreFakes.FakeCore)

						fakeOp := new(opFakes.FakeOp)
						fakeCore.OpReturns(fakeOp)

						expectedOpName := "dummyOpName"
						expectedOpDescription := "dummyOpDescription"
						expectedPath := ".opspec"

						objectUnderTest, _ := newCli(
							context.Background(),
							func(
								context.Context,
								clioutput.CliOutput,
								clioutput.OpFormatter,
								local.NodeCreateOpts,
							) (corePkg.Core, error) {
								return fakeCore, nil
							},
						)

						/* act */
						objectUnderTest.Run([]string{"opctl", "op", "create", "-d", expectedOpDescription, expectedOpName})

						/* assert */
						actualPath, actualOpDescription, actualOpName := fakeOp.CreateArgsForCall(0)
						Expect(actualOpName).To(Equal(expectedOpName))
						Expect(actualOpDescription).To(Equal(expectedOpDescription))
						Expect(actualPath).To(Equal(expectedPath))
					})
				})

				Context("w/out description", func() {
					It("should call coreFakes.Create w/ expected args", func() {
						/* arrange */
						fakeCore := new(coreFakes.FakeCore)

						fakeOp := new(opFakes.FakeOp)
						fakeCore.OpReturns(fakeOp)

						expectedName := "dummyOpName"
						expectedPath := ".opspec"

						objectUnderTest, _ := newCli(
							context.Background(),
							func(
								context.Context,
								clioutput.CliOutput,
								clioutput.OpFormatter,
								local.NodeCreateOpts,
							) (corePkg.Core, error) {
								return fakeCore, nil
							},
						)

						/* act */
						objectUnderTest.Run([]string{"opctl", "op", "create", expectedName})

						/* assert */
						actualPath, actualOpDescription, actualOpName := fakeOp.CreateArgsForCall(0)
						Expect(actualOpName).To(Equal(expectedName))
						Expect(actualOpDescription).To(BeEmpty())
						Expect(actualPath).To(Equal(expectedPath))
					})
				})
			})

			Context("install", func() {
				It("should call coreFakes.Install w/ expected args", func() {
					/* arrange */
					fakeCore := new(coreFakes.FakeCore)

					fakeOp := new(opFakes.FakeOp)
					fakeCore.OpReturns(fakeOp)

					expectedPath := "dummyPath"
					expectedOpRef := "dummyOpRef"
					expectedUsername := "dummyUsername"
					expectedPassword := "dummyPassword"

					objectUnderTest, _ := newCli(
						context.Background(),
						func(
							context.Context,
							clioutput.CliOutput,
							clioutput.OpFormatter,
							local.NodeCreateOpts,
						) (corePkg.Core, error) {
							return fakeCore, nil
						},
					)

					/* act */
					objectUnderTest.Run([]string{
						"opctl",
						"op",
						"install",
						"--path",
						expectedPath,
						"-u",
						expectedUsername,
						"-p",
						expectedPassword,
						expectedOpRef,
					})

					/* assert */
					actualCtx,
						actualPath,
						actualOpRef,
						actualUsername,
						actualPassword := fakeOp.InstallArgsForCall(0)

					Expect(actualCtx).To(Equal(context.TODO()))
					Expect(actualPath).To(Equal(expectedPath))
					Expect(actualOpRef).To(Equal(expectedOpRef))
					Expect(actualUsername).To(Equal(expectedUsername))
					Expect(actualPassword).To(Equal(expectedPassword))
				})
			})

			Context("validate", func() {

				It("should call coreFakes.OpValidate w/ expected args", func() {
					/* arrange */
					fakeCore := new(coreFakes.FakeCore)

					fakeOp := new(opFakes.FakeOp)
					fakeCore.OpReturns(fakeOp)

					opRef := ".opspec/dummyOpName"

					objectUnderTest, _ := newCli(
						context.Background(),
						func(
							context.Context,
							clioutput.CliOutput,
							clioutput.OpFormatter,
							local.NodeCreateOpts,
						) (corePkg.Core, error) {
							return fakeCore, nil
						},
					)

					/* act */
					objectUnderTest.Run([]string{"opctl", "op", "validate", opRef})

					/* assert */
					actualCtx,
						actualOpRef := fakeOp.ValidateArgsForCall(0)

					Expect(actualCtx).To(Equal(context.TODO()))
					Expect(actualOpRef).To(Equal(opRef))
				})

			})

		})

		Context("run", func() {
			Context("with two op run args & an arg-file", func() {
				It("should call coreFakes.Run w/ expected args", func() {
					/* arrange */
					fakeCore := new(coreFakes.FakeCore)

					expectedRunOpts := &corePkg.RunOpts{
						Args:    []string{"arg1Name=arg1Value", "arg2Name=arg2Value"},
						ArgFile: "dummyArgFile",
					}
					expectedOpRef := ".opspec/dummyOpName"

					objectUnderTest, _ := newCli(
						context.Background(),
						func(
							context.Context,
							clioutput.CliOutput,
							clioutput.OpFormatter,
							local.NodeCreateOpts,
						) (corePkg.Core, error) {
							return fakeCore, nil
						},
					)

					/* act */
					objectUnderTest.Run([]string{
						"opctl",
						"run",
						"-a",
						expectedRunOpts.Args[0],
						"-a",
						expectedRunOpts.Args[1],
						"--arg-file",
						expectedRunOpts.ArgFile,
						expectedOpRef,
					})

					/* assert */
					actualCtx,
						actualOpUrl,
						actualRunOpts, _ := fakeCore.RunArgsForCall(0)

					Expect(actualCtx).To(Equal(context.TODO()))
					Expect(actualOpUrl).To(Equal(expectedOpRef))
					Expect(actualRunOpts).To(Equal(expectedRunOpts))
				})
			})

			Context("with zero op run args", func() {
				It("should call coreFakes.Run w/ expected args", func() {
					/* arrange */
					fakeCore := new(coreFakes.FakeCore)

					expectedOpRef := ".opspec/dummyOpName"

					objectUnderTest, _ := newCli(
						context.Background(),
						func(
							context.Context,
							clioutput.CliOutput,
							clioutput.OpFormatter,
							local.NodeCreateOpts,
						) (corePkg.Core, error) {
							return fakeCore, nil
						},
					)

					/* act */
					objectUnderTest.Run([]string{"opctl", "run", expectedOpRef})

					/* assert */
					actualCtx,
						actualOpRef,
						actualRunOpts, _ := fakeCore.RunArgsForCall(0)

					Expect(actualCtx).To(Equal(context.TODO()))
					Expect(actualOpRef).To(Equal(expectedOpRef))
					Expect(actualRunOpts.Args).To(BeEmpty())
				})
			})
		})
	})

	Context("self-update", func() {

		Context("with channel flag", func() {

			It("should call coreFakes.SelfUpdate with expected releaseChannel", func() {
				/* arrange */
				expectedChannel := "beta"

				fakeCore := new(coreFakes.FakeCore)

				objectUnderTest, _ := newCli(
					context.Background(),
					func(
						context.Context,
						clioutput.CliOutput,
						clioutput.OpFormatter,
						local.NodeCreateOpts,
					) (corePkg.Core, error) {
						return fakeCore, nil
					},
				)

				/* act */
				objectUnderTest.Run([]string{"opctl", "self-update", "-c", expectedChannel})

				/* assert */
				actualChannel := fakeCore.SelfUpdateArgsForCall(0)
				Expect(actualChannel).To(Equal(expectedChannel))
			})

		})

		Context("without channel flag", func() {

			It("should call coreFakes.SelfUpdate with expected releaseChannel", func() {
				/* arrange */
				expectedChannel := "stable"

				fakeCore := new(coreFakes.FakeCore)

				objectUnderTest, _ := newCli(
					context.Background(),
					func(
						context.Context,
						clioutput.CliOutput,
						clioutput.OpFormatter,
						local.NodeCreateOpts,
					) (corePkg.Core, error) {
						return fakeCore, nil
					},
				)

				/* act */
				objectUnderTest.Run([]string{"opctl", "self-update"})

				/* assert */
				actualChannel := fakeCore.SelfUpdateArgsForCall(0)
				Expect(actualChannel).To(Equal(expectedChannel))
			})
		})

	})

})
