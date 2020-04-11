package isabellebot

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/yiping-allison/isabelle/cmd"
	"github.com/yiping-allison/isabelle/models"
)

// Bot represents a daisymae bot instance
type Bot struct {
	// AdminRole contains the role ID of the discord server's admin role
	AdminRole string

	// DS represents the bot's current discord session
	DS *discordgo.Session

	// Service is the gateway to all Service interactions
	Service models.Services

	// Prefix is the user set bot prefix found in .config (default is ?)
	Prefix string

	// Commands is a map of command name to a closure of the func
	Commands map[string]Command
}

// New creates a new daisymae bot instance and loads bot commands.
//
// It will return the finished bot and nil upon success or
// empty bot and err upon failure
func New(botKey, admin string) (*Bot, error) {
	if botKey == "" {
		return nil, errors.New("isabellebot: you need to input a botKey in the .config file")
	}
	discord, err := discordgo.New("Bot " + botKey)
	if err != nil {
		return nil, errors.New("isabellebot: error connecting to discord")
	}
	// Commands Setup
	cmds := make(map[string]Command, 0)
	isa := &Bot{
		AdminRole: admin,
		Prefix:    "?",
		DS:        discord,
		Service:   models.Services{},
		Commands:  cmds,
	}
	isa.compileCommands()
	// Add Handlers
	isa.DS.AddHandler(isa.ready)
	isa.DS.AddHandler(isa.handleMessage)
	return isa, nil
}

// ready will update bot status after bot receives "ready" event from
// discord
func (b *Bot) ready(s *discordgo.Session, r *discordgo.Ready) {
	s.UpdateStatus(0, b.Prefix+"list")
}

// handleMessage handles all new discord messages which the bot uses to determine
// actions
func (b *Bot) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, b.Prefix) {
		b.processCmd(s, m)
	}
}

// Command represents a discord bot command
type Command struct {
	Cmd func(cmd.CommandInfo)
}

// processCmd attemps to process any string that is prefixed with bot notifier
//
// Valid commands will be run while invalid commands will be ignored
//
// Example bot commands:
//
// ?search
//
// ?search help
func (b *Bot) processCmd(s *discordgo.Session, m *discordgo.MessageCreate) {
	cmds := regexp.MustCompile("\\s+").Split(m.Content[len(b.Prefix):], -1)
	trim := strings.TrimPrefix(cmds[0], b.Prefix)
	res := b.find(trim)
	if res == nil {
		return
	}
	var commands []string
	for key := range b.Commands {
		commands = append(commands, key)
	}
	ci := cmd.CommandInfo{
		AdminRole: b.AdminRole,
		Ses:       s,
		Msg:       m,
		Service:   b.Service,
		Prefix:    b.Prefix,
		CmdName:   trim,
		CmdOps:    cmds,
		CmdList:   commands,
	}
	// Run command
	res.Cmd(ci)
}

// finds a command in the command map
//
// If it exists, it returns the Command
// If not, it returns nil
func (b *Bot) find(name string) *Command {
	if val, ok := b.Commands[name]; ok {
		return &val
	}
	return nil
}

// compileCommands contains all commands the bot should add to the bot command map
func (b *Bot) compileCommands() {
	b.addCommand("search", cmd.Search)
	b.addCommand("help", cmd.Help)
	b.addCommand("list", cmd.List)
	b.addCommand("event", cmd.Event)
	b.addCommand("queue", cmd.Queue)
	b.addCommand("close", cmd.Close)
	b.addCommand("unregister", cmd.Unregister)
	b.addCommand("trade", cmd.Trade)
	b.addCommand("offer", cmd.Offer)
	b.addCommand("rep", cmd.Rep)
	b.addCommand("accept", cmd.Accept)
	b.addCommand("reject", cmd.Reject)
}

// utility func to add command to bot command map
func (b *Bot) addCommand(name string, cmd func(cmd.CommandInfo)) {
	if _, ok := b.Commands[name]; ok {
		fmt.Printf("addCommand: %s already exists in the map\n", name)
		return
	}
	command := Command{
		Cmd: cmd,
	}
	b.Commands[name] = command
}

// SetPrefix sets user directed bot prefix from .config
func (b *Bot) SetPrefix(newPrefix string) {
	b.Prefix = newPrefix
}
