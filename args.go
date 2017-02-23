// Author hoenig

package main

import (
	"flag"
	"os"

	"github.com/pkg/errors"
)

type args struct {
	user       string
	hostexp    string
	scriptdir  string
	command    string
	pw         bool
	nopassword bool
	verbose    bool
}

func arguments() args {
	var args args

	flag.StringVar(&args.user, "user", os.Getenv("USER"), "ssh username")
	flag.StringVar(&args.hostexp, "hosts", "", "the list of hosts")
	flag.StringVar(&args.scriptdir, "scripts", "", "the directory full of scripts")
	flag.StringVar(&args.command, "command", "", "the command to run")
	flag.BoolVar(&args.pw, "pw", false, "send password on stdin after running --command")
	flag.BoolVar(&args.nopassword, "no-password", false, "no-password skips password prompt")
	flag.BoolVar(&args.verbose, "verbose", false, "verbose mode")

	flag.Parse()

	return args
}

func validate(args args) error {
	if args.hostexp == "" {
		return errors.Errorf("--hosts is required")
	}

	if args.user == "" {
		return errors.Errorf("--user or $USER must be set")
	}

	if args.scriptdir == "" && args.command == "" {
		return errors.Errorf("--scripts or --command is required")
	}

	if args.scriptdir != "" && args.command != "" {
		return errors.Errorf("only one of --scripts or --command allowed")
	}

	if args.command == "" && args.pw {
		return errors.Errorf("--pw only allowed in conjunction with --command")
	}

	return nil
}
