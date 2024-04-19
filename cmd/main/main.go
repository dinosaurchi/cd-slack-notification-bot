package main

import (
	"cd-slack-notification-bot/go/pkg/config"
	"cd-slack-notification-bot/go/pkg/matcher"
	"cd-slack-notification-bot/go/pkg/notifier"
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

	const lookBackDurationPRTracker = time.Hour * 500
	const lookBackDurationCDTracker = time.Hour * 24

	prTrackerState, err := prtracker.LoadInitialPRTrackerState(stateDirPath, lookBackDurationPRTracker)
	if err != nil {
		panic(err)
	}

	cdTrackerState, err := cdtracker.LoadInitialCDTrackerState(stateDirPath, lookBackDurationCDTracker)
	if err != nil {
		panic(err)
	}

	matcherState, err := matcher.LoadInitialMatcherState(stateDirPath)
	if err != nil {
		panic(err)
	}

	notifierState, err := notifier.LoadInitialNotifierState(stateDirPath)
	if err != nil {
		panic(err)
	}

	const waitTime = time.Minute * 3
	const waitTimeForError = time.Minute * 6

	for {
		curNow := time.Now()

		err := runAlls(prTrackerState, cdTrackerState, matcherState, notifierState, stateDirPath, curNow)
		if err != nil {
			logrus.Errorf("Error: %v\n", err)
			time.Sleep(waitTimeForError)
		} else {
			time.Sleep(waitTime)
		}
	}
}

func runAlls(
	prTrackerState *prtracker.State,
	cdTrackerState *cdtracker.State,
	matchState *matcher.State,
	notifierState *notifier.State,
	stateDirPath string,
	curNow time.Time,
) error {
	logrus.Infof("**** Start running at %v ****\n", curNow.Format(time.RFC3339))

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

	// Run matcher
	matchState, err = matcher.RunMatcher(matchState, cdTrackerState, prTrackerState)
	if err != nil {
		return errors.Errorf("error running Matcher: %v", err)
	}

	err = utils.DumpToFile(matcher.GetMatcherStatePath(stateDirPath), matchState)
	if err != nil {
		return errors.Errorf("dump Matcher file error: %v", err)
	}

	// Run notifier
	notifierState, err = notifier.RunNotifier(notifierState, matchState)
	if err != nil {
		logrus.Errorf("error running Notifier: %v", err)
	}

	if err == nil {
		// Only dump the notifier state if there is no error
		err = utils.DumpToFile(notifier.GetNotifierStatePath(stateDirPath), notifierState)
		if err != nil {
			return errors.Errorf("dump Notifier file error: %v", err)
		}
	}

	logrus.Infof("**** End running at %v ****\n\n", time.Now().Format(time.RFC3339))

	return nil
}
