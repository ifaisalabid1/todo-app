package main

import (
	"fmt"
	"os"

	"github.com/ifaisalabid1/todo-app/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(cfg)
}
