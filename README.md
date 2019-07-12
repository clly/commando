commando
========

Run commands in bulk on remote machines.

[![Go Report Card](https://goreportcard.com/badge/go.gophers.dev/cmds/commando)](https://goreportcard.com/report/go.gophers.dev/cmds/commando)
[![Build Status](https://travis-ci.com/shoenig/commando.svg?branch=master)](https://travis-ci.com/shoenig/commando)
[![GoDoc](https://godoc.org/go.gophers.dev/cmds/commando?status.svg)](https://godoc.org/go.gophers.dev/cmds/commando)
[![NetflixOSS Lifecycle](https://img.shields.io/osslifecycle/shoenig/regexplus.svg)](OSSMETADATA)
[![GitHub](https://img.shields.io/github/license/shoenig/regexplus.svg)](LICENSE)

# Project Overview

Module `go.gophers.dev/cmds/commando` provides a CLI utility for running commands
in bulk on a set of machines. It allows for setting a password that will be used
in prompts on the remote machine (e.g. using sudo).

# Getting Started

The `commando` command can be installed by running
```
$ go get go.gophers.dev/cmds/commando
```

### Example usage

#### No password required
```bash
$ commando --hosts gophers.dev --command "uname -a"
will execute command
uname -a
on hosts
[gophers.dev]
--- gophers.dev ---
executing command `uname -a`
Linux ubs1 3.10.0-957.21.3.el7.x86_64 #1 SMP Tue Jun 18 16:35:19 UTC 2019 x86_64 x86_64 x86_64 GNU/Linux
```

#### With password required
```bash
# running sudo whoami will cause a password prompt, so we pass --pw to have
# commando prompt for a password to send and answer the prompt with

$ ./commando --hosts gophers.dev --command "sudo whoami" --pw
will execute command
sudo whoami
on hosts
[gophers.dev]
  password for 'hoenig' -->
--- gophers.dev ---
executing command `sudo whoami`
[sudo] password for hoenig:
root
```

# Contributing

The `go.gophers.dev/cmds/commando` module is always improving with new features
and error corrections. For contributing bug fixes and new features please file an issue.

# License

The `go.gophers.dev/cmds/commando` module is open source under the [BSD-3-Clause](LICENSE) license.
