// Author hoenig

package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/pkg/errors"

	"golang.org/x/crypto/ssh/terminal"
)

func prompt(args args) (string, error) {
	if args.nopassword {
		tracef(args.verbose, "skipping password prompt")
		return "", nil
	}

	color.White("  password for '%s' --> ", args.user)
	bs, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", errors.Wrap(err, "failed to read password")
	}
	return string(bs), nil
}

func easyPrompt(user string) (string, error) {
	color.White("  password for '%s' --> ", user)
	bs, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", errors.Wrap(err, "failed to read password")
	}
	return string(bs), nil
}
