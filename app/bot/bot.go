package bot

// Interface is a bot interface.
type Interface interface {
	OnMessage(msg Message) (response Response)
}

// Message is primary record to pass data from/to bots
type Message struct {
	MessageID int    `json:"message_id"`
	From      User   `json:"from,omitempty"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text,omitempty"`
}

// Chat represents a chat.
type Chat struct {
	ID int64 `json:"id"`
}

// User represents a user or bot.
type User struct {
	ID int64 `json:"id"`
}

// Response describes bot's answer on particular message
type Response struct {
	Text    string
	Send    bool // status
	Preview bool // enable web preview
}

// MultiBot combines many bots to one virtual
type MultiBot []Interface

// OnMessage pass msg to all bots and collects reposnses (combining all of them)
//noinspection GoShadowedVar
func (b MultiBot) OnMessage(msg Message) (response Response) {
	for _, bot := range b {
		response = bot.OnMessage(msg)
	}
	return
}
