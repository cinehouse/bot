package bot

import (
	"context"
	"log"
	"sort"
	"strings"

	"github.com/go-pkgz/syncs"
)

//go:generate mockery --inpackage --name=Interface --case=snake

// Interface is a bot interface.
type Interface interface {
	Help() string
	ReactOn() []string
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
	ID       int64  `json:"id"`
	UserName string `json:"username"`
}

// Response describes bot's answer on particular message
type Response struct {
	Text    string
	Send    bool // status
	Preview bool // enable web preview
}

// MultiBot combines many bots to one virtual
type MultiBot []Interface

// Help returns help message
func (b MultiBot) Help() string {
	sb := strings.Builder{}
	for _, child := range b {
		help := child.Help()
		if help != "" {
			// WriteString always returns nil err
			if !strings.HasSuffix(help, "\n") {
				help += "\n"
			}
			_, _ = sb.WriteString(help)
		}
	}
	return sb.String()
}

// ReactOn returns combined list of all keywords
func (b MultiBot) ReactOn() (res []string) {
	for _, bot := range b {
		res = append(res, bot.ReactOn()...)
	}
	return res
}

// OnMessage pass msg to all bots and collects reposnses (combining all of them)
//noinspection GoShadowedVar
func (b MultiBot) OnMessage(msg Message) (response Response) {
	if contains([]string{"help", "/help", "help!"}, msg.Text) {
		return Response{
			Text: b.Help(),
			Send: true,
		}
	}

	resps := make(chan string)

	wg := syncs.NewSizedGroup(4)
	for _, bot := range b {
		bot := bot
		wg.Go(func(ctx context.Context) {
			if resp := bot.OnMessage(msg); resp.Send {
				resps <- resp.Text
			}
		})
	}

	go func() {
		wg.Wait()
		close(resps)
	}()

	var lines []string
	for r := range resps {
		log.Printf("[DEBUG] collect %q", r)
		lines = append(lines, r)
	}

	sort.Slice(lines, func(i, j int) bool {
		return lines[i] < lines[j]
	})

	log.Printf("[DEBUG] answers %d, send %v", len(lines), len(lines) > 0)
	return Response{
		Text: strings.Join(lines, "\n"),
		Send: len(lines) > 0,
	}
}

func generateHelpMessage(command []string, description string) string {
	return strings.Join(command, ", ") + " _– " + description + "_\n"
}

func contains(s []string, e string) bool {
	e = strings.TrimSpace(e)
	for _, a := range s {
		if strings.EqualFold(a, e) {
			return true
		}
	}
	return false
}
