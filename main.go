// Author hoenig

// commando is a command line tool for executing expect scripts
// via ssh on remote servers. It is as terrible as it sounds.
package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// typical example of running a basic command
// commando --command "uname -a" --hosts "tst-mexec{1..6}"

func main() {
	args := arguments()
	v := args.verbose

	tracef(v, "cliargs user: '%s'", args.user)
	tracef(v, "cliargs hosts: '%s'", args.hostexp)
	tracef(v, "cliargs scripts: '%s'", args.scriptdir)
	tracef(v, "cliargs nopassword: '%t'", args.nopassword)
	tracef(v, "cliargs verbose: '%t'", args.verbose)

	if err := validate(args); err != nil {
		dief("arguments are invalid: %v", err)
	}

	hosts := hosts(args.hostexp)
	if len(hosts) == 0 {
		dief("no hosts resolved from --host regex")
	}

	if args.command == "" {
		scripts, err := load(args)
		if err != nil {
			dief("failed to load scripts: %v", err)
		}

		color.Magenta("will execute scripts")
		color.Yellow(fmt.Sprintf("%v", scripts))
		color.Magenta("on hosts")
		color.Yellow(fmt.Sprintf("%v", hosts))

		pswd, err := prompt(args)
		if err != nil {
			dief("failed to read password: %v", err)
		}

		if err := run(args.user, pswd, hosts, scripts); err != nil {
			dief("failed to run scripts: %v", err)
		}
	} else {
		color.Magenta("will execute command")
		color.Yellow(args.command)
		color.Magenta("on hosts")
		color.Yellow(fmt.Sprintf("%v", hosts))

		pswd, err := prompt(args)
		if err != nil {
			dief("failed to read password: %v", err)
		}

		if err := runCmd(args.user, pswd, hosts, args.command); err != nil {
			dief("failed to run command: %v", err)
		}
	}
}

func dief(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func tracef(verbose bool, format string, args ...interface{}) {
	if verbose {
		color.Cyan(format+"\n", args...)
	}
}
