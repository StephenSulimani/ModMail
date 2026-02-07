package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func MessageCreate() func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		fmt.Println("=================================")
		fmt.Printf("Author: %s\n", m.Author.Username)
		fmt.Println(m.Content)
		fmt.Println("=================================")
	}
}
