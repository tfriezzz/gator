package main

import (
	"fmt"

	"github.com/tfriezzz/gator/internal/config"
)

func main() {
	config.Read()
	config.Read().SetUser("tfry")
	fmt.Printf("%s", config.Read())
}
