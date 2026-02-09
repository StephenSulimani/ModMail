package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/stephensulimani/modmail/database"
	"github.com/stephensulimani/modmail/handlers"
	"github.com/stephensulimani/modmail/internal"
	_ "modernc.org/sqlite"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	BOT_TOKEN := os.Getenv("BOT_TOKEN")

	if BOT_TOKEN == "" {
		fmt.Println("BOT_TOKEN environment variable is not set")
		return
	}

	// Generate Secret Key
	internal.SECRET_KEY = make([]byte, 32)
	rand.Read(internal.SECRET_KEY)

	// Generate Index Key
	internal.INDEX_KEY = make([]byte, 32)
	rand.Read(internal.INDEX_KEY)

	// Set up and configure database
	db, err := sql.Open("sqlite", "file:modmail?mode=memory&cache=shared")
	// db, err := sql.Open("sqlite", "./modmail.db")
	if err != nil {
		fmt.Println("Error opening database: ", err)
		return
	}

	err = database.CreateTables(db)
	if err != nil {
		fmt.Println("Table creation failed: ", err)
		return
	}

	fmt.Println("Database created successfully")

	// Set up and configure Discord Bot
	discord, err := discordgo.New("Bot " + BOT_TOKEN)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}
	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentsDirectMessages

	discord.AddHandler(handlers.Ready())

	discord.AddHandler(handlers.MessageCreate(db))

	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
		return
	}
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	err = discord.Close()
	if err != nil {
		fmt.Printf("could not close session gracefully: %s\n", err)
	}
}
