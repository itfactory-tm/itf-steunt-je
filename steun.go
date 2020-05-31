package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
	"math/rand"
	"regexp"
	"strconv"
	"sync"
	"time"
)

const totalFiles = 51

var seqMutex sync.Mutex
var seq int

func init() {
	rand.Seed(time.Now().UnixNano())
	registerCommand(command.Command{
		Name:        "shout",
		Category:    command.CategoryFun,
		Description: "",
		Hidden:      false,
		Handler:     sendSteun,
	})
	registerCommand(command.Command{
		Name:        "shoutxl",
		Category:    command.CategoryFun,
		Description: "",
		Hidden:      false,
		Handler:     sendSteunXL,
	})
}

var numRegex = regexp.MustCompile(`^tm!shout ([0-9]*)$`)

func sendSteun(s *discordgo.Session, m *discordgo.MessageCreate) {
	matches := numRegex.FindAllStringSubmatch(m.Message.Content, -1)
	if len(matches) > 0 && len(matches[0]) > 1 {
		i, err := strconv.ParseInt(matches[0][1], 10, 64)
		if err != nil {
			go queue(fmt.Sprintf("./audio/%02d.wav", i))
			return
		}
	}
	sendSequential(s)
}

func sendRandomMessage(s *discordgo.Session) {
	go connectVoice(s)
	i := rand.Intn(totalFiles)
	queue(fmt.Sprintf("./audio/%02d.wav", i))
}

func sendSequential(s *discordgo.Session) {
	go connectVoice(s)

	go queue(fmt.Sprintf("./audio/%02d.wav", seq))
	seqMutex.Lock()
	seq++
	if seq > totalFiles {
		seq = 0
	}
	seqMutex.Unlock()
}

func sendSteunXL(s *discordgo.Session, m *discordgo.MessageCreate) {
	randNrs := map[int]struct{}{}
	for {
		i := rand.Intn(totalFiles)
		if _, exists := randNrs[i]; !exists {
			randNrs[i] = struct{}{}
		}
		if len(randNrs) >= 10 {
			break
		}
	}
	for i := range randNrs {
		go queue(fmt.Sprintf("./audio/%02d.wav", i))
	}
}
