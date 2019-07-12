package main

import "testing"

func Test_hosts(t *testing.T) {
	tests := []struct {
		input string
		exp   []string
	}{
		{
			input: "qa-control2",
			exp:   []string{"qa-control2"},
		},
		{
			input: "qa-control2, qa-control4,qa-control5",
			exp: []string{
				"qa-control2",
				"qa-control4",
				"qa-control5",
			},
		},
		{
			input: "qa-control2, qa-control{4..6},qa-control9, qa-control13",
			exp: []string{
				"qa-control2",
				"qa-control4",
				"qa-control5",
				"qa-control6",
				"qa-control9",
				"qa-control13",
			},
		},
		{
			input: "qa-control1.zombo.com",
			exp:   []string{"qa-control1.zombo.com"},
		},
		{
			input: "qa-control1.zombo.com,qa-control3.zombo.com, qa-control5.zombo.com",
			exp: []string{
				"qa-control1.zombo.com",
				"qa-control3.zombo.com",
				"qa-control5.zombo.com",
			},
		},
		{
			input: "qa-control{3..6}.zombo.com",
			exp: []string{
				"qa-control3.zombo.com",
				"qa-control4.zombo.com",
				"qa-control5.zombo.com",
				"qa-control6.zombo.com",
			},
		},
		{
			input: "qa-control1.zombo.com,qa-control{3..6}.zombo.com,qa-executor{1..3}",
			exp: []string{
				"qa-control1.zombo.com",
				"qa-control3.zombo.com",
				"qa-control4.zombo.com",
				"qa-control5.zombo.com",
				"qa-control6.zombo.com",
				"qa-executor1",
				"qa-executor2",
				"qa-executor3",
			},
		},
	}

	for _, test := range tests {
		expanded := hosts(test.input)
		if len(expanded) != len(test.exp) {
			t.Fatal("expected:", test.exp, "got:", expanded)
		}
		for i, host := range expanded {
			if host != test.exp[i] {
				t.Fatal("expected:", test.exp, "got:", expanded)
			}
		}
	}
}

func Test_expand(t *testing.T) {
	tests := []struct {
		raw string
		exp []string
	}{
		{

			raw: "qa-control3",
			exp: []string{"qa-control3"},
		},
		{

			raw: "qa-control{1..2}",
			exp: []string{
				"qa-control1",
				"qa-control2",
			},
		},
		{

			raw: "qa-control{5..8}",
			exp: []string{
				"qa-control5",
				"qa-control6",
				"qa-control7",
				"qa-control8",
			},
		},
	}

	for _, test := range tests {
		expanded := expand(test.raw)
		if len(expanded) != len(test.exp) {
			t.Fatal("expected:", test.exp, "got:", expanded)
		}
		for i, host := range expanded {
			if host != test.exp[i] {
				t.Fatal("expected:", test.exp, "got:", expanded)
			}
		}
	}
}
