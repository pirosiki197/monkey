package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/pirosiki197/monkey/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Hello, %s! This is the Monkey programming language!\n", user.Username)
	fmt.Println("Feel free to type in commands")
	repl.Start(os.Stdin, os.Stdout)
}
