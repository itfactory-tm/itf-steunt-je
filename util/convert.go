package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var latestNum int64
var wavRegex = regexp.MustCompile(`.*\.wav`)

func main() {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	unhandledFiles := []string{}
	for _, f := range files {
		if wavRegex.MatchString(f.Name()) {
			// is WAV
			fileNameParts := strings.Split(f.Name(), ".")
			if len(fileNameParts) <= 1 {
				// invalid file
				continue
			}
			i, _ := strconv.ParseInt(fileNameParts[0], 10, 64)
			if i > latestNum {
				latestNum = i
			}
		} else if f.Name() != "convert.go" {
			// a file to convert
			unhandledFiles = append(unhandledFiles, f.Name())
		}
	}

	for _, file := range unhandledFiles {
		convert(file)
		fmt.Println("converted", file)
	}
}

func convert(file string) {
	latestNum++
	cmd := exec.Command("ffmpeg", "-i", file, "-ac", "1", "-ar", "48000", fmt.Sprintf("%02d.wav", latestNum))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	os.Remove(file)
}
