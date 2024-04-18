package main

import (
	prtracker "cd-slack-notification-bot/go/pkg/tracker/pr-tracker"
	"cd-slack-notification-bot/go/pkg/utils"
	"os"
	"path"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	stateDirPath := path.Join(".", "state")
	err = os.MkdirAll(stateDirPath, 0755)
	if err != nil {
		panic(err)
	}

	prTrackerState, err := prtracker.LoadInitialPRTrackerState(stateDirPath)
	if err != nil {
		panic(err)
	}

	const waitTime = time.Second * 10
	const waitTimeForError = time.Second * 120

	for {
		logrus.Println("Starting PR tracker")
		curNow := time.Now()
		prTrackerState, err = prtracker.RunPRTracker(prTrackerState, curNow)
		if err != nil {
			logrus.Printf("Error: %v\n", err)
			time.Sleep(waitTimeForError)
		} else {
			err = utils.DumpToFile(path.Join(stateDirPath, "pr-tracker.json"), prTrackerState)
			if err != nil {
				logrus.Printf("Dump file error: %v\n", err)
			}
			time.Sleep(waitTime)
		}
	}
}
