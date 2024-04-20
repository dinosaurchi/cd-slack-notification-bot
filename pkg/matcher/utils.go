package matcher

import (
	"cd-slack-notification-bot/go/pkg/utils"
	"os"
	"path"

	"github.com/pkg/errors"
)

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
