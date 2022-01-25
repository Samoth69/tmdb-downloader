package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/urfave/cli/v2"
)

// struct containing app settings
// this struct will be filled by reading program args
type structConf struct {
	TMDB_Id      int
	DownloadPath string
	KeyPath      string
}

//check if the provided path is valid
//it will convert the path to absolute if needed
//return true if valid,
//false otherwise with error
func checkPath(path *string) (bool, error) {
	if !filepath.IsAbs(*path) {
		absPath, err := filepath.Abs(*path)
		if err != nil {
			return false, err
		}
		*path = absPath
	}
	if _, err := os.Stat(*path); err != nil {
		return false, err
	}
	return true, nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app := &cli.App{
		Name:                   "MetadataDownloader",
		Usage:                  "Download images from tmdb",
		UseShortOptionHandling: true,
		UsageText:              "metadatadownloader [global options] tmdb_id",
		HideHelpCommand:        true,
		Version:                "1.0",
		Authors: []*cli.Author{
			&cli.Author{
				Name: "Samoth69",
			},
		},
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Required: false,
				Usage:    "Output folder",
				//Destination: &appConfig.DownloadPath,
			},
			&cli.PathFlag{
				Name:     "keyfile",
				Aliases:  []string{"key", "k"},
				Required: false,
				Usage:    "Path to file containing tmdb api key",
				//Destination: &appConfig.KeyPath,
			},
		},
		Action: func(c *cli.Context) error {
			appConfig := new(structConf)

			// checking download path
			if dlPath := c.String("output"); dlPath == "" {
				appConfig.DownloadPath, _ = os.Getwd()
			} else {
				appConfig.DownloadPath = dlPath
			}
			_, err := checkPath(&appConfig.DownloadPath)
			if err != nil {
				return cli.Exit(fmt.Sprintf("Error on download path: %s", err), -1)
			}

			// checking keyfile path
			if kp := c.String("keyfile"); kp == "" {
				path, _ := os.Getwd()
				appConfig.KeyPath = filepath.Join(path, "api_keys.json")
			} else {
				appConfig.KeyPath = kp
			}
			_, err = checkPath(&appConfig.KeyPath)
			if err != nil {
				return cli.Exit(fmt.Sprintf("Error on key path: %s", err), -1)
			}

			//checking tmdb id value
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

			LoadKeys(appConfig.KeyPath)

			item := GetLinks(appConfig.TMDB_Id, Keys.TmdbKeyV4)
			DownloadFiles(item, appConfig.DownloadPath)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		return
	}
}
