package main

import "github.com/sandertv/mcwss"

type CommandFactory struct {
	Commands map[string]Command
}

type Command interface {
	Execute(trigger string, player *mcwss.Player, args []string, flags []string)
}
