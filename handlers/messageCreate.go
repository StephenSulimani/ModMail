package handlers

import (
	"database/sql"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/stephensulimani/modmail/database"
)

func MessageCreate(db *sql.DB) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Member == nil {
			// Message sent in DM
			ticket, err := database.FindTicketByUserId(db, m.Author.ID)
			if err != nil {
				fmt.Println(err)
				if err == sql.ErrNoRows {
					fmt.Println("no tickets found")
					ticket, err := database.CreateTicket(db, &database.Ticket{
						UserID:     m.Author.ID,
						ChannelID:  "test",
						WebhookURL: "test",
					})
					if err != nil {
						fmt.Println(err)
						return
					}
					fmt.Println("New Ticket:", ticket)
					return
				}
				fmt.Println(err)
				return
			}
			fmt.Println("Existing Ticket:", ticket)

		} else {
			// Message sent in channel
			alias, err := database.FindAliasByChannelId(db, m.ChannelID)
			if err != nil {
				if err == sql.ErrNoRows {
					return
				}
				fmt.Println(err)
				return
			}
			fmt.Println(alias)

		}

		fmt.Println("=================================")
		fmt.Printf("Author: %s\n", m.Author.Username)
		fmt.Println(m.Content)
		fmt.Println("=================================")
	}
}

func HandleDM(s *discordgo.Session, m *discordgo.MessageCreate) {
}

func HandleTicket(s *discordgo.Session, m *discordgo.MessageCreate) {
}
