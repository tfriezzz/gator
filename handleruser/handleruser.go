// Package handleruser handles user specific actions
package handleruser

import (
	"fmt"
	// "github.com/tfriezzz/gator/internal/commands"
	"github.com/tfriezzz/gator/commands"
	"github.com/tfriezzz/gator/internal/config"
)

func HandlerLogin(s *config.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("the login handler expects a single argument")
	}
	username := cmd.Args[0]
	err := s.Config.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("the user has been set to %s", username)

	return nil
}
