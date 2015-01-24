package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	var bindCommands []*exec.Cmd

	bindings := GetAllBindings()
	for _, binding := range bindings {
		bindCommands = append(bindCommands, CreateCommand(binding))
	}

	err := RunCommands(bindCommands)
	if err != nil {
		panic(err)
	}
}

type Credentials struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"password"`
}

type Binding struct {
	Credentials Credentials `json:"credentials"`
	Plan        string      `json:"plan"`
	Name        string      `json:"name"`
}

func GetAllBindings() []Binding {
	sshfsBindings := make(map[string][]Binding)
	jsonString := os.Getenv("VCAP_SERVICES")

	err := json.Unmarshal([]byte(jsonString), &sshfsBindings)
	if err != nil {
		panic("Unable to decode JSON in env var VCAP_SERVICES")
	}
	return sshfsBindings["sshfs"]
}

func CreateCommand(binding Binding) *exec.Cmd {
	cmd := exec.Command(
		"sshfs",
		fmt.Sprintf("%s@%s:", binding.Credentials.User, binding.Credentials.Host),
		"-p", strconv.Itoa(binding.Credentials.Port),
		"-o", "password_stdin",
		"-o", "StrictHostKeyChecking=false",
		binding.Name,
	)

	cmd.Stdin = strings.NewReader(binding.Credentials.Pass + "\n")

	return cmd
}

func RunCommands(commands []*exec.Cmd) error {
	for _, command := range commands {
		commandOutput, err := command.CombinedOutput()
		if err != nil {
			return errors.New(fmt.Sprintf(
				"Failed while running %s : %s : %s",
				command,
				commandOutput,
				err.Error()))
		}
	}

	return nil
}
