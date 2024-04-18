package prtracker

import (
	"cd-slack-notification-bot/go/pkg/config"
	"cd-slack-notification-bot/go/pkg/github"
	"cd-slack-notification-bot/go/pkg/slack"
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
	logrus.Infof("Fetching new PR statuses up to %v", upTo.String())
	logrus.Infof("Last fetched timestamp: %v", state.LastFetchedTimestamp.String())

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

	for _, message := range messages {
		githubPRInfo, err := slack.ParseGithubPRInfoFromPROpenedMessage(message)
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

		// Initialize PR info if not exists
		if _, ok := state.PRs[githubPRInfo.ThreadTimestamp]; !ok {
			state.PRs[githubPRInfo.ThreadTimestamp] = &PRInfo{
				PRNumber: githubPRInfo.PRNumber,
				Statuses: []statusInfo{},
			}
		}

		curPR := state.PRs[githubPRInfo.ThreadTimestamp]
		if len(curPR.Statuses) > 0 {
			// The PR info is already resolved as it has statuses
			continue
		}

		cdInfo, err := githubClient.GetPullRequestCDInfo(githubPRInfo.PRNumber)
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

	logrus.Infof("Total PRs: %v / %v", countResolvedPRs(state.PRs), len(state.PRs))
	logrus.Infof("-------------")

	return state, nil
}

func countResolvedPRs(
	prs map[string]*PRInfo,
) int {
	count := 0
	for _, pr := range prs {
		if len(pr.Statuses) > 0 {
			count++
		}
	}
	return count
}
