package main

import (
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type Auth struct {

}

func sshAgentAuth() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func PasswordCallback(user string) ssh.AuthMethod {
	return ssh.PasswordCallback(func() (string, error) {
		return easyPrompt(user)
	})
}
