/*

dls.go - docker ls
-John Taylor
Apr-16-2020

List files within a running docker container

To statically compile the program on Linux:
go build -tags netgo -ldflags '-extldflags "-static" -s -w'

*/

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/olekukonko/tablewriter"
)

const version = "1.0.0"

type stats struct {
	fileCount int
	dirCount  int
	errCount  int
}

var excludeDirs = map[string]int{
	"dev":  0,
	"proc": 0,
	"sys":  0,
	".git": 0,
}

func outputTable(colHeaders []string, tblData [][]string) {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader(colHeaders)
	table.SetAutoWrapText(false)
	table.AppendBulk(tblData)
	table.Render()
}

func matchesExclude(entry string) bool {
	prefix := ""
	for key, _ := range excludeDirs {

		max := len(key)
		if max > len(entry) {
			max = len(entry)
		}
		prefix = entry[:max]
		if key == prefix {
			return true
		}
	}
	return false
}

func getMetadata(dirName string, showAll bool) ([][]string, [][]string, stats) {
	var allEntries [][]string
	var allErrors [][]string
	errCount := 0
	fileCount := 0
	dirCount := 0
	fileSize := ""
	modTime := ""

	filepath.Walk(dirName,
		func(path string, info os.FileInfo, err error) error {
			if !showAll && (matchesExclude(path) || "." == path) {
				return nil
			}
			if err != nil {
				allErrors = append(allErrors, []string{fmt.Sprintf("%v", err)})
				errCount++
				return nil
			}

			fileSize = fmt.Sprintf("%9d", info.Size())
			modTime = fmt.Sprintf("%v", info.ModTime())[:19]
			objType := "F"
			if info.IsDir() {
				objType = "D"
				dirCount++
			} else {
				fileCount++
			}
			allEntries = append(allEntries, []string{fileSize, modTime, objType, path})

			return nil
		})

	myStats := new(stats)
	myStats.fileCount = fileCount
	myStats.dirCount = dirCount
	myStats.errCount = errCount

	return allEntries, allErrors, *myStats
}

func main() {
	argsShowAll := flag.Bool("a", false, "show all files, including .git, dev, proc, and sys")
	argsShowErrors := flag.Bool("e", false, "show file/directory errors")
	argsVersion := flag.Bool("v", false, "show version and then exit")
	flag.Parse()
	if *argsVersion {
		fmt.Fprintf(os.Stderr, "version %s\n", version)
		os.Exit(1)
	}

	var allEntries, allErrors [][]string
	var myStats stats
	allEntries, allErrors, myStats = getMetadata(".", *argsShowAll)

	var colHeader []string
	if *argsShowErrors && len(allErrors) > 0 {
		colHeader = []string{fmt.Sprintf("Errors: %d", myStats.errCount)}
		outputTable(colHeader, allErrors)
	}

	colHeader = []string{"Size", "Mod Time", "Type", fmt.Sprintf("Name (Files:%d Dirs:%d)", myStats.fileCount, myStats.dirCount)}
	outputTable(colHeader, allEntries)
}
