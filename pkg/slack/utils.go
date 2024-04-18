package slack

import (
	"cd-slack-notification-bot/go/pkg/aws"
	"cd-slack-notification-bot/go/pkg/utils"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

func ParseRunIDFromMessage(
	message slack.Message,
) (string, error) {
	blockString, err := utils.ToJSONString(message.Attachments[0].Blocks.BlockSet)
	if err != nil {
		return "", errors.WithStack(err)
	}

	blocks := []any{}
	err = json.Unmarshal([]byte(blockString), &blocks)
	if err != nil {
		return "", errors.WithStack(err)
	}

	const minRunIDLength = 10
	for _, block := range blocks {
		blockType, err := getBlockType(block)
		if err != nil {
			return "", errors.WithStack(err)
		}

		if blockType == "section" {
			runID, err := getRunIDFromSectionBlock(block)
			if err != nil {
				return "", errors.WithStack(err)
			}
			if len(runID) > minRunIDLength {
				return runID, nil
			}
		} else if blockType == "actions" {
			// Skip, not supported for now
			continue
		}
	}

	// No run ID found
	return "", nil
}

func getBlockType(
	block any,
) (string, error) {
	type blockSet struct {
		Type string `json:"type"`
	}
	blockString, err := utils.ToJSONString(block)
	if err != nil {
		return "", errors.WithStack(err)
	}

	res := blockSet{}
	err = json.Unmarshal([]byte(blockString), &res)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return res.Type, nil
}

func getRunIDFromSectionBlock(
	block any,
) (string, error) {
	type blockSection struct {
		Text struct {
			// Type string `json:"type"`
			Text string `json:"text"`
		} `json:"text"`
		// Type string `json:"type"`
		// BlockID string `json:"block_id"`
	}

	blockString, err := utils.ToJSONString(block)
	if err != nil {
		return "", errors.WithStack(err)
	}

	res := blockSection{}
	err = json.Unmarshal([]byte(blockString), &res)
	if err != nil {
		return "", errors.WithStack(err)
	}

	runID, err := aws.GetAWSCodeSuiteRunIDFromMessage(res.Text.Text)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return runID, nil
}
