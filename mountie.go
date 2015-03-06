package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

func main() {
	for _, binding := range GetAllBindings() {
		mountPoint, err := MakeMountPoint(binding)
		if err != nil {
			panic(err)
		}
		RunCommand(CreateCommand(binding))
		fmt.Printf("Mounted SSHFS filesystem for service instance %s into %s\n", binding.Name, mountPoint)
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
	sshfsBindings := map[string][]Binding{}
	jsonString := os.Getenv("VCAP_SERVICES")

	err := json.Unmarshal([]byte(jsonString), &sshfsBindings)
	if err != nil {
		panic("Unable to decode JSON in env var VCAP_SERVICES")
	}

	return sshfsBindings["sshfs"]
}

func MakeMountPoint(binding Binding) (mountPath string, err error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	mountPath = path.Join(currentPath, binding.Name)
	err = os.MkdirAll(mountPath, 0777)
	return mountPath, err
}

func CreateCommand(binding Binding) *exec.Cmd {
	cmd := exec.Command(
		"sshfs",
		fmt.Sprintf("%s@%s:", binding.Credentials.User, binding.Credentials.Host),
		"-p", strconv.Itoa(binding.Credentials.Port),
		"-o", "password_stdin",
		"-o", "StrictHostKeyChecking=false",
		"-o", "reconnect",
		"-o", "sshfs_debug",
		binding.Name,
	)

	cmd.Stdin = strings.NewReader(binding.Credentials.Pass + "\n")

	return cmd
}

func RunCommand(command *exec.Cmd) error {
	commandOutput, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed while running %s : %s : %s", command, commandOutput, err)
	}

	return nil
}
