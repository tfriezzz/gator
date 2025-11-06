package main

import (
	"fmt"
	"os"

	// "github.com/tfriezzz/gator/handleruser"
	"github.com/tfriezzz/gator/commands"
	"github.com/tfriezzz/gator/handleruser"
	"github.com/tfriezzz/gator/internal/config"
)

func main() {
	var testState config.State
	testConfig := config.Read()
	testState.Config = testConfig
	var testCommands commands.Commands
	testCommands.Handlers = make(map[string]func(*config.State, commands.Command) error)
	var testCmd commands.Command
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "no args")
		os.Exit(1)
	}
	testCmd.Name = os.Args[1]
	testCmd.Args = os.Args[2:]
	if len(testCmd.Args) < 1 {
		fmt.Fprintln(os.Stderr, "username required")
		os.Exit(1)
	}
	testCommands.Register("login", handleruser.HandlerLogin)
	// if err != nil {
	// 	fmt.Errorf("%w", err)
	// }
	err := testCommands.Run(&testState, testCmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	// config.Read().SetUser("tfry")
	// fmt.Printf("%s", config.Read())
}
