package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMountie(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mountie Suite")
}

const VCAP_SERVICES = `
{
    "sshfs": [
      {
        "credentials": {
          "host": "sshfs-server.example.com",
          "password": "b1b88dc6a9ffe81f96033c6b67c977",
          "port": 49159,
          "user": "aca29c23876c4138b4499aef1c4177b9"
        },
        "label": "sshfs",
        "name": "my-sshfs-instance",
        "plan": "1gb",
        "tags": []
      },
	  {
        "credentials": {
          "host": "different-fs-server.example.com",
          "password": "some-password",
          "port": 499,
          "user": "some-user"
        },
        "label": "sshfs",
        "name": "a-different-service-instance",
        "plan": "one-million-gigabytes",
        "tags": []
      }
    ],
    "v1-test-n/a": [
      {
        "credentials": {
          "host": "10.244.1.6",
          "login": "binding",
          "name": "2c7f2a8d-a923-4437-9a70-87da84b4813c",
          "port": 38828,
          "secret": "e373d9d1-0ed7-4dc2-82b0-2eb9e54b7f97",
          "url": "dummy-node.10.244.0.34.xip.io/2c7f2a8d-a923-4437-9a70-87da84b4813c"
        },
        "label": "v1-test-n/a",
        "name": "dummy-instance",
        "plan": "free",
        "tags": []
      }
    ]
 }`
