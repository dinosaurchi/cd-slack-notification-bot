package notifier

import (
	"cd-slack-notification-bot/go/pkg/utils"
	"os"
	"path"

	"github.com/pkg/errors"
)

func LoadInitialNotifierState(
	stateDirPath string,
) (*State, error) {
	statePath := GetNotifierStatePath(stateDirPath)
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		// If the file does not exist, create a new state
		return &State{
			CDNotified: map[string]bool{},
			PRNotified: map[string]bool{},
		}, nil
	}

	// Otherwise, load the state from the file
	notifierState := &State{}
	err := utils.LoadFromFile(statePath, notifierState)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return notifierState, nil
}

func GetNotifierStatePath(
	stateDirPath string,
) string {
	return path.Join(stateDirPath, "notifier.json")
}
