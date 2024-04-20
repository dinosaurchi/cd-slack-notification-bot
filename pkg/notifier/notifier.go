package notifier

import (
	"cd-slack-notification-bot/go/pkg/config"
	"cd-slack-notification-bot/go/pkg/matcher"
	"cd-slack-notification-bot/go/pkg/slack"
	prtracker "cd-slack-notification-bot/go/pkg/tracker/pr-tracker"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type State struct {
	// Map from RunID to whether the thread has been notified
	CDNotified map[string]bool `json:"CDNotified"`
	PRNotified map[string]bool `json:"PRNotified"`
}

func RunNotifier(
	state *State,
	matcherState *matcher.State,
	prTrackerState *prtracker.State,
) (*State, error) {
	return RunNotifierCustom(
		state,
		matcherState,
		prTrackerState,
		config.GetConfigDefault().Slack.GithubPRNotificationChannelID,
		config.GetConfigDefault().Slack.CodeBuildNotificationChannelID,
		slack.NewClientDefault(),
	)
}

func RunNotifierCustom(
	state *State,
	matcherState *matcher.State,
	prTrackerState *prtracker.State,
	githubPRNotificationChannelID string,
	codeBuildNotificationChannelID string,
	slackClient slack.Client,
) (*State, error) {
	logrus.Infof("=== Run notifier ====")

	logrus.Infof("Total resolved RunIDs: %v", len(matcherState.ResolvedRunIDs))

	resolvedNotifiedCount := 0

	for runID, matchedResult := range matcherState.ResolvedRunIDs {
		_, ok := state.CDNotified[runID]
		if !ok {
			// Notify the CD
			err := notifyCDThread(slackClient, matchedResult, githubPRNotificationChannelID, codeBuildNotificationChannelID)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			state.CDNotified[runID] = true
		}

		_, ok = state.PRNotified[runID]
		if !ok {
			// Notify the PR
			err := notifyPRThread(slackClient, matchedResult, githubPRNotificationChannelID, codeBuildNotificationChannelID)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			resolvedNotifiedCount++
			state.PRNotified[runID] = true
		}
	}

	logrus.Infof("Notified %v new resolved RunIDs", resolvedNotifiedCount)

	successfulPRs, err := prTrackerState.GetPRsWithSuccessfulCD()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Separate the successful PRs from the PRs that have been notified
	logrus.Infof("")

	logrus.Infof("Total successful PRs: %v", len(successfulPRs))

	successfulPRsNotifiedCount := 0
	for _, prInfo := range successfulPRs {
		if !prInfo.IsSuccessfulCD() {
			// As we assume that the PR is successful, if it does not have a successful CD, we should return an error
			return nil, errors.Errorf("PR %v does not have a successful CD", prInfo.PRNumber)
		}

		codeBuildURL, err := prInfo.GetCodeBuildURL()
		if err != nil {
			return nil, errors.WithStack(err)
		}

		runID, err := slack.GetAWSCodeBuildRunID(codeBuildURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if len(runID) == 0 {
			return nil, errors.Errorf("RunID is empty for PR %v", prInfo.PRNumber)
		}

		_, ok := state.PRNotified[runID]
		if !ok {
			// Notify the PR
			err := notifySuccessfulPRThread(
				slackClient,
				codeBuildURL,
				prInfo.ThreadTimestamp,
				githubPRNotificationChannelID,
			)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			successfulPRsNotifiedCount++
			state.PRNotified[runID] = true
		}
	}

	logrus.Infof("Notified %v new successful PRs", successfulPRsNotifiedCount)

	logrus.Infof("-------------")

	return state, nil
}

func notifyCDThread(
	slackClient slack.Client,
	matchedResult *matcher.MatchedResult,
	githubPRNotificationChannelID string,
	codeBuildNotificationChannelID string,
) error {
	prThreadLink, err := slackClient.GetMessageLink(
		githubPRNotificationChannelID,
		matchedResult.PRThreadTimestamp,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	message, err := GetCDMessage(
		prThreadLink,
		matchedResult.Statuses,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	timestamp, err := slackClient.ReplyThread(
		codeBuildNotificationChannelID,
		matchedResult.CDThreadTimestamp,
		message,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	logrus.Infof("Notified CD thread at timestamp %v of channel %v", timestamp, codeBuildNotificationChannelID)

	return nil
}

func notifyPRThread(
	slackClient slack.Client,
	matchedResult *matcher.MatchedResult,
	githubPRNotificationChannelID string,
	codeBuildNotificationChannelID string,
) error {
	cdThreadLink, err := slackClient.GetMessageLink(
		codeBuildNotificationChannelID,
		matchedResult.CDThreadTimestamp,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	message, err := GetCDMessage(
		cdThreadLink,
		matchedResult.Statuses,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	timestamp, err := slackClient.ReplyThread(
		githubPRNotificationChannelID,
		matchedResult.PRThreadTimestamp,
		message,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	logrus.Infof("Notified PR thread at timestamp %v of channel %v", timestamp, githubPRNotificationChannelID)

	return nil
}

func notifySuccessfulPRThread(
	slackClient slack.Client,
	codeBuildURL string,
	prThreadTimestamp string,
	githubPRNotificationChannelID string,
) error {
	message, err := GetSuccessfulCDMessage(codeBuildURL)
	if err != nil {
		return errors.WithStack(err)
	}
	timestamp, err := slackClient.ReplyThread(
		githubPRNotificationChannelID,
		prThreadTimestamp,
		message,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	logrus.Infof("Notified PR thread at timestamp %v of channel %v", timestamp, githubPRNotificationChannelID)

	return nil
}
