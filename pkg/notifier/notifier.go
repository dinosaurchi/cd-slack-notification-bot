package notifier

import (
	"cd-slack-notification-bot/go/pkg/config"
	"cd-slack-notification-bot/go/pkg/matcher"
	"cd-slack-notification-bot/go/pkg/slack"

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
) (*State, error) {
	return RunNotifierCustom(
		state,
		matcherState,
		config.GetConfigDefault().Slack.GithubPRNotificationChannelID,
		config.GetConfigDefault().Slack.CodeBuildNotificationChannelID,
		slack.NewClientDefault(),
	)
}

func RunNotifierCustom(
	state *State,
	matcherState *matcher.State,
	githubPRNotificationChannelID string,
	codeBuildNotificationChannelID string,
	slackClient slack.Client,
) (*State, error) {
	logrus.Infof("=== Run notifier ====")

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
			state.PRNotified[runID] = true
		}
	}

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
