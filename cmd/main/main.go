package main

import (
	"cd-slack-notification-bot/go/pkg/config"
	cdtracker "cd-slack-notification-bot/go/pkg/tracker/cd-tracker"
	prtracker "cd-slack-notification-bot/go/pkg/tracker/pr-tracker"
	"cd-slack-notification-bot/go/pkg/utils"
	"os"
	"path"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	// Use repo owner and repo name to create a state directory
	// to avoid conflicts between different repositories
	stateDirPath := path.Join(
		".",
		"state",
		config.GetConfigDefault().Github.RepoOwner,
		config.GetConfigDefault().Github.RepoName,
	)
	err = os.MkdirAll(stateDirPath, 0755)
	if err != nil {
		panic(err)
	}

	const lookBackDuration = time.Hour * 500

	prTrackerState, err := prtracker.LoadInitialPRTrackerState(stateDirPath, lookBackDuration)
	if err != nil {
		panic(err)
	}

	cdTrackerState, err := cdtracker.LoadInitialCDTrackerState(stateDirPath, lookBackDuration)
	if err != nil {
		panic(err)
	}

	const waitTime = time.Minute * 3
	const waitTimeForError = time.Minute * 6

	for {
		curNow := time.Now()
		err := runTrackers(prTrackerState, cdTrackerState, stateDirPath, curNow)
		if err != nil {
			logrus.Errorf("Error: %v\n", err)
			time.Sleep(waitTimeForError)
		} else {
			time.Sleep(waitTime)
		}
	}
}

func runTrackers(
	prTrackerState *prtracker.State,
	cdTrackerState *cdtracker.State,
	stateDirPath string,
	curNow time.Time,
) error {
	// Run PR tracker
	prTrackerState, err := prtracker.RunPRTracker(prTrackerState, curNow)
	if err != nil {
		return errors.Errorf("error running PR Tracker: %v", err)
	}

	err = utils.DumpToFile(prtracker.GetPRTrackerStatePath(stateDirPath), prTrackerState)
	if err != nil {
		return errors.Errorf("dump PR Tracker file error: %v", err)
	}

	// Run CD tracker
	cdTrackerState, err = cdtracker.RunCDTracker(cdTrackerState, curNow)
	if err != nil {
		return errors.Errorf("error running CD Tracker: %v", err)
	}

	err = utils.DumpToFile(cdtracker.GetCDTrackerStatePath(stateDirPath), cdTrackerState)
	if err != nil {
		return errors.Errorf("dump CD Tracker file error: %v", err)
	}

	return nil
}
