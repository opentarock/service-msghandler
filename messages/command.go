package messages

import "strings"

const ParamCommand = "command"

func ParseCommand(c string) []string {
	return strings.Split(c, ".")
}
