package main

import (
	"fmt"
	"sync"

	"github.com/itfactory-tm/itf-steunt-je/pkg/mixer"

	"github.com/bwmarrin/discordgo"
)

// TODO: automate these
const fnDiscord = "687565213943332875"
const audioChannel = "715889803937185812"

var audioConnected = false
var voiceQueueChan = make(chan string)

var encoderMutex sync.Mutex
var encoder *mixer.Encoder

func connectVoice(dg *discordgo.Session) {
	audioConnected = true
	dgv, err := dg.ChannelVoiceJoin(fnDiscord, audioChannel, false, true)
	if err != nil {
		fmt.Println(err)
		audioConnected = false
		return
	}

	encoderMutex.Lock()
	encoder = mixer.NewEncoder()
	encoder.VC = dgv
	go encoder.Run()
	encoderMutex.Unlock()

	doneChan := make(chan struct{})
	go func() {
		var i uint64
		for {
			select {
			case f := <-voiceQueueChan:
				fmt.Println(f)
				go encoder.Queue(uint64(i), f)
				// i++
			case <-doneChan:
				return
			}
		}
	}()

	//time.Sleep(5 * time.Second) // waiting for first audio
	//for !encoder.HasFinishedAll() {
	//	time.Sleep(15 * time.Second)
	//}

	// Close connections once all are played
	//doneChan <- struct{}{}
	//dgv.Disconnect()
	//dgv.Close()
	//encoder.Stop()
	//audioConnected = false
}
