package channelserver

import (
	"fmt"

	"github.com/matterbridge/discordgo"
)

// onDiscordMessage handles receiving messages from discord and forwarding them ingame.
func (s *Server) onDiscordMessage(ds *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from our bot, or ones that are not in the correct channel.
	if m.Author.ID == ds.State.User.ID || m.ChannelID != s.erupeConfig.Discord.ChannelID {
		return
	}

	// Broadcast to the game clients.
	message := fmt.Sprintf("[DISCORD] %s: %s", m.Author.Username, m.Content)
	s.BroadcastChatMessage(message)
}
