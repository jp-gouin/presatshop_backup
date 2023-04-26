package main

import (
	"fmt"

	"github.com/mylittleboxy/backup/pkg/configType"
	"github.com/mylittleboxy/backup/pkg/dump"
	"github.com/mylittleboxy/backup/pkg/s3"
	"github.com/mylittleboxy/backup/pkg/slack"
	"github.com/robfig/cron"
)

func main() {

	config := configType.Config{
		DB:    configType.ConfigDB{},
		S3:    configType.ConfigS3{},
		Slack: configType.ConfigSlack{},
	}
	c := cron.New()

	// Schedule the task to run every day at midnight
	c.AddFunc("@midnight", func() {
		filename, err := dump.Dump(config)
		if err != nil {
			err = slack.SendSlackMessage(config, fmt.Sprintf("Error while saving dump : %v", err))
		} else {
			err = s3.SendFile(config, filename)
			if err != nil {
				err = slack.SendSlackMessage(config, fmt.Sprintf("Error while uploading dump to s3 : %v", err))
			} else {
				slack.SendSlackMessage(config, "Dump save into S3")
			}
		}
	})
	c.Start()

	select {}

}
