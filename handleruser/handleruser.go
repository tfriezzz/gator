// Package handleruser handles user specific actions
package handleruser

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	// "github.com/tfriezzz/gator/internal/commands"
	"github.com/google/uuid"
	"github.com/tfriezzz/gator/commands"
	"github.com/tfriezzz/gator/internal/config"
	"github.com/tfriezzz/gator/internal/database"
)

func requireOneArg(cmd commands.Command, usage string) {
	if len(cmd.Args) < 1 {
		fmt.Fprintln(os.Stderr, "usage:", usage)
		os.Exit(1)
	}
	// cmd.Name = os.Args[1]
	// cmd.Args = os.Args[2:]
	// if len(cmd.Name) == 0 {
	// 	fmt.Fprintln(os.Stderr, "username required")
	// 	os.Exit(1)
	// }
}

func HandlerLogin(s *config.State, cmd commands.Command) error {
	requireOneArg(cmd, "gator login <username>")
	// if len(cmd.Args) != 1 {
	// 	return fmt.Errorf("the login handler expects a single argument")
	// }
	userName := cmd.Args[0]
	_, err := s.DB.GetUser(context.Background(), userName)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Fprintln(os.Stderr, "user not found")
		os.Exit(1)
	}
	err = s.Config.SetUser(userName)
	if err != nil {
		return err
	}
	fmt.Printf("the user has been set to %s", userName)

	return nil
}

func HandlerRegister(s *config.State, cmd commands.Command) error {
	requireOneArg(cmd, "gator register <username>")
	if len(cmd.Args) == 0 {
		return fmt.Errorf("the register handler expects a single argument")
	}

	userID := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()
	userName := cmd.Args[0]

	userArgs := database.CreateUserParams{
		userID, createdAt,
		updatedAt, userName,
	}

	_, err := s.DB.GetUser(context.Background(), userName)
	if err == nil {
		fmt.Fprintln(os.Stderr, "user already exists")
		os.Exit(1)
	}

	user, err := s.DB.CreateUser(context.Background(), userArgs)
	if err != nil {
		return err
	}

	if err = s.Config.SetUser(user.Name); err != nil {
		os.Exit(1)
		return err
	}

	fmt.Printf("user %s was created: %+v", user.Name, user)
	return nil
}
