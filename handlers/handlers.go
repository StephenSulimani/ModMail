package handlers

import "github.com/bwmarrin/discordgo"

type Handler[T any] func(s *discordgo.Session, event T)
