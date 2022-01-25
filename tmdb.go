package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//Represent a file to download
type DownloadableItem struct {
	//web address to download
	//ex: https://google.com/...
	Address string

	//file name on disk
	//ex: image.png
	FileName string
}

type tmdbItem struct {
	FilePath string `json:"file_path"`
}

type tmdbItems *[]tmdbItem

type tmdbGetImagesAnswer struct {
	Backdrops tmdbItems `json:"backdrops"`
	Logos     tmdbItems `json:"logos"`
	Posters   tmdbItems `json:"posters"`
}

//send a request to tmdb api and then create a list of DownloadableItem
//tmdbId: tmdbid to download files
//bearerToken: tmdb api token for auth
func GetLinks(tmdbId int, bearerToken string) *[]DownloadableItem {
	//--------------------
	//creating http client
	//--------------------

	//url to fetch
	url := fmt.Sprintf("https://api.themoviedb.org/3/tv/%d/images?include_image_language=fr,en,null", tmdbId)

	//bearer token
	bearer := "Bearer " + bearerToken

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}

	//sending http request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on api call:", err)
		os.Exit(-1)
	}
	defer resp.Body.Close()

	//reading http answer
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
		os.Exit(-1)
	}

	//parsing to json
	var parsed_data tmdbGetImagesAnswer
	err = json.Unmarshal(body, &parsed_data)
	if err != nil {
		log.Println("Error parsing api data:", err)
		os.Exit(-1)
	}

	//--------------------
	//processing datas
	//--------------------

	//total number of items
	itemSize := len(*parsed_data.Backdrops) + len(*parsed_data.Logos) + len(*parsed_data.Posters)

	//array to be returned is initialised with the needed size
	var ret = make([]DownloadableItem, itemSize)

	//count between all elements
	globalIndex := 0

	var elemIndex int
	var element tmdbItem

	//list of tmdbItems
	currentList := &[]tmdbItems{
		parsed_data.Backdrops,
		parsed_data.Logos,
		parsed_data.Posters,
	}

	for bigListIndex, list := range *currentList {
		for elemIndex, element = range *list {
			var fileName string
			switch bigListIndex {
			case 0:
				fileName = "Backdrop "
			case 1:
				fileName = "Logo "
			case 2:
				fileName = "Posters "
			default:
				panic("Unknown index, check code")
			}
			ret[globalIndex] = DownloadableItem{
				Address:  "https://www.themoviedb.org/t/p/original" + element.FilePath,
				FileName: fileName + strconv.Itoa(elemIndex) + getExtensionFromFilePath(element.FilePath),
			}
			globalIndex++
		}
	}
	return &ret
}

//search and return the file extension in the string
//return "" if not found
func getExtensionFromFilePath(name string) string {
	index := strings.LastIndex(name, ".")
	if index == -1 {
		return ""
	}
	return name[index:]
}
