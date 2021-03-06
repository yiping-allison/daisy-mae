package cmd

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Unregister allows a queue user to remove themselves from the queue
func Unregister(cmdInfo CommandInfo) {
	// cmdInfo.CmdOps[1:] starts after ;unregister
	if len(cmdInfo.CmdOps) != 3 {
		// Error - not enough arguments
		msg := cmdInfo.createMsgEmbed(
			"Error: Wrong Arguments", errThumbURL, "Try checking your syntax.", errColor,
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"unregister event 1234", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.BotChID, msg)
		return
	}
	args := cmdInfo.CmdOps[1:]
	switch strings.ToLower(args[0]) {
	case "event":
		cmdInfo.removeFromEvent(args[1], cmdInfo.Msg.Author)
	case "trade":
		cmdInfo.removeFromTrade(args[1], cmdInfo.Msg.Author)
	}
}

// helper func which removes queue users from an event queue
func (c CommandInfo) removeFromEvent(eventID string, user *discordgo.User) {
	if !c.Service.Event.EventExists(eventID) {
		// event does not exist
		msg := c.createMsgEmbed(
			"Error: Event does not exist", errThumbURL, "Event ID: "+eventID, errColor,
			format(
				createFields("Suggestion", "Try checking if you supplied the correct Event ID", false),
			))
		c.Ses.ChannelMessageSendEmbed(c.BotChID, msg)
		return
	}

	// remove user
	c.Service.Event.Remove(eventID, user)
	// Remove tracking on user
	c.Service.User.RemoveQueue(eventID, user)

	// successfully removed user
	msg := c.createMsgEmbed(
		"Removed from Event", checkThumbURL, "Queue ID: "+c.CmdOps[2],
		successColor, format(
			createFields("User", user.Mention(), true),
			createFields("Suggestion", "Feel free to queue for any other events or create your own.", false),
		))
	c.Ses.ChannelMessageSendEmbed(c.BotChID, msg)
}

// helper func to remove user's offer from trade event
func (c CommandInfo) removeFromTrade(tradeID string, user *discordgo.User) {
	if !c.Service.Trade.Exists(tradeID) {
		// trade does not exist
		msg := c.createMsgEmbed(
			"Error: Trade does not exist", errThumbURL, "Trade ID: "+tradeID, errColor,
			format(
				createFields("Suggestion", "Try checking if you supplied the correct Trade ID", false),
			))
		c.Ses.ChannelMessageSendEmbed(c.BotChID, msg)
		return
	}

	offer := c.Service.Trade.GetOffer(tradeID, user.ID)
	// remove user
	c.Service.Trade.Remove(tradeID, user)
	// remove tracking on user
	c.Service.User.RemoveOffer(tradeID, user)

	// successfully removed user
	msg := c.createMsgEmbed(
		"User Withdrew From Trade", checkThumbURL, "Trade ID: "+tradeID,
		successColor, format(
			createFields("User", user.Mention(), true),
			createFields("Offer", offer, true),
			createFields("Suggestion", "Feel free to offer for any other trades or create your own.", false),
		))
	c.Ses.ChannelMessageSendEmbed(c.BotChID, msg)
}
