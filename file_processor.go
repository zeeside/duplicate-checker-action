package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sethvargo/go-githubactions"
	"github.com/sirupsen/logrus"
)

// Encapsulate duplicates metadata for duplicates found
type FileMatch struct {
	LineNumber      int
	FileName        string
	MatchingStrings []string
}

// A wrapper for the file processor
type FileProcessor struct {
	Config                *Config
	AllMatches            map[string][]FileMatch
	DuplicateKeys         map[string]int
	Logger                *logrus.Logger
	LargeFiles            []string
	FilesChecked          int32
	InvalidExtensionFiles int
	MatchingValues        map[string]struct{}
	action                *githubactions.Action
	StartTime             time.Time
	Duration              time.Duration
	MU                    sync.RWMutex
}

// Instantiates a new file processor
func NewFileProcessor(c *Config, logger *logrus.Logger) *FileProcessor {
	logger.Debugf("Content regex is %s", c.ContentRegex)
	return &FileProcessor{
		Config:         c,
		Logger:         logger,
		AllMatches:     make(map[string][]FileMatch),
		DuplicateKeys:  make(map[string]int),
		MatchingValues: make(map[string]struct{}),
		action:         githubactions.New(),
	}
}

// Runs the core logic
func (f *FileProcessor) Run() {
	f.StartTime = time.Now()

	var wg sync.WaitGroup

	wg.Add(1)
	go f.checkDirectory(f.Config.DirectoryContext, &wg)
	wg.Wait()

	f.Duration = time.Since(f.StartTime)
}

// recursively walks the directory tree starting with the
// configured scope
func (f *FileProcessor) checkDirectory(path string, wg *sync.WaitGroup) {
	defer wg.Done()

	for k := range f.Config.IgnorePaths {
		if len(k) > 0 && strings.Contains(path, k) {
			f.Logger.Debugf("Path %s is ignore list. Skipping", path)
			return
		}
	}

	dir, err := os.Open(path)
	if err != nil {
		f.action.Fatalf("Error opening directory: %s\n", err)
		return
	}
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		f.action.Fatalf("Error reading directory: %s\n", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			wg.Add(1)
			go f.checkDirectory(fmt.Sprintf("%s/%s", path, file.Name()), wg)
		} else if filepath.Ext(file.Name()) == f.Config.FileExtension {
			_, excludeExtension := f.Config.ExcludedExtensions[filepath.Ext(file.Name())]
			_, ignore := f.Config.IgnoreFiles[file.Name()]
			if !excludeExtension && !ignore {
				f.MU.RLock()
				if len(f.AllMatches) >= f.Config.MaxFilesToProcess {
					break
				}
				f.MU.RUnlock()
				wg.Add(1)
				f.checkFile(path, file, wg)
			} else {
				f.Logger.Debugf("file %s is excluded", file.Name())
			}
		}
	}
}

// Checks the file for duplicates if it is under the max file-size threshold
func (f *FileProcessor) checkFile(directory string, baseFile os.FileInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	atomic.AddInt32(&f.FilesChecked, 1)

	path := fmt.Sprintf("%s/%s", directory, baseFile.Name())
	file, err := os.Open(path)
	if err != nil {
		f.action.Fatalf("Error opening file: %s\n", err)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		f.action.Fatalf("Error getting file info: %s\n", err)
		return
	}

	if info.Size() > int64(f.Config.MaxFileSizeKB) {
		f.LargeFiles = append(f.LargeFiles, path)
		return
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		f.action.Fatalf("Error reading file: %s\n", err)
	}

	contents := string(bytes)

	if f.Config.ContentRegex.MatchString(contents) {
		location := f.Config.ContentRegex.FindStringSubmatchIndex(contents)
		matches := f.Config.ContentRegex.FindStringSubmatch(contents)
		matchingString := matches[len(matches)-1]

		f.MU.RLock()
		matchList, isMatch := f.AllMatches[matchingString]
		f.MU.RUnlock()

		if isMatch {
			f.DuplicateKeys[matchingString] = len(matchList) + 1
		}

		f.MU.Lock()
		f.AllMatches[matchingString] = append(matchList, FileMatch{
			LineNumber:      location[0],
			FileName:        path,
			MatchingStrings: matches,
		})
		f.MU.Unlock()

	}
}

// Prints a formatted output showing
// which files contain duplicates,
// what the dupliates were and excerpts
// from the files
func (f *FileProcessor) PrintOutput() {

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Total files Scanned: %d\nDirectory Scanned: %s\nExtension Checked: %s\nDuplicated strings: %d\nFiles with Duplicates: %d\nDuration: %s\n\n",
		f.FilesChecked,
		f.Config.DirectoryContext,
		f.Config.FileExtension,
		len(f.DuplicateKeys),
		getSumOfValues(f.DuplicateKeys),
		f.Duration))

	if len(f.DuplicateKeys) > 0 {
		sb.WriteString("Caution!! Potentially harmful duplicates found in files:\n\n")
		for key := range f.DuplicateKeys {
			duplicateList := f.AllMatches[key]

			for _, i := range duplicateList {
				sb.WriteString(fmt.Sprintf("Matching String: %s\nMatch Path: %s\nMatch Location: %d\nMatching Block: %s\n...\n\n",
					key,
					i.FileName,
					i.LineNumber,
					i.MatchingStrings[0]))
			}
		}
	}

	if len(f.LargeFiles) > 0 {
		sb.WriteString("Duplicate checker skipped the following large files:\n")
		for _, lf := range f.LargeFiles {
			sb.WriteString(fmt.Sprintf("%s\n", lf))
		}
	}

	f.action.AddStepSummary(sb.String())
	f.action.SetOutput("result", sb.String())
	f.action.SetEnv("result", sb.String())
	jsonString, _ := json.Marshal(sb.String())
	f.action.SetEnv("result_escaped", string(jsonString))
	f.action.SetOutput("result_escaped", string(jsonString))

	if len(sb.String()) > 0 {
		f.action.Errorf(sb.String())
	}
}

func getSumOfValues(numMap map[string]int) int {
	returnVal := 0
	for _, value := range numMap {
		returnVal += value
	}
	return returnVal
}
