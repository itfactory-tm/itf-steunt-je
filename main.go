package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/itfactory-tm/thomas-bot/pkg/command"

	"github.com/bwmarrin/discordgo"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Token  string
	Prefix string `default:"tm"`
}

var c config
var handlers = map[string]command.Command{}
var commandRegex *regexp.Regexp

func main() {
	err := envconfig.Process("itfsteun", &c)
	if err != nil {
		log.Fatal(err)
	}
	if c.Token == "" {
		log.Fatal("No token specified")
	}

	commandRegex = regexp.MustCompile(c.Prefix + `!(\w*)\b`)

	dg, err := discordgo.New("Bot " + c.Token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	// Register handlers
	dg.AddHandler(onMessage)

	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}
	// TODO: add connection error handlers

	go func() {
		go connectVoice(dg)
		for {
			i := rand.Intn(11)
			queue(fmt.Sprintf("./audio/%02d.wav", i))
			time.Sleep(10 * time.Second)
		}
	}()
	log.Println("Thomas Bot is now supporting students.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func queue(f string) {
	voiceQueueChan <- f
	time.Sleep(500 * time.Millisecond)
	for {
		encoderMutex.Lock()
		if encoder == nil {
			encoderMutex.Unlock()
			time.Sleep(100 * time.Millisecond)
			continue
		} else {
			encoderMutex.Unlock()
		}

		if encoder.HasFinishedAll() {
			break //MS80
		} else {
			time.Sleep(200 * time.Millisecond)
		}
	}

}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if commandRegex.MatchString(m.Content) {
		if c, exists := handlers[commandRegex.FindStringSubmatch(m.Content)[1]]; exists {
			c.Handler(s, m)
		}
	}
}

func registerCommand(c command.Command) {
	handlers[c.Name] = c
}
