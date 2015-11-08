package checklist

import (
	"log"
	"os/exec"
	"regexp"
	"sync"

	"github.com/karolgorecki/goprove/util"
)

var wg sync.WaitGroup

// checkItem is the basic sturct for the test cases
type checkItem struct {
	name     string
	function func(string) (message string, success bool)
}

func get() (checkList map[string]checkItem) {

	checkList = map[string]checkItem{
		"builds": {
			name:     "Compiles: Does the project build?",
			function: builds,
		},
		"isFormated": {
			name:     "gofmt Correctness: Is the code formatted correctly?",
			function: isFormatted,
		},
		"hasLicense": {
			name:     "Licensed: Does the project have a license?",
			function: hasLicense,
		},
		"isLinted": {
			name:     "golint Correctness: Is the linter satisfied?",
			function: isLinted,
		},
		"isVet": {
			name:     "go tool vet Correctness: Is the Go vet satisfied?",
			function: isVet,
		},
		"hasReadme": {
			name:     "README Presence: Does the project's include a documentation entrypoint?",
			function: hasReadme,
		},
		"hasContribution": {
			name:     "Contribution Process: Does the project document a contribution process?",
			function: hasContribution,
		},
		"testPassing": {
			name:     "Are the tests passing?",
			function: testPassing,
		},
	}

	return checkList
}

func (checkItem checkItem) Run() (message string, success bool) {
	return checkItem.function(checkItem.name)
}

// RunTasks is a wrapper for running all tasks from the list
func RunTasks() (successTasks []string, failedTasks []string) {
	tasks := get()

	wg.Add(len(tasks))
	for _, task := range tasks {

		go func(task checkItem) {
			desc, isSuccess := task.Run()

			if isSuccess {
				successTasks = append(successTasks, desc)

			} else {
				failedTasks = append(failedTasks, desc)
			}

			wg.Done()

		}(task)
	}

	wg.Wait()

	return successTasks, failedTasks
}

// -----------------------------------------------------------------------------
// CHECKLIST FUNCTIONS
// -----------------------------------------------------------------------------
func builds(taskName string) (message string, success bool) {
	_, err := exec.Command("go", "build").Output()

	if err != nil {
		return util.GetFailMessage(taskName), false
	}
	return util.GetSuccessMessage(taskName), true
}

func isFormatted(taskName string) (message string, success bool) {
	output, _ := exec.Command("gofmt", "-l", ".").Output()

	if len(output) > 0 {
		return util.GetFailMessage(taskName), false
	}

	return util.GetSuccessMessage(taskName), true
}

func testPassing(taskName string) (message string, success bool) {
	output, _ := exec.Command("go", "test", "./...").Output()

	if testFails, _ := regexp.Match(`--- FAIL`, output); testFails {
		return util.GetFailMessage(taskName), false
	}

	return util.GetSuccessMessage(taskName), true
}

func hasLicense(taskName string) (message string, success bool) {
	hasLicense, err := util.FileExists("license", "licensing")

	if err != nil {
		log.Fatal(err)
	}

	if hasLicense {
		return util.GetSuccessMessage(taskName), true
	}

	return util.GetFailMessage(taskName), false
}

func hasReadme(taskName string) (message string, success bool) {
	hasReadme, err := util.FileExists("readme")
	if err != nil {
		log.Fatal(err)
	}

	if hasReadme {
		return util.GetSuccessMessage(taskName), true
	}

	return util.GetFailMessage(taskName), false
}

func hasContribution(taskName string) (message string, success bool) {
	hasContribution, err := util.FileExists("contribution", "contribute", "contributing")
	if err != nil {
		log.Fatal(err)
	}

	if hasContribution {
		return util.GetSuccessMessage(taskName), true
	}

	return util.GetFailMessage(taskName), false
}

func isLinted(taskName string) (string, bool) {

	output, _ := exec.Command("golint").Output()

	if len(output) > 0 {
		return util.GetFailMessage(taskName), false
	}
	return util.GetSuccessMessage(taskName), true
}

func isVet(taskName string) (message string, success bool) {
	_, err := exec.Command("go", "vet").Output()

	if err != nil {
		return util.GetFailMessage(taskName), false
	}
	return util.GetSuccessMessage(taskName), true
}
