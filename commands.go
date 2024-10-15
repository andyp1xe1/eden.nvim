package main

import (
	"fmt"
)

type CommandAction string
const (
    Navigate  CommandAction = "navigate"
    ScrollTo                = "scrollTo"
    Build                   = "build"
    Reload                  = "reload"
)

type Command struct {
	Action CommandAction
	Value  string
}

var (
	JsTable  = map[CommandAction]string{ 
	  Navigate: "",
		ScrollTo: "",
		Reload: "",
	}
)

func parseCommand(message string) (Command, error) {
	var cmd Command
	if n, err := fmt.Sscanf(message, "%[^:]:%s", &cmd.Action, &cmd.Value); n < 2 || err != nil {
		return Command{}, fmt.Errorf("Invalid Command")
	}
	return cmd, nil
}
