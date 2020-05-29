package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
	"math/rand"
)

func init() {
	registerCommand(command.Command{
		Name:        "steun",
		Category:    command.CategoryFun,
		Description: "",
		Hidden:      false,
		Handler:     sendSteun,
	})
}

func sendSteun(s *discordgo.Session, m *discordgo.MessageCreate) {
	sendRandomMessage(s)
}

func sendRandomMessage(s *discordgo.Session) {
	go connectVoice(s)
	i := rand.Intn(11)
	queue(fmt.Sprintf("./audio/%02d.wav", i))
}
