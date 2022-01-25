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
	TMDB_Id        int
	AnswerYes      bool
	AnswerNo       bool
	MaxThreadCount int
}

//load os provided args
func LoadArgs(settings *structConf) (bool, error) {
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
			if settings.AnswerYes && settings.AnswerNo {
				return cli.Exit("Only one of --yes, --no or none should be provided, not both --yes and --no", -2)
			}

			inputTmdbIdValue := c.Args().Get(0)
			//check if first arg can be converted to an int
			if tmdbId, err := strconv.Atoi(inputTmdbIdValue); err == nil {
				if tmdbId > 0 {
					settings.TMDB_Id = tmdbId
				} else {
					return cli.Exit("invalid tmdb_id, should be a positive number", -3)
				}
			} else {
				return cli.Exit(fmt.Sprintf("%s is invalid, should be a positive number", inputTmdbIdValue), -1)
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "yes", Aliases: []string{"y"}, Value: false, Usage: "Automatically answer Yes to questions", Destination: &settings.AnswerYes},
			&cli.BoolFlag{Name: "no", Aliases: []string{"n"}, Value: false, Usage: "Automatically answer No to questions", Destination: &settings.AnswerNo},
			&cli.IntFlag{Name: "threads", Aliases: []string{"t"}, Value: 4, Usage: "Max number of downloading threads", Destination: &settings.MaxThreadCount},
		},
	}

	err := app.Run(os.Args)
	success := err == nil

	return success, err
}

func main() {
	appConfig := new(structConf)
	_, err := LoadArgs(appConfig)
	if err != nil {
		log.Fatal(err)
		return
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	//fmt.Printf("%#v\n", appConfig)
	//fmt.Printf("%#v\n", Keys)

	item := GetLinks(appConfig.TMDB_Id, Keys.TmdbKeyV4)
	DownloadFiles(item, appConfig.MaxThreadCount)
}
