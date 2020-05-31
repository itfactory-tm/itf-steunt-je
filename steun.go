package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
	"math/rand"
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
}

func sendSteun(s *discordgo.Session, m *discordgo.MessageCreate) {
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
