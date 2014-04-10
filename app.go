package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// App is the main structure of a cli application.
// New App variables should be created with the cli.NewApp() function.
type App struct {

	// Name of the program. Defaults to os.Args[0]
	Name string

	// Description of the program.
	Usage string

	// Version of the program
	Version string

	// List of commands to execute
	Commands []Command

	// List of flags to parse
	Flags []Flag

	// An action to execute before any subcommands are run, but after the context is ready.
	// If a non-nil error is returned, no subcommands are run.
	Before func(context *Context) error

	// The action to execute when no subcommands are specified
	Action func(context *Context)

	// Compilation date
	Compiled time.Time

	// Author
	Author string

	// Author e-mail
	Email string
}

// compileTime tries to find out when this binary was compiled.
// Returns the current time if it fails to find it.
func compileTime() time.Time {
	info, err := os.Stat(os.Args[0])
	if err != nil {
		return time.Now()
	}
	return info.ModTime()
}

// NewApp creates a new cli Application with some reasonable defaults for Name, Usage, Version and Action.
func NewApp() *App {
	return &App{
		Name:     os.Args[0],
		Usage:    "A new cli application",
		Version:  "0.0.0",
		Action:   helpCommand.Action,
		Compiled: compileTime(),
		Author:   "Author",
		Email:    "unknown@email",
	}
}

// Run provides an entry point to the cli app.
// It parses the slice of arguments and routes to the proper flag/args combination.
func (a *App) Run(arguments []string) error {

	// append help to commands
	if a.Command(helpCommand.Name) == nil {
		a.Commands = append(a.Commands, helpCommand)
	}

	// append version/help flags
	a.appendFlag(BoolFlag{"version, v", "print the version"})
	a.appendFlag(BoolFlag{"help, h", "show help"})

	// parse flags
	set := flagSet(a.Name, a.Flags)
	set.SetOutput(ioutil.Discard)
	err := set.Parse(arguments[1:])
	nerr := normalizeFlags(a.Flags, set)
	if nerr != nil {
		fmt.Println(nerr)
		context := NewContext(a, set, set)
		ShowAppHelp(context)
		fmt.Println()
		return nerr
	}
	context := NewContext(a, set, set)

	if err != nil {
		fmt.Println("Incorrect Usage.")
		fmt.Println()
		ShowAppHelp(context)
		fmt.Println()
		return err
	}

	if checkHelp(context) {
		return nil
	}

	if checkVersion(context) {
		return nil
	}

	if a.Before != nil {
		err := a.Before(context)
		if err != nil {
			return err
		}
	}

	args := context.Args()
	if args.Present() {
		name := args.First()
		c := a.Command(name)
		if c != nil {
			return c.Run(context)
		}
	}

	// Run default Action
	a.Action(context)

	return nil
}

// Command returns the named command on App. If the command does not exist, nil is returned.
func (a *App) Command(name string) *Command {
	for _, c := range a.Commands {
		if c.HasName(name) {
			return &c
		}
	}
	return nil
}

// hasFlag checks for the presence of a flag.
func (a *App) hasFlag(flag Flag) bool {
	for _, f := range a.Flags {
		if flag == f {
			return true
		}
	}
	return false
}

// appendFlag appends a flag if it does not already exist.
func (a *App) appendFlag(flag Flag) {
	if !a.hasFlag(flag) {
		a.Flags = append(a.Flags, flag)
	}
}
