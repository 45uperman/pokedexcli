package main

import (
	"bufio"
	"os"
	"testing"
)

func init() {
	supportedCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
	}
}
func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input: "    diglet  DIG!   diglet DIG!  TRIO!   TRIO!!   TRIO!!!     ",
			expected: []string{
				"diglet", "dig!", "diglet", "dig!", "trio!", "trio!!", "trio!!!",
			},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf(
				"Output of unexpected length\nEXPECTED: %v\nACTUAL: %v",
				c.expected,
				actual,
			)
			t.FailNow()
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf(
					"Found incorrect element!\nEXPECTED: %v\nACTUAL: %v",
					c.expected,
					actual,
				)
				t.FailNow()
			}
		}
	}
}

func TestCommands(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input: "help",
			expected: []string{
				"Welcome to the Pokedex!",
				"Usage:",
				"",
				"exit: Exit the Pokedex",
				"help: Displays a help message",
			},
		},
	}

	for _, c := range cases {
		var actual []string
		r, w, _ := os.Pipe()
		os.Stdout = w
		supportedCommands[c.input].callback()
		w.Close()
		scanner := bufio.NewScanner(r)
		i := 0
		for scanner.Scan() {
			actual = append(actual, scanner.Text())
			if c.expected[i] != actual[i] {
				t.Errorf(
					"Incorrect help message\nEXPECTED: %v\nACTUAL: %v",
					c.expected,
					actual,
				)
				t.FailNow()
			}
			i++
		}
		if err := scanner.Err(); err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
}
