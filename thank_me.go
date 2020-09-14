package tusk

import (
	"github.com/lrstanley/girc"
	"strings"
)

// NarwhalThankMePlugin is our Thank Me plugin
// The goal is just to reply with something random like "No problem" when thanked
type NarwhalThankMePlugin struct {
	Messages []string
}

// NarwhalThank is our thank me plugin
var NarwhalThank NarwhalThankMePlugin

func init() {
	NarwhalThank = NarwhalThankMePlugin{
		Messages: []string{
			"No problem.",
			"You're welcome.",
			"Anytime.",
			"I'm just a bot, no need to thank me.",
			"I'm just doing what I can.",
			"Sure thing.",
		},
	}
}

func (thank *NarwhalThankMePlugin) Parse(c *girc.Client, e girc.Event, m NarwhalMessage) {
	if !strings.HasPrefix(m.Message, Config.User) { // If our message doesn't start with a name of our bot
		return // We don't care about it
	}

	strippedMessage := strings.Trim(m.Message, Config.User) // Remove our username reference
	strippedMessage = strings.Trim(strippedMessage, ",")    // Remove any starting ,
	strippedMessage = strings.Trim(strippedMessage, ":")    // Remove any : sometimes used by people for ref
	strippedMessage = strings.TrimSpace(strippedMessage)    // Remove any whitespace
	strippedMessage = strings.ToLower(strippedMessage)      // Lowercase to make search simpler

	if strings.HasPrefix(strippedMessage, "ty") || // Starts with ty
		strings.HasPrefix(strippedMessage, "thanks") || // thanks
		strings.HasPrefix(strippedMessage, "thank you") || // thank you
		strings.HasPrefix(strippedMessage, "kiitos") { // Torille!
		message := GetRandomString(thank.Messages) // Get a random thank you
		c.Cmd.ReplyTo(e, message)
	}
}
