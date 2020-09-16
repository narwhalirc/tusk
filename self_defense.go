package tusk

import (
	"github.com/lrstanley/girc"
)

// This contains a basic self-defense mechanism.
// Current functionality: Re-joins when kicked and kicks the bad actor

// OnKick will handle when the bot gets kicked
func OnKick(c *girc.Client, e girc.Event) {
	m := ParseMessage(c, e) // Parse our message

	if e.Params[1] == Config.User { // If the bot is being kicked
		c.Cmd.Join(m.Channel)
		c.Cmd.Mode(m.Channel, "+o") // Attempt to op self
		c.Cmd.Reply(e, "Kick of bot detected. Enforcing countermeasure.")
		KickUser(c, e, m, m.Issuer)
	}
}
