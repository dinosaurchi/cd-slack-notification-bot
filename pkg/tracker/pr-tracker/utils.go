package prtracker

import (
	"cd-slack-notification-bot/go/pkg/slack"
	"cd-slack-notification-bot/go/pkg/utils"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
)

func LoadInitialPRTrackerState(
	stateDirPath string,
	lookBackDuration time.Duration,
) (*State, error) {
	statePath := GetPRTrackerStatePath(stateDirPath)
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		// If the file does not exist, create a new state
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

// Try to get the run ID from the PR info.
//   - If the PR is not merged yet, it wont have status info, thus
//     we return an empty string.
func GetRunIDFromPRInfo(
	prInfo *PRInfo,
) (string, error) {
	for _, status := range prInfo.Statuses {
		if status.CodeBuildURL != "" {
			runID, err := slack.GetAWSCodeBuildRunID(status.CodeBuildURL)
			if err != nil {
				return "", errors.WithStack(err)
			}
			return runID, nil
		}
	}
	return "", nil
}

func GetPRTrackerStatePath(
	stateDirPath string,
) string {
	return path.Join(stateDirPath, "pr-tracker.json")
}
