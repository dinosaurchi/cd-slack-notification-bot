package matcher

import (
	"cd-slack-notification-bot/go/pkg/slack"
	matcher "cd-slack-notification-bot/go/pkg/tracker/pr-tracker"
	"cd-slack-notification-bot/go/pkg/utils"
	"os"
	"path"

	"github.com/pkg/errors"
)

// Try to get the run ID from the PR info.
//   - If the PR is not merged yet, it wont have status info, thus
//     we return an empty string.
func GetRunIDFromPRInfo(
	prInfo *matcher.PRInfo,
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

func LoadInitialMatcherState(
	stateDirPath string,
) (*State, error) {
	statePath := GetMatcherStatePath(stateDirPath)
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		// If the file does not exist, create a new state
		return &State{
			ResolvedRunIDs: map[string]*MatchedResult{},
		}, nil
	}

	// Otherwise, load the state from the file
	matcherState := &State{}
	err := utils.LoadFromFile(statePath, matcherState)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return matcherState, nil
}

func GetMatcherStatePath(
	stateDirPath string,
) string {
	return path.Join(stateDirPath, "matcher.json")
}
