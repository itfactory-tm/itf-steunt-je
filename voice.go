package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

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
	if audioConnected {
		// we're done here
		return
	}
	audioConnected = true
	dgv, err := dg.ChannelVoiceJoin(fnDiscord, audioChannel, false, true)
	if err != nil {
		log.Println(err)
		audioConnected = false
		dgv.Disconnect()
		// keep hitting yourself till you connect
		connectVoice(dg)
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
				dg.UpdateStreamingStatus(0, f, "")
				go encoder.Queue(uint64(i), f)
				// i++
			case <-doneChan:
				return
			}
		}
	}()

	go func() {
		failCount := 0
		for {
			time.Sleep(5 * time.Second)
			if !dgv.Ready {
				failCount++
			} else {
				failCount = 0
			}

			if failCount > 5 {
				log.Println("Connection check failed 5 times")
				os.Exit(1)
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
