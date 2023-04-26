package slack

import (
	"fmt"

	"github.com/mylittleboxy/backup/pkg/configType"
	"github.com/slack-go/slack"
)

func SendSlackMessage(conf configType.Config, message string) error {
	api := slack.New(conf.Slack.APIToken)
	channelID := conf.Slack.ChannelID // replace this with the ID of the channel you want to send a message to

	_, _, err := api.PostMessage(channelID, slack.MsgOptionText(message, false))
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}
	return nil
}
