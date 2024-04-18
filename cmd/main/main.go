package main

import (
	cdtracker "cd-slack-notification-bot/go/pkg/tracker/cd-tracker"
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

	cdTrackerState, err := cdtracker.LoadInitialCDTrackerState(stateDirPath)
	if err != nil {
		panic(err)
	}

	const waitTime = time.Second * 10
	const waitTimeForError = time.Second * 120

	for {
		curNow := time.Now()

		// Run PR tracker
		prTrackerState, err = prtracker.RunPRTracker(prTrackerState, curNow)
		if err != nil {
			logrus.Printf("Error: %v\n", err)
			time.Sleep(waitTimeForError)
		} else {
			err = utils.DumpToFile(prtracker.GetPRTrackerStatePath(stateDirPath), prTrackerState)
			if err != nil {
				logrus.Printf("Dump file error: %v\n", err)
			}
			time.Sleep(waitTime)
		}

		// Run CD tracker
		cdTrackerState, err = cdtracker.RunCDTracker(cdTrackerState, curNow)
		if err != nil {
			logrus.Printf("Error: %v\n", err)
			time.Sleep(waitTimeForError)
		} else {
			err = utils.DumpToFile(cdtracker.GetCDTrackerStatePath(stateDirPath), cdTrackerState)
			if err != nil {
				logrus.Printf("Dump file error: %v\n", err)
			}
			time.Sleep(waitTime)
		}
	}
}
