package reporter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cinehouse/bot/app/bot"
)

// Reporter collects all messages and saves to plain file
type Reporter struct {
	logsPath string
	messages chan string
}

// NewLogger makes new reporter bot
func NewLogger(logs string) (result Reporter) {
	log.Printf("[INFO] new reporter, path=%s", logs)
	_ = os.MkdirAll(logs, 0750)
	result = Reporter{logsPath: logs, messages: make(chan string, 1000)}
	go result.activate()
	return result
}

// Save to log channel, non-blocking and skip if needed
func (l Reporter) Save(msg *bot.Message) {
	if msg.Text == "" {
		log.Print("[DEBUG] Message not saved to log: no text")
		return
	}

	bdata, err := json.Marshal(&msg)
	if err != nil {
		log.Printf("[WARN] failed to log, error %v", err)
		return
	}

	select {
	case l.messages <- string(bdata) + "\n":
	default:
		log.Printf("[WARN] can't buffer log entry %v", msg)
	}
}

func (l Reporter) activate() {
	log.Print("[INFO] activate reporter")
	buffer := make([]string, 0, 100)

	writeBuff := func() error {
		if len(buffer) == 0 {
			return nil
		}
		// nolint
		fh, err := os.OpenFile(fmt.Sprintf("%s/%s.log", l.logsPath, time.Now().Format("20060102")),
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)

		if err != nil {
			log.Printf("[WARN] failed to log, %v", err)
			return err
		}
		defer fh.Close() //nolint
		for _, rec := range buffer {
			if _, err = fh.WriteString(rec); err != nil {
				log.Printf("[WARN] failed to write log, %v", err)
			}
		}

		log.Printf("[DEBUG] wrote %d log entries", len(buffer))
		buffer = buffer[:0]
		return nil
	}

	for {
		select {
		case entry := <-l.messages:
			buffer = append(buffer, entry)
			if len(buffer) >= 100 { // forced flush every 100 records
				if err := writeBuff(); err != nil {
					log.Printf("[WARN] failed to write reporter buffer, %v", err)
				}
			}
		case <-time.After(time.Second * 5): // flush on 5 seconds inactivity
			if err := writeBuff(); err != nil {
				log.Printf("[WARN] failed to write reporter buffer, %v", err)
			}
		}
	}
}
