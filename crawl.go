package fsfileprocessor

import (
	"os"
	"path/filepath"
)

// Crawl walks through the file system and processes files based on the Processor provided to the Config
func Crawl(config Config) error {
	pathChannel := make(chan string)
	errChannel := make(chan error, 1)
	go func() {
		defer close(pathChannel)
		filepath.Walk(config.Controller.Rootdir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				errChannel <- err
				return err
			}

			if info.IsDir() && !config.Controller.Recursive {
				return nil
			}

			pathChannel <- path
			return nil
		})
	}()

	go func() {
		defer close(errChannel)
		config.Processor(pathChannel, errChannel)
	}()

	crawlErr := <-errChannel

	if crawlErr != nil {
		return crawlErr
	}

	return nil

}
