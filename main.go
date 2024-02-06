package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

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
		archiveName := fmt.Sprintf("%s/%s-%d.tar.gz", config.DB.DumpDir, "Dump_prestashop", time.Now().Unix())
		var filenamesToArchive []string

		// Dump mysql database
		filename, err := dump.Dump(config)

		if err != nil {
			err = slack.SendSlackMessage(config, fmt.Sprintf("Error while saving dump : %v", err))
		} else {
			filenamesToArchive = append(filenamesToArchive, filename)
			err := filepath.Walk("/prestashop",
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !info.IsDir() {
						filenamesToArchive = append(filenamesToArchive, path)
					}
					return nil
				})
			err = s3.SendFiles(config, filenamesToArchive, archiveName)
			if err != nil {
				err = slack.SendSlackMessage(config, fmt.Sprintf("Error while uploading dump to s3 : %v", err))
			} else {
				slack.SendSlackMessage(config, fmt.Sprintf("Dump save into S3 %s", archiveName))
			}
		}
	})
	c.Start()
	c.Run()
	select {}
}
