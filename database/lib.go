package database

import (
	"database/sql"
	"errors"

	"github.com/stephensulimani/modmail/internal"
)

// CreateTables creates the database tables
func CreateTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + "`tickets`" + `(
		` + "`id`" + ` INTEGER PRIMARY KEY AUTOINCREMENT,	
		` + "`user_id`" + ` VARCHAR(255) NOT NULL UNIQUE,
		` + "`blind_user_id`" + ` VARCHAR(255) NOT NULL UNIQUE,
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

// FindTicketByUserId finds a ticket by the user's ID
func FindTicketByUserId(db *sql.DB, userId string) (*Ticket, error) {
	var ticket Ticket

	err := db.QueryRow("SELECT id, user_id, channel_id, webhook_url, created_at FROM tickets WHERE blind_user_id = ?",
		internal.GenerateBlindIndex(userId)).Scan(
		&ticket.ID,
		&ticket.UserID,
		&ticket.ChannelID,
		&ticket.WebhookURL,
		&ticket.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	decrypted, err := internal.Decrypt(ticket.UserID)
	if err != nil {
		return nil, err
	}

	ticket.UserID = decrypted

	return &ticket, nil
}

// CreateTicket creates a new ticket
func CreateTicket(db *sql.DB, ticket *Ticket) (*Ticket, error) {
	if ticket.UserID == "" || ticket.ChannelID == "" || ticket.WebhookURL == "" {
		return nil, errors.New("user_id, channel_id, and webhook_url are required")
	}
	encryptedUserId, err := internal.Encrypt(ticket.UserID)
	if err != nil {
		return nil, err
	}

	res, err := db.Exec("INSERT INTO tickets (user_id, blind_user_id, channel_id, webhook_url) VALUES (?, ?, ?, ?)", encryptedUserId, internal.GenerateBlindIndex(ticket.UserID), ticket.ChannelID, ticket.WebhookURL)
	if err != nil {
		return nil, err
	}

	ticket_id, _ := res.LastInsertId()

	ticket.ID = int(ticket_id)

	return ticket, nil
}

// FindTicketByChannelId finds a ticket by the channel's ID
func FindTicketByChannelId(db *sql.DB, channelId string) (*Ticket, error) {
	var ticket Ticket

	err := db.QueryRow("SELECT * FROM tickets WHERE channel_id = ?", channelId).Scan(
		&ticket.ID,
		&ticket.UserID,
		&ticket.ChannelID,
		&ticket.WebhookURL,
		&ticket.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	decrypted, err := internal.Decrypt(ticket.UserID)
	if err != nil {
		return nil, err
	}

	ticket.UserID = decrypted

	return &ticket, nil
}

// FindAliasByChannelId finds an alias by the channel's ID
func FindAliasByChannelId(db *sql.DB, channelId string) (*Alias, error) {
	ticket, err := FindTicketByChannelId(db, channelId)
	if err != nil {
		return nil, err
	}

	var alias Alias

	err = db.QueryRow("SELECT staff_id, alias, created_at FROM aliases WHERE ticket_id = ?", ticket.ID).Scan(
		&alias.StaffID,
		&alias.Alias,
		&alias.CreatedAt,
	)

	alias.Ticket = ticket

	return &alias, err
}
