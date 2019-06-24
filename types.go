package tusk

import (
	"github.com/lrstanley/girc"
	"net/url"
)

// Our structs and interfaces

// NarwhalConfig is our primary Narwhal configuration
type NarwhalConfig struct {
	// Network is the IRC network to connection to.
	Network string `toml:"Network"`

	// Port is the port on the network we're connecting to. Likely 6667.
	Port int `toml:"Port,omitempty"`

	// User is the IRC Bot username
	User string

	// Name is the IRC Bot name
	Name string

	// FallbackNick is the IRC bot fallback nickname if first nick is registered to someone else
	FallbackNick string `toml:"FallbackNick,omitempty"`

	// Password is the IRC bot password for authentication
	Password string

	// Plugins is a list of plugin configurations
	Plugins NarwhalPluginsConfig `toml:"Plugins,omitempty"`

	// Channels is a list of channels to join
	Channels []string

	// Users is our users configuration
	Users NarwhalUsersConfig `toml:"Users,omitempty"`
}

// NarwhalLink is a struct containing information related to an HTTP resource
type NarwhalLink struct {
	// IsReddit designates whether this resource is a Reddit URL
	IsReddit bool

	// IsYoutube designates whether this resource is a Youtube URL
	IsYoutube bool

	// Link is our net URL struct
	Link url.URL

	// Title is the page title
	Title string

	// Votes is the Reddit votes (if IsReddit)
	Votes NarwhalRedditVotes
}

// NarwhalMessage is a custom message
type NarwhalMessage struct {
	Channel      string
	Command      string
	Host         string
	FullIssuer   string
	Issuer       string
	Message      string
	MessageNoCmd string
	Params       []string
}

// NarwhalPlugin is a plugin interface
type NarwhalPlugin interface {
	Parse(c *girc.Client, e girc.Event, m NarwhalMessage)
}

// NarwhalRedditVotes is the total votes for a reddit thread
type NarwhalRedditVotes struct {
	Dislikes string
	Likes    string
	Score    string
}

// NarwhalUsersConfig is our configuration for blacklisting users, administrative users, and autokicking
type NarwhalUsersConfig struct {
	// Admins is an array of users authorized to perform admin actions
	Admins []string

	// Blacklist is an array of users blacklisted from performing Plugins
	Blacklist []string
}

// NarwhalPluginsConfig is a list of command configurations
type NarwhalPluginsConfig struct {
	Enabled []string

	Admin    NarwhalAdminConfig      `toml:"Admin,omitempty"`
	AutoKick NarwhalAutoKickerConfig `toml:"AutoKick,omitempty"`
	Replacer NarwhalReplacerConfig   `toml:"Replacer,omitempty"`
	Slap     NarwhalSlapConfig       `toml:"Slap,omitempty"`
}

////////////////////////////////////////////////////////
// PLUGINS
////////////////////////////////////////////////////////

// NarwhalAdminConfig is our configuration for the Narwhal admin plugin
type NarwhalAdminConfig struct {
	// DisabledCommands is a list of admin commands to disable
	DisabledCommands []string
}

// NarwhalAdminPlugin is our Admin plugin
type NarwhalAdminPlugin struct{}

// NarwhallAutoKickerConfig is our configuration for the Narwhal autokicker
type NarwhalAutoKickerConfig struct {
	// EnabledAutoban determines whether to enable the automatic banning of users which exceed our MinimumKickToBanCount
	EnabledAutoban bool `json:",omitempty"`

	// Hosts to kick. Matches from end.
	Hosts []string

	// MessageMatches is a list of messages that will result in kicks
	MessageMatches []string

	// MinimumKickToBanCount is a minimum amount of times a user should be kicked before being automatically banned. Only enforced when EnabledAutoban is set
	MinimumKickToBanCount int `json:",omitempty"`

	// Users to kick. Matches from beginning.
	Users []string
}

// NarwhalAutoKickerPlugin is our Autokick plugin
type NarwhalAutoKickerPlugin struct {
	// Tracker is a map of usernames to the amount of times they've been kicked
	Tracker map[string]int
}

// NarwhalReplacerConfig is our configuration for the Narwhal replacer plugin
type NarwhalReplacerConfig struct {
	// CachedMessageLimit is our limit of how many messages to cache
	CachedMessageLimit int
}

// NarwhalReplacerPlugin is our Replacer plugin
type NarwhalReplacerPlugin struct{}

// NarwhalSlapConfig is our configuration for the Narwhal autokicker
type NarwhalSlapConfig struct {
	// CustomActions is a list of custom actions on how to slap a user
	CustomActions []string
}

// NarwhalSlapPlugin is our slap plugin
type NarwhalSlapPlugin struct {
	Objects []string
}

// NarwhalSong is our Song plugin
type NarwhalSongPlugin struct{}

//NarwhalUrlParserPlugin is our URL plugin
type NarwhalUrlParserPlugin struct{}
