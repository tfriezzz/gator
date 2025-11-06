package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/tfriezzz/gator/commands"
	"github.com/tfriezzz/gator/handleruser"
	"github.com/tfriezzz/gator/internal/config"
	"github.com/tfriezzz/gator/internal/database"
)

func main() {
	var testState config.State
	testConfig := config.Read()
	testState.Config = testConfig
	var testCommands commands.Commands
	testCommands.Handlers = make(map[string]func(*config.State, commands.Command) error)
	testCmd := commands.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	db, err := sql.Open("postgres", testConfig.DBURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	dbQueries := database.New(db)

	testState.DB = dbQueries
	testCommands.Register("login", handleruser.HandlerLogin)

	// err := testCommands.Run(&testState, testCmd)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err)
	// }

	testCommands.Register("register", handleruser.HandlerRegister)

	err = testCommands.Run(&testState, testCmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
