// The packer package contains the core components of Packer.
package packer

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

// The environment struct contains all the state necessary for a single
// instance of Packer.
//
// It is *not* a singleton, but generally a single environment is created
// when Packer starts running to represent that Packer run. Technically,
// if you're building a custom Packer binary, you could instantiate multiple
// environments and run them in parallel.
type Environment struct {
	builderFactory BuilderFactory
	command map[string]Command
	ui      Ui
}

// This struct configures new environments.
type EnvironmentConfig struct {
	BuilderFactory BuilderFactory
	Command map[string]Command
	Ui      Ui
}

// DefaultEnvironmentConfig returns a default EnvironmentConfig that can
// be used to create a new enviroment with NewEnvironment with sane defaults.
func DefaultEnvironmentConfig() *EnvironmentConfig {
	config := &EnvironmentConfig{}
	config.BuilderFactory = new(NilBuilderFactory)
	config.Ui = &ReaderWriterUi{os.Stdin, os.Stdout}
	return config
}

// This creates a new environment
func NewEnvironment(config *EnvironmentConfig) (env *Environment, err error) {
	if config == nil {
		err = errors.New("config must be given to initialize environment")
		return
	}

	env = &Environment{}
	env.builderFactory = config.BuilderFactory
	env.command = make(map[string]Command)
	env.ui = config.Ui

	for k, v := range config.Command {
		env.command[k] = v
	}

	// TODO: Should "version" be allowed to be overriden?
	if _, ok := env.command["version"]; !ok {
		env.command["version"] = new(versionCommand)
	}

	return
}

// Returns the BuilderFactory associated with this Environment.
func (e *Environment) BuilderFactory() BuilderFactory {
	return e.builderFactory
}

// Executes a command as if it was typed on the command-line interface.
// The return value is the exit code of the command.
func (e *Environment) Cli(args []string) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		e.PrintHelp()
		return 1
	}

	command, ok := e.command[args[0]]
	if !ok {
		// The command was not found. In this case, let's go through
		// the arguments and see if the user is requesting the version.
		for _, arg := range args {
			if arg == "--version" || arg == "-v" {
				command = e.command["version"]
				break
			}
		}

		// If we still don't have a command, show the help.
		if command == nil {
			e.PrintHelp()
			return 1
		}
	}

	return command.Run(e, args[1:])
}

// Prints the CLI help to the UI.
func (e *Environment) PrintHelp() {
	// Created a sorted slice of the map keys and record the longest
	// command name so we can better format the output later.
	commandKeys := make([]string, len(e.command))
	i := 0
	maxKeyLen := 0
	for key, _ := range e.command {
		commandKeys[i] = key
		if len(key) > maxKeyLen {
			maxKeyLen = len(key)
		}

		i++
	}

	// Sort the keys
	sort.Strings(commandKeys)

	e.ui.Say("usage: packer [--version] [--help] <command> [<args>]\n\n")
	e.ui.Say("Available commands are:\n")
	for _, key := range commandKeys {
		command := e.command[key]

		// Pad the key with spaces so that they're all the same width
		key = fmt.Sprintf("%v%v", key, strings.Repeat(" ", maxKeyLen-len(key)))

		// Output the command and the synopsis
		e.ui.Say("    %v     %v\n", key, command.Synopsis())
	}
}

// Returns the UI for the environment. The UI is the interface that should
// be used for all communication with the outside world.
func (e *Environment) Ui() Ui {
	return e.ui
}
