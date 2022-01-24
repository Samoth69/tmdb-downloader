package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

// struct containing app settings
// this struct will be filled by reading program args
type structConf struct {
	TMDB_ID    int
	ANILIST_ID int
	Verbose    bool
	AnswerYes  bool
	AnswerNo   bool
}

//load os provided args
func LoadArgs(settings *structConf) (success bool, err error) {
	app := &cli.App{
		Name:                   "MetadataDownloader",
		Usage:                  "Download metadata for specified media id, metadata include images, banner...",
		UseShortOptionHandling: true,
		Compiled:               time.Now(),
		UsageText:              "metadatadownloader [global options]",
		HideHelpCommand:        true,
		Authors: []*cli.Author{
			&cli.Author{
				Name: "Samoth69",
			},
		},
		Action: func(c *cli.Context) error {
			if settings.ANILIST_ID == -1 && settings.TMDB_ID == -1 {
				return cli.Exit("--anilist-id and/or --tmdb-id should be provided", -1)
			}

			if settings.AnswerYes && settings.AnswerNo {
				return cli.Exit("Only one of --yes, --no or none should be provided, not both --yes and --no", -2)
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "tmdb-id",
				Aliases:     []string{"tid"},
				Value:       -1,
				Usage:       "https://www.themoviedb.org/ - TMDB show id",
				Destination: &settings.TMDB_ID,
			},
			&cli.IntFlag{
				Name:        "anilist-id",
				Aliases:     []string{"aid"},
				Value:       -1,
				Usage:       "https://anilist.co/ - AniList show id",
				Destination: &settings.ANILIST_ID,
			},
			&cli.BoolFlag{Name: "verbose", Aliases: []string{"v"}, Value: false, Usage: "App will be more chatty", Destination: &settings.Verbose},
			&cli.BoolFlag{Name: "yes", Aliases: []string{"y"}, Value: false, Usage: "Automatically answer Yes to questions", Destination: &settings.AnswerYes},
			&cli.BoolFlag{Name: "no", Aliases: []string{"n"}, Value: false, Usage: "Automatically answer No to questions", Destination: &settings.AnswerNo},
		},
	}

	err = app.Run(os.Args)
	success = err == nil

	return
}

func main() {
	appConfig := new(structConf)
	_, err := LoadArgs(appConfig)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("%#v", appConfig)
}
