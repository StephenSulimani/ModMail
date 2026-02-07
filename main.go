package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/stephensulimani/modmail/handlers"
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
	secret_key := make([]byte, 32)
	rand.Read(secret_key)

	// Set up and configure database
	db, err := sql.Open("sqlite", "file:modmail?mode=memory&cache=shared")
	if err != nil {
		fmt.Println("Error opening database: ", err)
		return
	}

	err = CreateTables(db)
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
	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers

	discord.AddHandler(handlers.Ready())

	discord.AddHandler(handlers.MessageCreate())

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

func CreateTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + "`tickets`" + `(
		` + "`id`" + ` INTEGER PRIMARY KEY AUTOINCREMENT,	
		` + "`user_id`" + ` VARCHAR(255) NOT NULL UNIQUE,
        ` + "`channel_id`" + ` VARCHAR(255) NOT NULL UNIQUE,
		` + "`webhook_url`" + ` VARCHAR(255) NOT NULL,
		` + "`created_at`" + ` DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + "`aliases`" + `(
		` + "`ticket_id`" + ` INTEGER NOT NULL,
		` + "`staff_id`" + ` VARCHAR(255) NOT NULL,
		` + "`alias`" + ` INTEGER NOT NULL,
        ` + "`created_at`" + ` DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(` + "`ticket_id`" + `) REFERENCES tickets(` + "`id`" + `) ON DELETE CASCADE,
		PRIMARY KEY(` + "`ticket_id`" + `, ` + "`staff_id`" + `)
		)
	`)

	return err
}
