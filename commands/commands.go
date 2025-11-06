// Package commands handles routing
package commands

import (
	"fmt"

	"github.com/tfriezzz/gator/internal/config"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Handlers map[string]func(*config.State, Command) error
}

func (c *Commands) Run(s *config.State, cmd Command) error {
	Handler, ok := c.Handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command")
	}
	return Handler(s, cmd)
}

func (c *Commands) Register(name string, f func(*config.State, Command) error) error {
	_, ok := c.Handlers[name]
	if ok {
		return fmt.Errorf("command: '%s' already exists", name)
	}
	c.Handlers[name] = f
	return nil
}
