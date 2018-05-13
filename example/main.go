package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/krumbot/fsfileprocessor"
)

// example FileProcessor type function. This function simply pulls incoming
// values off of the path channel and prints the path to the console
func process(fileReceiver <-chan fsfileprocessor.WalkInfo, errorChannel chan<- error) error {
	for filewalkinfo := range fileReceiver {
		fmt.Println(filewalkinfo.Path)
	}

	return nil
}

func main() {
	//In this example, we will walk through this github package's file directory recursively

	//We will only processf files with a .go or .md file extension
	fe, _ := regexp.Compile(".(go|md)")

	//We will add a custom conditional function to our Crawler. In this example, we will exclude
	//any directories containing the string "example"
	exampleConditionFunc := func(conditionChannel chan<- bool, info fsfileprocessor.WalkInfo) {
		if strings.Contains(info.Path, "example") && info.Info.IsDir() {
			fmt.Println("Failed custom conditional. Not processing for : ", info.Path)
			conditionChannel <- false
		} else {
			conditionChannel <- true
		}
	}
	exampleConditionSlice := []fsfileprocessor.ConditionFunc{exampleConditionFunc}

	//We will be recursively crawling over this github package and processing all .go and .md
	// files that have been modified after May 15, 2016
	controller := fsfileprocessor.Controller{
		Rootdir:              "../",
		Recursive:            true,
		EarliestTimeModified: time.Date(2016, time.May, 15, 0, 0, 0, 0, time.UTC),
		FileExt:              fe,
	}

	//This example sets up the configuration instructions for the crawler
	config := fsfileprocessor.Crawler{
		Processor:    process,
		Controller:   controller,
		Conditionals: exampleConditionSlice,
	}

	//We are now crawling through the file directory and concurrently processing the filepaths
	crawlErr := config.Crawl()

	if crawlErr != nil {
		fmt.Println(crawlErr)
		os.Exit(1)
	}
}
