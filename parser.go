package tusk

import (
	"fmt"
	"github.com/JoshStrobl/trunk"
	"github.com/lrstanley/girc"
	"time"
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

func AnnounceLibera(c *girc.Client) {
	for _, channel := range Config.Channels { // For each channel we are
		c.Cmd.Message(channel, "We are available on Libera IRC. Same channels, different server!")
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

	if !Config.LiberaAnnounceMigration { // Enable migration announcement
		return
	}

	time.Sleep(5 * time.Second) // We don't know when we've joined, so just waitY

	AnnounceLibera(c)
	reminder := time.NewTicker(1 * time.Hour)
	go func(c *girc.Client) {
		for {
			select {
			case <-reminder.C:
				AnnounceLibera(c)
			}
		}
	}(c)
}

// OnJoin will handle when a user joins a channel
func OnJoin(c *girc.Client, e girc.Event) {
	m := ParseMessage(c, e)

	m.Channel = e.Params[0] // Set channel to first param

	fmt.Println(MessageBreaker)
	fmt.Printf("[%s] Joined %s: %s (%s)\n", GetNowAsISO8601(), m.Channel, m.Issuer, m.FullIssuer)

	if PluginManager.IsEnabled("AutoKick") { // AutoKick enabled
		NarwhalAutoKicker.Parse(c, e, m) // Run through auto-kicker first
	}
}

// OnLeave will handle when a user leaves a channel
func OnLeave(c *girc.Client, e girc.Event) {
	m := ParseMessage(c, e)
	fmt.Printf("Left: %v\n", m)
}

// OnInvite will handle a request to invite an IRC channel
func OnInvite(c *girc.Client, e girc.Event) {
	msg := ParseMessage(c, e) // Parse our message

	if !msg.Admin { // User is not trusted
		trunk.LogErr(fmt.Sprintf("Rejecting invite by non-admin %s to %s", msg.Issuer, msg.Message))
		return
	}

	trunk.LogInfo(fmt.Sprintf("Joining channel %s invited by admin %s", msg.Message, msg.Issuer))
	c.Cmd.Join(msg.Message)

	Config.Channels = append(Config.Channels, msg.Message)
	Config.Channels = DeduplicateList(Config.Channels)
	SaveConfig()
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

		if !userInBlacklist && (m.Issuer != Config.User) { // Ensure we aren't parsing our own bot messages
			PrintPrettyMessage(m)

			if PluginManager.IsEnabled("Admin") { // Admin Management enabled
				NarwhalAdminManager.Parse(c, e, m) // Run through management
			}

			if PluginManager.IsEnabled("AutoKick") { // AutoKick enabled
				NarwhalAutoKicker.Parse(c, e, m) // Run through auto-kicker first
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
