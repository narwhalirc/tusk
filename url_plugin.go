package tusk

import (
	"fmt"
	"github.com/TryStreambits/sauron"
	"github.com/lrstanley/girc"
	"net/url"
	"strings"
)

//NarwhalUrlParserPlugin is our URL plugin
type NarwhalUrlParserPlugin struct{}

// NarwhalUrlParser is our url parser
var NarwhalUrlParser NarwhalUrlParserPlugin

func (parser *NarwhalUrlParserPlugin) Parse(c *girc.Client, e girc.Event, m NarwhalMessage) {
	urls := []url.URL{}
	splitMessage := strings.Split(m.Message, " ") // Split on space, URLs should not contain whitespace

	for _, subMessage := range splitMessage {
		if url, parseErr := url.Parse(subMessage); parseErr == nil { // If we successfully parsed this URL
			urls = append(urls, *url) // Add this url to our urls
		}
	}

	if len(urls) > 0 { // If we have URLs
		for _, url := range urls {
			link, linkGetErr := sauron.GetLink(url.String()) // Get the link for this URL

			if linkGetErr != nil { // Failed to get the link
				continue // Skip it
			}

			var message string // The message we'll send

			if (link.Extras["IsImageLink"] == "true") || (link.Extras["IsVideoLink"] == "true") { // Direct image or video link
				continue // Skip making a message for it
			}

			if link.Extras["IsRedditLink"] == "true" { // Is a Reddit Link
				message = fmt.Sprintf("[ %s ][Score: %s, %d%% upvotes]", link.Title, link.Extras["Score"], link.Extras["Percentage"])
			} else if link.Extras["IsYouTubeLink"] == "true" && link.Extras["IsVideo"] == "true" { // Is a YouTube video specifically
				desktopYT := fmt.Sprintf("https://youtube.com/watch?v=%s", link.Extras["Video"])
				mobileYT := fmt.Sprintf("https://m.youtube.com/watch?v=%s", link.Extras["Video"])
				message = fmt.Sprintf("[ %s | Desktop: %s | Mobile: %s ]", link.Title, desktopYT, mobileYT)
			} else {
				message = fmt.Sprintf("[ %s ]", link.Title)
			}

			if message != "" { // Not an empty title
				c.Cmd.Reply(e, message)
			}
		}
	}
}
