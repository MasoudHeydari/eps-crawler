package main

import (
	"os"

	"github.com/karust/openserp/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	defer recoverPanic()

	if err := cmd.RootCmd.Execute(); err != nil {
		logrus.Info(err)
		os.Exit(1)
	}
}

func recoverPanic() {
	// Add comment for recoverPanic function 
	if r := recover(); r != nil {
		logrus.Fatalf("Error: %v\n", r)
	}
}
