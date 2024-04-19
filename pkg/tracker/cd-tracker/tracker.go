package cdtracker

import (
	"cd-slack-notification-bot/go/pkg/config"
	"cd-slack-notification-bot/go/pkg/slack"
	"cd-slack-notification-bot/go/pkg/utils"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type State struct {
	LastFetchedTimestamp time.Time `json:"lastFetchedTimestamp"`
	// Mapping from runID to CD info
	RunIDToCDs map[string]*CDInfo `json:"runIDToCDs"`
}

type CDInfo struct {
	ThreadTimestamp string `json:"threadTimestamp"`
	RunID           string `json:"runID"`
}

func RunCDTracker(
	state *State,
	upTo time.Time,
) (*State, error) {
	logrus.Infof("=== Run CD tracker ====")
	logrus.Infof("Last fetched timestamp: %v", state.LastFetchedTimestamp.String())
	logrus.Infof("Fetching new CD statuses up to %v", upTo.String())
	slackClient := slack.NewClientDefault()
	messages, err := slackClient.RetrieveChannelHistory(
		config.GetConfigDefault().Slack.CodeBuildNotificationChannelID,
		state.LastFetchedTimestamp,
		upTo,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	logrus.Infof("Fetched %v messages", len(messages))
	hasRunIDCount := 0

	allTimestamps := []time.Time{}
	for _, message := range messages {
		var runID string
		runID, err = slack.ParseRunIDFromCodeBuildMessage(message)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		threadTimestamp := slack.GetThreadTimestamp(message.ThreadTimestamp, message.Timestamp)
		allTimestamps = append(allTimestamps, utils.ConvertTimestampStringToTime(threadTimestamp))

		if runID != "" {
			if _, ok := state.RunIDToCDs[runID]; ok {
				continue
			}

			hasRunIDCount++
			state.RunIDToCDs[runID] = &CDInfo{
				ThreadTimestamp: threadTimestamp,
				RunID:           runID,
			}
		}
	}

	logrus.Infof("Found %v new messages with runID", hasRunIDCount)
	logrus.Infof("-------------")

	maxTimestamp, err := utils.MaxSlice[time.Time](
		allTimestamps,
		func(a, b time.Time) (bool, error) {
			return a.After(b), nil
		},
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Update the last fetched timestamp
	state.LastFetchedTimestamp = maxTimestamp

	return state, nil
}
