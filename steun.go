package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
	"io/ioutil"
	"log"
	"math/rand"
	"regexp"
	"sync"
	"time"
)

var seqMutex sync.Mutex
var seq int

var hadRandomMutex sync.Mutex
var hadRandom = map[string]struct{}{}

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
		go queue(fmt.Sprintf("./audio/%s.wav", matches[0][1]))
		return
	}
	sendSequential(s)
}

func sendRandomMessage(s *discordgo.Session) {
	go connectVoice(s)
	var i int

	hadRandomMutex.Lock()
	files, err := ioutil.ReadDir("./audio")
	if err != nil {
		log.Fatal(err)
	}

	if (len(hadRandom)) >= len(files) {
		// reset counter
		hadRandom = map[string]struct{}{}
	}
	try := 0
	for {
		i = rand.Intn(len(files))
		if _, exists := hadRandom[files[i].Name()]; !exists {
			hadRandom[files[i].Name()] = struct{}{}
			break
		}
		try++
		if try > 500 {
			// give up
			break
		}
	}
	hadRandomMutex.Unlock()
	queue(fmt.Sprintf("./audio/%s", files[i].Name()))
}

func sendSequential(s *discordgo.Session) {
	go connectVoice(s)

	files, err := ioutil.ReadDir("./audio")
	if err != nil {
		log.Fatal(err)
	}

	go queue(fmt.Sprintf("./audio/%s", files[seq].Name()))
	seqMutex.Lock()
	seq++
	if seq > len(files) {
		seq = 0
	}
	seqMutex.Unlock()
}

func sendSteunXL(s *discordgo.Session, m *discordgo.MessageCreate) {
	files, err := ioutil.ReadDir("./audio")
	if err != nil {
		log.Fatal(err)
	}

	randNrs := map[int]struct{}{}
	for {
		i := rand.Intn(len(files))
		if _, exists := randNrs[i]; !exists {
			randNrs[i] = struct{}{}
		}
		if len(randNrs) >= 10 {
			break
		}
	}
	for i := range randNrs {
		go queue(fmt.Sprintf("./audio/%s", files[i].Name()))
	}
}
