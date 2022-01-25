package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
)

// struct containing app settings
// this struct will be filled by reading program args
type structConf struct {
	TMDB_Id int
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app := &cli.App{
		Name:                   "MetadataDownloader",
		Usage:                  "Download images from tmdb",
		UseShortOptionHandling: true,
		Compiled:               time.Now(),
		UsageText:              "metadatadownloader [global options] tmdb_id",
		HideHelpCommand:        true,
		Authors: []*cli.Author{
			&cli.Author{
				Name: "Samoth69",
			},
		},
		Action: func(c *cli.Context) error {
			appConfig := new(structConf)

			inputTmdbIdValue := c.Args().Get(0)
			if len(inputTmdbIdValue) == 0 {
				return cli.Exit("please provide a tmdb_id", -1)
			}
			//check if first arg can be converted to an int
			if tmdbId, err := strconv.Atoi(inputTmdbIdValue); err == nil {
				if tmdbId > 0 {
					appConfig.TMDB_Id = tmdbId
				} else {
					return cli.Exit("invalid tmdb_id, should be a positive number", -1)
				}
			} else {
				return cli.Exit(fmt.Sprintf("%s is invalid, should be a positive number", inputTmdbIdValue), -1)
			}

			item := GetLinks(appConfig.TMDB_Id, Keys.TmdbKeyV4)
			DownloadFiles(item)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		return
	}
}
