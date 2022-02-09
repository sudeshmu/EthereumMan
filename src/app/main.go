package main

import (
	"etherman/src/config"
	"etherman/src/greet"
	"etherman/src/logger"
	"fmt"
)

func main() {
	fmt.Println("Welcome")
	greet.Hello()
	logger.InfoLogger.Print("Test")
	fmt.Print(config.Users())
}
