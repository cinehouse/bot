package bot

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// SysBot implements basic bot function to respond on ping and others from basic.data file.
// also, reacts on say! with keys/values from say.data file
type SysBot struct {
	dataLocation string
	commands     []sysCommand
}

// sysCommand hold one type triggers from basic.data
type sysCommand struct {
	triggers    []string
	description string
	message     string
}

// NewSys makes new sys bot and load data to []say and basic map
func NewSys(dataLocation string) (*SysBot, error) {
	log.Printf("[INFO] created sys bot, data location=%s", dataLocation)
	res := SysBot{dataLocation: dataLocation}
	if err := res.loadBasicData(); err != nil {
		return nil, err
	}
	rand.Seed(0)
	return &res, nil
}

// Help returns help message
func (p SysBot) Help() (line string) {
	for _, c := range p.commands {
		line += generateHelpMessage(c.triggers, c.description)
	}
	return line
}

// ReactOn keys
func (p SysBot) ReactOn() []string {
	res := make([]string, 0)
	for _, bot := range p.commands {
		res = append(bot.triggers, res...)
	}
	return res
}

// OnMessage implements bot.Interface
func (p SysBot) OnMessage(msg Message) (response Response) {
	if !contains(p.ReactOn(), msg.Text) {
		return Response{}
	}

	for _, bot := range p.commands {
		if found := contains(bot.triggers, strings.ToLower(msg.Text)); found {
			return Response{Text: bot.message, Send: true}
		}
	}

	return Response{}
}

func (p *SysBot) loadBasicData() error {
	basicData, err := readLines(p.dataLocation + "/basic.data")
	if err != nil {
		return errors.Wrap(err, "can't load basic.data")
	}

	for _, line := range basicData {
		elems := strings.Split(line, "|")
		if len(elems) != 3 {
			log.Printf("[DEBUG] bad format %s, ignored", line)
			continue
		}
		sysCommand := sysCommand{
			description: elems[1],
			message:     elems[2],
			triggers:    strings.Split(elems[0], ";"),
		}
		p.commands = append(p.commands, sysCommand)
		log.Printf("[DEBUG] loaded basic response, %v, %s", sysCommand.triggers, sysCommand.message)
	}
	return nil
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, errors.Wrapf(err, "can't open %s", path)
	}
	defer f.Close() //nolint

	result := make([]string, 0)
	s := bufio.NewScanner(f)
	for s.Scan() {
		result = append(result, s.Text())
	}

	return result, nil
}
