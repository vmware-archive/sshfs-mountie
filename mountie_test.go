package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	. "github.com/pivotal-cf/sshfs-mountie"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mountie", func() {
	BeforeEach(func() {
		os.Setenv("VCAP_SERVICES", VCAP_SERVICES)
	})

	AfterEach(func() {
		os.Setenv("VCAP_SERVICES", `{}`)
	})

	Describe("Getting bindings from the environment", func() {
		Context("when there are multiple services in the VCAP_SERVICES env var", func() {
			It("Returns only the sshfs bindings", func() {
				binding0 := Binding{
					Plan: "1gb",
					Name: "my-sshfs-instance",
					Credentials: Credentials{
						Host: "sshfs-server.example.com",
						Port: 49159,
						User: "aca29c23876c4138b4499aef1c4177b9",
						Pass: "b1b88dc6a9ffe81f96033c6b67c977",
					},
				}
				binding1 := Binding{
					Plan: "one-million-gigabytes",
					Name: "a-different-service-instance",
					Credentials: Credentials{
						Host: "different-fs-server.example.com",
						Port: 499,
						User: "some-user",
						Pass: "some-password",
					},
				}

				bindings := GetAllBindings()

				Expect(bindings).To(ConsistOf([]Binding{binding0, binding1}))
			})
		})
	})

	Describe("Generating shell command to call sshfs client binary", func() {
		It("returns the correct command struct when given a Binding", func() {
			binding := Binding{
				Plan: "one-million-gigabytes",
				Name: "a-different-service-instance",
				Credentials: Credentials{
					Host: "different-fs-server.example.com",
					Port: 499,
					User: "some-user",
					Pass: "some-password",
				},
			}

			sshfsMountCommand := CreateCommand(binding)
			password, err := ioutil.ReadAll(sshfsMountCommand.Stdin)
			if err != nil {
				panic(err)
			}

			Expect(string(password)).To(Equal("some-password\n"))
			Expect(sshfsMountCommand.Args).To(Equal([]string{"sshfs", fmt.Sprintf("%s@%s:", binding.Credentials.User, binding.Credentials.Host), "-p", "499",
				"-o", "password_stdin", "-o", "StrictHostKeyChecking=false", binding.Name}))
		})
	})

	Describe("running commands", func() {
		Context("when all commands succeed", func() {
			It("returns a nil error", func() {
				cmds := []*exec.Cmd{
					exec.Command("echo", "hello"),
					exec.Command("echo", "world"),
				}

				err := RunCommands(cmds)

				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("if one of the commands failed", func() {
			It("returns the error", func() {
				cmds := []*exec.Cmd{
					exec.Command("/bin/bash", "-c", "echo success && exit 0"),
					exec.Command("/bin/bash", "-c", "echo mounting failed because your network is down && exit 3"),
				}

				err := RunCommands(cmds)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("mounting failed because your network is down"))
			})
		})
	})
})
