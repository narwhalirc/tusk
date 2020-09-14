package tusk

import (
	"fmt"
	"github.com/JoshStrobl/trunk"
	"github.com/lrstanley/girc"
	"strconv"
)

// adminCommands is a list of admin commands
var adminCommands []string

// fixedIgnoreUsers is a list of users we'll ignore messages from no matter what
var fixedIgnoreUsers []string

// ignoreCommands is a list of numerical command codes to ignore
var ignoreCommands []string

func init() {
	ignoreCommands = []string{
		"002", // RPL_YOURHOST
		"003", // RPL_CREATED
		"004", // RPL_MYINFO
		"005", // RPL_BOUNCE
		"251", // RPL_LUSERCLIENT
		"252", // RPL_LUSEROP
		"253", // RPL_LUSERUNKNOWN
		"254", // RPL_LUSERCHANNELS
		"255", // RPL_LUSERME
		"265", // RPL_LOCALUSERS
		"266", // RPL_GLOBALUSERS
		"331", // RPL_NOTOPIC
		"332", // RPL_TOPIC
		"333", // RPL_TOPICWHOTIME
		"372", // RPL_MOTD
		"375", // RPL_MOTDSTART
		"376", // RPL_ENDOFMOTD
	}
}

// OnConnected will handle connection to an IRC network
func OnConnected(c *girc.Client, e girc.Event) {
	trunk.LogSuccess("Successfully connected to " + Config.Network + " as " + Config.User)
	if len(Config.Channels) > 0 { // If we have channels set to join
		for _, channel := range Config.Channels { // For each channel to join
			c.Cmd.Join(channel)
			c.Cmd.Mode(channel, "+o") // Attempt to op self
			trunk.LogInfo("Joining " + channel)
		}
	}
}

// OnInvite will handle a request to invite an IRC channel
func OnInvite(c *girc.Client, e girc.Event) {
	msg := ParseMessage(c, e)              // Parse our message
	clientUser := c.LookupUser(msg.Issuer) // Attempt to look up the user

	if clientUser == nil { // Failed to look up the user
		trunk.LogErr("Failed to look up user: " + msg.Issuer)
		return
	}

	fmt.Println("%v", clientUser)

	channel := c.LookupChannel(msg.Message)
	fmt.Println("%v", channel)

	channelPerms, permsOk := clientUser.Perms.Lookup(msg.Message) // Get the channel the invite is being issued from

	if !permsOk {
		trunk.LogErr(fmt.Sprintf("Failed to get permissions from %s for %s", msg.Message, msg.Issuer))
		return
	}

	if !channelPerms.IsTrusted() { // User is not trusted
		trunk.LogInfo(fmt.Sprintf("Rejecting invite by non-admin %s to %s", msg.Issuer, msg.Message))
		return
	}

	trunk.LogInfo(fmt.Sprintf("Joining channel %s from channel trusted user: %s", msg.Message, msg.Issuer))
	c.Cmd.Join(msg.Message)

	Config.Channels = append(Config.Channels, msg.Message)
	Config.Channels = DeduplicateList(Config.Channels)

	msg = ParseMessage(c, e) // Re-parse our message for Msg.Admin check

	if msg.Admin { // Bot admin
		SaveConfig()
	}
}

// Parser will handle the majority of incoming messages, user joins, etc.
func Parser(c *girc.Client, e girc.Event) {
	m := ParseMessage(c, e)

	var ignoreMessage bool
	command := e.Command

	for _, ignoreCommand := range ignoreCommands { // For each ignore command
		if ignoreCommand == command {
			ignoreMessage = true
			break
		}
	}

	if !ignoreMessage {
		var userInBlacklist bool

		for _, blacklistUser := range Config.Users.Blacklist { // For each user
			userInBlacklist = Matches(blacklistUser, m.Issuer) // Check against issuer

			if !userInBlacklist { // Didn't match based on nick
				userInBlacklist = Matches(blacklistUser, m.Host) // Check against host
			}

			if userInBlacklist { // Matched
				break
			}
		}

		if PluginManager.IsEnabled("AutoKick") { // AutoKick enabled
			NarwhalAutoKicker.Parse(c, e, m) // Run through auto-kicker first
		}

		if !userInBlacklist && (m.Issuer != Config.User) { // Ensure we aren't parsing our own bot messages
			trunk.LogInfo("Allowed: " + m.Issuer)
			trunk.LogInfo("Authenticated: " + strconv.FormatBool(m.Authenticated))
			trunk.LogInfo("Full Issuer: " + m.FullIssuer)
			trunk.LogInfo("Received: " + m.Message)
			trunk.LogInfo("Host: " + m.Host)
			trunk.LogInfo("Possible Channel: " + m.Channel)

			if PluginManager.IsEnabled("Admin") { // Admin Management enabled
				NarwhalAdminManager.Parse(c, e, m) // Run through management
			}

			if PluginManager.IsEnabled("Replacer") { // Replacer enabled
				NarwhalReplacer.AddToCache(m)
				NarwhalReplacer.Parse(c, e, m) // Run through replacer
			}

			if PluginManager.IsEnabled("Song") { // Song enabled
				NarwhalSong.Parse(c, e, m) // Run through song
			}

			if PluginManager.IsEnabled("Slap") { // Slap enabled
				NarwhalSlap.Parse(c, e, m) // Run through slap
			}

			if PluginManager.IsEnabled("ThankMe") { // Thank The Bot enabled
				NarwhalThank.Parse(c, e, m) // Run through thank me
			}

			if PluginManager.IsEnabled("UrlParser") && (m.Channel != "") { // Url Parser enabled and provided URL over a channel
				NarwhalUrlParser.Parse(c, e, m) // Run through URL parser
			}

			for _, parseFunc := range PluginManager.Modules { // For each parser function in our Plugin Modules
				parseFunc.(func(*girc.Client, girc.Event, NarwhalMessage))(c, e, m)
			}
		}
	}
}
