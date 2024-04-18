package cdtracker

import (
	"cd-slack-notification-bot/go/pkg/utils"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
)

func LoadInitialCDTrackerState(
	stateDirPath string,
) (*State, error) {
	statePath := GetCDTrackerStatePath(stateDirPath)
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		// If the file does not exist, create a new state
		const lookBackDuration = time.Hour * 300
		return &State{
			LastFetchedTimestamp: time.Now().Add(-lookBackDuration),
			RunIDToCDs:           map[string]*CDInfo{},
		}, nil
	}

	// Otherwise, load the state from the file
	cdTrackerState := &State{}
	err := utils.LoadFromFile(statePath, cdTrackerState)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return cdTrackerState, nil
}

func GetCDTrackerStatePath(
	stateDirPath string,
) string {
	return path.Join(stateDirPath, "cd-tracker.json")
}
