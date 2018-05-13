package fsfileprocessor

import (
	"os"
	"path/filepath"
)

// Crawl walks through the file system and processes files based on the Processor provided to the Crawler
func (config Crawler) Crawl() error {
	conditionChannel := make(chan WalkInfo)
	errChannel := make(chan error, 1)

	conditionFunc, validConditionChannel := config.generateConditionFunction(errChannel)

	go func() {
		defer close(conditionChannel)
		filepath.Walk(config.Controller.Rootdir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				errChannel <- err
				return err
			}

			conditionChannel <- WalkInfo{path, info}
			return nil
		})
	}()

	go conditionFunc(conditionChannel)

	go func() {
		defer close(errChannel)
		config.Processor(validConditionChannel, errChannel)
	}()

	crawlErr := <-errChannel

	if crawlErr != nil {
		return crawlErr
	}

	return nil
}
