package common

import (
	"strings"
)

type RegisterCommand struct {
	ServiceName string
	ServiceUrl  string
}

func CreateRegisterCommand(name, url string) string {
	return "R " + name + " " + url
}

func ParseCommand(command string) interface{} {
	split := strings.Split(command, " ")

	if split[0] == "R" {
		return RegisterCommand{split[1], split[2]}
	}

	return nil
}
