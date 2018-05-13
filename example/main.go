package main

import (
	"fmt"
	"os"

	"github.com/krumbot/fsfileprocessor"
)

// example FileProcessor type function. This function simply pulls incoming
// values off of the path channel and prints the path to the console
func process(fileReceiver <-chan string, errorChannel chan<- error) error {
	for filepath := range fileReceiver {
		fmt.Println(filepath)
	}

	return nil
}

func main() {
	//In this example, we will walk through this github package's file directory recursively
	controller := fsfileprocessor.Controller{
		Rootdir:   "../",
		Recursive: true,
	}

	//This example sets up the configuration instructions for the crawler
	config := fsfileprocessor.Config{
		Processor:  process,
		Controller: controller,
	}

	//We are now crawling through the file directory and concurrently processing the filepaths
	crawlErr := fsfileprocessor.Crawl(config)

	if crawlErr != nil {
		fmt.Println(crawlErr)
		os.Exit(1)
	}
}
