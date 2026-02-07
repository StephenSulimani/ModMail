package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Ready() func(s *discordgo.Session, e *discordgo.Ready) {
	return func(s *discordgo.Session, e *discordgo.Ready) {
		fmt.Printf("Successfully logged in as %s\n", s.State.User.Username)
	}
}
