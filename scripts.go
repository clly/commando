// Author hoenig

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

type script struct {
	command string
	stdin   []string
}

// A scriptfile contains one or more scripts to be executed. the first line is the command and the second line is stdin
type scriptfile struct {
	name    string
	scripts []script
	sudo    bool
}

func (s scriptfile) String() string {
	return s.name
}

func load(cfg args) ([]scriptfile, error) {
	scripts := []scriptfile{}

	err := filepath.Walk(cfg.scriptdir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "failed to read scripts")
		}

		// skip directories
		if info.IsDir() {
			return nil
		}

		script, err := read(info.Name(), path)
		if err != nil {
			return errors.Wrapf(err, "failed to read script file %s", info.Name())
		}

		scripts = append(scripts, script)
		return nil
	})

	if len(scripts) == 0 {
		return nil, errors.Errorf("no scripts found")
	}

	return scripts, err
}

func read(name, path string) (scriptfile, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return scriptfile{}, errors.Wrap(err, "failed to read script")
	}
	s := strings.TrimSpace(string(bs))
	return parse(name, s)
}

func parse(name, content string) (scriptfile, error) {
	// --- is a "key" string which causes us to create a new "script" from a single file
	parts := strings.Split(content, "---")
	scriptfile := scriptfile{name: name}

	for _, part := range parts {
		if strings.Contains(part, "PASSWORD") {
			scriptfile.sudo = true
		}
		lines := cleanup(strings.Split(part, "\n"))
		if len(lines) == 0 {
			return scriptfile, errors.Errorf("no command in script %s", name)
		}
		s := script{lines[0], lines[1:]}
		scriptfile.scripts = append(scriptfile.scripts, s)
	}
	return scriptfile, nil
}

func cleanup(lines []string) []string {
	cleansed := make([]string, 0, len(lines))
	for _, dirty := range lines {
		clean := strings.TrimSpace(dirty)
		switch {
		case clean == "":
		case clean[0] == '#':
		default:
			cleansed = append(cleansed, clean)
		}
	}
	return cleansed
}

func run(user, pass string, hosts []string, files []scriptfile) error {
	for _, host := range hosts {

		client, err := makeClient(user, pass, host)
		if err != nil {
			return errors.Wrap(err, "failed to dial host")
		}

		for _, file := range files {
			if err := executeScriptfile(client, user, pass, host, file); err != nil {
				return errors.Wrapf(err, "failed to run %s on %s", file, host)
			}
			fmt.Println("")
		}
	}
	return nil
}

func runCmd(user string, hosts []string, command string, pw bool) error {
	var pass string

	if pw || strings.Contains(command, "sudo") {
		var err error
		pass, err = easyPrompt(user)
		if err != nil {
			return err
		}
		pw = true
	}
	for _, host := range hosts {
		client, err := makeClient(user, pass, host)
		if err != nil {
			return errors.Wrap(err, "failed to dial host")
		}

		if err := executeCommand(client, user, pass, host, command, pw); err != nil {
			return errors.Wrapf(err, "failed to run %s on %s", command, host)
		}
		fmt.Println("")
	}

	return nil
}

func substitute(stdin []string, substitutions map[string]string) []string {
	replaced := []string{}
	for _, line := range stdin {
		for old, new := range substitutions {
			line = strings.Replace(line, old, new, -1)
		}
		replaced = append(replaced, line)
	}
	return replaced
}

func combine(stdin []string) string {
	var b bytes.Buffer
	for _, line := range stdin {
		line = strings.TrimSpace(line)
		b.WriteString(line)
		b.WriteString("\n")
	}
	return b.String()
}

func executeScriptfile(client *ssh.Client, user, pass, host string, sf scriptfile) error {
	color.Magenta(fmt.Sprintf("--- %s ---", host))

	for _, script := range sf.scripts {
		if err := executeScript(client, user, pass, host, script); err != nil {
			return err
		}
	}

	return nil
}

func executeCommand(client *ssh.Client, user, pass, host, command string, pw bool) error {
	color.Magenta(fmt.Sprintf("--- %s ---", host))

	sc := script{command: command}
	if pw {
		sc.stdin = []string{"PASSWORD"}
	}

	return executeScript(client, user, pass, host, sc)
}

func executeScript(client *ssh.Client, user, pass, host string, sc script) error {
	color.Yellow("executing command `%s`\n", sc.command)

	session, err := client.NewSession()
	if err != nil {
		return errors.Wrap(err, "asdf")
	}

	stdin := combine(substitute(sc.stdin, map[string]string{
		"PASSWORD": pass,
	}))

	session.Stdin = strings.NewReader(stdin)

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		return errors.Wrap(err, "request pty failed")
	}

	bs, err := session.CombinedOutput(sc.command)

	// print the output regardless of err
	output := strings.TrimSpace(string(bs))
	if len(output) == 0 {
		color.Magenta("<no output>")
	} else {
		color.Blue(output)
	}

	return err
}

func makeClient(user, pass, host string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            newSSHAuth(user, pass),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	address := fmt.Sprintf("%s:22", host)
	return ssh.Dial("tcp", address, config)
}

func newSSHAuth(user, pass string) []ssh.AuthMethod {
	authMethods := make([]ssh.AuthMethod, 0)
	sshAgent := sshAgentAuth()
	if sshAgent != nil {
		authMethods = append(authMethods, sshAgent)
	}

	if pass == "" {
		authMethods = append(authMethods, PasswordCallback(user))
	} else {
		authMethods = append(authMethods, ssh.Password(pass))
	}
	return authMethods
}
