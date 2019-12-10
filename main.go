package main

import (
	"github.com/fidellr/edu_malay/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.RootCMD.Execute(); err != nil {
		logrus.Fatalf("Fail init root cmd with error : %s", err.Error())
	}
}
