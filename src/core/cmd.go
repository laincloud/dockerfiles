package core

import (
	"fmt"
)

// Command denotes build/pull/push/retag
type Command int

const (
	// Build denotes build images
	Build Command = iota
	// Pull denotes pull images
	Pull
	// Push denotes push images
	Push
	// Retag denotes retag images
	Retag
	// RetagSingle denotes retag an image
	RetagSingle
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
	case RetagSingle:
		return "retag-single"
	default:
		panic(fmt.Sprintf("unexpected Command %v", c))
	}
}
