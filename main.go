package main

import (
	"Coins/cmd"
	"fmt"
)
import _ "Coins/cmd"

func main() {
	fmt.Println("main backend of Coins")
	cmd.Execute()
}
