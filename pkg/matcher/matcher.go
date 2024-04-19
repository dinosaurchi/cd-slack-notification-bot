package matcher

import (
	cdtracker "cd-slack-notification-bot/go/pkg/tracker/cd-tracker"
	prtracker "cd-slack-notification-bot/go/pkg/tracker/pr-tracker"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type State struct {
	ResolvedRunIDs map[string]*MatchedResult `json:"resolvedRunIDs"`
}

type MatchedResult struct {
	CDThreadTimestamp string   `json:"cdThreadTimestamp"`
	PRThreadTimestamp string   `json:"prThreadTimestamp"`
	PRNumber          uint64   `json:"prNumber"`
	Statuses          []string `json:"statuses"`
}

func RunMatcher(
	state *State,
	cdTrackerState *cdtracker.State,
	prTrackerState *prtracker.State,
) (*State, error) {
	runIDsToPRs, err := mapRunIDsToPRs(prTrackerState)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for runID, cdInfo := range cdTrackerState.RunIDToCDs {
		_, ok := state.ResolvedRunIDs[runID]
		if ok {
			// The RunID is already resolved, skip it
			continue
		}

		prInfo, ok := runIDsToPRs[runID]
		if !ok {
			// We do not return the error here because there maybe some deplay between
			// the PR and the CD, so we just log the error and continue, wait for the next
			// run to see if the PR thread is created.
			logrus.Warnln("RunID", runID, "does not have a corresponding PR")
			continue
		}

		state.ResolvedRunIDs[runID] = &MatchedResult{
			CDThreadTimestamp: cdInfo.ThreadTimestamp,
			PRThreadTimestamp: prInfo.ThreadTimestamp,
			PRNumber:          prInfo.PRNumber,
			Statuses:          toStateStrings(prInfo),
		}
	}

	return state, nil
}

func toStateStrings(prInfo *prtracker.PRInfo) []string {
	statuses := []string{}
	for _, status := range prInfo.Statuses {
		statuses = append(statuses, status.State)
	}
	return statuses
}

func mapRunIDsToPRs(prTrackerState *prtracker.State) (map[string]*prtracker.PRInfo, error) {
	runIDsToPRs := map[string]*prtracker.PRInfo{}
	for _, prInfo := range prTrackerState.PRs {
		prRunID, err := GetRunIDFromPRInfo(prInfo)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		_, ok := runIDsToPRs[prRunID]
		if ok {
			return nil, errors.Errorf("runID %s is duplicated", prRunID)
		}
		runIDsToPRs[prRunID] = prInfo
	}
	return runIDsToPRs, nil
}