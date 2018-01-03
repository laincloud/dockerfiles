package util

import (
	"fmt"
)

// Command denotes build/pull/push/retag
type Command int

const (
	// Build denotes build image
	Build Command = iota
	// Pull denotes pull image
	Pull
	// Push denotes push image
	Push
	// Retag denotes retag image
	Retag
)

// ShowCommand return string representation of Command
func ShowCommand(c Command) string {
	switch c {
	case Build:
		return "build"
	case Pull:
		return "pull"
	case Push:
		return "push"
	case Retag:
		return "retag"
	default:
		panic(fmt.Sprintf("unexpected Command %v", c))
	}
}
