package tusk

import (
	"fmt"
	"github.com/lrstanley/girc"
)

// This contains a basic self-defense mechanism.
// Current functionality: Re-joins when kicked and kicks the bad actor

// OnKick will handle when the bot gets kicked
func OnKick(c *girc.Client, e girc.Event) {
	m := ParseMessage(c, e) // Parse our message
	fmt.Printf("Kicked: %v\n", m)
}
