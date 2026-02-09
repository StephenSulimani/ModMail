package database

type Ticket struct {
	ID         int
	UserID     string
	ChannelID  string
	WebhookURL string
	CreatedAt  string
}

type Alias struct {
	Ticket    *Ticket
	StaffID   string
	Alias     int
	CreatedAt string
}
