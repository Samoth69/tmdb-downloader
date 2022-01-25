package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func DownloadFiles(toDownload *[]DownloadableItem, downloadLocation string) {
	arrSize := len(*toDownload)
	bar := progressbar.NewOptions(arrSize,
		//progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		//progressbar.OptionSetWidth(15),
		//progressbar.OptionSetDescription("[cyan][1/3][reset] Writing moshable file..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	for index := range *toDownload {
		elem := &(*toDownload)[index]
		bar.Describe(fmt.Sprintf("[cyan][%d/%d][reset] Downloading %40s", index+1, arrSize, elem.FileName))
		err := downloadFile(elem.FileName, elem.Address, downloadLocation)
		if err != nil {
			fmt.Println(err)
		}
		bar.Add(1)
	}
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func downloadFile(filepath string, url string, downloadLocation string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if os.IsExist(err) {
		if !askConfirmation("File %s already exists, confirm overwrite ?", false) {
			return nil
		}
	} else if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

//ask the user the question
func askConfirmation(question string, defaultAnswer bool) (result bool) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(question)
		if defaultAnswer {
			fmt.Print(" [Y/n] ")
		} else {
			fmt.Print(" [y/N] ")
		}
		text, err := reader.ReadString('\n')

		if err != nil {
			return defaultAnswer
		}

		text = strings.ToLower(text)
		if text[0] == 'y' {
			return true
		} else if text[0] == 'n' {
			return false
		} else if text[0] == '\r' || text[0] == '\n' {
			return defaultAnswer
		}
	}
}
