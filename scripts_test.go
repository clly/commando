package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const file1 = `
sudo whoami
PASSWORD
`

const file2 = `

herp derp
PASSWORD
foo
PASSWORD

`

const file3 = `
echo alpha
---
sudo whoami
PASSWORD
---
echo beta
bar
PASSWORD
`

const file4 = `
# comment1
echo alpha
---
# comment 2
sudo whoami
PASSWORD
#comment3
#comment4
PASSWORD
---
# comment 0
# comment 1
# comment 2

echo beta
# comment 3
bar
PASSWORD

# comment 4
`

func Test_parseScriptWithSudo(t *testing.T) {
	tests := []struct {
		content    string
		name       string
		expScripts []script
	}{
		{
			content: file1,
			name:    "0-script1",
			expScripts: []script{
				{
					command: "sudo whoami",
					stdin:   []string{"PASSWORD"},
				},
			},
		},
		{
			content: file2,
			name:    "1-script2",
			expScripts: []script{
				{
					command: "herp derp",
					stdin:   []string{"PASSWORD", "foo", "PASSWORD"},
				},
			},
		},
		{
			content: file3,
			name:    "2-script3",
			expScripts: []script{
				{
					command: "echo alpha",
					stdin:   []string{},
				},
				{
					command: "sudo whoami",
					stdin:   []string{"PASSWORD"},
				},
				{
					command: "echo beta",
					stdin:   []string{"bar", "PASSWORD"},
				},
			},
		},
		{
			content: file4,
			name:    "3-script4",
			expScripts: []script{
				{
					command: "echo alpha",
					stdin:   []string{},
				},
				{
					command: "sudo whoami",
					stdin:   []string{"PASSWORD", "PASSWORD"},
				},
				{
					command: "echo beta",
					stdin:   []string{"bar", "PASSWORD"},
				},
			},
		},
	}

	for _, test := range tests {
		scriptFile, err := parse(test.name, test.content)
		require.NoError(t, err)
		require.Equal(t, test.name, scriptFile.name)
		require.Equal(t, len(test.expScripts), len(scriptFile.scripts))
		require.True(t, scriptFile.sudo, "Failed to parse whether or not we need a password")
		for i := 0; i < len(test.expScripts); i++ {
			expScript := test.expScripts[i]
			script := scriptFile.scripts[i]
			require.Equal(t, expScript.command, script.command)
			require.Equal(t, len(expScript.stdin), len(script.stdin))
			for j := 0; j < len(expScript.stdin); j++ {
				require.Equal(t, expScript.stdin[j], script.stdin[j])
			}
		}
	}
}

const file5 = `
whoami
`

const file6 = `
echo alpha
---
whoami
---
echo beta
bar
whatup
`

func Test_parseScriptWithoutSudo(t *testing.T) {
	tests := []struct {
		content    string
		name       string
		expScripts []script
	}{
		{
			content: file5,
			name:    "1-script5-no-sudo",
			expScripts: []script{
				{
					command: "whoami",
					stdin:   []string{},
				},
			},
		},
		{
			content: file6,
			name:    "2-script6-no-sudo",
			expScripts: []script{
				{
					command: "echo alpha",
					stdin:   []string{},
				},
				{
					command: "whoami",
					stdin:   []string{},
				},
				{
					command: "echo beta",
					stdin:   []string{"bar", "whatup"},
				},
			},
		},
	}

	for _, test := range tests {
		scriptFile, err := parse(test.name, test.content)
		require.NoError(t, err)
		require.Equal(t, test.name, scriptFile.name)
		require.Equal(t, len(test.expScripts), len(scriptFile.scripts))
		require.False(t, scriptFile.sudo, "%q %q %q", scriptFile.scripts, scriptFile.sudo, test.name)
		for i := 0; i < len(test.expScripts); i++ {
			expScript := test.expScripts[i]
			script := scriptFile.scripts[i]
			require.Equal(t, expScript.command, script.command)
			require.Equal(t, len(expScript.stdin), len(script.stdin))
			for j := 0; j < len(expScript.stdin); j++ {
				require.Equal(t, expScript.stdin[j], script.stdin[j])
			}
		}
	}
}
