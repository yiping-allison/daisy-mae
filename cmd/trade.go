package cmd

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/yiping-allison/isabelle/models"
)

const (
	// MaxTrade defines the max amount of trade events a user can make
	// at a time
	MaxTrade int = 2
)

// Trade will handle trade options within the server
//
// Trade is meant to be paired along with offer.go in order
// for members to trade and offer items among each other
func Trade(cmdInfo CommandInfo) {
	user := cmdInfo.Msg.Author
	if !cmdInfo.Service.Rep.Exists(user.ID) {
		// user does not exist in database
		// initialize user to 0
		rep := models.Rep{
			DiscordID: user.ID,
			RepNum:    0,
		}
		err := cmdInfo.Service.Rep.Create(&rep)
		if err != nil {
			fmt.Println("Trade(): error creating new user in database...")
			return
		}
	}

	// If the user currently doesn't exist in server tracking, make a new one
	if !cmdInfo.Service.User.UserExists(user) {
		// Create a user
		cmdInfo.Service.User.AddUser(user)
	}

	if cmdInfo.Service.User.LimitTrade(user.ID) {
		// Max trades created - can't make anymore
		return
	}

	// generate trade id
	id := generateID(1000, 9999)

	// TODO: Add trade tracking in models...
	if cmdInfo.Service.Trade.Exists(id) {
		// error - id exists
		return
	}
	cmdInfo.Service.Trade.AddTrade(id, user)

	// retrieve reps from database
	reps := cmdInfo.Service.Rep.GetRep(user.ID)

	// Print Trade Offer
	offer := strings.Join(cmdInfo.CmdOps[1:], " ")
	msg := cmdInfo.createMsgEmbed(
		"Trade Offer", tradeThumbURL, user.String(), tradeColor,
		format(
			createFields("Trade ID", id, true),
			createFields("Reputation", strconv.Itoa(reps), true),
			createFields("Offer", offer, false),
		),
	)
	cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
}

// generateID will come up with a pseudo random number with 4 digits
// and return it in string format
//
// This is used to generate trade IDs
func generateID(min, max int) string {
	id := min + rand.Intn(max-min)
	return strconv.Itoa(id)
}
