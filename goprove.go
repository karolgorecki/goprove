// Package goprove contains lib for checking the Golang best practises
package goprove

import (
	"os"
	"strings"
	"sync"

	"github.com/fatih/structs"
)

const (
	minimumCriteria itemCategory = iota
	goodCitizen
	extraCredit
)

var (
	sourcePath   string
	sourceGoPath string
	checkList    []checkItem
)

//go:generate jsonenums -type=itemCategory
type itemCategory byte

type checkItem struct {
	Name     string       `structs:"name"`
	Desc     string       `structs:"desc"`
	Category itemCategory `structs:"category"`
	fn       func() bool
}

func (ci checkItem) run() (success bool) {
	return ci.fn()
}

func init() {
	checkList = []checkItem{
		{
			Name:     "projectBuilds",
			Category: minimumCriteria,
			Desc:     "Compiles: Does the project build?",
			fn:       projectBuilds,
		},
		{
			Name:     "isFormatted",
			Category: minimumCriteria,
			Desc:     "gofmt Correctness: Is the code formatted correctly?",
			fn:       isFormatted,
		},
		{
			Name:     "hasLicense",
			Category: minimumCriteria,
			Desc:     "Licensed: Does the project have a license?",
			fn:       hasLicense,
		},
		{
			Name:     "isLinted",
			Category: minimumCriteria,
			Desc:     "golint Correctness: Is the linter satisfied?",
			fn:       isLinted,
		},
		{
			Name:     "isVetted",
			Category: minimumCriteria,
			Desc:     "go tool vet Correctness: Is the Go vet satisfied?",
			fn:       isVetted,
		},
		{
			Name:     "hasReadme",
			Category: minimumCriteria,
			Desc:     "README Presence: Does the project's include a documentation entrypoint?",
			fn:       hasReadme,
		},
		{
			Name:     "testPassing",
			Category: minimumCriteria,
			Desc:     "Are the tests passing?",
			fn:       testPassing,
		},
		{
			Name:     "isDirMatch",
			Category: minimumCriteria,
			Desc:     "Directory Names and Packages Match: Does each package <pkg> statement's package name match the containing directory name?",
			fn:       isDirMatch,
		},
		{
			Name:     "hasContributing",
			Category: goodCitizen,
			Desc:     "Contribution Process: Does the project document a contribution process?",
			fn:       hasContributing,
		},
		{
			Name:     "hasBenches",
			Category: extraCredit,
			Desc:     "Benchmarks: In addition to tests, does the project have benchmarks?",
			fn:       hasBenches,
		},
		{
			Name:     "hasBlackboxTests",
			Category: extraCredit,
			Desc:     "Blackbox Tests: In addition to standard tests, does the project have blackbox tests?",
			fn:       hasBlackboxTests,
		},
	}
}

// RunTasks is a wrapper for running all tasks from the list
func RunTasks(path string, tasksToExlude []string) (successTasks []map[string]interface{}, failedTasks []map[string]interface{}) {
	var wg sync.WaitGroup
	sourcePath = path
	sourceGoPath = strings.Replace(sourcePath, os.Getenv("GOPATH")+"/src/", "", 1)

	excludeTasks(&checkList, tasksToExlude)

	wg.Add(len(checkList))
	for _, task := range checkList {
		go func(task checkItem) {
			if ok := task.run(); ok {
				successTasks = append(successTasks, structs.Map(task))
			} else {
				failedTasks = append(failedTasks, structs.Map(task))
			}
			wg.Done()
		}(task)
	}

	wg.Wait()
	return successTasks, failedTasks
}

func excludeTasks(chL *[]checkItem, tskE []string) {
	tempChL := *chL
	if !hasBuildPackage() {
		for idx, item := range tempChL {
			if item.Name == "projectBuilds" {
				tempChL = append(tempChL[:idx], tempChL[idx+1:]...)
			}
		}
	}

	if len(tskE) > 0 {
		for _, tskExc := range tskE {
			for idx, item := range tempChL {
				if item.Name == tskExc {
					tempChL = append(tempChL[:idx], tempChL[idx+1:]...)
				}
			}
		}
	}
	*chL = tempChL
}
