package fsfileprocessor

import (
	"os"
	"regexp"
	"time"
)

// Crawler provides both instructions and configurations to crawl the filesystem and process files
type Crawler struct {
	Processor    FileProcessor
	Conditionals []ConditionFunc
	Controller   Controller
}

//WalkInfo contains a filepath and os.FileInfo struct. It is passed through the channels of this applciation
type WalkInfo struct {
	Path string
	Info os.FileInfo
}

// FileProcessor function provides processing instructions for the crawler
type FileProcessor func(fileReceiver <-chan WalkInfo, errorChannel chan<- error) error

// ConditionFunc is a type with a required method signature matching that of the
// default filepath condition checks
type ConditionFunc func(conditionChannel chan<- bool, info WalkInfo)

// Controller exposes a set of configuration options for the crawler
type Controller struct {
	Rootdir              string
	Recursive            bool
	FileExt              *regexp.Regexp
	EarliestTimeModified time.Time
}
