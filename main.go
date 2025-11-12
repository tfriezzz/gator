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
	arg1 := ""
	arg2 := []string{}
	if len(os.Args) > 1 {
		arg1 = os.Args[1]
	}
	if len(os.Args) > 2 {
		arg2 = os.Args[2:]
	}
	testCmd := commands.Command{
		Name: arg1,
		Args: arg2,
	}
	// testCtx := context.Background()

	db, err := sql.Open("postgres", testConfig.DBURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	dbQueries := database.New(db)

	testState.DB = dbQueries
	testCommands.Register("login", handleruser.HandlerLogin)

	testCommands.Register("register", handleruser.HandlerRegister)

	testCommands.Register("reset", handleruser.HandlerReset)

	testCommands.Register("users", handleruser.HandlerList)

	testCommands.Register("agg", handleruser.HandlerAgg)

	testCommands.Register("addfeed", handleruser.HandlerFeed)

	err = testCommands.Run(&testState, testCmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// testFeed, _ := rss.FetchFeed(testCtx, "https://www.wagslane.dev/index.xml")
	// fmt.Printf("testFeed: %v", testFeed)
}
