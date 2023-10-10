package main

import (
	"fmt"

	"github.com/masonkmeyer/jet/cmd"
)

func main() {
	// it is useful to show the user a command that resulted in an exit.
	// for example, if the user tries to checkout a branch but has uncommitted changes
	exitOutput := <-cmd.Execute()

	fmt.Println()
	fmt.Println(exitOutput)
	fmt.Println()
}
