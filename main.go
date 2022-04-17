package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sandertv/mcwss"
	"github.com/sandertv/mcwss/mctype"
	"github.com/sandertv/mcwss/protocol/event"
)

const PREFIX = "§f[§4MCTools§f] "

func main() {

	server := mcwss.NewServer(nil)

	positions := make(map[string][2]mctype.Position)
	changePositions := func(clientId string, newPositions [2]mctype.Position) {
		positions[clientId] = newPositions
	}
	getPositions := func(clientId string) ([2]mctype.Position, bool) {
		value, exist := positions[clientId]
		return value, exist
	}

	helpMessageBytes, err := os.ReadFile("help.txt")

	if err != nil {
		panic(err)
	}

	commandFactory := CommandFactory{Commands: map[string]Command{
		"1":       &PositionCommand{getPositions: getPositions, setPosition: changePositions},
		"2":       &PositionCommand{getPositions: getPositions, setPosition: changePositions},
		"fill":    &FillCommand{getPositions: getPositions},
		"replace": &ReplaceCommand{getPositions: getPositions},
		"help":    &MessageCommand{message: strings.ReplaceAll(string(helpMessageBytes), "\r", "")},
		"loop":    &LoopCommand{},
		"say":     &SayCommand{},
	}}

	server.OnConnection(func(player *mcwss.Player) {
		commandFactory.Commands["help"].Execute("help", player, []string{}, []string{})
		player.OnPlayerMessage(func(event *event.PlayerMessage) {
			if event.MessageType == "chat" {
				if strings.HasPrefix(event.Message, ";") {
					lexed, flags, err := HandleQuotes(strings.TrimPrefix(event.Message, ";"))
					if err != nil {
						player.SendMessage(err.Error())
						return
					}
					if value, exist := commandFactory.Commands[lexed[0]]; exist {
						value.Execute(lexed[0], player, lexed[1:], flags)
					} else {
						player.SendMessage("Unknown command: " + lexed[0])
					}
				}
			}
		})
	})

	server.OnDisconnection(func(player *mcwss.Player) {
		delete(positions, player.ClientID)
	})
	server.Run()
}

type PositionCommand struct {
	setPosition  func(clientId string, positions [2]mctype.Position)
	getPositions func(clientId string) ([2]mctype.Position, bool)
}

type FillCommand struct {
	getPositions func(clientId string) ([2]mctype.Position, bool)
}

type ReplaceCommand struct {
	getPositions func(clientId string) ([2]mctype.Position, bool)
}

type MessageCommand struct {
	message string
}

type LoopCommand struct{}

type SayCommand struct{}

type Tellraw struct {
	Text string `json:"text"`
}

func (command *FillCommand) Execute(trigger string, player *mcwss.Player, args []string, flags []string) {
	value, exist := command.getPositions(player.ClientID)

	if (!exist || value[0] == mctype.Position{} || value[1] == mctype.Position{}) {
		player.SendMessage(PREFIX + "You have not set both positions yet. Use ;1 to set the first position and ;2 to set the second position.")
		return
	}

	if len(args) < 1 {
		player.SendMessage(PREFIX + "Usage: ;fill <block> [meta] [fillMode: replace, keep, destroy, hollow, outline]")
		return
	}

	cmd := ("fill " +
		strconv.FormatFloat(value[0].X, 'f', 0, 64) + " " +
		strconv.FormatFloat(value[0].Y, 'f', 0, 64) + " " +
		strconv.FormatFloat(value[0].Z, 'f', 0, 64) + " " +
		strconv.FormatFloat(value[1].X, 'f', 0, 64) + " " +
		strconv.FormatFloat(value[1].Y, 'f', 0, 64) + " " +
		strconv.FormatFloat(value[1].Z, 'f', 0, 64) + " " +
		args[0])
	if len(args) > 1 {
		cmd += " " + args[1]
	}
	if len(args) > 2 {
		cmd += " " + args[2]
	}
	player.Exec(
		cmd,
		nil,
	)
	player.SendMessage(PREFIX + "Filled blocks")

}

func (command *ReplaceCommand) Execute(trigger string, player *mcwss.Player, args []string, flags []string) {
	value, exist := command.getPositions(player.ClientID)

	if (!exist || value[0] == mctype.Position{} || value[1] == mctype.Position{}) {
		player.SendMessage(PREFIX + "You have not set both positions yet. Use ;1 to set the first position and ;2 to set the second position.")
		return
	}

	if len(args) < 2 {
		player.SendMessage(PREFIX + "Usage: ;replace <newBlock> <newBlockMeta> <oldBlock> [oldBlockMeta] ")
		return
	}

	newMeta := "0"
	if len(args) > 3 {
		newMeta = args[3]
	}

	cmd := ("fill " +
		strconv.FormatFloat(value[0].X, 'f', 0, 64) + " " +
		strconv.FormatFloat(value[0].Y, 'f', 0, 64) + " " +
		strconv.FormatFloat(value[0].Z, 'f', 0, 64) + " " +
		strconv.FormatFloat(value[1].X, 'f', 0, 64) + " " +
		strconv.FormatFloat(value[1].Y, 'f', 0, 64) + " " +
		strconv.FormatFloat(value[1].Z, 'f', 0, 64) + " " +
		args[0] + " " +
		args[1] + " replace " +
		args[2] + " " +
		newMeta)

	player.Exec(
		cmd,
		nil,
	)
	player.SendMessage(PREFIX + "Filled blocks")

}

func (command *PositionCommand) Execute(trigger string, player *mcwss.Player, args []string, flags []string) {
	index := 0
	if trigger == "2" {
		index = 1
	}
	player.Position(func(position mctype.Position) {
		position.Y -= 2
		if value, exist := command.getPositions(player.ClientID); exist {
			value[index] = position
			command.setPosition(player.ClientID, value)
		} else {
			positions := [2]mctype.Position{{}, {}}
			positions[index] = position
			command.setPosition(player.ClientID, positions)
		}

		fmt.Println(position)
		player.SendMessage(
			PREFIX + "Set position " +
				strconv.Itoa(index+1) +
				" to " +
				strconv.FormatFloat(
					position.X, 'f', 0, 64,
				) + " " +
				strconv.FormatFloat(
					position.Y, 'f', 0, 64,
				) + " " +
				strconv.FormatFloat(
					position.Z, 'f', 0, 64,
				) + ".",
		)
	})

}

func (command *MessageCommand) Execute(trigger string, player *mcwss.Player, args []string, flags []string) {
	player.SendMessage(command.message)
}

func (command *LoopCommand) Execute(trigger string, player *mcwss.Player, args []string, flags []string) {
	if len(args) < 2 {
		player.SendMessage(PREFIX + "Usage: ;loop <times> <command>")
		return
	}
	var i int64 = 0
	max, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		player.SendMessage(PREFIX + "Usage: ;loop <times> <command>")
		return
	}
	for i < max {
		player.Exec(strings.Join(args[1:], " "), nil)
		i++
	}
	player.SendMessage(PREFIX + "Ran command " + args[0] + " times")
}

func (command *SayCommand) Execute(trigger string, player *mcwss.Player, args []string, flags []string) {

	bytes, err := json.Marshal(Tellraw{Text: strings.Join(args, " ")})
	if err != nil {
		panic(err)
	}
	cmd := "tellraw @a {\"rawtext\":[" +
		string(bytes) + "]}"
	player.Exec(
		cmd,
		nil,
	)
}
