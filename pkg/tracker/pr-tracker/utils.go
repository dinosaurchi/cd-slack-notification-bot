package prtracker

import (
	"cd-slack-notification-bot/go/pkg/utils"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
)

func LoadInitialPRTrackerState(
	stateDirPath string,
) (*State, error) {
	statePath := path.Join(stateDirPath, "pr-tracker.json")
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		// If the file does not exist, create a new state
		const lookBackDuration = time.Hour * 200
		return &State{
			LastFetchedTimestamp: time.Now().Add(-lookBackDuration),
			PRs:                  map[string]*PRInfo{},
		}, nil
	}

	// Otherwise, load the state from the file
	prTrackerState := &State{}
	err := utils.LoadFromFile(statePath, prTrackerState)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return prTrackerState, nil
}