package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type ApiKeys struct {
	TmdbKeyV4 string
}

// contains api keys to be used for communicating with tmdb/anilist API
// thoses are sensible data, should be handled with care
var Keys ApiKeys

func LoadKeys(path string) {
	jsonFile, err := os.Open(path)
	if err != nil {
		text, _ := json.MarshalIndent(Keys, "", "  ")
		os.WriteFile(path, text, 0755)
		log.Printf("%s not found, creating default, please fill this file before using this tool\n", path)
		os.Exit(-1)
	}

	defer jsonFile.Close()

	data, _ := io.ReadAll(jsonFile)
	json.Unmarshal(data, &Keys)

	if len(Keys.TmdbKeyV4) <= 0 {
		log.Printf("Invalid TmdbKeyV4, check %s file\n", path)
		os.Exit(-1)
	}
}
