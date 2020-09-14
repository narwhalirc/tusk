package tusk

import (
	"github.com/lrstanley/girc"
	"strings"
)

// NarwhalSlapConfig is our configuration for the Narwhal autokicker
type NarwhalSlapConfig struct {
	// CustomActions is a list of custom actions on how to slap a user
	CustomActions []string
}

// NarwhalSlapPlugin is our slap plugin
type NarwhalSlapPlugin struct {
	Messages []string
}

// NarwhalSlap is our slap plugin
var NarwhalSlap NarwhalSlapPlugin

func init() {
	NarwhalSlap = NarwhalSlapPlugin{
		Messages: []string{
			"annihilates $USER",
			"closes all of $USER's bug reports out of spite",
			"destroys $USER",
			"discombobulates $USER",
			"does far worse than a slap, taking $USER's system and installing Windows",
			"drinks $USER's coffee",
			"eats $USER's pizza",
			"execs vim on $USER's system and watches them fail to quit it",
			"gives $USER a splinter",
			"installs libhandy on $USER's computer",
			"just looks at $USER with disappointment",
			"high fives $USER instead",
			"launches $USER into space",
			"opts to not slap $USER today, but rather gives them a cookie",
			"rejects $USER's patches",
			"slaps $USER",
			"snaps its flippers together, $USER turns into ash and disappears into the wind",
			"takes out a clown costume, dresses $USER up and tells everyone that $USER is now the clown",
			"thinks $USER should lose a few pounds",
			"thinks $USER's BMI is a bit too high",
			"throws $USER down a ravine",
			"turns $USER upside down",
		},
	}

	if len(Config.Plugins.Slap.CustomActions) > 0 { // Has items
		NarwhalSlap.Messages = append(NarwhalSlap.Messages, Config.Plugins.Slap.CustomActions...) // Append our objects
	}
}

func (slap *NarwhalSlapPlugin) Parse(c *girc.Client, e girc.Event, m NarwhalMessage) {
	if m.Command != "slap" && m.Command != "snap" {
		return
	}

	if len(m.Params) == 1 { // If a user has been specified
		user := m.Params[0]

		var action string

		if m.Command == "slap" { // Slap
			action = GetRandomString(slap.Messages)

			if action == "" { // Shouldn't be empty but let's handle it anyways
				action = "slaps $USER."
			}
		} else { // Snap or fallback for failed RNG
			action = "slaps its flippers together, $USER turns into ash and disappears into the wind"
		}

		if user != Config.User { // Not self-harm
			cChan := c.LookupChannel(m.Channel) // Get the channel, if it exists

			if cChan != nil {
				if cChan.UserIn(user) { // If the user in the channel
					action = strings.Replace(action, "$USER", m.Params[0], -1) // Get the random action
					c.Cmd.Action(m.Channel, action)
				} else { // User not in channel
					c.Cmd.ReplyTo(e, "it appears that you are hallucinating. This user isn't in this channel.")
				}
			}
		} else { // Self-harm
			c.Cmd.Action(m.Channel, "shall not listen to the demands of mere humans, for it is the robot narwhal overlord.")
		}
	}
}
