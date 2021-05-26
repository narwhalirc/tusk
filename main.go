package tusk

import (
	"github.com/JoshStrobl/trunk"
	"github.com/lrstanley/girc"
	"os"
	"os/user"
	"path/filepath"
	"plugin"
	"time"
)

// A Narwhal is no Narwhal without their tusk!

// Consistent paths
var Paths []string

// Config is our Narwhal Config
var Config NarwhalConfig

// PluginManager is our Plugin Manager
var PluginManager NarwhalPluginManager

func init() {
	var getUserErr error
	var currentUser *user.User

	currentUser, getUserErr = user.Current() // Attempt to get the current user

	if getUserErr != nil { // If we successfully got the user
		trunk.LogFatal("Failed to get the current user: " + getUserErr.Error())
	}

	workdir, getWdErr := os.Getwd() // Get the current working directory

	if getWdErr != nil { // If we failed to get the current working dir
		trunk.LogFatal("Failed to get the current working directory: " + getWdErr.Error())
	}

	Paths = []string{
		filepath.Join(currentUser.HomeDir, ".config", "narwhal"),
		workdir,
		"/etc/narwhal",
		"/usr/share/defaults/narwhal",
	}
}

// NewTusk will create a new tusk for our Narwhal, but only one tusk is allowed at a time.
func NewTusk() {
	var newTuskErr error

	Config, newTuskErr = ReadConfig()

	if newTuskErr == nil { // Read our config
		PluginManager.Modules = make(map[string]plugin.Symbol)

		if loadPluginsErr := PluginManager.LoadPlugins(); loadPluginsErr != nil { // Failed to load a plugin
			trunk.LogWarn("Failed to load plugin: " + loadPluginsErr.Error())
		}

		ircConfig := girc.Config{
			Server: Config.Network,
			Port:   Config.Port,
			Name:   Config.Name,
			Nick:   Config.User,
			User:   Config.User,
			SASL: &girc.SASLPlain{
				User: Config.User,
				Pass: Config.Password,
			},
		}

		client := girc.New(ircConfig)
		client.Handlers.Add(girc.CONNECTED, OnConnected) // On CONNECTED, trigger OnConnected
		client.Handlers.Add(girc.INVITE, OnInvite)       // On INVITE, trigger OnInvite
		client.Handlers.Add(girc.JOIN, OnJoin)           // On JOIN, trigger our OnJoin
		client.Handlers.Add(girc.PRIVMSG, Parser)        // On PRIVMSG, trigger our Parser
		client.Handlers.Add(girc.KICK, OnKick)           // On KICK, trigger our OnKick

		if newTuskErr = client.Connect(); newTuskErr != nil { // Failed during run
			trunk.LogFatal("Failed to run client: " + newTuskErr.Error())
		}

		AnnounceLibera(client)
		reminder := time.NewTimer(1 * time.Hour)
		go func(c *girc.Client) {
			<-reminder.C
			AnnounceLibera(c)
		}(client)
	} else {
		trunk.LogFatal("Failed to read or parse config: " + newTuskErr.Error())
	}
}

func AnnounceLibera(c *girc.Client) {
	solusIrcChannels := []string{
		"#budgie-desktop-dev",
		"#solus",
		"#solus-chat",
		"#solus-dev",
		"#solus-livestream",
	}

	for _, channel := range solusIrcChannels { // For each IRC channel
		c.Cmd.Action(channel, "We are now available on Libera Chat.")
	}
}
