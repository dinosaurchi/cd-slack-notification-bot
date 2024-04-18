package prtracker

import (
	"cd-slack-notification-bot/go/pkg/config"
	"cd-slack-notification-bot/go/pkg/github"
	"cd-slack-notification-bot/go/pkg/slack"
	"cd-slack-notification-bot/go/pkg/utils"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type State struct {
	LastFetchedTimestamp time.Time `json:"lastFetchedTimestamp"`
	// Mapping from threadTimestamp to PR info
	PRs map[string]*PRInfo `json:"PRs"`
}

type PRInfo struct {
	PRNumber uint64       `json:"prNumber"`
	Statuses []statusInfo `json:"statuses"`
}

type statusInfo struct {
	State        string `json:"state"`
	CodeBuildURL string `json:"codeBuildURL"`
}

func RunPRTracker(
	state *State,
	upTo time.Time,
) (*State, error) {
	logrus.Infof("Last fetched timestamp: %v", state.LastFetchedTimestamp.String())

	logrus.Infof("Fetching new PR statuses up to %v", upTo.String())
	state, err := fetchNewPRs(state, upTo)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	logrus.Infof("Updating fetched and not resolved PRs")
	var resolvedCount int
	state, resolvedCount, err = updateFetchedNotResolvedPRs(state)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	logrus.Infof("Total PRs: %v / %v", resolvedCount, len(state.PRs))
	logrus.Infof("-------------")

	return state, nil
}

func isResolved(
	pr *PRInfo,
) bool {
	return len(pr.Statuses) > 0
}

func fetchNewPRs(
	state *State,
	upTo time.Time,
) (*State, error) {
	// Fetch new messages to check for new PR's threads
	slackClient := slack.NewClientDefault()
	messages, err := slackClient.RetrieveChannelHistory(
		config.GetConfigDefault().Slack.GithubPRNotificationChannelID,
		state.LastFetchedTimestamp,
		upTo,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	githubClient := github.NewClientDefault()

	newTimestamps := []time.Time{}
	for _, message := range messages {
		var githubPRInfo *slack.GithubPRInfo
		githubPRInfo, err = slack.ParseGithubPRInfoFromPROpenedMessage(message)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if githubPRInfo == nil {
			continue
		}
		if githubPRInfo.RepoName != config.GetConfigDefault().Github.RepoName {
			continue
		}
		if githubPRInfo.RepoOwner != config.GetConfigDefault().Github.RepoOwner {
			continue
		}

		newTimestamps = append(newTimestamps, utils.ConvertTimestampStringToTime(githubPRInfo.ThreadTimestamp))

		// Initialize PR info if not exists
		if _, ok := state.PRs[githubPRInfo.ThreadTimestamp]; !ok {
			state.PRs[githubPRInfo.ThreadTimestamp] = &PRInfo{
				PRNumber: githubPRInfo.PRNumber,
				Statuses: []statusInfo{},
			}
		}

		curPR := state.PRs[githubPRInfo.ThreadTimestamp]
		if isResolved(curPR) {
			// The PR info is already resolved as it has statuses
			continue
		}

		var cdInfo *github.CDInfo
		cdInfo, err = githubClient.GetPullRequestCDInfo(githubPRInfo.PRNumber)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		for _, status := range cdInfo.Statuses {
			curPR.Statuses = append(
				curPR.Statuses,
				statusInfo{
					State:        status.State,
					CodeBuildURL: status.TargetURL,
				},
			)
		}

		// Update PR info
		state.PRs[githubPRInfo.ThreadTimestamp] = curPR
	}

	maxThreadTimestamp, err := utils.MaxSlice[time.Time](
		newTimestamps,
		func(a, b time.Time) (bool, error) {
			return a.After(b), nil
		},
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Update the last fetched timestamp
	state.LastFetchedTimestamp = maxThreadTimestamp

	return state, nil
}

func updateFetchedNotResolvedPRs(
	state *State,
) (*State, int, error) {
	githubClient := github.NewClientDefault()

	resolvedCount := 0

	// Update the not resolved PRs
	for threadTimestamp, pr := range state.PRs {
		if isResolved(pr) {
			resolvedCount++
			continue
		}

		cdInfo, err := githubClient.GetPullRequestCDInfo(pr.PRNumber)
		if err != nil {
			return nil, -1, errors.WithStack(err)
		}

		newStatuses := []statusInfo{}
		for _, status := range cdInfo.Statuses {
			newStatuses = append(
				newStatuses,
				statusInfo{
					State:        status.State,
					CodeBuildURL: status.TargetURL,
				},
			)
		}

		if len(newStatuses) > 0 {
			pr.Statuses = newStatuses
			state.PRs[threadTimestamp] = pr
		}
	}

	return state, resolvedCount, nil
}
