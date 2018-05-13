package fsfileprocessor

// Config provides both instructions and configurations to crawl the filesystem and process files
type Config struct {
	Processor  FileProcessor
	Controller Controller
}

// FileProcessor function provides processing instructions for the crawler
type FileProcessor func(fileReceiver <-chan string, errorChannel chan<- error) error

// Controller exposes a set of configuration options for the crawler
type Controller struct {
	Rootdir   string
	Recursive bool
}
