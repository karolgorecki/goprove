// Goprove is the command line tool
// for checking if your project follows
// the good practises
package main

import (
	"fmt"

	"github.com/karolgorecki/goprove/checkList"

	"github.com/karolgorecki/goprove/util"
)

func main() {
	// Print the Goprove banner
	util.PrintBanner()

	// Start the benchmark
	util.BenchmarkStart()

	// Create 2 arrays with passed/failed checklist items
	okTasks, nokTasks := checkList.RunTasks()

	// Print a list of completed & failed tasks
	fmt.Println("Passed tests:", len(okTasks), "of", len(okTasks)+len(nokTasks))
	fmt.Println("---------------------------------------------------------------")
	for _, taskDone := range okTasks {
		fmt.Println(taskDone)
	}

	fmt.Println("---------------------------------------------------------------")
	for _, taskFailed := range nokTasks {
		fmt.Println(taskFailed)
	}

	// Write the execution time
	util.ExecutionTime()
}
